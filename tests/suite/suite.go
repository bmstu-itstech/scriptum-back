package suite

import (
	"context"
	"math/rand/v2"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/api"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/mock"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server/auth"
)

type Suite struct {
	BoxService  apiv2.BoxServiceClient
	FileService apiv2.FileServiceClient
	JobService  apiv2.JobServiceClient
}

func suiteConfig() config.Config {
	return config.Config{
		Docker: config.Docker{
			ImagePrefix:   "sc-box",
			RunnerTimeout: time.Minute,
		},
		GRPC: config.GRPC{}, // Не используется
		Logging: config.Logging{
			Level: "debug",
		},
		Postgres: config.Postgres{
			URI: os.Getenv("POSTGRES_URI"),
		},
		Storage: config.Storage{
			BasePath: "../uploads",
		},
	}
}

func New(t *testing.T) (context.Context, *Suite) {
	cfg := suiteConfig()
	l := logs.NewLogger(cfg.Logging)

	repos := postgres.MustNewRepository(cfg.Postgres, l)
	runner := docker.MustNewRunner(cfg.Docker, l)
	storage := local.MustNewStorage(cfg.Storage, l)
	mockIAP := mock.NewIsAdminProvider()

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	infra := app.Infra{
		BoxProvider:     repos,
		BoxRepo:         repos,
		FileReader:      storage,
		FileUploader:    storage,
		IsAdminProvider: mockIAP,
		JobProvider:     repos,
		JobPublisher:    jPub,
		JobRepository:   repos,
		Runner:          runner,
	}
	a := app.NewApp(infra, l)

	ctx := context.Background()

	go func() {
		err := jSub.Listen(ctx, func(ctx2 context.Context, jobID string) error {
			return a.Commands.RunJob.Handle(ctx2, request.RunJob{JobID: jobID})
		})
		if err != nil {
			t.Error(err)
		}
	}()

	lis := bufconn.Listen(1024 * 1024) // 1 Mb
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		auth.UnaryServerInterceptor(),
	))

	api.RegisterBoxService(s, a, l)
	api.RegisterFileService(s, a, l)
	api.RegisterJobService(s, a, l)

	go func() {
		err := s.Serve(lis)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	boxClient := apiv2.NewBoxServiceClient(conn)
	fileClient := apiv2.NewFileServiceClient(conn)
	jobClient := apiv2.NewJobServiceClient(conn)

	uid := rand.Int64() //nolint:gosec // В рамках тестов допустимо использование такого рода генерации чисел
	ctx = auth.ClientOutgoingContext(ctx, uid)

	return ctx, &Suite{
		BoxService:  boxClient,
		FileService: fileClient,
		JobService:  jobClient,
	}
}

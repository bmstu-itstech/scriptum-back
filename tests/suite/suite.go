package suite

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/api"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/command"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/query"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/mock"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server/auth"
)

const runnerTimeout = 15 * time.Second

type Suite struct {
	BoxService  apiv2.BoxServiceClient
	FileService apiv2.FileServiceClient
	JobService  apiv2.JobServiceClient
}

func connectDB() (*sqlx.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		return nil, errors.New("DATABASE_URI must be set")
	}
	return sqlx.Connect("postgres", uri)
}

func New(t *testing.T) (context.Context, *Suite) {
	l := logs.DefaultLogger()

	db, err := connectDB()
	if err != nil {
		l.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repos := postgres.NewRepository(db, l)
	runner := docker.MustNewRunner(l)
	storage := local.NewStorage("../uploads", l)
	mockIAP := mock.NewIsAdminProvider()

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	a := &app.App{
		Commands: app.Commands{
			CreateBox:  command.NewCreateBoxHandler(repos, mockIAP, l),
			DeleteBox:  command.NewDeleteBoxHandler(repos, l),
			RunJob:     command.NewRunJobHandler(runner, repos, storage, l),
			StartJob:   command.NewStartJobHandler(repos, repos, jPub, l),
			UploadFile: command.NewUploadFileHandler(storage, l),
		},
		Queries: app.Queries{
			GetBox:      query.NewGetBoxHandler(repos, l),
			GetBoxes:    query.NewGetBoxesHandler(repos, l),
			GetJob:      query.NewGetJobHandler(repos, l),
			GetJobs:     query.NewGetJobsHandler(repos, l),
			SearchBoxes: query.NewSearchBoxesHandler(repos, l),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), runnerTimeout)
	t.Cleanup(cancel)

	go func() {
		err2 := jSub.Listen(ctx, func(ctx2 context.Context, jobID string) error {
			ctx2, cancel2 := context.WithTimeout(ctx2, runnerTimeout)
			defer cancel2()
			return a.Commands.RunJob.Handle(ctx2, request.RunJob{JobID: jobID})
		})
		if err != nil {
			t.Error(err2)
		}
	}()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		auth.UnaryServerInterceptor(),
	))

	api.RegisterBoxService(s, a, l)
	api.RegisterFileService(s, a, l)
	api.RegisterJobService(s, a, l)

	go func() {
		err = s.Serve(lis)
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
	md := metadata.New(map[string]string{"x-user-id": strconv.FormatInt(uid, 10)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx, &Suite{
		BoxService:  boxClient,
		FileService: fileClient,
		JobService:  jobClient,
	}
}

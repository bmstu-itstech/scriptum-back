package api

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/pkg/server/auth"
)

type jobService struct {
	apiv2.UnimplementedJobServiceServer

	app *app.App
	l   *slog.Logger
}

func RegisterJobService(s *grpc.Server, app *app.App, l *slog.Logger) {
	apiv2.RegisterJobServiceServer(s, &jobService{app: app, l: l})
}

func (s jobService) GetJob(ctx context.Context, req *apiv2.GetJobRequest) (*apiv2.GetJobResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	job, err := s.app.Queries.GetJob.Handle(ctx, request.GetJob{
		UID:   uid,
		JobID: req.JobId,
	})
	switch {
	case errors.Is(err, ports.ErrJobNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrPermissionDenied):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.GetJobResponse{Job: jobToAPI(job)}, nil
}

func (s jobService) GetJobs(ctx context.Context, req *apiv2.GetJobsRequest) (*apiv2.GetJobsResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	jobs, err := s.app.Queries.GetJobs.Handle(ctx, request.GetJobs{
		UID:   uid,
		State: jobStateFromAPIOrNil(req.State),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.GetJobsResponse{Jobs: jobsToAPI(jobs)}, nil
}

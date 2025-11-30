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

type boxService struct {
	apiv2.UnimplementedBoxServiceServer

	app *app.App
	l   *slog.Logger
}

func RegisterBoxService(s *grpc.Server, app *app.App, l *slog.Logger) {
	apiv2.RegisterBoxServiceServer(s, &boxService{app: app, l: l})
}

func (s boxService) CreateBox(ctx context.Context, req *apiv2.CreateBoxRequest) (*apiv2.CreateBoxResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	r, err := createBoxRequestFromAPI(req, uid)
	var iiErr domain.InvalidInputError
	if errors.As(err, &iiErr) {
		return nil, statusInvalidInput(iiErr)
	} else if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	boxID, err := s.app.Commands.CreateBox.Handle(ctx, r)
	if errors.As(err, &iiErr) {
		return nil, statusInvalidInput(iiErr)
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.CreateBoxResponse{BoxId: boxID}, nil
}

func (s boxService) DeleteBox(ctx context.Context, req *apiv2.DeleteBoxRequest) (*apiv2.DeleteBoxResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	err := s.app.Commands.DeleteBox.Handle(ctx, request.DeleteBox{
		UID:   uid,
		BoxID: req.BoxId,
	})
	switch {
	case errors.Is(err, ports.ErrBoxNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrPermissionDenied):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.DeleteBoxResponse{}, nil
}

func (s boxService) GetBox(ctx context.Context, req *apiv2.GetBoxRequest) (*apiv2.GetBoxResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	box, err := s.app.Queries.GetBox.Handle(ctx, request.GetBox{
		UID:   uid,
		BoxID: req.BoxId,
	})
	switch {
	case errors.Is(err, ports.ErrBoxNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrPermissionDenied):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.GetBoxResponse{Box: boxToAPI(box)}, nil
}

func (s boxService) GetBoxes(ctx context.Context, _ *apiv2.GetBoxesRequest) (*apiv2.GetBoxesResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	boxes, err := s.app.Queries.GetBoxes.Handle(ctx, request.GetBoxes{UID: uid})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.GetBoxesResponse{Boxes: boxesToAPI(boxes)}, nil
}

func (s boxService) SearchBoxes(
	ctx context.Context, req *apiv2.SearchBoxesRequest,
) (*apiv2.SearchBoxesResponse, error) {
	uid, ok := auth.ExtractUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "'x-user-id' header is missing")
	}

	boxes, err := s.app.Queries.SearchBoxes.Handle(ctx, request.SearchBoxes{
		UID:  uid,
		Name: req.Name,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv2.SearchBoxesResponse{Boxes: boxesToAPI(boxes)}, nil
}

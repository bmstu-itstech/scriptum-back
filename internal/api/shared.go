package api

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

func statusInvalidInput(iiErr domain.InvalidInputError) error {
	st := status.New(codes.InvalidArgument, iiErr.Message)
	errInfo := &errdetails.ErrorInfo{
		Reason: iiErr.Code,
	}
	st, err := st.WithDetails(errInfo)
	if err != nil {
		return fmt.Errorf("st.WithDetails: %w", err)
	}
	return st.Err()
}

// internal/worker/launch_handler.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type LaunchHandler struct {
	usecase *LaunchUC
}

func NewLaunchHandler(usecase *LaunchUC) *LaunchHandler {
	return &LaunchHandler{usecase: usecase}
}

func (h *LaunchHandler) Handler(msg *message.Message) error {
	var req scripts.LaunchRequest
	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		fmt.Printf("failed to decode launch request: %v\n", err)
		return err
	}

	ctx := context.Background()

	if err := h.usecase.ProcessLaunchRequest(ctx, req); err != nil {
		fmt.Printf("failed to process launch request: %v\n", err)
		return err
	}

	return nil
}

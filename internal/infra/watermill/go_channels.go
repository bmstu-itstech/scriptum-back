package watermill

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/bmstu-itstech/scriptum-back/pkg/logs/sl"
)

func NewJobPubSubGoChannels(l *slog.Logger) (JobPublisher, JobSubscriber) {
	logger := sl.NewWatermillLoggerAdapter(l)
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	return NewJobPublisher(pubSub), NewJobSubscriber(pubSub, l)
}

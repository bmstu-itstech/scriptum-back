package main

// import (
// 	"context"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"

// 	httpapi "github.com/bmstu-itstech/scriptum-back/internal/api/http"
// 	"github.com/bmstu-itstech/scriptum-back/internal/app"
// 	"github.com/bmstu-itstech/scriptum-back/internal/service"
// 	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
// 	"github.com/bmstu-itstech/scriptum-back/pkg/server"
// )

// func main() {
// 	l := logs.DefaultLogger()
// 	context := context.Background()

// 	emailNotifier, err := service.NewEmailNotifier()
// 	jobRepo, err := service.NewJobRepo(context)
// 	pythonLauncher, err := service.PythonLauncher("")

// 	botRepository := service.NewPgBotsRepository(db)
// 	participantRepository := service.NewPgParticipantsRepository(db)
// 	botMessageSender := service.NewTelegramMessageSender()

// 	process := ProcessHandlerAdapter{app.NewProcessHandler(botRepository, participantRepository, botMessageSender, l, mc)}
// 	entry := EntryHandlerAdapter{app.NewEntryHandler(botRepository, participantRepository, botMessageSender, l, mc)}

// 	telegramService := service.NewTelegramService(l, process, entry)

// 	a := app.Application{
// 		Commands: app.Commands{
// 			CreateBot:     app.NewCreateBotHandler(botRepository, l, mc),
// 			DeleteBot:     app.NewDeleteBotHandler(botRepository, telegramService, l, mc),
// 			StartBot:      app.NewStartBotHandler(botRepository, telegramService, l, mc),
// 			StopBot:       app.NewStopBotHandler(botRepository, telegramService, l, mc),
// 			UpdateStatus:  app.NewUpdateStatusHandler(botRepository, l, mc),
// 			Entry:         app.NewEntryHandler(botRepository, participantRepository, botMessageSender, l, mc),
// 			Process:       app.NewProcessHandler(botRepository, participantRepository, botMessageSender, l, mc),
// 			CreateMailing: app.NewCreateMailingHandler(botRepository, l, mc),
// 			StartMailing:  app.NewStartMailingHandler(botRepository, participantRepository, botMessageSender, l, mc),
// 		},
// 		Queries: app.Queries{
// 			AllAnswers:  app.NewGetAnswersTableHandler(botRepository, participantRepository, l, mc),
// 			GetBot:      app.NewGetBotHandler(botRepository, l, mc),
// 			GetBots:     app.NewGetBotsHandler(botRepository, l, mc),
// 			StartedBots: app.NewGetStartedBotsHandler(botRepository, l, mc),
// 		},
// 	}

// 	server.RunHTTPServer(func(router chi.Router) http.Handler {
// 		return httpapi.HandlerFromMux(httpapi.NewHTTPServer(&a), router)
// 	})
// }

// type ProcessHandlerAdapter struct {
// 	H app.ProcessHandler
// }

// func (a ProcessHandlerAdapter) Process(ctx context.Context, botId string, userId int64, msg bots.Message) error {
// 	return a.H.Handle(ctx, app.Process{
// 		BotUUID: botId,
// 		UserID:  userId,
// 		Text:    msg.Text,
// 	})
// }

// type EntryHandlerAdapter struct {
// 	H app.EntryHandler
// }

// func (a EntryHandlerAdapter) Entry(ctx context.Context, botId string, userId int64, key string) error {
// 	return a.H.Handle(ctx, app.Entry{
// 		BotUUID: botId,
// 		UserID:  userId,
// 		Key:     key,
// 	})
// }

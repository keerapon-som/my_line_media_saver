package service

// import (
// 	"context"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"bot_message_collector/api"
// 	"bot_message_collector/config"
// 	httpCatchup "bot_message_collector/http"
// 	"bot_message_collector/repository"

// 	"github.com/gofiber/fiber/v2"

// 	"golang.org/x/sync/errgroup"
// )

// type Service struct {
// 	fiberApp *fiber.App
// }

// func New(LineWebhook *api.LineWebhookService, jsonArchive *repository.LineJsonfileArchive) *Service {

// 	return &Service{

// 		fiberApp: httpCatchup.NewHTTPRouter(LineWebhook, jsonArchive),
// 	}

// }

// func (s *Service) Run(ctx context.Context) error {

// 	errgroup, ctx := errgroup.WithContext(ctx)

// 	wg := sync.WaitGroup{}

// 	wg.Wait()

// 	errgroup.Go(func() error {

// 		err := s.fiberApp.Listen(":" + config.GetConfig().ServerConfig.HTTP.Port)
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			log.Printf("Failed to start HTTP server: %v", err)
// 			return err
// 		}

// 		return nil
// 	})

// 	errgroup.Go(func() error {
// 		<-ctx.Done()
// 		log.Printf("Received shutdown signal, shutting down HTTP server...")

// 		return s.fiberApp.Shutdown()
// 	})

// 	return errgroup.Wait()
// }

package commands

import (
	"context"
	"errors"
	"github.com/akbariandev/pacassistant/internal/app/bot"
	"github.com/akbariandev/pacassistant/pkg/telegram"
	"log"

	"github.com/akbariandev/pacassistant/transport/client"

	"github.com/akbariandev/pacassistant/config"
	"github.com/akbariandev/pacassistant/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run bot service",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New[config.ExtraData](configPath)
		if err != nil {
			return err
		}

		logging, err := defaultLogging()
		if cfg.Logging != nil {
			logOpt := logger.Options{
				Development:  cfg.Development,
				Debug:        cfg.Logging.Debug,
				EnableCaller: cfg.Logging.EnableCaller,
				SkipCaller:   3,
			}

			if len(cfg.Logging.SentryDSN) != 0 {
				logOpt.Sentry = &logger.SentryConfig{
					DSN:              cfg.Logging.SentryDSN,
					AttachStacktrace: true,
					Environment:      logger.DEVELOPMENT,
					EnableTracing:    true,
					Debug:            true,
					TracesSampleRate: 1.0,
				}
				if !cfg.Development {
					logOpt.Sentry.Environment = logger.PRODUCTION
				}
			}

			logging, err = logger.New(logger.HandleType(cfg.Logging.Handler), logOpt)
		}
		if err != nil {
			return err
		}

		if len(cfg.GrpcClients) == 0 {
			return errors.New("failed to get pactus grpc client form gprc clients")
		}

		pactusClient, err := client.NewPactusClient(cfg.GrpcClients[0])
		if err != nil {
			return err
		}

		telegramBot := telegram.NewTelegramBot(cfg.Telegram.ApiKey)

		application, err := bot.New(cmd.Context(), pactusClient, telegramBot, cfg, logging)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.WithValue(cmd.Context(), "development", cfg.Development)
		application.Run(ctx)

		return nil
	},
}

func defaultLogging() (logger.Logger, error) {
	return logger.New(logger.ConsoleHandler, logger.Options{
		Development:  false,
		Debug:        false,
		EnableCaller: true,
		SkipCaller:   3,
	})
}

package commands

import (
	"log"

	"github.com/akbariandev/pacassistant/internal/app/bot"

	"github.com/akbariandev/pacassistant/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("start database migration.")
		log.Println("loading config file...")
		cfg, err := config.New[config.ExtraData](configPath)
		if err != nil {
			return err
		}
		log.Println("config has been loaded.")

		log.Println("creating new application...")
		application, err := bot.New(cmd.Context(),
			nil,
			nil,
			cfg,
			nil,
		)
		if err != nil {
			return err
		}
		log.Println("application created.")

		log.Println("starting migrate database...")
		if err := application.Migration(cmd.Context()); err != nil {
			return err
		}

		log.Println("successfully migrate done.")

		return nil
	},
}

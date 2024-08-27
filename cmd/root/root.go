package root

import (
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "MyApp is a REST API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		cfg, err := config.NewConfig()
		if err != nil {
			log.Err(err).Msg("Error loading configuration")
			return
		}

		// TODO: start server here
		log.Info().Interface("options", cfg).Msg("Server started")
	},
}

func init() {
	// Define command-line flags and bind them to Viper
	rootCmd.PersistentFlags().String("environment", "", "Application environment (e.g., production, staging, qa, development)")
	rootCmd.PersistentFlags().String("hostname", "", "Hostname of the application")
	rootCmd.PersistentFlags().String("port", "", "Port for the application")

	// Bind the flags to Viper
	viper.BindPFlag("ENVIRONMENT", rootCmd.PersistentFlags().Lookup("environment"))
	viper.BindPFlag("HOSTNAME", rootCmd.PersistentFlags().Lookup("hostname"))
	viper.BindPFlag("PORT", rootCmd.PersistentFlags().Lookup("port"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("Error executing cmd")
		// Handle errors appropriately in real applications
	}
}

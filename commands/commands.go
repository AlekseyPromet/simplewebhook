package commands

import (
	"AlekseyPromet/examples/simplewebhook/app"
	"AlekseyPromet/examples/simplewebhook/models"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	// Used for flags.
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "sw-cli [options]",
		Short: "Simple webhook server",
		Long:  "Complete documentation is available at https://github.com/AlekseyPromet/simplewebhook/blob/main/README.md",
		RunE: func(cmd *cobra.Command, args []string) error {

			config := &models.Config{
				Port:    viper.GetString("port"),
				Verbose: viper.GetBool("verbose"),
				Debug:   viper.GetBool("debug"),
			}

			service, err := app.NewService(config)
			if err != nil {
				return err
			}

			fx.New(
				fx.Provide(service.Run),
				fx.Invoke(func(*http.Server) {}),
				// fx.WithLogger(service.GetFxLogger),
			).Run()

			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("port", "p", "8088", "servise start on port")
	rootCmd.PersistentFlags().BoolP("verbose", "v", true, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", true, "verbose and debug output")
	viper.RegisterAlias("f", "config")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("author", "Aleksey Promet <promet.alex@gmail.com>")
	viper.SetDefault("license", "GPL-3.0")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name "sw-api.yaml"
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("sw-api.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

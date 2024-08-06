package cli

import (
	"os"
	"path"
	"strings"

	"github.com/ItsMyEyes/my-tools/internal/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logger.Get()

var (
	cfgFile string
	rootCmd = cobra.Command{
		Use:   "kiyora",
		Short: "kiyora - kiyora server cli agent",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig()
		},
	}
)

func init() {
	// initialize default values

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		`Config file path. If not given, it will look up the config file at following location:
		- /etc/kiyora/config.yml
		- $HOME/.kiyora/config.yaml
		- config.yaml in current directory (./config.yaml)`)

	rootCmd.AddCommand(NewBackupCLI())

}

// ExecuteRootCmd prepares all CLI commands
func ExecuteRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

// mkInitializer will create a Cobra initializer function given a `cfgFile`
func initConfig() {
	// Set config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath("/etc/kiyora")
		viper.AddConfigPath(path.Join(home, ".kiyora"))
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("kiyora")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msg("Config file used for kiyora: " + viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}
}

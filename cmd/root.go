package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Locale   string `mapstructure:"LOCALE"`
	Brand    string `mapstructure:"BRAND"`
	Model    string `mapstructure:"MODEL"`
	OSPF     string `mapstructure:"OSPF"`
	TermID   string `mapstructure:"TERM_ID"`
	TermName string `mapstructure:"TERM_NAME"`
	Token    string `mapstructure:"TOKEN"`
}

var rootCmd = &cobra.Command{}

func init() {
	cobra.OnInitialize(onInitialize)

	rootCmd.PersistentFlags().String("token", "", "Authentication token")
	rootCmd.PersistentFlags().String("term-id", "", "Terminal id")

	viper.BindPFlag("TOKEN", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("TERM_ID", rootCmd.PersistentFlags().Lookup("term-id"))

	viper.SetDefault("LOCALE", "en_US")
	viper.SetDefault("BRAND", "TP-Link")
	viper.SetDefault("MODEL", "Pixel 7")
	viper.SetDefault("OSPF", "Android 14")
}

func onInitialize() {

	homeDir, _ := os.UserHomeDir()

	viper.AddConfigPath(homeDir)
	viper.AddConfigPath(".")
	viper.SetConfigFile(".tapo-tools")
	viper.SetConfigType("dotenv")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

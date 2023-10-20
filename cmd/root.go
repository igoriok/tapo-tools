package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{}

func init() {
	cobra.OnInitialize(onInitialize)

	rootCmd.PersistentFlags().String("token", "", "Tapo Care token")
	rootCmd.PersistentFlags().String("term-id", "", "Tapo Care term id")

	viper.BindPFlag("TOKEN", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("TERM_ID", rootCmd.PersistentFlags().Lookup("term-id"))
}

func onInitialize() {

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
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

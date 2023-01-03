/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tieubaoca/go-chat-server/app"
	"github.com/tieubaoca/go-chat-server/services"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if godotenv.Load() != nil {
			log.Fatal("Error loading .env file")
		}
		services.InitDbClient(
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT"),
			os.Getenv("MONGO_USERNAME"),
			os.Getenv("MONGO_PASSWORD"),
			os.Getenv("MONGO_DB"))
		services.InitWebSocket()
		app.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
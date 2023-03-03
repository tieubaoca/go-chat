/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tieubaoca/go-chat-server/app"
	"github.com/tieubaoca/go-chat-server/utils/log"
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
		tls, _ := cmd.Flags().GetBool("security")
		app := app.NewApp()
		if tls {
			app.RunTLS()
		} else {
			app.Run()
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().BoolP("security", "s", false, "Run with TLS")

	log.InfoLogger.Println("Starting app")

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.New(io.MultiWriter(os.Stdout, file))
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	if godotenv.Load() != nil {
		log.FatalLogger.Fatal("Error loading .env file")
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

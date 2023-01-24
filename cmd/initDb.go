/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tieubaoca/go-chat-server/services"
)

// initDbCmd represents the initDb command
var initDbCmd = &cobra.Command{
	Use:   "initDb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.InitDbClient(
			os.Getenv("MONGO_CONNECTION_STRING"),
			os.Getenv("MONGO_DB"),
		)

		services.InitCollections()
	},
}

func init() {
	rootCmd.AddCommand(initDbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initDbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initDbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

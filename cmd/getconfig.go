/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"

	"github.com/cory-evans/what-did-i-work-on/config"
	"github.com/spf13/cobra"
)

// getconfigCmd represents the getconfig command
var getconfigCmd = &cobra.Command{
	Use:   "getconfig",
	Short: "prints the current config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig()

		if err != nil {
			return
		}

		output, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			return
		}

		cmd.Print(string(output))

	},
}

func init() {
	rootCmd.AddCommand(getconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

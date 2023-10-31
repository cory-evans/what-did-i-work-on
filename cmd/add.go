/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/cory-evans/what-did-i-work-on/config"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [directories]",
	Short: "Add a directory to the config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		maxSearchDepth, err := cmd.Flags().GetInt("max-search-depth")
		if err != nil {
			fmt.Println(err)
			return
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, p := range args {
			// get absolute path
			if !path.IsAbs(p) {
				p = path.Join(cwd, path.Clean(p))
			}

			cfg.AddPath(p, maxSearchDepth)
		}

		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.
	// addCmd.PersistentFlags().StringP("path", "p", "", "Path to add to config")
	addCmd.PersistentFlags().IntP("max-search-depth", "d", 0, "Max search depth for path")
}

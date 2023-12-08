package cmd

import (
	"fmt"
	"strconv"

	"github.com/cory-evans/what-did-i-work-on/config"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [number(s)]",
	Short: "remove a directory from the config",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()

		if err != nil {
			return
		}

		for _, nStr := range args {
			n, err := strconv.Atoi(nStr)
			if err != nil {
				continue
			}

			cfg.RemoveNumber(n)
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
	rootCmd.AddCommand(removeCmd)
}

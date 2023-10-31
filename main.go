package main

import (
	"log/slog"
	"os"

	"github.com/cory-evans/what-did-i-work-on/common"
	"github.com/cory-evans/what-did-i-work-on/config"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var logger *slog.Logger
var cfg *config.Config

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	c, err := config.LoadConfig()
	cfg = c
	common.CheckError(err)
}

func main() {
	defer config.SaveConfig(cfg)

	r, err := git.PlainOpen("./")
	common.CheckError(err)

	author, err := common.GetAuthorName(r)
	common.CheckError(err)

	cIter, err := r.Log(&git.LogOptions{
		All: true,
	})
	common.CheckError(err)

	commits, err := common.GetCommitsByAuthor(cIter, author)
	common.CheckError(err)

	for _, c := range commits {
		logger.Info("commit", "date", c.Author.When.Format("2006-01-02 15:04:05"), "message", c.Message)
	}

	home, err := os.UserHomeDir()
	common.CheckError(err)

	logger.Info("user home dir", "dir", home)

	configDir, err := os.UserConfigDir()
	common.CheckError(err)

	logger.Info("user config dir", "dir", configDir)
}

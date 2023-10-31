/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cory-evans/what-did-i-work-on/common"
	"github.com/cory-evans/what-did-i-work-on/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "what-did-i-work-on",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig()

		if err != nil {
			return
		}

		var commits []*CommitForRepo
		var headsChecked []*plumbing.Reference

		for _, d := range conf.Directories {

			absDepth := strings.Count(d.Path, string(os.PathSeparator))
			filepath.WalkDir(d.Path, func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					return nil
				}

				// get depth to see if we should skip
				depth := strings.Count(path, string(os.PathSeparator))
				if (depth - absDepth) > (d.MaxSearchDepth + 1) {
					return filepath.SkipDir
				}

				// skip these
				toSkip := []string{"node_modules", "vendor"}
				for _, skip := range toSkip {
					if info.Name() == skip {
						return filepath.SkipDir
					}
				}

				if info.Name() == ".git" {
					r, headRef, err := openRepoAndGetHeadRef(path)

					// cl, headRef, err := getCommits(path)
					if err != nil {
						return filepath.SkipDir
					}

					// check to see if we've already checked this head
					for _, h := range headsChecked {
						if h.Hash() == headRef.Hash() {
							return filepath.SkipDir
						}
					}

					cl, err := getCommits(path, r, headRef)
					if err != nil {
						return filepath.SkipDir
					}

					commits = append(commits, cl...)

					headsChecked = append(headsChecked, headRef)

					return filepath.SkipDir
				}

				return nil
			})
		}

		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Commit.Author.When.Before(commits[j].Commit.Author.When)
		})

		printLogs(commits)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.what-did-i-work-on.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printLogs(commits []*CommitForRepo) {
	for _, c := range commits {
		fmt.Printf(
			"%s - %s - %s\n",
			c.Commit.Author.When.Format("02 Jan 03:04 pm"),
			c.RepoName,
			strings.TrimRight(c.Commit.Message, "\n"),
		)
	}
}

func getCommits(gitFolder string, r *git.Repository, ref *plumbing.Reference) ([]*CommitForRepo, error) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	commits, err := r.Log(&git.LogOptions{
		All:   true,
		Since: &yesterday,
		From:  ref.Hash(),
	})

	if err != nil {
		return nil, err
	}

	var commitList []*CommitForRepo

	me, err := common.GetAuthorName(r)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(gitFolder, string(os.PathSeparator))

	var repoName string
	if len(parts) > 2 {
		repoName = parts[len(parts)-3] + "/" + parts[len(parts)-2]
	} else {
		repoName = parts[len(parts)-2]
	}

	commits.ForEach(func(c *object.Commit) error {

		if c.Author.Name != me {
			return nil
		}

		commitList = append(commitList, NewCommitForRepo(repoName, c))
		return nil
	})

	return commitList, nil
}

type CommitForRepo struct {
	RepoName string
	Commit   *object.Commit
}

func NewCommitForRepo(repoName string, commit *object.Commit) *CommitForRepo {
	return &CommitForRepo{
		RepoName: repoName,
		Commit:   commit,
	}
}

func openRepoAndGetHeadRef(gitFolder string) (*git.Repository, *plumbing.Reference, error) {
	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, nil, err
	}

	return r, head, nil
}

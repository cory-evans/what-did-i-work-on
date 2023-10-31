package common

import (
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetAuthorName(r *git.Repository) (string, error) {

	globalConfig, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return "", err
	}

	localConfig, err := r.Config()
	if err != nil {
		return "", err
	}

	if localConfig.User.Name != "" {
		return localConfig.User.Name, nil
	}

	if globalConfig.User.Name != "" {
		return globalConfig.User.Name, nil
	}

	return "", nil
}

func GetCommitsByAuthor(iter object.CommitIter, author string) ([]*object.Commit, error) {
	commits := []*object.Commit{}
	iter.ForEach(func(c *object.Commit) error {
		if c.Author.Name == author {
			commits = append(commits, c)
		}
		return nil
	})

	return commits, nil
}

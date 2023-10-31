package main

import (
	"log"

	"github.com/cory-evans/what-did-i-work-on/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	r, err := git.PlainOpen("./")
	common.CheckError(err)

	ref, err := r.Head()
	common.CheckError(err)

	cIter, err := r.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	common.CheckError(err)

	cIter.ForEach(func(c *object.Commit) error {
		// common.PrintCommit(c)

		log.Println(c.Author.Name, c.Author.Email, c.Author.When, c.Message)
		return nil
	})

}

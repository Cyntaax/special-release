package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	repo, _ := git.PlainOpen("../")

	head, _ := repo.Head()
	iter, _ := repo.Log(&git.LogOptions{From: head.Hash()})
	var currentCommit *object.Commit
	iter.ForEach(func(c *object.Commit) error {
		if currentCommit == nil {
			currentCommit = c
		}
		return nil
	})

	fmt.Println("Current commit", currentCommit.Message)

	parents := currentCommit.NumParents()

	if parents >= 1 {
		parent, _ := currentCommit.Parent(0)
		fmt.Println(parent.Hash)

		pt, _ := currentCommit.Patch(parent)

		for _, patch := range pt.FilePatches() {
			fr, to := patch.Files()
			fmt.Println(fr.Path(), to.Path())
		}
	}
}
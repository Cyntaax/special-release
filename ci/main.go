package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"os"
	"path"
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
	tr, _ := currentCommit.Tree()

	if parents >= 1 {
		parent, _ := currentCommit.Parent(0)
		fmt.Println(parent.Hash)

		pt, err := currentCommit.Patch(parent)
		if err != nil {
			panic(err)
		}


		zipBuff := new(bytes.Buffer)

		zipWriter := zip.NewWriter(zipBuff)

		for _, patch := range pt.FilePatches() {
			fr, _ := patch.Files()
			fmt.Println(path.Dir(fr.Path()), path.Base(fr.Path()))

			if matched, _ := path.Match("resources/*", fr.Path()); matched == false {
				continue
			}


			file, _ := tr.File(fr.Path())

			fmt.Println(file.Name)

			f, err := zipWriter.Create(fr.Path())
			if err != nil {
				fmt.Println("ERROR CREATING FILE", err.Error())
				continue
			}

			fileReader, err := file.Reader()

			copied, err := io.Copy(f, fileReader)
			if err != nil {
				fmt.Println("ERROR", err.Error())
			}
			fmt.Println("Wrote", copied)
		}

		output, _ := os.Create("release.zip")
		zipWriter.Close()
		output.Write(zipBuff.Bytes())
	}
}
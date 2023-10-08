package main

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetBranches() []string {

	var branches []string
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	branchIter, err := repo.Branches()
	if err != nil {
		log.Print(err)
	}

	myFunc := func(ref *plumbing.Reference) error {
		trim := strings.TrimPrefix(ref.Name().String(), "refs/heads/")
		branches = append(branches, trim)
		return nil
	}

	err = branchIter.ForEach(myFunc)
	if err != nil {
		log.Print(err)
	}

	return branches
}

// An alternative to fetching branches in the tests
func GetBranchesHelper(r *git.Repository, branchNames []string) []string {
	var branches []string

	for _, b := range branchNames {
		b, err := r.Branch(b)
		if err != nil {
			log.Print("Error", err)
		}
		branches = append(branches, b.Name)
	}

	return branches
}

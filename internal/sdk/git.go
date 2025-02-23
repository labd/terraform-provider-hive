package sdk

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type CommitInfo struct {
	Author string
	Hash   string
}

// GetLatestCommitInfo returns the latest commit's author (name <email>) and commit hash.
func GetLatestCommitInfo() (*CommitInfo, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference.
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commit object pointed by HEAD.
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	info := &CommitInfo{
		Author: fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email),
		Hash:   commit.Hash.String(),
	}
	return info, nil
}

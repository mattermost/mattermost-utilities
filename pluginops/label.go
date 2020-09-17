package main

import (
	"github.com/google/go-github/v32/github"
)

// Label respresents a GitHub label.
type Label struct {
	Name        string
	Description string
	Color       string
}

// ToGithubLabel converts a label to a github library compatible type.
func (l Label) ToGithubLabel() *github.Label {
	return &github.Label{
		Name:        &l.Name,
		Description: &l.Description,
		Color:       &l.Color,
	}
}

// Equal compare to a github library compatible type.
func (l Label) Equal(gh *github.Label) bool {
	return gh.GetName() == l.Name &&
		gh.GetDescription() == l.Description &&
		gh.GetColor() == l.Color
}

// MergeLabels return the union of two slices of labels.
func MergeLabels(l1 []Label, l2 ...[]Label) []Label {
	for _, l := range l2 {
		for _, e2 := range l {
			found := false

			for _, e1 := range l1 {
				if e1.Name == e2.Name {
					found = true
					break
				}
			}

			if !found {
				l1 = append(l1, e2)
			}
		}
	}

	return l1
}

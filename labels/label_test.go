package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeLabels(t *testing.T) {
	l1 := []Label{
		{"1: PM Review", "Requires review by a product manager", "#006b75"},
		{"2: Dev Review", "Requires review by a core committer", "#eb6420"},
		{"Work In Progress", "Not yet ready for review", "e11d21"},
	}

	l2 := []Label{
		{"Bug", "Something isn't working", "d73a4a"},
		{"Enhancement", "New feature or request", "a2eeef"},
	}

	l3 := []Label{
		{"Changelog/Done", "Required changelog entry has been written", "0e8a16"},
		{"Changelog/Not Needed", "Does not require a changelog entry", "d4c5f9"},
	}

	for name, test := range map[string]struct {
		L1             []Label
		L2             []Label
		L3             []Label
		ExpectedOutput []Label
	}{
		"Equal lists": {
			L1:             l1,
			L2:             l1,
			ExpectedOutput: l1,
		},
		"First list empty": {
			L1:             []Label{},
			L2:             l1,
			ExpectedOutput: l1,
		},
		"Second list empty": {
			L1:             l1,
			L2:             []Label{},
			ExpectedOutput: l1,
		},
		"merge": {
			L1: l1,
			L2: l2,
			ExpectedOutput: []Label{
				{"1: PM Review", "Requires review by a product manager", "#006b75"},
				{"2: Dev Review", "Requires review by a core committer", "#eb6420"},
				{"Work In Progress", "Not yet ready for review", "e11d21"},
				{"Bug", "Something isn't working", "d73a4a"},
				{"Enhancement", "New feature or request", "a2eeef"},
			},
		},
		"merge with duplicates": {
			L1: append(l1, Label{"Enhancement", "New feature or request", "a2eeef"}),
			L2: l2,
			ExpectedOutput: []Label{
				{"1: PM Review", "Requires review by a product manager", "#006b75"},
				{"2: Dev Review", "Requires review by a core committer", "#eb6420"},
				{"Work In Progress", "Not yet ready for review", "e11d21"},
				{"Enhancement", "New feature or request", "a2eeef"},
				{"Bug", "Something isn't working", "d73a4a"},
			},
		},
		"merge with duplicates list": {
			L1:             append(l1, l2...),
			L2:             l2,
			ExpectedOutput: append(l1, l2...),
		},
		"merge with three lists": {
			L1: l1,
			L2: l2,
			L3: l3,
			ExpectedOutput: []Label{
				{"1: PM Review", "Requires review by a product manager", "#006b75"},
				{"2: Dev Review", "Requires review by a core committer", "#eb6420"},
				{"Work In Progress", "Not yet ready for review", "e11d21"},
				{"Bug", "Something isn't working", "d73a4a"},
				{"Enhancement", "New feature or request", "a2eeef"},
				{"Changelog/Done", "Required changelog entry has been written", "0e8a16"},
				{"Changelog/Not Needed", "Does not require a changelog entry", "d4c5f9"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedOutput, MergeLabels(test.L1, test.L2, test.L3))
		})
	}
}

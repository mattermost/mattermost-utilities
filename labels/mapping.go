package main

var mapping = map[string][]Label{
	"mattermost-plugin-giphy": MergeLabels(core, plugin, helpWanted),
}

var core = []Label{
	{"1: PM Review", "Requires review by a product manager", "006b75"},
	{"2: Dev Review", "Requires review by a core committer", "eb6420"},
	{"2: QA Review", "Requires review by a QA tester", "7cdfe2"},
	{"3: Reviews Complete", "All reviewers have approved the pull request", "0e8a16"},
	{"Awaiting Submitter Action", "Blocked on the author", "b60205"},
	{"Do Not Merge/Awaiting PR", "Awaiting another pull request before merging (e.g. server changes)", "a32735"},
	{"Do Not Merge", "Should not be merged until this label is removed", "a32735"},
	{"Invalid", "This doesn't seem right", "e4e669"},
	{"Lifecycle/frozen", "", "d3e2f0"},
	{"Lifecycle/1:stale", "", "5319e7"},
	{"Triage", "", "efcb6e"},
	{"Work In Progress", "Not yet ready for review", "e11d21"},
}

var plugin = []Label{
	{"Bug", "Something isn't working", "d73a4a"},
	{"Duplicate", "This issue or pull request already exists", "cfd3d7"},
	{"Enhancement", "New feature or request", "a2eeef"},
	{"Needs Mattermost Changes", "Requires changes to the Mattermost Plugin tookit", "9c14c9"},
	{"Question", "Further information is requested", "d876e3"},
	{"Wontfix", "This will not be worked on", "ffffff"},
}

var helpWanted = []Label{
	{"Difficulty/1:Easy", "Easy ticket", "c2e0c6"},
	{"Difficulty/2:Medium", "Medium ticket", "bfdadc"},
	{"Difficulty/3:Hard", "Hard ticket", "f9d0c4"},
	{"Good First Issue", "Suitable for first-time contributors", "7057ff"},
	{"Hacktoberfest", "", "dc7d02"},
	{"Help Wanted", "Community help wanted", "33aa3f"},
	{"Needs spec", "Needs further specification to be a good help wanted ticket", "88acea"},
	{"Tech/Go", "", "0e8a16"},
	{"Tech/ReactJS", "", "1d76db"},
	{"Up For Grabs", "Ready for help from the community. Removed when someone volunteers", "8B4500"},
}

var docs = []Label{
	{"Docs/Done", "Required documentation has been written", "0e8a16"},
	{"Docs/Needed", "Requires documentation", "b60205"},
	{"Docs/Not Needed", "Does not require documentation", "d4c5f9"},
}

var changelog = []Label{
	{"Changelog/Done", "Required changelog entry has been written", "0e8a16"},
	{"Changelog/Not Needed", "Does not require a changelog entry", "d4c5f9"},
}

var tests = []Label{
	{"Tests/Done", "Required tests have been written", "0e8a16"},
	{"Tests/Not Needed", "Does not require tests", "d4c5f9"},
}

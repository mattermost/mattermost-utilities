package main

var defaultLabels = coreLabels
var coreLabels = MergeLabels(pullRequest)
var pluginLabels = MergeLabels(coreLabels, issue, helpWanted, docs, plugin)

var defaultMapping = map[string][]Label{
	"mattermost-oembed-plugin":            pluginLabels,
	"mattermost-plugin-agenda":            pluginLabels,
	"mattermost-plugin-antivirus":         pluginLabels,
	"mattermost-plugin-autolink":          pluginLabels,
	"mattermost-plugin-autotranslate":     coreLabels,
	"mattermost-plugin-aws-SNS":           pluginLabels,
	"mattermost-plugin-channel-export":    coreLabels,
	"mattermost-plugin-community":         pluginLabels,
	"mattermost-plugin-demo":              coreLabels,
	"mattermost-plugin-custom-attributes": pluginLabels,
	"mattermost-plugin-email-reply":       pluginLabels,
	"mattermost-plugin-giphy":             pluginLabels,
	"mattermost-plugin-github":            pluginLabels,
	"mattermost-plugin-gitlab":            pluginLabels,
	"mattermost-plugin-google-calendar":   pluginLabels,
	"mattermost-plugin-incident-response": MergeLabels(coreLabels, issue, docs, plugin),
	"mattermost-plugin-jenkins":           pluginLabels,
	"mattermost-plugin-jira":              pluginLabels,
	"mattermost-plugin-memes":             pluginLabels,
	"mattermost-plugin-mscalendar":        pluginLabels,
	"mattermost-plugin-msoffice":          pluginLabels,
	"mattermost-plugin-nop":               coreLabels,
	"mattermost-plugin-nps":               coreLabels,
	"mattermost-plugin-profanity-filter":  coreLabels,
	"mattermost-plugin-skype4business":    pluginLabels,
	"mattermost-plugin-solar-lottery":     pluginLabels,
	"mattermost-plugin-starter-template":  coreLabels,
	"mattermost-plugin-suggestions":       pluginLabels,
	"mattermost-plugin-todo":              pluginLabels,
	"mattermost-plugin-webex":             pluginLabels,
	"mattermost-plugin-welcomebot":        pluginLabels,
	"mattermost-plugin-workflow":          MergeLabels(coreLabels, issue, docs, plugin),
	"mattermost-plugin-workflow-client":   coreLabels,
	"mattermost-plugin-zoom":              pluginLabels,
}

// PR is the list of labels typically used on PRs. Use --
var pullRequest = []Label{
	{"1: PM Review", "Requires review by a product manager", "006b75"},
	{"1: UX Review", "Requires review by a UX Designer", "7cdfe2"},
	{"2: Dev Review", "Requires review by a core committer", "eb6420"},
	{"2: QA Review", "Requires review by a QA tester", "7cdfe2"},
	{"3: Reviews Complete", "All reviewers have approved the pull request", "0e8a16"},
	{"AutoMerge", "Used by Mattermod to merge PR automatically", "b74533"},
	{"Awaiting Submitter Action", "Blocked on the author", "b60205"},
	{"Do Not Merge/Awaiting PR", "Awaiting another pull request before merging (e.g. server changes)", "a32735"},
	{"Do Not Merge", "Should not be merged until this label is removed", "a32735"},
	{"Lifecycle/frozen", "", "d3e2f0"},
	{"Lifecycle/1:stale", "", "5319e7"},
	{"Work In Progress", "Not yet ready for review", "e11d21"},
}

var issue = []Label{
	{"Bug", "Something isn't working", "d73a4a"},
	{"Duplicate", "This issue or pull request already exists", "cfd3d7"},
	{"Enhancement", "New feature or request", "a2eeef"},
	{"Invalid", "This doesn't seem right", "e4e669"},
	{"Question", "Further information is requested", "d876e3"},
	{"Triage", "", "efcb6e"},
	{"Wontfix", "This will not be worked on", "ffffff"},
}

var plugin = []Label{
	{"Needs Mattermost Changes", "Requires changes to the Mattermost Plugin tookit", "9c14c9"},
	{"Setup Cloud Test Server", "Setup a test server using Mattermost Cloud", "0052cc"},
	{"Setup HA Cloud Test Server", "Setup an HA test server using Mattermost Cloud", "00efff"},
}

var helpWanted = []Label{
	{"Difficulty/1:Easy", "Easy ticket", "c2e0c6"},
	{"Difficulty/2:Medium", "Medium ticket", "bfdadc"},
	{"Difficulty/3:Hard", "Hard ticket", "f9d0c4"},
	{"Good First Issue", "Suitable for first-time contributors", "7057ff"},
	{"Hacktoberfest", "", "dc7d02"},
	{"Help Wanted", "Community help wanted", "33aa3f"},
	{"Help Wanted Candidate", "Needs further specification to be a good help wanted ticket", "88acea"},
	{"Tech/Go", "", "0e8a16"},
	{"Tech/JavaScript", "", "f9d0c4"},
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

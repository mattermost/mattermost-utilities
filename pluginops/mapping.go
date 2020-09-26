package main

var coreLabels = pullRequest
var communityPlugins = MergeLabels(coreLabels, []Label{securityReview})
var pluginLabelsWithoutHW = MergeLabels(coreLabels, issue, docs, plugin)
var pluginLabels = MergeLabels(pluginLabelsWithoutHW, helpWanted)

var defaultMapping = map[string][]Label{
	"mattermost-icebreaker-plugin":        communityPlugins,
	"mattermost-plugin-agenda":            pluginLabels,
	"mattermost-plugin-antivirus":         pluginLabels,
	"mattermost-plugin-api":               pluginLabels,
	"mattermost-plugin-autolink":          pluginLabels,
	"mattermost-plugin-autotranslate":     pluginLabels,
	"mattermost-plugin-aws-SNS":           pluginLabels,
	"mattermost-plugin-canary":            pluginLabels,
	"mattermost-plugin-channel-export":    pluginLabels,
	"mattermost-plugin-cloud":             pluginLabels,
	"mattermost-plugin-community":         pluginLabels,
	"mattermost-plugin-confluence":        pluginLabels,
	"mattermost-plugin-custom-attributes": pluginLabels,
	"mattermost-plugin-demo-creator":      pluginLabels,
	"mattermost-plugin-demo":              pluginLabels,
	"mattermost-plugin-dice-roller":       communityPlugins,
	"mattermost-plugin-digitalocean":      communityPlugins,
	"mattermost-plugin-docup":             pluginLabels,
	"mattermost-plugin-email-reply":       pluginLabels,
	"mattermost-plugin-giphy-moussetc":    communityPlugins,
	"mattermost-plugin-github":            pluginLabels,
	"mattermost-plugin-gitlab":            pluginLabels,
	"mattermost-plugin-gmail":             communityPlugins,
	"mattermost-plugin-google-calendar":   pluginLabels,
	"mattermost-plugin-incident-response": pluginLabels,
	"mattermost-plugin-jenkins":           pluginLabels,
	"mattermost-plugin-jira":              pluginLabels,
	"mattermost-plugin-jitsi":             pluginLabels,
	"mattermost-plugin-memes":             pluginLabels,
	"mattermost-plugin-mscalendar":        pluginLabels,
	"mattermost-plugin-msteams-meetings":  pluginLabels,
	"mattermost-plugin-nop":               pluginLabels,
	"mattermost-plugin-nps":               pluginLabels,
	"mattermost-plugin-oembed":            pluginLabels,
	"mattermost-plugin-profanity-filter":  pluginLabels,
	"mattermost-plugin-recommend":         communityPlugins,
	"mattermost-plugin-skype4business":    pluginLabels,
	"mattermost-plugin-solar-lottery":     pluginLabels,
	"mattermost-plugin-starter-template":  pluginLabels,
	"mattermost-plugin-suggestions":       pluginLabels,
	"mattermost-plugin-todo":              pluginLabels,
	"mattermost-plugin-walltime":          pluginLabels,
	"mattermost-plugin-webex":             pluginLabels,
	"mattermost-plugin-webrtc-video":      communityPlugins,
	"mattermost-plugin-welcomebot":        pluginLabels,
	"mattermost-plugin-workflow-client":   pluginLabels,
	"mattermost-plugin-workflow":          pluginLabelsWithoutHW,
	"mattermost-plugin-zoom":              pluginLabels,
	"standup-raven":                       communityPlugins,
}

var securityReview = Label{
	"3: Security Review", "Review requested from Security Team", "1d76db",
}

// PR is the list of labels typically used on PRs. Use --
var pullRequest = []Label{
	{"1: PM Review", "Requires review by a product manager", "006b75"},
	{"1: UX Review", "Requires review by a UX Designer", "7cdfe2"},
	{"2: Dev Review", "Requires review by a core committer", "eb6420"},
	{"3: QA Review", "Requires review by a QA tester", "7cdfe2"},
	{"4: Reviews Complete", "All reviewers have approved the pull request", "0e8a16"},
	{"AutoMerge", "Used by Mattermod to merge PR automatically", "b74533"},
	{"Awaiting Submitter Action", "Blocked on the author", "b60205"},
	{"Do Not Merge/Awaiting PR", "Awaiting another pull request before merging (e.g. server changes)", "a32735"},
	{"Do Not Merge", "Should not be merged until this label is removed", "a32735"},
	{"Lifecycle/frozen", "", "d3e2f0"},
	{"Lifecycle/1:stale", "", "5319e7"},
	{"Lifecycle/2:inactive", "", "a34523"},
	{"Lifecycle/3:orphaned", "", "111111"},
	{"QA Review Done", "PR has been approved by QA", "0e8a16"},
	{"Work In Progress", "Not yet ready for review", "e11d21"},
}

var issue = []Label{
	{"Duplicate", "This issue or pull request already exists", "cfd3d7"},
	{"Invalid", "This doesn't seem right", "e4e669"},
	{"Triage", "", "efcb6e"},
	{"Type/Bug", "Something isn't working", "d73a4a"},
	{"Type/Enhancement", "New feature or improvement of existing feature", "a2eeef"},
	{"Type/Task", "A general task", "6698d1"},
	{"Type/Question", "Further information is requested", "d876e3"},
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
	{"Needs Spec", "Needs further specification to be a good (help wanted) ticket", "e25de0"},
	{"Tech/Go", "", "0e8a16"},
	{"Tech/ReactJS", "", "1d76db"},
	{"Tech/TypeScript", "", "c9ffff"},
	{"Up For Grabs", "Ready for help from the community. Removed when someone volunteers", "8B4500"},
}

var docs = []Label{
	{"Docs/Done", "Required documentation has been written", "0e8a16"},
	{"Docs/Needed", "Requires documentation", "b60205"},
	{"Docs/Not Needed", "Does not require documentation", "d4c5f9"},
}

var migrateMap = map[string]string{
	"Bug":             "Type/Bug",
	"Enhancement":     "Type/Enhancement",
	"Question":        "Type/Question",
	"Tech/JavaScript": "",
	// Migrate QA labels: https://community.mattermost.com/core/pl/eegcso8xr3bqzr7giftc3dgdka
	"2: QA Review":        "3: QA Review",
	"3: Reviews Complete": "4: Reviews Complete",
}

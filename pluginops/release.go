package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/manifoldco/promptui"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bumpVersionCmd)
}

var bumpVersionCmd = &cobra.Command{
	Use:   "release",
	Short: "Prepare a plugin release by creating a version bump PR",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		err := releaseVersion()
		if err != nil {
			log.WithError(err).Fatal("Release failed")
		}

		log.Info("Release complete!")
	},
}

const (
	remoteName = "origin"
	gitPath    = "."
)

// bumpVersion
func releaseVersion() error {
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return errors.Wrapf(err, "failed to access git repository at %q", gitPath)
	}

	log.Info("Checking if repository is clean")

	worktree, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	status, err := worktree.Status()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree status")
	}

	if !status.IsClean() {
		return errors.New("Repository is not clean")
	}

	newVersionString, err := bumpManifestVersion()
	if err != nil {
		return err
	}

	log.Info("Running \"make apply\"")

	cmd := exec.Command("make", []string{"apply"}...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		return err
	}
	fmt.Print(string(out))

	// TODO: Use library
	cmd = exec.Command("git", []string{"diff"}...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		return err
	}
	fmt.Print(string(out))

	ok, err := confirmPrompt("Does the diff look good")
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	ok, err = confirmPrompt(fmt.Sprintf("Do you want to create a new branch for the release and push to %q", remoteName))
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	branchName := fmt.Sprintf("release_v%s", newVersionString)

	log.Infof("Creating branch %q", branchName)

	remote, err := repo.Remote(remoteName)
	if err != nil {
		return errors.Wrapf(err, "Remote %q not found", remoteName)
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	nameRef := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(nameRef, headRef.Hash())

	branchConfig := &config.Branch{
		Name:   branchName,
		Remote: remoteName,
		Merge:  nameRef,
	}
	err = repo.CreateBranch(branchConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to create branch %q", branchName)
	}

	err = repo.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Keep:   true,
	})
	if err != nil {
		return err
	}

	log.Infof("Commiting changes")

	commit, err := worktree.Commit(fmt.Sprintf("Bump version to %s", newVersionString), &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return err
	}
	log.Info("Create commit:")
	fmt.Println(obj)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("master"),
	})
	if err != nil {
		return err
	}

	log.Infof("Pushing to %q", remoteName)

	pushOptions := &git.PushOptions{
		RemoteName: remoteName,
	}
	err = repo.Push(pushOptions)
	if err != nil {
		if errors.Is(err, git.ErrNonFastForwardUpdate) {
			return errors.New("Can't push to remote branch. Please check that the branch names doesn't collide.")
		}

		return err
	}

	orgName, repoName, err := GetRepoURL(remote.Config().URLs[0])
	if err != nil {
		return errors.Wrap(err, "failed to extract repository information out of origin URL")
	}

	u, err := url.Parse(fmt.Sprintf("https://github.com/%s/%s/compare/master...%s:%s", org, repoName, orgName, branchName))
	if err != nil {
		return err
	}

	// newVersion might be a pre release
	newVersion := semver.MustParse(newVersionString)

	var ticketLink string
	client, err := getGitHubClient()
	if err != nil {
		log.WithError(err).Debug("Failed to create GitHub client. Won't try to automatically set ticket link")
	} else {
		issueTitle := fmt.Sprintf("Release v%s", newVersion.FinalizeVersion())
		query := fmt.Sprintf("state:open is:issue org:mattermost repo:repoName %s", issueTitle)
		result, _, err := client.Search.Issues(context.Background(), query, nil)
		if err != nil {
			return err
		}

		len := len(result.Issues)
		if len == 0 {
			log.Debug("No issues found")
		} else {
			if len > 1 {
				log.Warn("More then one issue found. Picking the first one")
			}
			ticketLink = fmt.Sprintf("Part of %s", result.Issues[0].GetHTMLURL())
		}
	}

	title := fmt.Sprintf("Release v%s", newVersionString)

	values := url.Values{}
	values.Add("quick_pull", "1")
	values.Add("body", fmt.Sprintf("#### Summary\n\n#### Ticket Link\n%s\n", ticketLink))
	values.Add("title", title)
	values.Add("labels", "2: Dev Review,3: QA Review")
	values.Add("milestone", fmt.Sprintf("v%s", newVersion.FinalizeVersion()))
	u.RawQuery = values.Encode()

	log.Infof("You can open a PR by clicking %s", u.String())

	ok, err = confirmPrompt(fmt.Sprintf("Do you want to delete your local branch %q", branchName))
	if err != nil {
		return err
	}

	if ok {
		// TODO: This doesn't work right now
		err = repo.DeleteBranch(branchName)
		if err != nil {
			return err
		}
	}

	return nil
}

func bumpManifestVersion() (string, error) {
	manifest, err := findManifest()
	if err != nil {
		return "", errors.Wrap(err, "failed to find manifest")
	}

	oldVersion, err := semver.Parse(manifest.Version)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse version in manifest")
	}

	newVersionMajor := oldVersion
	err = newVersionMajor.IncrementMajor()
	if err != nil {
		return "", errors.Wrap(err, "failed to increment major version")
	}

	newVersionMinor := oldVersion
	err = newVersionMinor.IncrementMinor()
	if err != nil {
		return "", errors.Wrap(err, "failed to increment minor version")
	}

	newVersionPatch := oldVersion
	err = newVersionPatch.IncrementPatch()
	if err != nil {
		return "", errors.Wrap(err, "failed to increment patch version")
	}

	log.Infof("Current version is %s", oldVersion.String())

	items := []string{
		fmt.Sprintf("%s (Next Major)", newVersionMajor.FinalizeVersion()),
		fmt.Sprintf("%s (Next Minor)", newVersionMinor.FinalizeVersion()),
		fmt.Sprintf("%s (Next Patch)", newVersionPatch.FinalizeVersion()),
	}

	if len(oldVersion.Pre) > 0 {
		items = append(items, fmt.Sprintf("%s (Next Stable Version)", oldVersion.FinalizeVersion()))
	}

	var result string
	index := -1

	for index < 0 {
		v := func(o string) error {
			_, err := semver.Parse(o)
			if err != nil {
				return err
			}

			return nil
		}

		prompt := promptui.SelectWithAdd{
			Label:    "To which version do you want to bump",
			Items:    items,
			AddLabel: "Custom version",
			Validate: v,
		}
		index, result, err = prompt.Run()

		if err != nil {
			return "", errors.Wrap(err, "prompt failed")
		}

		if index == -1 {
			items = append(items, result)
		}
	}

	var newVersionString string
	switch index {
	// Custom version
	case 0:
		newVersionString = newVersionMajor.FinalizeVersion()
	case 1:
		newVersionString = newVersionMinor.FinalizeVersion()
	case 2:
		newVersionString = newVersionPatch.FinalizeVersion()
	case 3:
		if len(oldVersion.Pre) > 0 {
			newVersionString = oldVersion.FinalizeVersion()
		} else {
			newVersionString = result
		}
	default:
		newVersionString = result
	}

	manifest.Version = newVersionString

	// Patch ReleaseNotesURL if needed
	if manifest.ReleaseNotesURL != "" {
		manifest.ReleaseNotesURL = strings.ReplaceAll(manifest.ReleaseNotesURL, oldVersion.String(), newVersionString)
	}

	log.Infof("Bumping version to %v", newVersionString)

	err = writeManifest(manifest)
	if err != nil {
		return "", errors.Wrap(err, "failed to writing manifest after bumping version")
	}

	return newVersionString, nil
}

func confirmPrompt(msg string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			return false, nil
		}

		return false, errors.Wrap(err, "prompt failed")
	}

	return true, nil
}

func GetRepoURL(originURL string) (string, string, error) {
	switch {
	case strings.HasPrefix(originURL, "https://github.com/"):
		originURL = strings.TrimPrefix(originURL, "https://github.com/")
	case strings.HasPrefix(originURL, "git@github.com:"):
		originURL = strings.TrimPrefix(originURL, "git@github.com:")
	default:
		return "", "", errors.Errorf("unknown prefix of origin URL %s", originURL)
	}

	originURL = strings.TrimSuffix(originURL, ".git")

	split := strings.Split(originURL, "/")
	if len(split) != 2 {
		return "", "", errors.Errorf("malformed origin URL")
	}

	return split[0], split[1], nil
}

func findManifest() (*model.Manifest, error) {
	_, manifestFilePath, err := model.FindManifest(".")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find manifest in current working directory")
	}

	manifestFile, err := os.Open(manifestFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", manifestFilePath)
	}
	defer manifestFile.Close()

	// Re-decode the manifest, disallowing unknown fields. When we write the manifest back out,
	// we don't want to accidentally clobber anything we won't preserve.
	var manifest model.Manifest
	decoder := json.NewDecoder(manifestFile)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "failed to parse manifest")
	}

	return &manifest, nil
}

// writeManifest writes a given manifest back to file
func writeManifest(manifest *model.Manifest) error {
	_, manifestFilePath, err := model.FindManifest(".")
	if err != nil {
		return errors.Wrap(err, "failed to find manifest in current working directory")
	}

	file, err := os.OpenFile(manifestFilePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", manifestFilePath)
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err = encoder.Encode(manifest); err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to encode manifest into manifest file")
	}

	return nil
}

package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(patchVersionCmd)
}

var patchVersionCmd = &cobra.Command{
	Use:   "patch",
	Short: "TODO",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		err := patchPlugin()
		if err != nil {
			log.WithError(err).Fatal("Release failed")
		}

		log.Info("Release complete!")
	},
}

func patchPlugin() error {
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

	/*

		os.ReadFile()

		log.Info("Running \"make apply\"")

		cmd := exec.Command("make", []string{"apply"}...)
		out, _ := cmd.CombinedOutput()
		// This is allowed to fail as make apply isn't available on all plugins any longer
		fmt.Print(string(out))

		// TODO: Use library
		cmd = exec.Command("git", []string{"am"}...)
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
	*/

	return nil
}

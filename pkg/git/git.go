package git

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	argogit "github.com/argoproj/argo-cd/util/git"
	"github.com/argoproj/argo-cd/util/rand"

	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
)

// Package git implements write of changes back to a Git repository.

//
func CommitChanges(app *v1alpha1.Application) error {
	root := filepath.Join(os.TempDir(), rand.RandString(20))
	defer func() {
		log.Log().Debugf("Removing temp git checkout at '%s' for app '%s'", root, app.GetName())
		os.RemoveAll(root)
	}()

	gitC, _ := argogit.NewClientExt(app.Spec.Source.RepoURL, root, argogit.NopCreds{}, false, false)
	err := gitC.Init()
	if err != nil {
		log.Errorf("Could not initialize local Git repo at '%s' for '%s'", root, app.GetName())
		return err
	}

	err = gitC.Fetch()
	if err != nil {
		log.Errorf("Could not fetch from remote: %v", err)
		return err
	}

	rev, err := gitC.LsRemote(app.Spec.Source.TargetRevision)
	if err != nil {
		log.Errorf("Could not list remote revision: %v", err)
	}

	err = gitC.Checkout(rev)
	if err != nil {
		log.Errorf("Could not checkout revision %v", err)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(root, ".argocd-source.yaml"), []byte("# Test data yeehaw"), 0644)
	if err != nil {
		log.Errorf("Could not write .argocd-source.yaml: %v", err)
		return err
	}

	err = gitC.Commit("Update image to %s", ".argocd-source.yaml")
	if err != nil {
		log.Errorf("Could not commit file: %v", err)
		return err
	}

	err = gitC.Push("master")
	if err != nil {
		log.Errorf("Could not push changes: %v", err)
		return err
	}

	return nil
}

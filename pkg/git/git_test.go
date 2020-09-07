package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/util/rand"
	"github.com/plus3it/gorecurcopy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CommitChanges(t *testing.T) {
	repoPath := filepath.Join(os.TempDir(), "testrepo-"+rand.RandString(20))
	err := os.Mkdir(repoPath, 0755)
	require.NoError(t, err)
	err = gorecurcopy.CopyDirectory("../../test/testdata/git", repoPath)
	require.NoError(t, err)
	defer func() {
		os.RemoveAll(repoPath)
	}()
	app := v1alpha1.Application{
		Spec: v1alpha1.ApplicationSpec{
			Source: v1alpha1.ApplicationSource{
				RepoURL:        repoPath,
				Path:           "/app",
				TargetRevision: "HEAD",
			},
		},
	}
	err = CommitChanges(&app)
	assert.NoError(t, err)
}

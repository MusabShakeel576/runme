package project

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/require"
)

func Test_selectRemoteURL(t *testing.T) {
	storer := memory.NewStorage()
	remotes := []*git.Remote{
		git.NewRemote(storer, &config.RemoteConfig{
			Name: "custom",
			URLs: []string{"git@my-server.com:stateful/cli.git"},
		}),
		git.NewRemote(storer, &config.RemoteConfig{
			Name: "origin",
			URLs: []string{"git@github.com:stateful/cli.git"},
		}),
	}
	url := selectRemoteURL(remotes)
	require.Equal(t, "git@github.com:stateful/cli.git", url)
}

var branchFixtures = `
Merge pull request #333 from stateful:seb/cal-edits--||--More text edits etc
Merge pull request #220 from stateful/admc/status-vscode-button--||--Add open in vscode button
Merge pull request #132 from stateful/jgee/feat/cli-instructions-platform-specific--||--Use accordion like component to contain CLI instructions
Merge branch 'main' into jgee/feat/cli-instructions-platform-specific--||--
Merge branch 'main' into admc/standup-ux-refactor--||--
Merging--||--
Merge branch 'admc/slack-attack'--||--
Merge pull request #7 from activecove/seb/file-cycles--||--Move to file sessions`

func Test_branchNamePreNonGreedy(t *testing.T) {
	expected := []Branch{
		{Name: "seb/cal-edits", Description: "More text edits etc"},
		{Name: "admc/status-vscode-button", Description: "Add open in vscode button"},
		{Name: "jgee/feat/cli-instructions-platform-specific", Description: "Use accordion like component to contain CLI instructions"},
		{Name: "seb/file-cycles", Description: "Move to file sessions"},
	}
	actual := getBranchNamesFromStdout(branchFixtures, false)
	require.Equal(t, expected, actual)
}

func Test_branchNamePreGreedy(t *testing.T) {
	expected := []Branch{
		{Name: "cal-edits", Description: "More text edits etc"},
		{Name: "status-vscode-button", Description: "Add open in vscode button"},
		{Name: "cli-instructions-platform-specific", Description: "Use accordion like component to contain CLI instructions"},
		{Name: "file-cycles", Description: "Move to file sessions"},
	}
	actual := getBranchNamesFromStdout(branchFixtures, true)
	require.Equal(t, expected, actual)
}

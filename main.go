package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/you06/releaser/manager"
	"github.com/you06/releaser/pkg/types"
)

const (
	nmToken           = "token"
	nmReleaseNoteRepo = "repo"
	nmVersion         = "version"
)

var (
	// common args
	token           string
	releaseNoteRepo string
	version         string
	subCommand      string
	// PRList sub command
	cmdPRList                = flag.NewFlagSet(types.SubCmdPRList, flag.ExitOnError)
	cmdPRListToken           = cmdPRList.String(nmToken, "", "github token")
	cmdPRListReleaseNoteRepo = cmdPRList.String(nmReleaseNoteRepo, "", "release note repo")
	cmdPRListVersion         = cmdPRList.String(nmVersion, "", "release version")
	// ReleaseNotes sub command
	cmdReleaseNotes                = flag.NewFlagSet(types.SubCmdReleaseNotes, flag.ExitOnError)
	cmdReleaseNotesToken           = cmdReleaseNotes.String(nmToken, "", "github token")
	cmdReleaseNotesReleaseNoteRepo = cmdReleaseNotes.String(nmReleaseNoteRepo, "", "release note repo")
	cmdReleaseNotesVersion         = cmdReleaseNotes.String(nmVersion, "", "release version")
	// Release sub command
	cmdRelease                = flag.NewFlagSet(types.SubCmdRelease, flag.ExitOnError)
	cmdReleaseToken           = cmdRelease.String(nmToken, "", "github token")
	cmdReleaseReleaseNoteRepo = cmdRelease.String(nmReleaseNoteRepo, "", "release note repo")
	cmdReleaseVersion         = cmdRelease.String(nmVersion, "", "release version")
	// Build sub command
	cmdBuild                = flag.NewFlagSet(types.SubCmdBuild, flag.ExitOnError)
	cmdBuildToken           = cmdBuild.String(nmToken, "", "github token")
	cmdBuildReleaseNoteRepo = cmdBuild.String(nmReleaseNoteRepo, "", "release note repo")
	cmdBuildVersion         = cmdBuild.String(nmVersion, "", "release version")
)

func init() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Printf("expected %s, %s, %s or %s subcommands\n",
			types.SubCmdPRList,
			types.SubCmdReleaseNotes,
			types.SubCmdRelease,
			types.SubCmdBuild)
		os.Exit(1)
	}

	switch os.Args[1] {
	case types.SubCmdPRList:
		subCommand = types.SubCmdPRList
		cmdPRList.Parse(os.Args[2:])
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
	case types.SubCmdReleaseNotes:
		subCommand = types.SubCmdReleaseNotes
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
	case types.SubCmdRelease:
		subCommand = types.SubCmdRelease
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
	case types.SubCmdBuild:
		subCommand = types.SubCmdBuild
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("expected %s, %s, %s or %s subcommands\n",
			types.SubCmdPRList,
			types.SubCmdReleaseNotes,
			types.SubCmdRelease,
			types.SubCmdBuild)
		os.Exit(1)
	}

	m, err := manager.New(&manager.Config{
		SubCommand:      subCommand,
		Version:         version,
		GithubToken:     token,
		ReleaseNoteRepo: releaseNoteRepo,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Run(); err != nil {
		log.Fatal(err)
	}
}

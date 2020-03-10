package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/you06/releaser/config"
	"github.com/you06/releaser/manager"
	"github.com/you06/releaser/pkg/types"
)

const (
	nmToken           = "token"
	nmReleaseNoteRepo = "repo"
	nmVersion         = "version"
	nmConfig          = "config"
)

var (
	// common args
	token           string
	releaseNoteRepo string
	version         string
	configPath      string
	subCommand      string
	// PRList sub command
	cmdPRList                = flag.NewFlagSet(types.SubCmdPRList, flag.ExitOnError)
	cmdPRListToken           = cmdPRList.String(nmToken, "", "github token")
	cmdPRListReleaseNoteRepo = cmdPRList.String(nmReleaseNoteRepo, "", "release note repo")
	cmdPRListVersion         = cmdPRList.String(nmVersion, "", "release version")
	cmdPRListConfig          = cmdPRList.String(nmConfig, "", "config path")
	// ReleaseNotes sub command
	cmdReleaseNotes                = flag.NewFlagSet(types.SubCmdReleaseNotes, flag.ExitOnError)
	cmdReleaseNotesToken           = cmdReleaseNotes.String(nmToken, "", "github token")
	cmdReleaseNotesReleaseNoteRepo = cmdReleaseNotes.String(nmReleaseNoteRepo, "", "release note repo")
	cmdReleaseNotesVersion         = cmdReleaseNotes.String(nmVersion, "", "release version")
	cmdReleaseNotesConfig          = cmdReleaseNotes.String(nmConfig, "", "config path")
	// Release sub command
	cmdRelease                = flag.NewFlagSet(types.SubCmdRelease, flag.ExitOnError)
	cmdReleaseToken           = cmdRelease.String(nmToken, "", "github token")
	cmdReleaseReleaseNoteRepo = cmdRelease.String(nmReleaseNoteRepo, "", "release note repo")
	cmdReleaseVersion         = cmdRelease.String(nmVersion, "", "release version")
	cmdReleaseConfig          = cmdRelease.String(nmConfig, "", "config path")
	// Build sub command
	cmdBuild                = flag.NewFlagSet(types.SubCmdBuild, flag.ExitOnError)
	cmdBuildToken           = cmdBuild.String(nmToken, "", "github token")
	cmdBuildReleaseNoteRepo = cmdBuild.String(nmReleaseNoteRepo, "", "release note repo")
	cmdBuildVersion         = cmdBuild.String(nmVersion, "", "release version")
	cmdBuildConfig          = cmdBuild.String(nmConfig, "", "config path")
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
		cmdPRList.Parse(os.Args[2:])
		subCommand = types.SubCmdPRList
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
		configPath = *cmdPRListConfig
	case types.SubCmdReleaseNotes:
		cmdReleaseNotes.Parse(os.Args[2:])
		subCommand = types.SubCmdReleaseNotes
		token = *cmdPRListToken
		releaseNoteRepo = *cmdPRListReleaseNoteRepo
		version = *cmdPRListVersion
		configPath = *cmdPRListConfig
	case types.SubCmdRelease:
		cmdRelease.Parse(os.Args[2:])
		subCommand = types.SubCmdRelease
		token = *cmdReleaseToken
		releaseNoteRepo = *cmdReleaseReleaseNoteRepo
		version = *cmdReleaseVersion
		configPath = *cmdReleaseConfig
	case types.SubCmdBuild:
		cmdBuild.Parse(os.Args[2:])
		subCommand = types.SubCmdBuild
		token = *cmdBuildToken
		releaseNoteRepo = *cmdBuildReleaseNoteRepo
		version = *cmdBuildVersion
		version = *cmdBuildVersion
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

	cfg := config.New()
	if err := cfg.Read(configPath); err != nil {

	}

	m, err := manager.New(cfg, &manager.Option{
		Version: version,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Run(subCommand); err != nil {
		log.Fatal(err)
	}
}

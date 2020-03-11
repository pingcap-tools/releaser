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
	version    string
	configPath string
	subCommand string
	// PRList sub command
	cmdPRList        = flag.NewFlagSet(types.SubCmdPRList, flag.ExitOnError)
	cmdPRListVersion = cmdPRList.String(nmVersion, "", "release version")
	cmdPRListConfig  = cmdPRList.String(nmConfig, "", "config path")
	// ReleaseNotes sub command
	cmdReleaseNotes        = flag.NewFlagSet(types.SubCmdReleaseNotes, flag.ExitOnError)
	cmdReleaseNotesVersion = cmdReleaseNotes.String(nmVersion, "", "release version")
	cmdReleaseNotesConfig  = cmdReleaseNotes.String(nmConfig, "", "config path")
	// Release sub command
	cmdCheckModule        = flag.NewFlagSet(types.SubCmdCheckModule, flag.ExitOnError)
	cmdCheckModuleVersion = cmdCheckModule.String(nmVersion, "", "release version")
	cmdCheckModuleConfig  = cmdCheckModule.String(nmConfig, "", "config path")
	// GenerateReleaseNote sub command
	cmdGenerateReleaseNote        = flag.NewFlagSet(types.SubCmdGenerateReleaseNote, flag.ExitOnError)
	cmdGenerateReleaseNoteVersion = cmdGenerateReleaseNote.String(nmVersion, "", "release version")
	cmdGenerateReleaseNoteConfig  = cmdGenerateReleaseNote.String(nmConfig, "", "config path")
)

func init() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Printf("expected %s, %s or %s subcommands\n",
			types.SubCmdPRList,
			types.SubCmdReleaseNotes,
			types.SubCmdCheckModule)
		os.Exit(1)
	}

	switch os.Args[1] {
	case types.SubCmdPRList:
		cmdPRList.Parse(os.Args[2:])
		subCommand = types.SubCmdPRList
		version = *cmdPRListVersion
		configPath = *cmdPRListConfig
	case types.SubCmdReleaseNotes:
		cmdReleaseNotes.Parse(os.Args[2:])
		subCommand = types.SubCmdReleaseNotes
		version = *cmdReleaseNotesVersion
		configPath = *cmdReleaseNotesConfig
	case types.SubCmdCheckModule:
		cmdCheckModule.Parse(os.Args[2:])
		subCommand = types.SubCmdCheckModule
		version = *cmdCheckModuleVersion
		configPath = *cmdCheckModuleConfig
	case types.SubCmdGenerateReleaseNote:
		cmdGenerateReleaseNote.Parse(os.Args[2:])
		subCommand = types.SubCmdGenerateReleaseNote
		version = *cmdGenerateReleaseNoteVersion
		configPath = *cmdGenerateReleaseNoteConfig
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("expected %s, %s or %s subcommands\n",
			types.SubCmdPRList,
			types.SubCmdReleaseNotes,
			types.SubCmdCheckModule)
		os.Exit(1)
	}

	cfg := config.New()
	if err := cfg.Read(configPath); err != nil {
		log.Fatalf("%+v", err)
	}
	// cfg.Print()

	m, err := manager.New(cfg, &manager.Option{
		Version: version,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := m.Run(subCommand); err != nil {
		log.Fatal(err)
	}
}

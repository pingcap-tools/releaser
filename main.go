package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/manager"
	"github.com/you06/releaser/pkg/types"
)

const (
	nmVersion = "version"
	nmConfig  = "config"
)

var (
	// common args
	version    string
	configPath string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "releaser",
		Short: "Releaser is a tool which helps you with your release notes",
		Long: "Releaser is a tool which helps you with your release notes." +
			"\nsee more from https://github.com/you06/releaser",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("expected %s, %s or %s subcommands\n",
				types.SubCmdPRList,
				types.SubCmdReleaseNotes,
				types.SubCmdCheckModule)
		},
	}

	var subCmdPRListCmd = &cobra.Command{
		Use:   types.SubCmdPRList,
		Short: "Generate release notes from milestone",
		Run: func(cmd *cobra.Command, args []string) {
			runWithSubCommand(types.SubCmdPRList)
		},
	}

	var generateReleaseNoteCmd = &cobra.Command{
		Use:   types.SubCmdGenerateReleaseNote,
		Short: "Generate release notes from milestone",
		Run: func(cmd *cobra.Command, args []string) {
			runWithSubCommand(types.SubCmdGenerateReleaseNote)
		},
	}

	rootCmd.AddCommand(subCmdPRListCmd)
	rootCmd.AddCommand(generateReleaseNoteCmd)

	rootCmd.PersistentFlags().StringVar(&configPath, nmConfig, "./config.toml", "config file")
	rootCmd.PersistentFlags().StringVar(&version, nmVersion, "", "release version")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runWithSubCommand(cmd string) {
	cfg := config.New()
	if err := cfg.Read(configPath); err != nil {
		log.Fatalf("%+v", err)
	}

	m, err := manager.New(cfg, &manager.Option{
		Version: version,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := m.Run(cmd); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:   "breachguard",
	Short: "BreachGuard checks email addresses against breach data",
	Long: `BreachGuard is a CLI tool for scanning email addresses
and reporting known data breach exposure.

Built by Dzounz Eedriz`,
	Version: version,
}

func Execute() {
	rootCmd.SetVersionTemplate("BreachGuard {{.Version}}\nBuilt by Dzounz Eedriz\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

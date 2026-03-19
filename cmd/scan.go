package cmd

import (
	"breachguard/internal/config"
	"breachguard/internal/hibp"
	"breachguard/internal/output"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	scanDelay    time.Duration
	onlyBreached bool
	dateFormat   string
	noColor      bool
	demoMode     bool
	outputFormat string
)

func init() {
	scanCmd := &cobra.Command{
		Use:   "scan <file>",
		Short: "Scan a file of email addresses",
		Args:  cobra.ExactArgs(1),
		RunE:  runScan,
	}

	scanCmd.Flags().DurationVar(&scanDelay, "delay", 7*time.Second, "Delay between requests")
	scanCmd.Flags().BoolVar(&onlyBreached, "only-breached", false, "Show only breached emails")
	scanCmd.Flags().StringVar(&dateFormat, "date-format", "year", "Date format: year, month, full")
	scanCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	scanCmd.Flags().BoolVar(&demoMode, "demo", false, "Use demo mode instead of live API")
	scanCmd.Flags().StringVar(&outputFormat, "output", "table", "Output format: table, json, csv")

	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	cfg := config.Load()
	client := hibp.NewClient(cfg.APIKey, demoMode)

	scanner := bufio.NewScanner(f)
	seen := make(map[string]struct{})
	var results []hibp.Result
	lineNo := 0

	for scanner.Scan() {
		email := strings.TrimSpace(scanner.Text())
		if email == "" {
			continue
		}

		email = strings.ToLower(email)
		if _, exists := seen[email]; exists {
			continue
		}
		seen[email] = struct{}{}
		lineNo++

		fmt.Printf("[%d] Checking %s...\n", lineNo, email)

		result, err := client.CheckEmail(email, dateFormat)
		if err != nil {
			fmt.Printf("  error: %v\n", err)
			time.Sleep(scanDelay)
			continue
		}

		if onlyBreached && !result.Breached {
			time.Sleep(scanDelay)
			continue
		}

		results = append(results, result)
		time.Sleep(scanDelay)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	switch strings.ToLower(outputFormat) {
	case "table":
		output.PrintTable(results, !noColor)
	case "json":
		return output.PrintJSON(results)
	case "csv":
		return output.PrintCSV(results)
	default:
		return fmt.Errorf("invalid output format: %s (use table, json, or csv)", outputFormat)
	}

	return nil
}

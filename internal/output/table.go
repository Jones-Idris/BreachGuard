package output

import (
	"breachguard/internal/hibp"
	"fmt"
	"strings"
)

func PrintTable(results []hibp.Result, useColor bool) {
	headers := []string{"EMAIL", "BREACHED", "COUNT", "BREACHES"}

	emailW := len(headers[0])
	breachedW := len(headers[1])
	countW := len(headers[2])
	breachesW := len(headers[3])

	breachedCount := 0
	rows := make([][]string, 0, len(results))

	for _, r := range results {
		statusPlain := "NO"
		statusDisplay := "NO"
		breachList := "-"

		if r.Breached {
			statusPlain = "YES"
			statusDisplay = "YES"
			breachList = shortenBreaches(r.Breaches, 3)
			breachedCount++
		}

		if useColor {
			if r.Breached {
				statusDisplay = red(statusDisplay)
			} else {
				statusDisplay = green(statusDisplay)
			}
		}

		emailW = max(emailW, len(r.Email))
		breachedW = max(breachedW, len(statusPlain))
		countW = max(countW, len(fmt.Sprintf("%d", r.Count)))
		breachesW = max(breachesW, len(breachList))

		rows = append(rows, []string{
			r.Email,
			statusDisplay,
			fmt.Sprintf("%d", r.Count),
			breachList,
		})
	}

	widths := []int{emailW, breachedW, countW, breachesW}

	printSeparator(widths)
	printRow(headers, widths)
	printSeparator(widths)

	for _, row := range rows {
		printRow(row, widths)
	}

	printSeparator(widths)
	fmt.Println()
	fmt.Printf("Summary: Checked: %d | Breached: %d | Safe: %d\n",
		len(results), breachedCount, len(results)-breachedCount)
}

func shortenBreaches(breaches []string, maxShown int) string {
	if len(breaches) == 0 {
		return "-"
	}
	if len(breaches) <= maxShown {
		return strings.Join(breaches, ", ")
	}

	shown := strings.Join(breaches[:maxShown], ", ")
	remaining := len(breaches) - maxShown
	return fmt.Sprintf("%s (+%d more)", shown, remaining)
}

func printSeparator(widths []int) {
	fmt.Print("+")
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2) + "+")
	}
	fmt.Println()
}

func printRow(cols []string, widths []int) {
	fmt.Print("|")
	for i, c := range cols {
		fmt.Printf(" %s |", pad(c, widths[i]))
	}
	fmt.Println()
}

func pad(s string, width int) string {
	if len(stripANSI(s)) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(stripANSI(s)))
}

func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == 0x1b {
			inEsc = true
			continue
		}
		if inEsc {
			if ch == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}

func green(text string) string {
	return "\033[32m" + text + "\033[0m"
}

func red(text string) string {
	return "\033[31m" + text + "\033[0m"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

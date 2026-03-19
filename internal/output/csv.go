package output

import (
	"breachguard/internal/hibp"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

func PrintCSV(results []hibp.Result) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if err := w.Write([]string{"email", "breached", "count", "breaches"}); err != nil {
		return err
	}

	for _, r := range results {
		record := []string{
			r.Email,
			strconv.FormatBool(r.Breached),
			strconv.Itoa(r.Count),
			strings.Join(r.Breaches, "; "),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return w.Error()
}

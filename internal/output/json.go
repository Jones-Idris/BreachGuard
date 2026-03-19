package output

import (
	"breachguard/internal/hibp"
	"encoding/json"
	"os"
)

func PrintJSON(results []hibp.Result) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

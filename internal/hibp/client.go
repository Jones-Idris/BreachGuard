package hibp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Breach struct {
	Name       string `json:"Name"`
	BreachDate string `json:"BreachDate"`
}

type Result struct {
	Email    string
	Breached bool
	Count    int
	Breaches []string
}

type Client struct {
	apiKey string
	demo   bool
	http   *http.Client
}

func NewClient(apiKey string, demo bool) *Client {
	return &Client{
		apiKey: apiKey,
		demo:   demo,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) CheckEmail(email string, dateFormat string) (Result, error) {
	return c.checkWithRetry(email, dateFormat, 2) // max 2 retries
}

func (c *Client) checkWithRetry(email string, dateFormat string, retries int) (Result, error) {
	if c.demo {
		return checkEmailDemo(email), nil
	}

	if strings.TrimSpace(c.apiKey) == "" {
		return Result{}, fmt.Errorf("HIBP_API_KEY is not set")
	}

	endpoint := "https://haveibeenpwned.com/api/v3/breachedaccount/" +
		url.PathEscape(email) + "?truncateResponse=false"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return Result{}, err
	}

	req.Header.Set("hibp-api-key", c.apiKey)
	req.Header.Set("user-agent", "BreachGuard")

	resp, err := c.http.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		var breaches []Breach
		if err := json.NewDecoder(resp.Body).Decode(&breaches); err != nil {
			return Result{}, err
		}

		names := make([]string, 0, len(breaches))
		for _, b := range breaches {
			names = append(names, formatBreach(b, dateFormat))
		}

		return Result{
			Email:    email,
			Breached: true,
			Count:    len(names),
			Breaches: names,
		}, nil

	case http.StatusNotFound:
		io.Copy(io.Discard, resp.Body)
		return Result{
			Email:    email,
			Breached: false,
			Count:    0,
			Breaches: nil,
		}, nil

	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")

		wait := 5 * time.Second
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				wait = time.Duration(seconds) * time.Second
			}
		}

		if retries > 0 {
			fmt.Printf("Rate limited. Waiting %v then retrying...\n", wait)
			time.Sleep(wait)
			return c.checkWithRetry(email, dateFormat, retries-1)
		}

		return Result{}, fmt.Errorf("rate limited after retries")

	case http.StatusUnauthorized:
		body, _ := io.ReadAll(resp.Body)
		return Result{}, fmt.Errorf("unauthorized: %s", strings.TrimSpace(string(body)))

	default:
		body, _ := io.ReadAll(resp.Body)
		return Result{}, fmt.Errorf("unexpected HIBP response %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}
}

func formatBreach(b Breach, mode string) string {
	date := b.BreachDate

	switch mode {
	case "month":
		if len(date) >= 7 {
			return fmt.Sprintf("%s (%s)", b.Name, date[:7])
		}
	case "full":
		if len(date) >= 10 {
			return fmt.Sprintf("%s (%s)", b.Name, date[:10])
		}
	default:
		if len(date) >= 4 {
			return fmt.Sprintf("%s (%s)", b.Name, date[:4])
		}
	}

	return b.Name
}

func checkEmailDemo(email string) Result {
	l := strings.ToLower(email)

	if strings.Contains(l, "test") || strings.Contains(l, "hack") || strings.Contains(l, "pwn") {
		return Result{
			Email:    email,
			Breached: true,
			Count:    2,
			Breaches: []string{"Adobe (2013)", "LinkedIn (2012)"},
		}
	}

	return Result{
		Email:    email,
		Breached: false,
		Count:    0,
		Breaches: nil,
	}
}

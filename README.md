# BreachGuard

BreachGuard is a Go CLI tool for checking email addresses against known data breaches using the Have I Been Pwned API.

Built by Dzounz Eedriz.

## Features

- Batch email scanning from file
- Table, JSON, and CSV output
- Automatic retry on rate limits
- Breach date formatting (year/month/full)
- Colored terminal output
- Version command

## Installation

```bash
git clone https://github.com/Jones-Idris/BreachGuard.git
cd BreachGuard
go build -o breachguard
```

## Usage

 You must provide your Have I Been Pwned API KEY:

```bash
export HIBP_API_KEY='your_api_key' 
```

Run a scan:

./breachguard scan <file> [flags]

Example: 
 
```bash
./breachguard scan emails.txt --delay 7s --output table
```

## Flags

flag				Description			Default

--delay			Delay between API requests		7s
--output		Output format (table,json,csv)		table
--only-breached		Show only breached emails 		false
--demo 			Use only demo mode (no api calls)	false

## Output Format 

Table (Default):

Human readable format with summarized breach data.

Example:

+----------------------+----------+-------+-----------------------------------------------+
| EMAIL                | BREACHED | COUNT | BREACHES                                      |
+----------------------+----------+-------+-----------------------------------------------+
| user@example.com     | YES      | 3     | Adobe (2013), LinkedIn (2012), Dropbox (+1)   |
| safe@email.com       | NO       | 0     | -                                             |
+----------------------+----------+-------+-----------------------------------------------+







^G Help          ^O Write Out     ^F Where Is      ^K Cut           ^T Execute       ^C Location      M-U Undo         M-A Set Mark     M-] To Bracket   M-B Previous     ◂ Back
^X Exit          ^R Read File     ^\ Replace       ^U Paste         ^J Justify       ^/ Go To Line    M-E Redo         M-6 Copy         ^B Where Was 





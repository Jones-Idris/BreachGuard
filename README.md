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

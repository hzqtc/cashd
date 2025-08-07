package ledger

import (
	"cashd/internal/data"
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

type LedgerDataSource struct{}

var ledgerFilePath string = func() string {
	if ledgerFileFlag != "" {
		return ledgerFileFlag
	} else if env := os.Getenv("LEDGER_FILE"); env != "" {
		return env
	} else if env := os.Getenv("HLEDGER_FILE"); env != "" {
		return env
	} else {
		return ""
	}
}()

var ledgerFileFlag string

func init() {
	pflag.StringVar(&ledgerFileFlag, "ledger", "", "Ledger file path")
}

func (l LedgerDataSource) LoadTransactions() ([]*data.Transaction, error) {
	commands := []string{"ledger", "hledger"}

	for _, cmd := range commands {
		cmd := exec.Command(cmd, "-f", ledgerFilePath, "print")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		if err := cmd.Start(); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				continue
			} else {
				return nil, err
			}
		}

		// Stream output to parser
		transactions, parseErr := parseJournal(stdout)
		// Wait for command to complete
		if err := cmd.Wait(); err != nil {
			return nil, err
		} else if parseErr != nil {
			return nil, parseErr
		} else {
			return transactions, nil
		}
	}

	return nil, os.ErrNotExist
}

func (l LedgerDataSource) Preferred() bool {
	return ledgerFileFlag != ""
}

func (l LedgerDataSource) Enabled() bool {
	return ledgerFilePath != ""
}

package ledger

import (
	"errors"
	"lledger/internal/data"
	"os"
	"os/exec"
)

type LedgerDataSource struct{}

func (l LedgerDataSource) LoadTransactions() ([]data.Transaction, error) {
	commands := []string{"ledger", "hledger"}

	for _, cmd := range commands {
		cmd := exec.Command(cmd, "print")

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

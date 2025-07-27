package data

import (
	"os/exec"
)

func LoadTransactions() ([]Transaction, error) {
	cmd := exec.Command("hledger", "print")

	// Get a pipe to the command's standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	transactions, parseErr := parseJournal(stdout)
	if err := cmd.Wait(); err != nil {
		return nil, err
	} else if parseErr != nil {
		return nil, parseErr
	} else {
		return transactions, nil
	}
}

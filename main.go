package main

import (
	"cashd/internal/model"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

var flagShowHelp = pflag.BoolP("help", "h", false, "Show help message")

func main() {
	pflag.Parse()

	if *flagShowHelp {
		pflag.Usage()
		os.Exit(0)
	}

	f, err := os.OpenFile("/tmp/cashd.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("failed to create log file: %v", err)
	}
	defer f.Close()
	// Send log output to the file
	log.SetOutput(f)

	p := tea.NewProgram(model.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

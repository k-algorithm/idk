package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/search"
	"github.com/k-algorithm/idk/internal/tui"
)

var supportedOS = map[string]bool{
	"linux":   true,
	"freebsd": true,
	"openbsd": true,
	"netbsd":  true,
	"darwin":  true,
}

func main() {
	result := search.Google(search.GoogleParam{
		Query:    "golang",
		PageSize: 10,
	})
	if len(result.QuestionIDs) == 0 {
		log.Println("No results..")
		return
	}

	log.Println(runtime.GOOS)
	if !supportedOS[runtime.GOOS] {
		fmt.Println("Current OS not supported. Refer to the documentation for more information.")
		os.Exit(0)
	}
	p := tea.NewProgram(tui.InitializeModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

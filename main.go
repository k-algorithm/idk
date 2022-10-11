package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/tui"
)

var supportedOS = map[string]bool{
	"linux":   true,
	"freebsd": true,
	"openbsd": true,
	"netbsd":  true,
	"darwin":  true,
	"windows": true,
}

func main() {
	// parse cli flags
	dbgModePtr := flag.Bool("dbg_mode", false, "Run debug mode.")
	flag.Parse()
	if *dbgModePtr {
		log.Println("DebugMode entered.")
	}

	if !supportedOS[runtime.GOOS] {
		fmt.Printf("Current OS (%s) not supported. Refer to the documentation for more information.", runtime.GOOS)
		os.Exit(0)
	}
	p := tea.NewProgram(tui.InitializeModel(*dbgModePtr), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

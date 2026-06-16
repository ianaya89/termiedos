package main

import (
	"fmt"
	"os"

	"termiedos/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Println("termiedos", version)
			return
		case "--help", "-h", "help":
			fmt.Println("termiedos — resultados de fútbol en la terminal (promiedos)")
			fmt.Println("\nUso: termiedos")
			fmt.Println("\nFlags:\n  -v, --version   muestra la versión\n  -h, --help      muestra esta ayuda")
			fmt.Println("\nTeclas: ↑↓ mover · ←→ día · t hoy · tab panel · enter abrir · b atrás · r recargar · q salir")
			return
		}
	}
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

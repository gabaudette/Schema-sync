package terminal

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"log"
	"os"
	"path/filepath"
)

func SetupTerminal() {
	_, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	ClearTerminal()

	err = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Schema Sync", pterm.FgLightBlue.ToStyle())).Render()
	pterm.Println()

	if err != nil {
		panic(err)
	}

	pterm.Println()
}

func ClearTerminal() {
	fmt.Println("\033[2J")
}

func BoolToText(value bool) string {
	if value {
		return pterm.Green("Yes")
	}

	return pterm.Red("No")
}

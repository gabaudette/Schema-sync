package main

import (
	"github.com/pterm/pterm"
	"ssync/internal/terminal"
	"ssync/internal/tsgenerator/api"
)

func main() {
	terminal.SetupTerminal()
	api.ResetTypeFile()
	api.ListFiles()

	pterm.Success.Prefix = pterm.Prefix{Text: "SYNC COMPLETE", Style: pterm.NewStyle(pterm.BgLightBlue, pterm.FgBlack)}
	pterm.Success.Println("All files have been synced successfully!")

	// TODO Re-add the pause function when the program is ready for production
	//// Pause the program to prevent the terminal from closing useful for Iterm2
	//_, err = fmt.Scanln()
	//if err != nil {
	//	return
	//}
}

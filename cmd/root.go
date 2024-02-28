/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomodoro-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

const (
	WORK_TIMER        = 25 * 60
	SHORT_BREAK_TIMER = 5 * 60
	LONG_BREAK_TIMER  = 15 * 60
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	activeTimer := WORK_TIMER

	ticker := time.NewTicker(time.Second)

	if err != nil {
		fmt.Printf("Error setting raw terminal: %v\n", err)
		return
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	input := make(chan byte)

	go func() {
		paused := true
		var formattedTime string

		for _, ok := <-input; ok; fmt.Print("\r\x1B[0K") {
			formattedTime = fmt.Sprintf("%2d:%02d", activeTimer/60, activeTimer%60)
			fmt.Print("\r\x1B[0K", formattedTime)
			if paused {
				fmt.Print("    (paused)")
			}

			select {
			case <-input:
				paused = !paused
			case <-ticker.C:
				if !paused {
					activeTimer -= 1
					if activeTimer == 0 {
						ticker.Stop()
					}
				}
			}
		}
	}()

	input <- 's' // display the timer right away, without user input

	var buf [1]byte

	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			return
		}

		if buf[0] == 's' || buf[0] == 'S' {
			input <- buf[0]
		} else {
			close(input)
			fmt.Print("\r\x1B[0K")
			return
		}
		time.Sleep(time.Millisecond)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pomodoro-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

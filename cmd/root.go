/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ian-tang/pomodoro-cli/cmd/timer"

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

var initialTimerState = timer.PausedTimerState{
	Timer: timer.Timer{
		TimerType:     timer.FOCUS_TIMER,
		TimeRemaining: timer.TimerDuration[timer.FOCUS_TIMER],
		PomodoroCount: 0,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Printf("Error setting raw terminal: %v\n", err)
		return
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	timerState := timer.TimerState(initialTimerState)
	ticker := time.NewTicker(time.Second)

	input := make(chan byte)

	go func() {
		var formattedTime string

		for ; ; fmt.Print("\r\x1B[0K") {
			formattedTime = timerState.GetFormattedTimeString()
			fmt.Print("\r\x1B[0K", formattedTime)

			select {
			case <-input:
				timerState = timerState.Pause()
			case <-ticker.C:
				timerState = timerState.Tick()
			}
		}
	}()

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

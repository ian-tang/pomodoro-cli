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
		PomodoroCount: 1,
	},
}

const (
	LOWERCASE_S = iota
	LOWERCASE_R
	LOWERCASE_F
	LOWERCASE_Q
)

var validInputKeys = map[byte]int{
	's': LOWERCASE_S,
	'r': LOWERCASE_R,
	'f': LOWERCASE_F,
	'q': LOWERCASE_Q,
}

const inputHelpMessage = "[s] start/stop [t] adjust timers [a] add task\n\r[r] reset current timer [R] reset progress [f] skip current timer\n\r[q] quit\r"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Printf("Error setting raw terminal: %v\n", err)
		return
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)
	// make cursor visible and clear screen of text
	defer fmt.Print("\x1B[?25h" + "\r\x1B[3A\x1B[J")

	timerState := timer.TimerState(initialTimerState)
	ticker := time.NewTicker(time.Second)

	input := make(chan byte)

	go func() {
		var formattedTime string
		// make cursor invisible and insert 3 new lines
		fmt.Print("\x1B[?25l\n\n\n")

		for {
			formattedTime = timerState.GetFormattedTimeString()
			// move cursor to left of screen and up 3 rows, then erase from cursor to end of screen
			fmt.Print("\r\x1B[3A\x1B[J", formattedTime, "\n\r", inputHelpMessage)

			select {
			case input := <-input:
				timerState = handleUserInput(timerState, input)
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
		if buf[0] == 'q' {
			return
		}
		input <- buf[0]

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

func handleUserInput(t timer.TimerState, input byte) timer.TimerState {
	inputKey := validInputKeys[input]

	switch inputKey {
	case LOWERCASE_S:
		return t.Pause()
	case LOWERCASE_F:
		return t.SkipCurrentTimer()
	case LOWERCASE_R:
		return t.ResetCurrentTimer()
	default:
		return t
	}
}

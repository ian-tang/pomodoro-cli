/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
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

var defaultTimer = timer.Timer{
	TimerType:     timer.FOCUS_TIMER,
	TimeRemaining: timer.DEFAULT_FOCUS_TIMER_DURATION,
	PomodoroCount: 1,
}
var defaultTimerState timer.TimerState = timer.TSPool.Paused

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

type TimerSave struct {
	timer.Timer         `json:"Timer"`
	timer.TimerDuration `json:"TimerDuration"`
}

var NIL_TIMER = timer.Timer{}

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

	var timerSaveFile TimerSave
	t := defaultTimer
	timerState := defaultTimerState
	saveFileName := "./data.json"

	func() {
		var data []byte
		var err error
		data, err = os.ReadFile(saveFileName)
		if err != nil {
			fmt.Println("Error opening saved session data:", err)
			return
		}
		err = json.Unmarshal(data, &timerSaveFile)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		if timerSaveFile.Timer != NIL_TIMER {
			t = timerSaveFile.Timer
			timerState = timer.TSPool.Paused
		}
		timer.SetTimerDuration(timer.FOCUS_TIMER, timerSaveFile.Focus)
		timer.SetTimerDuration(timer.SHORT_BREAK_TIMER, timerSaveFile.ShortBreak)
		timer.SetTimerDuration(timer.LONG_BREAK_TIMER, timerSaveFile.LongBreak)
	}()

	ticker := time.NewTicker(time.Second / timer.TICKS_PER_SECOND)
	input := make(chan byte)

	go func() {
		var formattedTime string
		// make cursor invisible and insert 3 new lines
		fmt.Print("\x1B[?25l\n\n\n")

		for {
			formattedTime = timerState.GetFormattedTimeString(&t)
			// move cursor to left of screen and up 3 rows, then erase from cursor to end of screen
			fmt.Print("\r\x1B[3A\x1B[J", formattedTime, "\n\r", inputHelpMessage)

			select {
			case inputKey := <-input:
				if inputKey == 'q' {
					timerDurations := timer.GetTimerDurations()
					saveTimerDurations := timer.TimerDuration{
						Focus:      timerDurations[timer.FOCUS_TIMER] / timer.TICKS_PER_MINUTE,
						ShortBreak: timerDurations[timer.SHORT_BREAK_TIMER] / timer.TICKS_PER_MINUTE,
						LongBreak:  timerDurations[timer.LONG_BREAK_TIMER] / timer.TICKS_PER_MINUTE,
					}

					buf, err := json.MarshalIndent(TimerSave{
						t,
						saveTimerDurations,
					}, "", "  ")
					if err != nil {
						fmt.Println("Error marshalling JSON:", err)
					}
					os.WriteFile(saveFileName, buf, 0644)
					input <- '0'
				}

				handleUserInput(&t, &timerState, inputKey)
			case <-ticker.C:
				timerState.Tick(&t, &timerState)
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
			input <- buf[0]
			<-input // wait for confirmation that shutdown procedure is complete
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

func handleUserInput(t *timer.Timer, ts *timer.TimerState, input byte) {
	inputKey, ok := validInputKeys[input]

	if !ok {
		return
	}

	switch inputKey {
	case LOWERCASE_S:
		(*ts).Pause(t, ts)
	case LOWERCASE_F:
		(*ts).SkipCurrentTimer(t, ts)
	case LOWERCASE_R:
		(*ts).ResetCurrentTimer(t, ts)
	default:
		return
	}
}

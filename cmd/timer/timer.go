package timer

import (
	"fmt"
)

const TICKS_PER_SECOND = 100
const TICKS_PER_MINUTE = TICKS_PER_SECOND * 60

const (
	FOCUS_TIMER = iota
	SHORT_BREAK_TIMER
	LONG_BREAK_TIMER
)

const (
	DEFAULT_FOCUS_TIMER_DURATION       = 25 * TICKS_PER_MINUTE
	DEFAULT_SHORT_FOCUS_TIMER_DURATION = 5 * TICKS_PER_MINUTE
	DEFAULT_LONG_BREAK_TIMER_DURATION  = 15 * TICKS_PER_MINUTE
)

var timerDuration = map[int]int{
	FOCUS_TIMER:       DEFAULT_FOCUS_TIMER_DURATION,
	SHORT_BREAK_TIMER: DEFAULT_SHORT_FOCUS_TIMER_DURATION,
	LONG_BREAK_TIMER:  DEFAULT_LONG_BREAK_TIMER_DURATION,
}

type TimerDuration struct {
	Focus      int
	ShortBreak int
	LongBreak  int
}

type TimerState interface {
	Tick() TimerState
	NextTimer() TimerState
	Pause() TimerState
	SkipCurrentTimer() TimerState
	ResetCurrentTimer() TimerState
	GetFormattedTimeString() string
	GetCurrentTimerState() Timer
}

type Timer struct {
	TimerType     int
	TimeRemaining int
	PomodoroCount int
}

type RunningTimerState struct {
	Timer
}

type PausedTimerState struct {
	Timer
}

type BetweenTimerState struct {
	Timer
}

func (t RunningTimerState) Tick() TimerState {
	t.TimeRemaining--

	if t.TimeRemaining <= 0 {
		return t.NextTimer()
	}

	return RunningTimerState{
		t.Timer,
	}
}

func (t RunningTimerState) NextTimer() TimerState {
	return BetweenTimerState(t)
}

func (t RunningTimerState) Pause() TimerState {
	return PausedTimerState(t)
}

func (t RunningTimerState) SkipCurrentTimer() TimerState {
	return t.NextTimer()
}

func (t RunningTimerState) ResetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: timerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t RunningTimerState) GetFormattedTimeString() string {
	return fmt.Sprintf("#%d %2d:%02d", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
}

func (t RunningTimerState) GetCurrentTimerState() Timer {
	return t.Timer
}

func (t PausedTimerState) Tick() TimerState {
	return t
}

func (t PausedTimerState) NextTimer() TimerState {
	return BetweenTimerState(t)
}

func (t PausedTimerState) Pause() TimerState {
	return RunningTimerState(t)
}

func (t PausedTimerState) SkipCurrentTimer() TimerState {
	return t.NextTimer()
}

func (t PausedTimerState) ResetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: timerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t PausedTimerState) GetFormattedTimeString() string {
	return fmt.Sprintf("#%d %2d:%02d    (paused)", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
}

func (t PausedTimerState) GetCurrentTimerState() Timer {
	return t.Timer
}

func (t BetweenTimerState) Tick() TimerState {
	return t
}

func (t BetweenTimerState) NextTimer() TimerState {
	timerType, timerDuration, pomodoroCount := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: pomodoroCount,
		},
	}
}

func (t BetweenTimerState) Pause() TimerState {
	timerType, timerDuration, pomodoroCount := t.getNextTimerType()

	return RunningTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: pomodoroCount,
		},
	}
}

func (t BetweenTimerState) SkipCurrentTimer() TimerState {
	return t
}

func (t BetweenTimerState) ResetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: timerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t BetweenTimerState) GetFormattedTimeString() string {
	timerType, _, pomodoroCount := t.getNextTimerType()

	switch timerType {
	case SHORT_BREAK_TIMER:
		if pomodoroCount == 1 {
			return "1 pomodoro done! Start short break?"
		}

		return fmt.Sprintf("%d pomodoros done! Start short break?", t.PomodoroCount)
	case LONG_BREAK_TIMER:
		return fmt.Sprintf("%d pomodoros done! Start long break?", t.PomodoroCount)
	default:
		return "Break over! Start pomodoro?"
	}
}

func (t BetweenTimerState) GetCurrentTimerState() Timer {
	return t.Timer
}

func (t Timer) getNextTimerType() (int, int, int) {
	if t.TimerType == FOCUS_TIMER && t.PomodoroCount%4 == 0 {
		return LONG_BREAK_TIMER, timerDuration[LONG_BREAK_TIMER], t.PomodoroCount
	}

	if t.TimerType == FOCUS_TIMER {
		return SHORT_BREAK_TIMER, timerDuration[SHORT_BREAK_TIMER], t.PomodoroCount
	}

	return FOCUS_TIMER, timerDuration[FOCUS_TIMER], t.PomodoroCount + 1
}

func SetTimerDuration(timerType int, duration int) error {
	if timerType > 2 {
		return fmt.Errorf("invalid timer type: %v", timerType)
	}
	if duration <= 0 {
		return fmt.Errorf("cannot set timer duration to %v", duration)
	}

	timerDuration[timerType] = duration * TICKS_PER_MINUTE
	return nil
}

func GetTimerDurations() map[int]int {
	return timerDuration
}

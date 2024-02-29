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

var TimerDuration = map[int]int{
	FOCUS_TIMER:       25 * 60 * TICKS_PER_SECOND,
	SHORT_BREAK_TIMER: 5 * 60 * TICKS_PER_SECOND,
	LONG_BREAK_TIMER:  15 * 60 * TICKS_PER_SECOND,
}

type TimerState interface {
	Tick() TimerState
	NextTimer() TimerState
	Pause() TimerState
	SkipCurrentTimer() TimerState
	ResetCurrentTimer() TimerState
	GetFormattedTimeString() string
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
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t RunningTimerState) GetFormattedTimeString() string {
	return fmt.Sprintf("#%d %2d:%02d", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
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
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t PausedTimerState) GetFormattedTimeString() string {
	return fmt.Sprintf("#%d %2d:%02d    (paused)", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
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
			TimeRemaining: TimerDuration[t.TimerType],
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

func (t Timer) getNextTimerType() (int, int, int) {
	if t.TimerType == FOCUS_TIMER && t.PomodoroCount%4 == 0 {
		return LONG_BREAK_TIMER, TimerDuration[LONG_BREAK_TIMER], t.PomodoroCount
	}

	if t.TimerType == FOCUS_TIMER {
		return SHORT_BREAK_TIMER, TimerDuration[SHORT_BREAK_TIMER], t.PomodoroCount
	}

	return FOCUS_TIMER, TimerDuration[FOCUS_TIMER], t.PomodoroCount + 1
}

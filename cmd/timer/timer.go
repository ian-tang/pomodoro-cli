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
	Tick(*Timer, *TimerState)
	NextTimer(*Timer, *TimerState)
	Pause(*Timer, *TimerState)
	SkipCurrentTimer(*Timer, *TimerState)
	ResetCurrentTimer(*Timer, *TimerState)
	GetFormattedTimeString(*Timer) string
}

type Timer struct {
	TimerType     int
	TimeRemaining int
	PomodoroCount int
}

type RunningTimerState struct{}

type PausedTimerState struct{}

type BetweenTimerState struct{}

type TimerStatePool struct {
	Running RunningTimerState
	Paused  PausedTimerState
	Between BetweenTimerState
}

var TSPool = TimerStatePool{
	Running: RunningTimerState{},
	Paused:  PausedTimerState{},
	Between: BetweenTimerState{},
}

func (rts RunningTimerState) Tick(t *Timer, ts *TimerState) {
	t.TimeRemaining--

	if t.TimeRemaining <= 0 {
		rts.NextTimer(t, ts)
	}
}

func (rts RunningTimerState) NextTimer(t *Timer, ts *TimerState) {
	*ts = TSPool.Between
}

func (rts RunningTimerState) Pause(t *Timer, ts *TimerState) {
	*ts = TSPool.Paused
}

func (rts RunningTimerState) SkipCurrentTimer(t *Timer, ts *TimerState) {
	rts.NextTimer(t, ts)
}

func (rts RunningTimerState) ResetCurrentTimer(t *Timer, ts *TimerState) {
	*ts = TSPool.Paused
	resetTimerValues(t)
}

func (rts RunningTimerState) GetFormattedTimeString(t *Timer) string {
	return fmt.Sprintf("#%d %2d:%02d", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
}

func (pts PausedTimerState) Tick(t *Timer, ts *TimerState) {}

func (pts PausedTimerState) NextTimer(t *Timer, ts *TimerState) {
	*ts = TSPool.Between
}

func (pts PausedTimerState) Pause(t *Timer, ts *TimerState) {
	*ts = TSPool.Running
}

func (pts PausedTimerState) SkipCurrentTimer(t *Timer, ts *TimerState) {
	pts.NextTimer(t, ts)
}

func (pts PausedTimerState) ResetCurrentTimer(t *Timer, ts *TimerState) {
	resetTimerValues(t)
}

func (pts PausedTimerState) GetFormattedTimeString(t *Timer) string {
	return fmt.Sprintf("#%d %2d:%02d    (paused)", t.PomodoroCount, t.TimeRemaining/TICKS_PER_MINUTE, (t.TimeRemaining%TICKS_PER_MINUTE)/TICKS_PER_SECOND)
}

func (bts BetweenTimerState) Tick(t *Timer, ts *TimerState) {}

func (bts BetweenTimerState) NextTimer(t *Timer, ts *TimerState) {}

func (bts BetweenTimerState) Pause(t *Timer, ts *TimerState) {
	setNextTimerValues(t)
	*ts = TSPool.Running
}

func (bts BetweenTimerState) SkipCurrentTimer(t *Timer, ts *TimerState) {}

func (bts BetweenTimerState) ResetCurrentTimer(t *Timer, ts *TimerState) {
	*ts = TSPool.Paused
	resetTimerValues(t)
}

func (bts BetweenTimerState) GetFormattedTimeString(t *Timer) string {
	switch {
	case t.TimerType == FOCUS_TIMER && t.PomodoroCount%4 == 0:
		return fmt.Sprintf("%d pomodoros done! Start long break?", t.PomodoroCount)
	case t.TimerType == FOCUS_TIMER:
		if t.PomodoroCount == 1 {
			return "1 pomodoro done! Start short break?"
		}

		return fmt.Sprintf("%d pomodoros done! Start short break?", t.PomodoroCount)
	default:
		return "Break over! Start pomodoro?"
	}
}

func resetTimerValues(t *Timer) {
	t.TimeRemaining = timerDuration[t.TimerType]
}

func setNextTimerValues(t *Timer) {
	switch {
	case t.TimerType == FOCUS_TIMER && t.PomodoroCount%4 == 0:
		t.TimerType = LONG_BREAK_TIMER
		t.TimeRemaining = timerDuration[LONG_BREAK_TIMER]
	case t.TimerType == FOCUS_TIMER:
		t.TimerType = SHORT_BREAK_TIMER
		t.TimeRemaining = timerDuration[SHORT_BREAK_TIMER]
	default:
		t.TimerType = FOCUS_TIMER
		t.TimeRemaining = timerDuration[FOCUS_TIMER]
		t.PomodoroCount++
	}
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

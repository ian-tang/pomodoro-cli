package timer

const (
	FOCUS_TIMER = iota
	SHORT_BREAK_TIMER
	LONG_BREAK_TIMER
)

var TimerDuration = map[int]int{
	FOCUS_TIMER:       25 * 60,
	SHORT_BREAK_TIMER: 5 * 60,
	LONG_BREAK_TIMER:  15 * 60,
}

type TimerState interface {
	tick()
	nextTimer() TimerState
	pause() TimerState
	skipCurrentTimer() TimerState
	resetCurrentTimer() TimerState
}

type Timer struct {
	TimerType     int
	TimeRemaining int
	PomodoroCount int
}

type FocusTimerState struct {
	Timer
}

type ShortBreakTimerState struct {
	Timer
}

type LongBreakTimerState struct {
	Timer
}

type PausedTimerState struct {
	Timer
}

func (t FocusTimerState) tick() {
	t.TimeRemaining--
	if t.TimeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t FocusTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t FocusTimerState) pause() TimerState {
	return PausedTimerState(t)
}

func (t FocusTimerState) skipCurrentTimer() TimerState {
	return t.nextTimer()
}

func (t Timer) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t ShortBreakTimerState) tick() {
	t.TimeRemaining--
	if t.TimeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t ShortBreakTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount + 1,
		},
	}
}

func (t ShortBreakTimerState) pause() TimerState {
	return PausedTimerState(t)
}

func (t ShortBreakTimerState) skipCurrentTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount + 1,
		},
	}
}

func (t ShortBreakTimerState) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t LongBreakTimerState) tick() {
	t.TimeRemaining--
	if t.TimeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t LongBreakTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount + 1,
		},
	}
}

func (t LongBreakTimerState) pause() TimerState {
	return PausedTimerState(t)
}

func (t LongBreakTimerState) skipCurrentTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     timerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount + 1,
		},
	}
}

func (t LongBreakTimerState) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t PausedTimerState) tick() {}

func (t PausedTimerState) nextTimer() TimerState {
	TimerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			TimerType:     TimerType,
			TimeRemaining: timerDuration,
			PomodoroCount: t.PomodoroCount + 1,
		},
	}
}

func (t PausedTimerState) pause() TimerState {
	switch t.TimerType {
	case FOCUS_TIMER:
		return FocusTimerState(t)
	case SHORT_BREAK_TIMER:
		return ShortBreakTimerState(t)
	case LONG_BREAK_TIMER:
		return LongBreakTimerState(t)
	default:
		return FocusTimerState(t)
	}
}

func (t PausedTimerState) skipCurrentTimer() TimerState {
	return t.nextTimer()
}

func (t PausedTimerState) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			TimerType:     t.TimerType,
			TimeRemaining: TimerDuration[t.TimerType],
			PomodoroCount: t.PomodoroCount,
		},
	}
}

func (t Timer) getNextTimerType() (int, int) {
	if t.TimerType == FOCUS_TIMER && t.PomodoroCount%4 == 0 {
		return LONG_BREAK_TIMER, TimerDuration[LONG_BREAK_TIMER]
	}

	if t.TimerType == FOCUS_TIMER {
		return SHORT_BREAK_TIMER, TimerDuration[SHORT_BREAK_TIMER]
	}

	return FOCUS_TIMER, TimerDuration[FOCUS_TIMER]
}

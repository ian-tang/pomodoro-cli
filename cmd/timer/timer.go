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
	timerType     int
	timeRemaining int
	pomodoroCount int
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
	t.timeRemaining--
	if t.timeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t FocusTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount,
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
			timerType:     t.timerType,
			timeRemaining: TimerDuration[t.timerType],
			pomodoroCount: t.pomodoroCount,
		},
	}
}

func (t ShortBreakTimerState) tick() {
	t.timeRemaining--
	if t.timeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t ShortBreakTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount + 1,
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
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount + 1,
		},
	}
}

func (t ShortBreakTimerState) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			timerType:     t.timerType,
			timeRemaining: TimerDuration[t.timerType],
			pomodoroCount: t.pomodoroCount,
		},
	}
}

func (t LongBreakTimerState) tick() {
	t.timeRemaining--
	if t.timeRemaining <= 0 {
		t.nextTimer()
	}
}

func (t LongBreakTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount + 1,
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
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount + 1,
		},
	}
}

func (t LongBreakTimerState) resetCurrentTimer() TimerState {
	return PausedTimerState{
		Timer{
			timerType:     t.timerType,
			timeRemaining: TimerDuration[t.timerType],
			pomodoroCount: t.pomodoroCount,
		},
	}
}

func (t PausedTimerState) tick() {}

func (t PausedTimerState) nextTimer() TimerState {
	timerType, timerDuration := t.getNextTimerType()

	return PausedTimerState{
		Timer{
			timerType:     timerType,
			timeRemaining: timerDuration,
			pomodoroCount: t.pomodoroCount + 1,
		},
	}
}

func (t PausedTimerState) pause() TimerState {
	switch t.timerType {
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
			timerType:     t.timerType,
			timeRemaining: TimerDuration[t.timerType],
			pomodoroCount: t.pomodoroCount,
		},
	}
}

func (t Timer) getNextTimerType() (int, int) {
	if t.timerType == FOCUS_TIMER && t.pomodoroCount%4 == 0 {
		return LONG_BREAK_TIMER, TimerDuration[LONG_BREAK_TIMER]
	}

	if t.timerType == FOCUS_TIMER {
		return SHORT_BREAK_TIMER, TimerDuration[SHORT_BREAK_TIMER]
	}

	return FOCUS_TIMER, TimerDuration[FOCUS_TIMER]
}

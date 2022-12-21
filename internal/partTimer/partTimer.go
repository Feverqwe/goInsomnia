package partTimer

import (
	"time"
)

const PART_DURATION = 5 * time.Minute

type PartTimer struct {
	timer     *time.Timer
	expiresAt time.Time
	onTimer   func()
}

func (s *PartTimer) Stop() {
	s.timer.Stop()
}

func (s *PartTimer) nextPart() {
	if s.expiresAt.Unix() < time.Now().Unix() {
		s.onTimer()
	} else {
		d := time.Until(s.expiresAt)
		if d > PART_DURATION {
			d = PART_DURATION
		}
		s.timer = time.AfterFunc(d, s.nextPart)
	}
}

func AfterFunc(d time.Duration, f func()) *PartTimer {
	ct := time.Now()
	ct = ct.Add(d)

	timer := &PartTimer{
		expiresAt: ct,
		onTimer:   f,
	}

	timer.nextPart()

	return timer
}

package core

import "time"

type Clock interface {
	Now() time.Time
}

func NewClock() *RealClock {
	return &RealClock{}
}

type RealClock struct {
}

func (c *RealClock) Now() time.Time {
	return time.Now().UTC()
}

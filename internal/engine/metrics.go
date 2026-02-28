package engine

import "time"

type Metrics struct {
	WPM        float64
	Accuracy   float64
	Errors     int
	TotalTyped int
	Correct    int
	TimeTaken  time.Duration
	Completed  bool
	TimedOut   bool
	Cancelled  bool
}

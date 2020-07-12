package simpletime

import "time"

// Schedule describes a job's duty cycle.
type Schedule interface {
	// Next returns the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	// If it is empty, it will be deleted.
	Next(time.Time) time.Time
}

// Job is an interface for submitted run jobs.
type Job interface {
	Run()
}

// WrapFuncJob is a wrapper that turns a func() into a Job.
type WrapFuncJob func()

func (fn WrapFuncJob) Run() { fn() }

// WrapFuncSchedule is a wrapper that turns a func(time.Time) time.Time into a Schedule.
type WrapFuncSchedule func(time.Time) time.Time

func (fn WrapFuncSchedule) Next(t time.Time) time.Time { return fn(t) }

// WrapSchedule is a wrapper that turns a Job and Schedule into a Schedule.
type WrapSchedule struct {
	Schedule Schedule
	Job      Job
}

// Get Next time.
// If necessary, Self-packaging: func(){ go run() }
func (wrap WrapSchedule) Next(t time.Time) time.Time {
	wrap.Job.Run()
	return wrap.Schedule.Next(t)
}

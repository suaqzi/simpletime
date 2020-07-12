package simpletime

// Enough Simple timer, Can easily handle hundreds of projects,
// Millions of projects, unreasonable slice structure.

// Linux cron format required:
// go git https://github.com/robfig/cron, extract: parser.go, spec.go, constantDelay.go
// Easy to finish Schedule interface.

// Timer does not do any error handling.
// Run in the specified time zone.
type Timer struct {
	entries []*Entry
	reset   chan *Entry
	running bool
}

// New returns a new Timer job runner, in the UTC time zone.
// Schedule return time.Time Automatic conversion to UTC.
func NewTimer() *Timer { return &Timer{reset: make(chan *Entry, 1)} }

// Wrap
func (timer *Timer) WrapFuncJob(s Schedule, fn func()) *Entry {
	return timer.WrapSchedule(s, WrapFuncJob(fn))
}
func (timer *Timer) WrapSchedule(s Schedule, cmd Job) *Entry {
	return timer.Schedule(WrapSchedule{Schedule: s, Job: cmd})
}

// Schedule adds a Job to the Timer to be run on the given schedule.
func (timer *Timer) Schedule(s Schedule) *Entry {
	entry := &Entry{schedule: s}
	// Add directly, next to get the new time
	timer.entries = append(timer.entries, entry)
	go timer.next(entry, timer.now())
	return entry
}

// Run scheduler.
func (timer *Timer) Run() {
	if !timer.running {
		timer.running = true
		timer.run()
	}
}

// If the scheduler is running, it stops; otherwise it does nothing.
func (timer *Timer) Stop() {
	if timer.running {
		timer.running = false
		timer.reset <- &Entry{} // Notify jump out run.
	}
}

// Reset all tasks
func (timer *Timer) Reset() {
	timer.Stop()
	for _, entry := range timer.entries {
		entry.Delete() // Notify on delete, If any.
	}
	timer.entries = nil
}

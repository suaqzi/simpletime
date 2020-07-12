package simpletime

import "time"

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	// Schedule on which this job should be run.
	schedule Schedule

	// The next time the job will run. This is the zero time if Tasks has not been
	next time.Time

	// Notify on delete
	notify Job
}

// Schedule nil or not.
func (entry *Entry) IsNil() bool {
	return entry.schedule == deleteSchedule
}

// Notify on delete, Follow up if necessary.
func (entry *Entry) SetNotifyFunc(fn func()) {
	entry.SetNotifyJob(WrapFuncJob(fn))
}
func (entry *Entry) SetNotifyJob(notify Job) {
	entry.notify = notify
}

// Delete self. Can't set time to 0, Probably in the middle of the first tier.
func (entry *Entry) Delete() {
	if entry.IsNil() {
		return
	}
	entry.schedule = deleteSchedule
	if entry.notify != nil {
		entry.notify.Run()
		entry.notify = nil // Release handle, GC.
	}
}

// Elegant delete task structure, size: 0. Not afraid of concurrency.
type deleteTask struct{}

func (d *deleteTask) Next(_ time.Time) time.Time {
	return time.Time{}
}

// Initialize null structure
var deleteSchedule = &deleteTask{}

# simpletime
  
// Schedule describes a job's duty cycle.  
		type Schedule interface {  
			// Next returns the next activation time, later than the given time.  
			// Next is invoked initially, and then each time the job is run.  
			// If it is empty, it will be deleted.  
			Next(time.Time) time.Time  
		}  
  
// WrapFuncSchedule is a wrapper that turns a func(time.Time) time.Time into a Schedule.  
type WrapFuncSchedule func(time.Time) time.Time  
  
timer := simpletime.NewTimer()  
entry := timer.Schedule(s Schedule) *Entry  
  
- Entry function:  
entry.IsNil() bool  
entry.SetNotifyFunc(fn func())  
entry.Delete()  
  
utils.go: various wrap.  


# simpletime
  
// Schedule describes a job's duty cycle.  
type Schedule interface {  
&emsp;&emsp;&emsp;&emsp;// Next returns the next activation time, later than the given time.  
&emsp;&emsp;&emsp;&emsp;// Next is invoked initially, and then each time the job is run.  
&emsp;&emsp;&emsp;&emsp;// If it is empty, it will be deleted.  
&emsp;&emsp;&emsp;&emsp;Next(time.Time) time.Time  
}  
  
// WrapFuncSchedule is a wrapper that turns a func(time.Time) time.Time into a Schedule.  
type WrapFuncSchedule func(time.Time) time.Time  
func (fn WrapFuncSchedule) Next(t time.Time) time.Time { return fn(t) }  
Various wrap: utils.go  
  
timer := simpletime.NewTimer()  
entry := timer.Schedule(s Schedule) *Entry  
  
- Entry function:  
entry.IsNil() bool  
entry.SetNotifyFunc(fn func())  
entry.Delete()  

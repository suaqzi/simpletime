package simpletime

import "time"

// tmp Duration
func (timer *Timer) tmpDuration() time.Duration {
	// If there are no entries yet,
	// just sleep - it still handles new entries and stop requests.
	if timer.entries == nil || timer.entries[0].next.IsZero() {
		return 1000 * time.Hour
	}
	// return time difference, New time is accurate.
	return timer.entries[0].next.Sub(timer.now())
}

// run the scheduler. this is private just due to the need to synchronize
// access to the 'running' state variable.
func (timer *Timer) run() {
	var (
		isSort   = true                     // Only meaningful for big data
		newTimer = time.NewTimer(time.Hour) // init, Restart sort time.
		nilTime  = time.Time{}
	)

	for timer.running {
		if isSort {
			// Sort priorities and remove obsolete.
			timer.minSortAndRemove()

			// Reset timer time
			newTimer.Reset(timer.tmpDuration())
		}

		select {
		case now := <-newTimer.C:
			// Run every Entry whose next time was less than now.
			for _, entry := range timer.entries {
				if entry.next.After(now) || entry.next.IsZero() {
					break
				}
				entry.next = nilTime
				go timer.next(entry, now)
			}
			isSort = true
		case entry := <-timer.reset:
			// Trigger the next round. Large amounts of data can reduce sorting
			if timer.entries != nil {
				isSort = timer.entries[0].next.IsZero() || !timer.entries[0].next.Before(entry.next)
			}
		} // continue
	}
	// end.
	newTimer.Stop()

	for {
		select {
		case <-newTimer.C:
		case <-timer.reset: // try to drain from the channel.
		default:
			return
		}
	}
}

// job is an interface for submitted timer jobs.
// afferent time - zone
func (timer *Timer) next(entry *Entry, now time.Time) {
	// defer entry.Recv()
	// Return time is empty, delete.
	entry.next = entry.schedule.Next(now).UTC()
	if entry.next.IsZero() {
		entry.Delete()
	} else if timer.running {
		timer.reset <- entry
	}
}

// New time based on region.
func (timer *Timer) now() time.Time {
	return time.Now().UTC()
}

// 1.When n is small, fast sorting is slow, and recursive sorting is fast.
// 2.When n is large and the degree of order is high, fast sorting is the fastest.
// 3.When n is large and the ordered program is low, the heap sorting is the fastest.
func (timer *Timer) minSortAndRemove() {
	// Reverse order, ensure no disorder.
	for i := len(timer.entries) - 1; i > -1; i-- {
		// is to delete?
		if timer.entries[i].IsNil() {
			// safe to delete, Insertion speed (slow->fast): slice -> list
			// append(, [i+1:]...), Don't worry about adding at any time
			timer.entries = append(timer.entries[:i], timer.entries[i+1:]...)
		} else {
			if timer.entries[i].next.IsZero() {
				continue
			}
			// Begin to exclude. Only one swap
			for j := 0; j < i; j++ {
				if timer.entries[j].next.Equal(timer.entries[i].next) {
					continue
				}
				// Time is different.
				// id: j min, i max; time i min j.
				if timer.entries[j].next.IsZero() || timer.entries[j].next.After(timer.entries[i].next) {
					timer.entries[i], timer.entries[j] = timer.entries[j], timer.entries[i]
				}
				break
			}
		}
	}
}

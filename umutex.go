package umutex

import (
	"sync"
	"sync/atomic"
)

// UMutex is simple implementation of Upgradable RWMutex.
type UMutex struct {
	rwmu sync.RWMutex
	u    int32
}

// RLock locks shared for multi reader.
func (m *UMutex) RLock() {
	m.rwmu.RLock()
}

// RUnlock unlocks reader lock.
func (m *UMutex) RUnlock() {
	m.rwmu.RUnlock()
}

// Lock locks exclusively for single writer.
func (m *UMutex) Lock() {
lock:
	m.rwmu.Lock()
	if atomic.LoadInt32(&m.u) > 0 {
		// Upgrade is given priority to Lock, retry lock.
		m.rwmu.Unlock()
		goto lock
	}
}

// Unlock unlocks writer lock.
func (m *UMutex) Unlock() {
	m.rwmu.Unlock()
}

// Upgrade converts reader lock to writer lock and returns success (true) or dead-lock (false).
// If Upgrade by multi reader locker at same time then dead-lock.
// Upgrade is given priority to Lock.
func (m *UMutex) Upgrade() bool {
	success := atomic.AddInt32(&m.u, 1) == 1
	if success {
		m.rwmu.RUnlock()
		m.rwmu.Lock()
	}
	atomic.AddInt32(&m.u, -1)
	return success
}

// Downgrade converts writer lock to reader lock.
// Downgrade() is given priority to Lock().
func (m *UMutex) Downgrade() {
	// no one does not have Lock because this one have the Lock
	atomic.StoreInt32(&m.u, 1)
	m.rwmu.Unlock()
	m.rwmu.RLock()
	atomic.AddInt32(&m.u, -1)
}

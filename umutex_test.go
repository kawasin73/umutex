package umutex

import (
	"testing"
	"time"
)

const waitTimeout = 100

func expectNotDone(t *testing.T, chDone chan struct{}, msg string) {
	t.Helper()
	select {
	case <-chDone:
		t.Error(msg)
	case <-time.After(waitTimeout * time.Millisecond):
		// not done
	}
}

func expectDone(t *testing.T, chDone chan struct{}, msg string) {
	t.Helper()
	select {
	case <-chDone:
		// done
	case <-time.After(waitTimeout * time.Millisecond):
		t.Error(msg)
	}
}

func TestUMutex_Upgrade(t *testing.T) {
	var mu UMutex
	chDone := make(chan struct{})

	mu.RLock()
	mu.RLock()
	go func() {
		mu.Upgrade()
		chDone <- struct{}{}
	}()
	expectNotDone(t, chDone, "RLock prevents Upgrade")

	// double Upgrade dead-locks.
	if mu.Upgrade() {
		t.Error("Upgrade dead-lock")
	}

	mu.RUnlock()
	expectDone(t, chDone, "RUnlock enables Upgrade")

	go func() {
		mu.RLock()
		chDone <- struct{}{}
	}()
	expectNotDone(t, chDone, "Upgraded mutex prevents RLock")

	mu.Unlock()
	expectDone(t, chDone, "Unlock enables RLock")

	// Upgrade is given priority to Lock.
	go func() {
		mu.Lock()
		chDone <- struct{}{}
	}()
	expectNotDone(t, chDone, "RLock prevents Lock")

	if !mu.Upgrade() {
		t.Error("failed to Upgrade")
	}

	expectNotDone(t, chDone, "Upgrade is given priority to Lock")

	mu.Unlock()

	expectDone(t, chDone, "Unlock enables Lock")
}

func TestUMutex_Downgrade(t *testing.T) {
	var mu UMutex
	chDone1 := make(chan struct{})

	mu.Lock()

	go func() {
		mu.Lock()
		chDone1 <- struct{}{}
	}()
	expectNotDone(t, chDone1, "Lock prevents Lock")

	mu.Downgrade()

	expectNotDone(t, chDone1, "Downgrade is given pritoriy to Lock")

	mu.RUnlock()

	expectDone(t, chDone1, "RUnlock enables Lock")
}

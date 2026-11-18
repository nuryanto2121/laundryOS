package bot

import (
	"testing"
	"time"
)

func TestMemorySessionStore(t *testing.T) {
	ttl := 100 * time.Millisecond
	store := NewMemorySessionStore(ttl)

	phone := "628123456789"

	// 1. Initial get should return false
	if _, ok := store.Get(phone); ok {
		t.Fatalf("expected empty store to return false")
	}

	// 2. Set session data
	data := &SessionData{
		State:        StateAwaitingAddress,
		CustomerName: "Budi",
	}
	store.Set(phone, data)

	// 3. Get session data immediately
	retrieved, ok := store.Get(phone)
	if !ok {
		t.Fatalf("expected session to be found")
	}
	if retrieved.State != StateAwaitingAddress || retrieved.CustomerName != "Budi" {
		t.Errorf("retrieved unexpected data: %+v", retrieved)
	}

	// 4. Update session data
	retrieved.State = StateAwaitingService
	retrieved.Address = "Jl. Merdeka No. 10"
	store.Set(phone, retrieved)

	updated, ok := store.Get(phone)
	if !ok || updated.State != StateAwaitingService || updated.Address != "Jl. Merdeka No. 10" {
		t.Errorf("failed to update session data: %+v", updated)
	}

	// 5. Expiration after TTL
	time.Sleep(150 * time.Millisecond)
	if _, ok := store.Get(phone); ok {
		t.Errorf("expected session to expire after TTL")
	}

	// 6. Delete explicitly
	store.Set(phone, &SessionData{State: StateIdle})
	store.Delete(phone)
	if _, ok := store.Get(phone); ok {
		t.Errorf("expected deleted session to return false")
	}
}

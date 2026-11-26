package main

import (
	"testing"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
)

func device(user string, id uint16) *store.Device {
	jid := types.JID{User: user, Device: id, Server: types.DefaultUserServer}
	return &store.Device{ID: &jid}
}

// After a re-link the store keeps the revoked companion row; reconnect must take the newest one.
func TestNewestDeviceByAccountKeyPrefersLatestLink(t *testing.T) {
	devices := []*store.Device{device("79990001122", 7), device("79990001122", 8), device("79990003344", 9)}
	got, matches := newestDeviceByAccountKey(devices, "79990001122")
	if got == nil || got.ID.Device != 8 {
		t.Fatalf("want device 8, got %v", got)
	}
	if matches != 2 {
		t.Fatalf("want 2 devices of the account, got %d", matches)
	}
}

func TestNewestDeviceByAccountKeyOrderIndependent(t *testing.T) {
	devices := []*store.Device{device("79990001122", 12), device("79990001122", 3)}
	if got, _ := newestDeviceByAccountKey(devices, "79990001122"); got == nil || got.ID.Device != 12 {
		t.Fatalf("want device 12, got %v", got)
	}
}

func TestNewestDeviceByAccountKeyNoMatch(t *testing.T) {
	devices := []*store.Device{device("79990003344", 2), {ID: nil}}
	if got, matches := newestDeviceByAccountKey(devices, "79990001122"); got != nil || matches != 0 {
		t.Fatalf("want no device, got %v (%d)", got, matches)
	}
}

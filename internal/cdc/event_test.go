package cdc

import "testing"

func TestOutboxRowToCloudEvent(t *testing.T) {
	row, err := FromBinlogRow([]any{uint64(7), []byte("event-7"), []byte("inventory.stock-moved.v1"), []byte("DEMO-001"), []byte(`{"delta":-1}`), "2026-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	event := row.CloudEvent()
	if event.ID() != "event-7" || event.Subject() != "DEMO-001" || event.Type() != "inventory.stock-moved.v1" {
		t.Fatalf("unexpected CloudEvent: %s", event.String())
	}
}

func TestMalformedOutboxRow(t *testing.T) {
	_, err := FromBinlogRow([]any{uint64(7), "id", "type", "sku", "not-json", "time"})
	if err == nil {
		t.Fatal("wanted invalid JSON error")
	}
}

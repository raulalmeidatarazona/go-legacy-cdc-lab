package main

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
)

// Run with INTEGRATION=1 after make up. It verifies the core atomicity claim.
func TestInvalidMovementRollsBackOutbox(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("set INTEGRATION=1 and run make up")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("mysql", "root:localdemo@tcp(127.0.0.1:3307)/legacy")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var beforeMovements, beforeOutbox, beforeBalance int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_movements").Scan(&beforeMovements); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events").Scan(&beforeOutbox); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT quantity FROM inventory_items WHERE sku='DEMO-001'").Scan(&beforeBalance); err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, `INSERT INTO stock_movements (sku, delta, location_code, actor_id)
		VALUES ('DEMO-001', -100000, 'A-01-02', 'test-operator')`)
	if err == nil {
		t.Fatal("expected nonnegative balance constraint to reject movement")
	}
	var afterMovements, afterOutbox, afterBalance int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_movements").Scan(&afterMovements); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events").Scan(&afterOutbox); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT quantity FROM inventory_items WHERE sku='DEMO-001'").Scan(&afterBalance); err != nil {
		t.Fatal(err)
	}
	if beforeMovements != afterMovements || beforeOutbox != afterOutbox || beforeBalance != afterBalance {
		t.Fatalf("failed movement changed state: movements %d→%d, outbox %d→%d, balance %d→%d",
			beforeMovements, afterMovements, beforeOutbox, afterOutbox, beforeBalance, afterBalance)
	}
}

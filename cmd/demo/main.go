package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/raulalmeidatarazona/go-legacy-cdc-lab/internal/cdc"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", "root:localdemo@tcp(127.0.0.1:3307)/legacy?parseTime=true")
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("MySQL unavailable; run make up: %w", err)
	}

	var file string
	var position uint32
	var ignored1, ignored2 any
	if err := db.QueryRowContext(ctx, "SHOW BINARY LOG STATUS").Scan(&file, &position, &ignored1, &ignored2, new(any)); err != nil {
		return fmt.Errorf("read binlog position: %w", err)
	}

	syncer := replication.NewBinlogSyncer(replication.BinlogSyncerConfig{
		ServerID: 101, Flavor: "mysql", Host: "127.0.0.1", Port: 3307,
		User: "cdc", Password: "localdemo",
	})
	defer syncer.Close()
	stream, err := syncer.StartSync(mysql.Position{Name: file, Pos: position})
	if err != nil {
		return fmt.Errorf("start binlog stream: %w", err)
	}

	// The movement, balance update, and outbox insert are one MySQL transaction.
	if _, err := db.ExecContext(ctx, `INSERT INTO stock_movements (sku, delta, location_code, actor_id)
		VALUES ('DEMO-001', -1, 'A-01-02', 'demo-operator')`); err != nil {
		return fmt.Errorf("write synthetic stock movement: %w", err)
	}

	for {
		entry, err := stream.GetEvent(ctx)
		if err != nil {
			return fmt.Errorf("wait for outbox binlog event: %w", err)
		}
		rows, ok := entry.Event.(*replication.RowsEvent)
		if !ok || string(rows.Table.Schema) != "legacy" || string(rows.Table.Table) != "outbox_events" {
			continue
		}
		if entry.Header.EventType != replication.WRITE_ROWS_EVENTv2 {
			continue
		}
		for _, values := range rows.Rows {
			row, err := cdc.FromBinlogRow(values)
			if err != nil {
				return err
			}
			encoded, err := row.CloudEvent().MarshalJSON()
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(encoded))
			return nil
		}
	}
}

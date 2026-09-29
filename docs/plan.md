# Implementation plan

## Existing vertical slice

MySQL 8.4 container, synthetic schema, trigger, Go binlog reader, CloudEvents conversion, runnable CLI and unit tests. `make demo` starts from the current binlog file position, inserts one movement and prints the generated event.

## Next vertical slice: reliable delivery

1. Introduce RabbitMQ to Compose and a Go publisher using persistent messages and publisher confirms.
2. Read from a durable GTID checkpoint on startup; checkpoint only after publish confirmation.
3. Add a fault injection mode that kills the process between confirmation and checkpoint.
4. Add an idempotent consumer with an `applied_events` uniqueness constraint and projection update in one transaction.
5. Test replay after crash and broker outage, measuring delivery lag and duplicate count.

## Final review slice

Add bounded queue/backpressure, OpenTelemetry metrics/traces, runbook, CI, dependency scanning, load fixture, benchmark methodology and a 90-second demo video. Report measured results with hardware, dataset and commands. Keep a separate section for unresolved risks.

## Key implementation rule

The checkpoint represents **confirmed delivery to the broker**, not merely a row read from MySQL. A crash after publish and before checkpoint can duplicate an event; the consumer's event ID constraint makes that safe.

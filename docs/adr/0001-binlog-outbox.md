# ADR 0001: Read a transactional outbox from the row binlog

**Status:** accepted for the lab; production adoption requires source-system review.

## Context

Downstream services need inventory changes from a legacy MySQL writer. Frequent polling adds repeated queries even when nothing changes. Direct dual writes to MySQL and a broker introduce a partial-failure window.

## Decision

Write a versioned outbox row in the same transaction as the stock movement, then read committed outbox inserts from MySQL's row binlog in Go. Standardize each row as a CloudEvent. In the next slice, publish with broker confirmation before persisting a GTID checkpoint. Consumers will deduplicate by event ID within their projection transaction.

## Consequences

- The source transaction pays for an extra insert and trigger execution; benchmark that cost.
- Binlog retention and replication privilege become operational dependencies.
- At-least-once delivery creates duplicates across failures; downstream idempotency is required.
- Schema evolution needs versioned event types and compatibility tests.
- A trigger keeps the demo's legacy writer unchanged, but source teams must approve it before production use.

## Alternatives

- **Polling the outbox:** simpler to operate at small scale, but repeatedly queries the database and introduces poll-interval latency.
- **Direct broker publish from writer:** avoids CDC infrastructure but requires changing every writer and leaves a dual-write failure window unless coordinated.
- **Managed CDC connector:** strongest operational alternative for a production team; compare staffing, control, and deployment constraints before building a custom streamer.

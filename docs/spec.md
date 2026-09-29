# Feature specification: inventory change stream

## Problem

A legacy inventory writer cannot be rewritten quickly. Downstream systems need committed stock changes without frequent source-table polling. They also need to survive disconnects, duplicates, and restarts without silently losing changes.

## Actors

- **Legacy writer:** inserts a stock movement.
- **CDC process:** reads committed outbox rows from the MySQL binlog.
- **Projection service:** eventually consumes events into a queryable read model.
- **Operator:** detects lag, repairs gaps, and replays safely.

## Functional requirements

- **FR-1:** A successful movement changes inventory balance and inserts exactly one outbox row in the same transaction.
- **FR-2:** An invalid movement rolls back the movement, balance and outbox row.
- **FR-3:** The CDC process emits a CloudEvent with stable ID, source, type, subject, versioned payload, and origin sequence.
- **FR-4:** After restart, the process resumes from a durable GTID checkpoint and may replay events without loss.
- **FR-5:** Broker publication is acknowledged before the checkpoint advances. A failed publish leaves the event replayable.
- **FR-6:** A duplicate event cannot apply a balance delta twice in the read model.
- **FR-7:** Operators can see lag, last checkpoint, retry count, and terminal errors.

## Quality and safety requirements

- Bound memory under a slow broker using a finite queue or synchronous flow control.
- Reject malformed records with a visible, recoverable policy; never silently skip them.
- Keep the MySQL source read-only for the CDC user except replication privileges; production credentials must be externalized.
- Produce a reproducible load test before claiming a latency or throughput target.
- Define behavior if a checkpoint points to a purged binlog or if GTID history diverges.

## Acceptance scenarios

1. Insert `DEMO-001` movement `-1`: source balance decreases by one; one outbox row is created; a CloudEvent with the matching SKU and delta appears.
2. Attempt a movement beyond available stock: the transaction fails; all three tables remain unchanged.
3. Restart after broker acknowledgement and before checkpoint: event is replayed; read model contains one logical application.
4. Stop the broker: CDC stops advancing its checkpoint, backlog becomes visible, and source writes still commit.
5. Feed an event twice: read model balance and event-applied count remain unchanged on the second delivery.

Scenarios 1 and 2 are implemented. Scenarios 3–5 are the next development gates.

## Out of scope for the first version

FEFO/FIFO allocation, scanning UX, operator accountability and cycle counts belong to a full WMS and are not implied by this synthetic CDC lab.

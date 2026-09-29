# Tasks and acceptance gates

## Done

- [x] Define source transaction, event contract and failure scenarios.
- [x] Add Docker MySQL with row binlog and GTID mode.
- [x] Write synthetic movement, balance and outbox schema.
- [x] Stream a committed outbox insert into a CloudEvent.
- [x] Add unit tests for mapping and malformed payload.
- [x] Add integration test for rollback on insufficient stock.
- [x] Configure CI for unit, race, integration, and demo checks (remote run pending).

## Next, in order

- [ ] Store and restore GTID checkpoint; test restart from a prior GTID.
- [ ] Publish CloudEvents to RabbitMQ with publisher confirms.
- [ ] Add an idempotent projection and API; duplicate delivery test.
- [ ] Add broker outage and crash-window integration tests.
- [ ] Add bounded backpressure and failure/runbook documentation.
- [ ] Record benchmark methodology and measured p50/p95 latency and source DB cost.
- [ ] Choose license and add threat model, architecture diagram and screencast.

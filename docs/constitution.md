# Engineering constitution

1. **Show evidence.** Every claimed reliability or speed property needs a repeatable test or measurement. No invented throughput or uptime numbers.
2. **Preserve source atomicity.** The movement, balance, and outbox row share one transaction. A failed movement must create no event.
3. **Make failure observable.** Every retry, gap, poison event, and checkpoint decision has a documented operator action.
4. **Expect duplicates.** Delivery is at least once; consumers use event ID to make projection updates idempotent.
5. **Keep client IP out.** All data, schema, names and code in this repo are synthetic and independently written.
6. **Optimize for review.** One command to run, one diagram to understand, tests for the hard cases, and explicit scope boundaries.

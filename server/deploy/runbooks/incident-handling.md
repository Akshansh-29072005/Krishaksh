# Incident Handling Notes

Severity levels:
- Sev1: payment flow outage, auth outage, DB unavailable.
- Sev2: scan/AI degraded with fallback still working.
- Sev3: non-critical dashboard/notification issues.

Flow:
1. Declare incident in ops channel.
2. Assign incident commander.
3. Stabilize user impact (rate limit/feature flag/degrade mode).
4. Execute rollback if needed.
5. Publish post-incident report with timeline and action items.

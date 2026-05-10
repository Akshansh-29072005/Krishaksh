# Rollback Checklist

- [ ] Identify failing component (`api`, `worker`, `nginx`, `redis`)
- [ ] Stop stack using `scripts/staging/rollback-staging.sh`
- [ ] Restore previous image tag/config
- [ ] If migration caused failure, apply corresponding `*.down.sql`
- [ ] Re-run health checks and smoke tests

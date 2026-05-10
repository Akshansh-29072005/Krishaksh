# Operational SOPs

## Queue backlog response
1. Check `/metrics` for queue enqueue patterns.
2. Check worker logs for `worker_errors_total:*` patterns.
3. Scale `worker` replicas or increase Asynq concurrency.

## AI degraded mode response
1. Check `ai_provider_fallback_total` and `ai_inference_failure_total`.
2. Validate external AI provider status.
3. Switch provider ordering if needed.

## Payment incident response
1. Verify webhook signature validation still active.
2. Check `transactions` and `webhook_events` tables for stuck events.
3. Requeue failed payment events if safe.

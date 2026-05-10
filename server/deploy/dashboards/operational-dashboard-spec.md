# Operational Dashboard Spec

Panels:
- API latency p95 (derived from `http_request_duration_ms_total` deltas)
- HTTP error rates (`errors_total` and status counters)
- Queue enqueue rates (`queue_enqueue_total_*`)
- Worker error rates (`worker_errors_total_*`)
- AI success/fallback/failure (`ai_inference_success_total_*`, `ai_provider_fallback_total`, `ai_inference_failure_total`)
- Payment event throughput (worker payment logs + transaction table)
- Ready probe status (`/ready` uptime)

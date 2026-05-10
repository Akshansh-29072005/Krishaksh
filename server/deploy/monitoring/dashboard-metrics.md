# Grafana-ready metrics reference

Use `/metrics` from API for dashboards:
- `http_requests_total_*`
- `http_request_duration_ms_total`
- `errors_total`
- `queue_enqueue_total_*`
- `ai_inference_success_total_*`
- `ai_provider_fallback_total`
- `ai_inference_failure_total`
- `worker_errors_total_*`

Payment and queue health:
- `queue_enqueue_total_payment`
- `worker_payment_event_processed` logs (structured) for derived counters.

AI latency dashboard:
- parse structured logs field `latency_ms` from `ai_inference_success`.

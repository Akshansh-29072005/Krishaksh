# Krisho Production Deployment (AWS)

## Architecture
- EC2 instance(s): run Docker Compose stack (`api`, `worker`, `redis`, `nginx`, `prometheus`).
- RDS PostgreSQL: managed DB, private subnet, `sslmode=require`.
- S3: scan image storage + backups.
- IAM Role on EC2/ECS: used for AWS API access (no static credentials).

## Containers
- `api`: serves Gin APIs, health, metrics.
- `worker`: Asynq processors (scan/payment/analytics).
- `redis`: queue backend and cache.
- `nginx`: TLS termination + reverse proxy + headers + limits.
- `prometheus`: scrape metrics.

## AWS setup checklist
1. RDS PostgreSQL
- Enable automated backups (7-35 days).
- Enable Performance Insights.
- Put in private subnet.
- Security group: only app nodes can connect 5432.

2. S3 bucket
- Enable versioning.
- Enable default encryption (SSE-S3 or SSE-KMS).
- Block public access.

3. IAM role strategy
- Attach least-privilege role to EC2 instance profile:
  - `s3:GetObject`, `s3:PutObject`, `s3:ListBucket` on required bucket/prefix.
  - Optional `kms:Decrypt` if KMS encryption.
- Do not set `AWS_ACCESS_KEY` / `AWS_SECRET_KEY` in env.

4. TLS certificates
- Place cert files in `deploy/nginx/certs/`:
  - `fullchain.pem`
  - `privkey.pem`
- Prefer AWS ACM + ALB for managed TLS in scaled setup.

## Environment and secrets
- Use `.env.prod` loaded by compose.
- Validate before deploy using `scripts/validate-env.sh`.
- Firebase service account mounted as Docker secret:
  - `deploy/secrets/firebase_service_account.json`

## Deploy
```bash
cd server
./scripts/deploy-prod.sh .env.prod
```

## Health validation
- `GET /health`
- `GET /ready`
- `GET /metrics`

## Backup strategy
- Run DB backups via `scripts/backup/rds-backup.sh` (cron or GitHub Action).
- Keep S3 versioning enabled for recovery of overwritten objects.
- Keep migration rollback SQL files (`*.down.sql`) under `db/migrations`.

## Disaster recovery notes
- RDS point-in-time restore enabled.
- Backup SQL dumps copied to dedicated backup bucket/prefix.
- Document recovery RTO/RPO targets (recommended: RPO <= 15m with PITR, RTO <= 60m).

## Monitoring and dashboards
- Prometheus scrapes `api:8080/metrics`.
- Build Grafana dashboards from metrics in `deploy/monitoring/dashboard-metrics.md`.
- Queue and AI behavior also available via structured logs.

## Security hardening
- Nginx security headers enabled.
- Request size cap in Nginx + app middleware.
- Rate limiting at Nginx and app levels.
- Secrets isolated from image layers.

## FCM production notes
- Set `FIREBASE_CREDENTIALS_PATH` to mounted secret JSON.
- Use separate Firebase project for production.
- Rotate service account keys periodically; prefer Workload Identity where possible.

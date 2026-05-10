# Krisho Production Deployment

## Architecture
- EC2/VPS instance(s): run Docker Compose stack (`api`, `worker`, `redis`, `nginx`, `prometheus`).
- Supabase PostgreSQL: Managed PostgreSQL (Session Pooler), `sslmode=require`.
- AWS S3: scan image storage.
- IAM Role on EC2/ECS: used for S3 access.

## Infrastructure Setup Checklist

### 1. Supabase PostgreSQL
- Create a Supabase project.
- Go to **Project Settings** -> **Database**.
- Use the **Transaction/Session Pooler** connection string (Port 5432).
- Ensure **sslmode=require** is appended to the connection string.
- Set `DATABASE_URL` in `.env.prod`.
- Configure `DB_MAX_CONNS` based on your Supabase tier (default: 50 for session pooler).

### 2. AWS S3 bucket
- Enable versioning.
- Enable default encryption (SSE-S3).
- Block public access.
- Attach least-privilege role to EC2 instance profile for S3 access.

### 3. TLS certificates
- Place cert files in `deploy/nginx/certs/`:
  - `fullchain.pem`
  - `privkey.pem`

## Supabase Pooling Notes
- The backend uses `pgxpool` with Supabase session pooling.
- Recommended Settings:
  - `DB_MAX_CONNS`: 50
  - `DB_MIN_CONNS`: 5
  - `DB_CONN_MAX_LIFETIME`: 1h
- Avoid "Statement" pooling; only "Transaction" or "Session" mode is supported for migration compatibility.

## Troubleshooting
- **Connection Timeout**: check if Supabase IP is restricted (Supabase allows all by default, but ensure no VPC blocks egress 5432).
- **SSL Error**: Ensure `?sslmode=require` is in the `DATABASE_URL`.
- **Pool Exhaustion**: Check transaction leaks or increase `DB_MAX_CONNS`.

## Deploy
```bash
cd server
cp .env.prod.example .env.prod # and configure
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

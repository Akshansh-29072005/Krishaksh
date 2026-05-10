# Staging Environment Checklist

- [ ] `.env.staging` created from `.env.staging.example`
- [ ] TLS certs in `deploy/nginx/certs`
- [ ] Firebase service account secret mounted
- [ ] RDS connectivity verified with `sslmode=require`
- [ ] S3 bucket exists and versioning enabled
- [ ] `scripts/staging/deploy-staging.sh` succeeded
- [ ] `/health`, `/ready`, `/metrics` green
- [ ] Queue worker process healthy

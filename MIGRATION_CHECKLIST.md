# Cloud Run + Upstash Migration Checklist

## 📋 Readiness Checklist

### Code Architecture Assessment
- [x] API is stateless (Gin handlers)
- [x] No persistent local filesystem needed
- [x] All file storage uses S3 (external)
- [x] Database is external (Supabase PostgreSQL)
- [x] No long-running daemon processes (except Asynq worker)
- [x] Query timeouts configured (30s)
- [x] All tasks complete in <10 seconds

### Current Async Jobs (All Migratable)
- [x] `TypeScanAnalyze` - Vision AI analysis (~10s)
- [x] `TypePaymentEvent` - Webhook processing (~5s)
- [x] `TypeRefundStart` - Refund trigger (~2s)
- [x] `TypeAnalyticsEvt` - Analytics logging (~1s)
- [x] `TypeOrderNotify` - Notifications (~3s)

### Configuration Points
- [x] Environment variables centralized
- [x] Redis connection string configurable
- [x] Database connection string from env
- [x] API port from env (8080)
- [x] All secrets in env vars (no hardcoded)

---

## 🚀 Migration Impact Summary

| Component | Current | Migration | Effort | Risk |
|-----------|---------|-----------|--------|------|
| **API** | Gin on VPS | Cloud Run | Minimal | 🟢 Low |
| **Database** | Supabase | Supabase | None | 🟢 None |
| **Job Queue** | Asynq + Redis | Cloud Tasks | Medium | 🟡 Medium |
| **Redis** | Self-managed | Upstash | Minimal | 🟢 Low |
| **Reverse Proxy** | Nginx | Cloud LB | Minimal | 🟢 Low |
| **Worker Daemon** | Asynq daemon | HTTP handlers | Medium | 🟡 Medium |

---

## 💰 Cost Comparison

### Fixed Costs (Current VPS)
```
VPS (2vCPU, 4GB RAM)    : $20-50/month
Redis (included)        : —
Supabase PostgreSQL     : $25-50/month
─────────────────────────────
Total                   : $65-120/month
```

### Variable Costs (Cloud Run + Upstash - typical)
```
Cloud Run (API)         : $40-60/month
Cloud Tasks             : $1-2/month
Upstash Redis           : $6-10/month
Supabase PostgreSQL     : $25-50/month
─────────────────────────────
Total                   : $72-122/month
```

**💡 Break-even traffic**: Normal usage. **Savings at <50% capacity**.

---

## 🏗️ Architecture Changes Required

### 1. Job Processing Migration

**BEFORE (Asynq on VPS)**:
```go
// In API handler
queueClient.Enqueue(asynq.NewTask(TypeScanAnalyze, payload))

// In worker daemon (long-running process)
mux.HandleFunc(TypeScanAnalyze, handler)
srv.Run(mux)  // Daemon runs forever
```

**AFTER (Cloud Tasks)**:
```go
// In API handler
tasksClient.CreateTask(ctx, &taskspb.CreateTaskRequest{...})

// In Cloud Run (HTTP endpoints)
router.POST("/internal/tasks/scan-analyze", handler)
// Cloud Tasks calls this endpoint, it processes and returns

// No daemon needed!
```

### 2. Redis Connection

**BEFORE**:
```go
redisClient := redis.NewClient(&redis.Options{
    Addr: "redis:6379",
})
```

**AFTER**:
```go
redisClient := redis.NewClient(&redis.Options{
    Addr: os.Getenv("UPSTASH_REDIS_URL"),
    Password: os.Getenv("UPSTASH_REDIS_TOKEN"),
})
```

### 3. Docker Composition

**BEFORE**:
```yaml
services:
  api:
    ...
  worker:        # ← DELETE THIS
    ...
  redis:         # ← DELETE THIS (use Upstash)
    ...
  nginx:         # ← DELETE THIS (use Cloud LB)
    ...
```

**AFTER**:
```yaml
# No docker-compose needed!
# Deploy directly to Cloud Run
# Configure via GCP console or Terraform
```

---

## 🔧 Implementation Steps

### Phase 1: Setup (1-2 days)
```bash
# 1. Create GCP Project
gcloud projects create krisho-prod

# 2. Enable APIs
gcloud services enable cloudrun.googleapis.com
gcloud services enable cloudtasks.googleapis.com
gcloud services enable artifactregistry.googleapis.com

# 3. Create Upstash Redis instance
# - Go to upstash.com → Create Redis database
# - Copy connection string

# 4. Create Cloud Tasks queue
gcloud tasks queues create scans
gcloud tasks queues create payments
gcloud tasks queues create refunds
```

### Phase 2: Code Changes (3-5 days)
```
src/
├── cmd/api/main.go
│   ├── Replace Asynq client → Cloud Tasks client
│   ├── Add task HTTP endpoints
│   └── Update Redis to Upstash
├── pkg/queue/
│   ├── Create cloud_tasks.go (Cloud Tasks client)
│   ├── Deprecate asynq.go
│   └── Update enqueue functions
├── internal/workers/
│   └── Move handlers to Cloud Run HTTP routes
├── cmd/worker/
│   └── DELETE (no longer needed)
└── docker-compose.prod.yml
    └── Simplify (only API service)
```

### Phase 3: Testing (3-5 days)
- Local Cloud Tasks emulator
- Integration tests
- Load tests
- Staging in Cloud Run

### Phase 4: Deployment (2-3 days)
- Deploy to Cloud Run
- Monitor
- Traffic shift
- Rollback if needed

---

## ✅ Pre-Requisites for Migration

Before starting, ensure:

- [ ] GCP Account with billing enabled
- [ ] `gcloud` CLI installed
- [ ] Docker installed (for Cloud Run image building)
- [ ] Terraform or IaC tool (optional, recommended)
- [ ] Staging environment (separate GCP project)
- [ ] 200-300 engineering hours allocated
- [ ] 4-week timeline approved
- [ ] Team trained on Cloud Run/Cloud Tasks (1 day)

---

## ⚠️ Migration Risks & Mitigations

| Risk | Severity | Mitigation |
|------|----------|-----------|
| Job loss during migration | High | Use Cloud Tasks DLQ, test retries |
| API downtime | High | Blue-green deployment, health checks |
| Performance regression | Medium | Load test before cutover |
| Cost overrun | Medium | Set GCP budget alerts |
| Vendor lock-in (GCP) | Low | Keep database external (Supabase) |
| Team unfamiliarity | Medium | Schedule training on Cloud Tasks |

---

## 🎯 Success Criteria

- [ ] All existing jobs process successfully
- [ ] API response time ≤ 200ms (p99)
- [ ] Error rate < 0.1%
- [ ] Auto-scaling works for 10x traffic
- [ ] Zero data loss during migration
- [ ] Rollback time < 30 minutes
- [ ] Cost within ±10% of estimate

---

## 📞 Discussion Topics

### 1. **Traffic Spike Handling**
- Current: VPS would crash, manual upgrade needed (hours)
- Proposed: Cloud Run scales in 30 seconds automatically
- **Decision**: Acceptable?

### 2. **Cost Predictability**
- Current: Fixed $65-120/month, predictable
- Proposed: Variable $50-150/month, depends on usage
- **Decision**: Preferred?

### 3. **Timeline & Resources**
- Requires: 200-300 engineering hours
- Timeline: 4-6 weeks with testing
- **Decision**: Feasible for your team?

### 4. **Deployment Strategy**
- Blue-green vs. canary deployment
- Rollback plan if issues arise
- **Decision**: Which approach preferred?

### 5. **Monitoring & Observability**
- Current: Prometheus on VPS
- Proposed: Google Cloud Monitoring + Cloud Trace
- **Decision**: Migration plan for existing dashboards?

---

## 📖 Resources

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud Tasks Documentation](https://cloud.google.com/tasks/docs)
- [Upstash Redis Documentation](https://upstash.com/docs)
- [Go Cloud Tasks Client](https://pkg.go.dev/cloud.google.com/go/cloudtasks/apiv2)
- [Go Upstash Client](https://github.com/upstash/redis-go)

---

## 📋 Sign-Off

**Ready for Architecture Discussion**: ✅ YES

**Next Steps**:
1. Review evaluation document
2. Discuss in team sync
3. Approve/modify migration plan
4. Allocate resources
5. Schedule implementation


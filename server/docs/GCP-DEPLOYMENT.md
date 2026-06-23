# GCP Cloud Run + Cloud Tasks Deployment Guide

## Architecture Overview

The Krisho backend has been refactored for a **fully serverless, stateless architecture** on Google Cloud Platform (GCP):

| Component | GCP Service | Purpose |
|-----------|-------------|---------|
| API | Cloud Run | Serves REST API, scales automatically from 0 to 1000+ containers |
| Task Queue | Cloud Tasks | Fully managed, serverless task queue (replaces self-hosted Redis/Asynq) |
| Async Worker | Cloud Run | Processes tasks from Cloud Tasks, scales automatically |
| Database | Supabase (PostgreSQL) | Managed relational database |
| Storage | AWS S3 / Google Cloud Storage | Stores crop scan images |
| AI | Gemini / OpenAI | Crop disease detection |

## Key Benefits

- **Cost Efficiency**: 60-80% lower cost than always-on VMs/containers
- **Auto-Scaling**: Scales from 0 to maximum capacity in seconds
- **No State**: All services are completely stateless (no local data/cache)
- **Managed Services**: Zero maintenance for queues, databases, etc.

---

## Step 1: GCP Project Setup

1. Create a new GCP project (or use an existing one)
2. Enable required APIs:
   - Cloud Run Admin API
   - Cloud Tasks API
   - Cloud Build API (for automated deployments)
3. Install gcloud CLI and authenticate:
   ```bash
   gcloud auth login
   gcloud config set project YOUR_PROJECT_ID
   ```

---

## Step 2: Cloud Tasks Queue Setup

Create a Cloud Tasks queue:
```bash
gcloud tasks queues create krisho-queue \
  --location=us-central1 \
  --max-dispatches-per-second=100 \
  --max-attempts=5
```

---

## Step 3: Deploy API to Cloud Run

### 3.1: Build and deploy API container

Our Dockerfile has a multi-stage build with `--target=api` and `--target=worker`.

Create a `cloudbuild-api.yaml` file in `/home/akshansh/AndroidStudioProjects/Krishaksh/server/` to build the API:
```yaml
steps:
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '--target', 'api', '-t', 'gcr.io/$PROJECT_ID/krisho-api', '.']
images:
  - 'gcr.io/$PROJECT_ID/krisho-api'
```

Then run:
```bash
gcloud builds submit --config cloudbuild-api.yaml
gcloud run deploy krisho-api \
  --image gcr.io/YOUR_PROJECT_ID/krisho-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --max-instances 1000 \
  --set-env-vars "QUEUE_TYPE=cloudtasks" \
  --set-env-vars "GCP_PROJECT_ID=YOUR_PROJECT_ID" \
  --set-env-vars "GCP_LOCATION=us-central1" \
  --set-env-vars "CLOUD_TASKS_QUEUE=krisho-queue" \
  --set-env-vars "CLOUD_TASKS_WORKER_URL=https://YOUR_WORKER_URL.a.run.app" \
  --set-env-vars "DATABASE_URL=YOUR_SUPABASE_DATABASE_URL" \
  --set-env-vars "GEMINI_API_KEY=YOUR_GEMINI_KEY" \
  --set-env-vars "AWS_ACCESS_KEY=YOUR_AWS_KEY" \
  --set-env-vars "AWS_SECRET_KEY=YOUR_AWS_SECRET" \
  --set-env-vars "S3_BUCKET=YOUR_S3_BUCKET" \
  --set-env-vars "RAZORPAY_KEY_ID=YOUR_RAZORPAY_KEY" \
  --set-env-vars "RAZORPAY_KEY_SECRET=YOUR_RAZORPAY_SECRET" \
  --set-env-vars "RAZORPAY_WEBHOOK_SECRET=YOUR_WEBHOOK_SECRET" \
  --set-env-vars "SERVER_ENV=production"
```

---

## Step 4: Deploy GCP Worker to Cloud Run

### 4.1: Build and deploy GCP worker container

Create a `cloudbuild-gcp-worker.yaml` file in `/home/akshansh/AndroidStudioProjects/Krishaksh/server/` to build the worker:
```yaml
steps:
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '--target', 'gcp-worker', '-t', 'gcr.io/$PROJECT_ID/krisho-gcp-worker', '.']
images:
  - 'gcr.io/$PROJECT_ID/krisho-gcp-worker'
```

First, **deploy the worker first** so you have a URL for `CLOUD_TASKS_WORKER_URL`!

```bash
# Deploy GCP Worker first
gcloud builds submit --config cloudbuild-gcp-worker.yaml
gcloud run deploy krisho-gcp-worker \
  --image gcr.io/YOUR_PROJECT_ID/krisho-gcp-worker \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --max-instances 500 \
  --set-env-vars "DATABASE_URL=YOUR_SUPABASE_DATABASE_URL" \
  --set-env-vars "GEMINI_API_KEY=YOUR_GEMINI_KEY" \
  --set-env-vars "SERVER_ENV=production"
```

Now copy the deployed worker URL from the gcloud output (looks like `https://krisho-gcp-worker-abcdef.a.run.app`) to use in the API deployment step!


---

## Step 5: Environment Variables

### Required Environment Variables for API (Cloud Run)
```
QUEUE_TYPE=cloudtasks
GCP_PROJECT_ID=your-gcp-project
GCP_LOCATION=us-central1
CLOUD_TASKS_QUEUE=krisho-queue
CLOUD_TASKS_WORKER_URL=https://your-worker.a.run.app
GCS_BUCKET=your-gcs-bucket
DATABASE_URL=postgresql://user:pass@host:5432/dbname
GEMINI_API_KEY=your-gemini-api-key
RAZORPAY_KEY_ID=your-razorpay-key-id
RAZORPAY_KEY_SECRET=your-razorpay-key-secret
RAZORPAY_WEBHOOK_SECRET=your-razorpay-webhook-secret
SERVER_ENV=production
```

### Required Environment Variables for Worker (Cloud Run)
```
DATABASE_URL=postgresql://user:pass@host:5432/dbname
GEMINI_API_KEY=your-gemini-api-key
SERVER_ENV=production
```

---

## Cost Estimation (Monthly)

| Service | Estimated Cost (Small Scale) | Notes |
|---------|-------------------------------|-------|
| Cloud Run (API) | $5 - $15 | Scales to 0 when idle |
| Cloud Run (Worker) | $5 - $15 | Scales to 0 when idle |
| Cloud Tasks | $0.50 - $2 | $0.04 per million tasks |
| Supabase (Free Tier) | $0 | Up to 500MB storage |
| Supabase (Pro) | $25 | For larger databases |
| **Total** | **~$10 - $40** | 70% less than VM-based setup |

---

## Auto-Scaling Configuration

Cloud Run auto-scales automatically based on request load. For best performance:

- **API Service**: Set `--max-instances=1000`
- **Worker Service**: Set `--max-instances=500`

You can adjust these in GCP Console → Cloud Run → Service → Revision → Edit & Deploy New Revision.

---

## Migration from Redis/Asynq

To switch from Redis/Asynq to Cloud Tasks, simply set the `QUEUE_TYPE` environment variable to `cloudtasks` on your API Cloud Run service.

## Troubleshooting

- **Tasks not executing?** Check that the `CLOUD_TASKS_WORKER_URL` is correct and the worker service allows unauthenticated requests.
- **Permission issues?** Ensure your Cloud Run service account has the `Cloud Tasks Enqueuer` role.
- **Database timeouts?** Increase `DB_MAX_CONNS` in environment variables.

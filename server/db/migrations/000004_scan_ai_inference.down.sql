ALTER TABLE scans
  DROP COLUMN IF EXISTS completed_at,
  DROP COLUMN IF EXISTS ai_metadata,
  DROP COLUMN IF EXISTS processing_error,
  DROP COLUMN IF EXISTS confidence_score,
  DROP COLUMN IF EXISTS ai_provider;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS is_suspended BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS analytics_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  event_type VARCHAR(80) NOT NULL,
  actor_user_id UUID,
  company_id UUID,
  entity_id UUID,
  order_id UUID,
  payment_id UUID,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  value NUMERIC(14,2),
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_type_time ON analytics_events(event_type, occurred_at);
CREATE INDEX IF NOT EXISTS idx_analytics_events_company_time ON analytics_events(company_id, occurred_at);
CREATE INDEX IF NOT EXISTS idx_analytics_events_actor_time ON analytics_events(actor_user_id, occurred_at);

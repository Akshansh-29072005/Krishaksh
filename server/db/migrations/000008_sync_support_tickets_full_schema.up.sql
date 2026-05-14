-- Ensure support_tickets has all required columns for the current application schema
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS title VARCHAR(200) NOT NULL;
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'open';
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'medium';
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS assigned_to UUID;
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS callback_requested BOOLEAN DEFAULT FALSE;
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS callback_status VARCHAR(20) DEFAULT 'none';
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMPTZ;

UPDATE support_tickets SET description = title WHERE description IS NULL OR description = '';
UPDATE support_tickets SET status = 'open' WHERE status IS NULL OR status = '';
UPDATE support_tickets SET priority = 'medium' WHERE priority IS NULL OR priority = '';
UPDATE support_tickets SET callback_requested = FALSE WHERE callback_requested IS NULL;
UPDATE support_tickets SET callback_status = 'none' WHERE callback_status IS NULL OR callback_status = '';
UPDATE support_tickets SET created_at = NOW() WHERE created_at IS NULL;
UPDATE support_tickets SET updated_at = NOW() WHERE updated_at IS NULL;

ALTER TABLE support_tickets ALTER COLUMN title SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN description SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN status SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN priority SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN callback_requested SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN callback_status SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE support_tickets ALTER COLUMN status DROP DEFAULT;
ALTER TABLE support_tickets ALTER COLUMN priority DROP DEFAULT;
ALTER TABLE support_tickets ALTER COLUMN callback_requested DROP DEFAULT;
ALTER TABLE support_tickets ALTER COLUMN callback_status DROP DEFAULT;
ALTER TABLE support_tickets ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE support_tickets ALTER COLUMN updated_at DROP DEFAULT;

-- Ensure support_tickets has all expected columns for the current application schema
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';
UPDATE support_tickets SET description = title WHERE description IS NULL OR description = '';
ALTER TABLE support_tickets ALTER COLUMN description SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN description DROP DEFAULT;

ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'medium';
UPDATE support_tickets SET priority = 'medium' WHERE priority IS NULL OR priority = '';
ALTER TABLE support_tickets ALTER COLUMN priority SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN priority DROP DEFAULT;

ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS callback_requested BOOLEAN DEFAULT FALSE;
ALTER TABLE support_tickets ALTER COLUMN callback_requested SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN callback_requested DROP DEFAULT;

ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS callback_status VARCHAR(20) DEFAULT 'none';
UPDATE support_tickets SET callback_status = 'none' WHERE callback_status IS NULL OR callback_status = '';
ALTER TABLE support_tickets ALTER COLUMN callback_status SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN callback_status DROP DEFAULT;

ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS assigned_to UUID;

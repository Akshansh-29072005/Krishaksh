-- Rollback for support_tickets schema sync migration
ALTER TABLE support_tickets DROP COLUMN IF EXISTS assigned_to;
ALTER TABLE support_tickets DROP COLUMN IF EXISTS callback_status;
ALTER TABLE support_tickets DROP COLUMN IF EXISTS callback_requested;
ALTER TABLE support_tickets DROP COLUMN IF EXISTS priority;
ALTER TABLE support_tickets DROP COLUMN IF EXISTS description;

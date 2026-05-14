-- Add missing description column to support_tickets table
ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS description TEXT;

-- Populate any missing description values
UPDATE support_tickets SET description = title WHERE description IS NULL OR description = '';

-- Ensure the column is not nullable and drop the default if present
ALTER TABLE support_tickets ALTER COLUMN description SET NOT NULL;
ALTER TABLE support_tickets ALTER COLUMN description DROP DEFAULT;

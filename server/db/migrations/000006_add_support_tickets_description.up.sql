-- Add missing description column to support_tickets table
ALTER TABLE support_tickets ADD COLUMN description TEXT NOT NULL DEFAULT '';

-- Update default for any existing rows
UPDATE support_tickets SET description = title WHERE description = '';

-- Remove the default constraint now that rows have been populated
ALTER TABLE support_tickets ALTER COLUMN description DROP DEFAULT;

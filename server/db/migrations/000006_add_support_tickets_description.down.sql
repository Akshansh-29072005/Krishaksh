-- Rollback: Remove description column from support_tickets
ALTER TABLE support_tickets DROP COLUMN description;

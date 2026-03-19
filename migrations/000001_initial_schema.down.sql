-- =============================================
-- LEAVE MANAGEMENT SYSTEM - ROLLBACK
-- =============================================

-- Drop triggers
DROP TRIGGER IF EXISTS update_leave_requests_updated_at ON leave_requests CASCADE;
DROP TRIGGER IF EXISTS update_leave_balances_updated_at ON leave_balances CASCADE;
DROP TRIGGER IF EXISTS update_leave_types_updated_at ON leave_types CASCADE;
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees CASCADE;
DROP TRIGGER IF EXISTS update_departments_updated_at ON departments CASCADE;
DROP TRIGGER IF EXISTS update_users_updated_at ON users CASCADE;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop views
DROP VIEW IF EXISTS pending_approvals CASCADE;
DROP VIEW IF EXISTS employee_leave_summary CASCADE;

-- Drop tables
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS leave_requests CASCADE;
DROP TABLE IF EXISTS leave_balances CASCADE;
DROP TABLE IF EXISTS leave_types CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
DROP EXTENSION IF EXISTS "pg_trgm" CASCADE;
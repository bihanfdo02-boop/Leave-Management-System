-- =============================================
-- LEAVE MANAGEMENT SYSTEM - INITIAL SCHEMA
-- =============================================

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================
-- 1. USERS TABLE (Authentication)
-- =============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee',
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 2. DEPARTMENTS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 3. EMPLOYEES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    manager_id UUID REFERENCES employees(id) ON DELETE SET NULL,
    hire_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 4. LEAVE TYPES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS leave_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    default_days_per_year INT NOT NULL DEFAULT 20,
    is_paid BOOLEAN DEFAULT true,
    requires_approval BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 5. LEAVE BALANCES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    total_days INT NOT NULL,
    used_days INT DEFAULT 0,
    remaining_days INT GENERATED ALWAYS AS (total_days - used_days) STORED,
    year INT NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_DATE),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, leave_type_id, year)
);

-- =============================================
-- 6. LEAVE REQUESTS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    number_of_days DECIMAL(5,2) NOT NULL CHECK (number_of_days > 0),
    reason TEXT NOT NULL,
    attachment_url VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    approved_by UUID REFERENCES employees(id) ON DELETE SET NULL,
    approval_date TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (end_date >= start_date)
);

-- =============================================
-- 7. AUDIT LOGS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 8. CREATE INDEXES (IF NOT EXISTS for safety)
-- =============================================

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Departments indexes
CREATE INDEX IF NOT EXISTS idx_departments_name ON departments(name);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments(is_active);

-- Employees indexes
CREATE INDEX IF NOT EXISTS idx_employees_user_id ON employees(user_id);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees(department_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees(manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_is_active ON employees(is_active);

-- Leave types indexes
CREATE INDEX IF NOT EXISTS idx_leave_types_name ON leave_types(name);
CREATE INDEX IF NOT EXISTS idx_leave_types_is_active ON leave_types(is_active);

-- Leave balances indexes
CREATE INDEX IF NOT EXISTS idx_leave_balances_employee_id ON leave_balances(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_balances_leave_type_id ON leave_balances(leave_type_id);
CREATE INDEX IF NOT EXISTS idx_leave_balances_year ON leave_balances(year);

-- Leave requests indexes
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);
CREATE INDEX IF NOT EXISTS idx_leave_requests_created_at ON leave_requests(created_at);
CREATE INDEX IF NOT EXISTS idx_leave_requests_start_date ON leave_requests(start_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_end_date ON leave_requests(end_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_status ON leave_requests(employee_id, status);

-- Audit logs indexes
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- =============================================
-- 9. CREATE VIEWS (IF NOT EXISTS for safety)
-- =============================================

DROP VIEW IF EXISTS pending_approvals CASCADE;
DROP VIEW IF EXISTS employee_leave_summary CASCADE;

-- View for leave summary
CREATE VIEW employee_leave_summary AS
SELECT 
    e.id as employee_id,
    e.first_name,
    e.last_name,
    e.email,
    d.name as department,
    lt.name as leave_type,
    lb.total_days,
    lb.used_days,
    lb.remaining_days,
    lb.year
FROM employees e
JOIN departments d ON e.department_id = d.id
JOIN leave_balances lb ON e.id = lb.employee_id
JOIN leave_types lt ON lb.leave_type_id = lt.id
WHERE e.is_active = true AND d.is_active = true;

-- View for pending approvals
CREATE VIEW pending_approvals AS
SELECT 
    lr.id,
    lr.employee_id,
    e.first_name || ' ' || e.last_name as employee_name,
    e.email,
    lt.name as leave_type,
    lr.start_date,
    lr.end_date,
    lr.number_of_days,
    lr.reason,
    lr.created_at,
    m.first_name || ' ' || m.last_name as manager_name,
    m.id as manager_id
FROM leave_requests lr
JOIN employees e ON lr.employee_id = e.id
JOIN leave_types lt ON lr.leave_type_id = lt.id
JOIN employees m ON e.manager_id = m.id
WHERE lr.status = 'pending'
ORDER BY lr.created_at DESC;

-- =============================================
-- 10. CREATE FUNCTIONS (IF NOT EXISTS for safety)
-- =============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- 11. CREATE TRIGGERS (DROP AND RECREATE)
-- =============================================

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_departments_updated_at ON departments;
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
DROP TRIGGER IF EXISTS update_leave_types_updated_at ON leave_types;
DROP TRIGGER IF EXISTS update_leave_balances_updated_at ON leave_balances;
DROP TRIGGER IF EXISTS update_leave_requests_updated_at ON leave_requests;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON departments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leave_types_updated_at BEFORE UPDATE ON leave_types
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leave_balances_updated_at BEFORE UPDATE ON leave_balances
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leave_requests_updated_at BEFORE UPDATE ON leave_requests
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
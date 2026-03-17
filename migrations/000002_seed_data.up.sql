-- Seed a default admin user
-- Password is 'password123'
-- The hash below is a valid bcrypt hash for 'password123'
INSERT INTO users (id, username, email, password_hash, full_name, role)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 
    'admin', 
    'admin@devpulse.io', 
    '$2y$10$LBBhR/.ByfO6jH8hG2r6.e2m2j9M2rOQ1O5Uo.9u5.M7hS7o/J.yG', 
    'DevPulse Administrator', 
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Seed a sample project
INSERT INTO projects (id, name, description, owner_id)
VALUES (
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
    'Project Apollo',
    'Next-generation synchronization engine for global teams.',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
) ON CONFLICT DO NOTHING;

-- Seed some tasks
INSERT INTO tasks (title, description, status, priority, project_id, assigned_to)
VALUES 
('Fix production API latency', 'Identify bottlenecks in the request pipeline', 'todo', 'urgent', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('Implement Redis caching', 'Layer 2 cache for frequent DB queries', 'in-progress', 'high', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT DO NOTHING;

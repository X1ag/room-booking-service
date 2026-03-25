INSERT INTO users (id, email, password_hash, role)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'admin@example.com', 'dummy', 'admin'),
  ('22222222-2222-2222-2222-222222222222', 'user@example.com', 'dummy', 'user')
ON CONFLICT (id) DO NOTHING;

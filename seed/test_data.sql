INSERT INTO rooms (id, name, description, capacity)
VALUES
  (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
    'Room A',
    'test room',
    6
  ),
  (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
    'Room B',
    'another test room',
    10
  ),
  (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
    'Room C',
    'room without schedule',
    4
  )
ON CONFLICT (id) DO NOTHING;

INSERT INTO schedules (id, room_id, start_time, end_time)
VALUES
  (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
    '09:00',
    '18:00'
  ),
  (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
    '10:00',
    '14:00'
  )
ON CONFLICT (room_id) DO NOTHING;

INSERT INTO schedule_days (schedule_id, day_of_week)
VALUES
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 1),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 2),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 3),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 4),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 5),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 2),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 4)
ON CONFLICT (schedule_id, day_of_week) DO NOTHING;

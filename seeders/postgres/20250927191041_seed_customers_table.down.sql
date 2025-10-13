-- Rollback seeder: Remove seeded customers
DELETE FROM customers WHERE email IN (
    'alice.j@example.com',
    'bob.w@example.com',
    'charlie.b@example.com',
    'diana.m@example.com',
    'ethan.d@example.com'
);

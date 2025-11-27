// Rollback seeder: Remove seeded customers
db.customers.deleteMany({
    email: {
        $in: [
            "alice.j@example.com",
            "bob.w@example.com",
            "charlie.b@example.com",
            "diana.m@example.com",
            "ethan.d@example.com"
        ]
    }
});

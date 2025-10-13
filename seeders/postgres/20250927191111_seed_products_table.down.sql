-- Rollback seeder: Remove seeded products
DELETE FROM products WHERE name IN (
    'Laptop Pro', 'Smartphone X', 'Wireless Headphones', 'Smartwatch 2', '4K Monitor',
    'Mechanical Keyboard', 'Webcam HD', 'USB-C Hub', 'External SSD 1TB', 'Gaming Mouse',
    'Standing Desk', 'Ergonomic Chair', 'LED Desk Lamp', 'Portable Charger', 'Bluetooth Speaker',
    'Graphics Tablet', 'VR Headset', 'Drone with Camera', 'E-reader', 'Fitness Tracker',
    'Electric Toothbrush', 'Coffee Maker', 'Blender', 'Air Fryer', 'Robot Vacuum',
    'Security Camera', 'Smart Thermostat', 'Wi-Fi Router', 'Projector', 'Soundbar',
    'Action Camera', 'Digital Photo Frame', 'Electric Kettle', 'Microwave Oven', 'Toaster',
    'Hair Dryer', 'Electric Shaver', 'Yoga Mat', 'Dumbbell Set', 'Resistance Bands',
    'Water Bottle', 'Backpack', 'Sunglasses', 'Wristwatch', 'Wallet',
    'Desk Organizer', 'Notebook and Pen Set', 'Wall Art', 'Scented Candle', 'Board Game'
);

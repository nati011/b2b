-- Seed data for products table
-- Development seed data
-- Note: supplier_id references suppliers table (assuming IDs start from 1)

INSERT INTO public.products (name, description, external_id, attributes, unit, is_active, supplier_id, price, total_quantity, reserved_quantity, created_date, last_modified, is_deleted)
VALUES
  -- Electronics (Supplier 1)
  ('Laptop Computer', 'High-performance laptop for business use', 'PROD-ELEC-001', '{"brand": "Dell", "model": "Latitude 5520", "processor": "Intel i7", "ram": "16GB", "storage": "512GB SSD"}'::jsonb, 'unit', TRUE, 1, 45000.00, 50, 5, NOW(), NOW(), FALSE),
  ('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 'PROD-ELEC-002', '{"brand": "Logitech", "model": "M705", "connectivity": "Wireless", "battery": "AA x 2"}'::jsonb, 'unit', TRUE, 1, 850.00, 200, 20, NOW(), NOW(), FALSE),
  ('USB Keyboard', 'Mechanical keyboard with backlight', 'PROD-ELEC-003', '{"brand": "Corsair", "model": "K70", "switch_type": "Cherry MX", "backlight": "RGB"}'::jsonb, 'unit', TRUE, 1, 3500.00, 100, 10, NOW(), NOW(), FALSE),
  
  -- Office Supplies (Supplier 2)
  ('A4 Paper Ream', 'Premium quality A4 paper, 500 sheets per ream', 'PROD-OFF-001', '{"size": "A4", "weight": "80gsm", "sheets": 500, "color": "White"}'::jsonb, 'ream', TRUE, 2, 250.00, 500, 50, NOW(), NOW(), FALSE),
  ('Ballpoint Pen Set', 'Pack of 12 blue ballpoint pens', 'PROD-OFF-002', '{"color": "Blue", "pack_size": 12, "ink_type": "Ballpoint"}'::jsonb, 'pack', TRUE, 2, 120.00, 1000, 100, NOW(), NOW(), FALSE),
  ('Stapler', 'Heavy-duty stapler for office use', 'PROD-OFF-003', '{"capacity": 210, "color": "Black", "material": "Metal"}'::jsonb, 'unit', TRUE, 2, 450.00, 300, 30, NOW(), NOW(), FALSE),
  
  -- Tools & Equipment (Supplier 3)
  ('Power Drill', 'Cordless power drill with battery and charger', 'PROD-TOOL-001', '{"brand": "Bosch", "voltage": "18V", "battery": "Lithium-ion", "chuck_size": "13mm"}'::jsonb, 'unit', TRUE, 3, 8500.00, 75, 8, NOW(), NOW(), FALSE),
  ('Hammer Set', 'Professional hammer set with 3 different sizes', 'PROD-TOOL-002', '{"pieces": 3, "weight_range": "500g-2kg", "handle": "Fiberglass"}'::jsonb, 'set', TRUE, 3, 1200.00, 150, 15, NOW(), NOW(), FALSE),
  ('Measuring Tape', '25-meter steel measuring tape', 'PROD-TOOL-003', '{"length": "25m", "width": "25mm", "material": "Steel", "case": "Plastic"}'::jsonb, 'unit', TRUE, 3, 350.00, 400, 40, NOW(), NOW(), FALSE),
  
  -- Furniture (Supplier 4)
  ('Office Chair', 'Ergonomic office chair with adjustable height', 'PROD-FURN-001', '{"material": "Mesh", "color": "Black", "weight_capacity": "120kg", "wheels": 5}'::jsonb, 'unit', TRUE, 4, 4500.00, 80, 8, NOW(), NOW(), FALSE),
  ('Desk Set', 'Modern office desk with drawers', 'PROD-FURN-002', '{"material": "MDF", "color": "Brown", "dimensions": "120x60x75cm", "drawers": 3}'::jsonb, 'unit', TRUE, 4, 8500.00, 60, 6, NOW(), NOW(), FALSE),
  ('Filing Cabinet', '4-drawer metal filing cabinet', 'PROD-FURN-003', '{"material": "Steel", "color": "Gray", "drawers": 4, "lock": "Yes"}'::jsonb, 'unit', TRUE, 4, 6500.00, 40, 4, NOW(), NOW(), FALSE),
  
  -- Safety & Security (Supplier 5)
  ('Safety Helmet', 'Industrial safety helmet with chin strap', 'PROD-SAFE-001', '{"color": "Yellow", "size": "Adjustable", "standard": "ANSI Z89.1", "material": "HDPE"}'::jsonb, 'unit', TRUE, 5, 850.00, 500, 50, NOW(), NOW(), FALSE),
  ('Safety Gloves', 'Cut-resistant safety gloves, pack of 12', 'PROD-SAFE-002', '{"size": "Large", "material": "Kevlar", "pack_size": 12, "protection_level": "Level 5"}'::jsonb, 'pack', TRUE, 5, 1200.00, 200, 20, NOW(), NOW(), FALSE),
  ('First Aid Kit', 'Comprehensive first aid kit for workplace', 'PROD-SAFE-003', '{"items": 150, "size": "Large", "case": "Plastic", "expiry": "2026-12-31"}'::jsonb, 'unit', TRUE, 5, 2500.00, 100, 10, NOW(), NOW(), FALSE),
  
  -- Industrial (Supplier 6)
  ('Steel Pipe', 'Galvanized steel pipe, 1 inch diameter, 6 meters', 'PROD-IND-001', '{"diameter": "1 inch", "length": "6m", "material": "Galvanized Steel", "thickness": "3mm"}'::jsonb, 'piece', TRUE, 6, 850.00, 300, 30, NOW(), NOW(), FALSE),
  ('Welding Machine', 'Arc welding machine 200A', 'PROD-IND-002', '{"type": "Arc Welder", "amperage": "200A", "voltage": "220V", "duty_cycle": "60%"}'::jsonb, 'unit', TRUE, 6, 15000.00, 25, 3, NOW(), NOW(), FALSE),
  ('Angle Grinder', 'Electric angle grinder 4.5 inch', 'PROD-IND-003', '{"power": "850W", "disc_size": "4.5 inch", "voltage": "220V", "speed": "11000 RPM"}'::jsonb, 'unit', TRUE, 6, 3200.00, 100, 10, NOW(), NOW(), FALSE),
  
  -- Construction Materials (Supplier 7)
  ('Cement Bag', 'Portland cement, 50kg bag', 'PROD-CONST-001', '{"weight": "50kg", "type": "Portland", "grade": "42.5R", "standard": "ES 1170"}'::jsonb, 'bag', TRUE, 7, 550.00, 2000, 200, NOW(), NOW(), FALSE),
  ('Steel Rebar', 'Reinforcement steel bar, 12mm diameter, 12 meters', 'PROD-CONST-002', '{"diameter": "12mm", "length": "12m", "grade": "Grade 60", "standard": "ASTM A615"}'::jsonb, 'piece', TRUE, 7, 450.00, 1500, 150, NOW(), NOW(), FALSE),
  ('Concrete Block', 'Hollow concrete block, standard size', 'PROD-CONST-003', '{"size": "40x20x20cm", "type": "Hollow", "strength": "7.5 MPa", "pieces_per_pallet": 80}'::jsonb, 'piece', TRUE, 7, 25.00, 5000, 500, NOW(), NOW(), FALSE),
  
  -- Cleaning Supplies (Supplier 8)
  ('Floor Cleaner', 'Concentrated floor cleaner, 5 liters', 'PROD-CLEAN-001', '{"volume": "5L", "type": "Concentrate", "scent": "Lemon", "dilution": "1:100"}'::jsonb, 'bottle', TRUE, 8, 450.00, 300, 30, NOW(), NOW(), FALSE),
  ('Mop Set', 'Professional mop with bucket and wringer', 'PROD-CLEAN-002', '{"pieces": 3, "bucket_size": "20L", "material": "Microfiber", "handle": "Aluminum"}'::jsonb, 'set', TRUE, 8, 1200.00, 150, 15, NOW(), NOW(), FALSE),
  ('Trash Bags', 'Heavy-duty trash bags, pack of 50', 'PROD-CLEAN-003', '{"size": "50L", "thickness": "30 microns", "pack_size": 50, "color": "Black"}'::jsonb, 'pack', TRUE, 8, 350.00, 500, 50, NOW(), NOW(), FALSE)
ON CONFLICT DO NOTHING;











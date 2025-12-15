-- PostgreSQL 手动测试数据库初始化脚本

-- 创建数据库 (需要连接到 postgres 数据库执行)
-- CREATE DATABASE testdb;

-- 连接到 testdb 后执行以下语句

-- 1. 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    age INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. 创建产品表
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. 创建订单详情表
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入测试数据

-- 插入用户数据
INSERT INTO users (username, email, age) VALUES
    ('alice', 'alice@example.com', 25),
    ('bob', 'bob@example.com', 30),
    ('charlie', 'charlie@example.com', 22),
    ('david', 'david@example.com', 28),
    ('eve', 'eve@example.com', 35)
ON CONFLICT (username) DO NOTHING;

-- 插入产品数据
INSERT INTO products (name, category, price, stock) VALUES
    ('iPhone 15', 'Electronics', 999.99, 50),
    ('MacBook Pro', 'Electronics', 2499.99, 30),
    ('AirPods Pro', 'Electronics', 249.99, 100),
    ('iPad Air', 'Electronics', 599.99, 40),
    ('Apple Watch', 'Electronics', 399.99, 60),
    ('Desk Chair', 'Furniture', 299.99, 20),
    ('Standing Desk', 'Furniture', 599.99, 15),
    ('Monitor', 'Electronics', 449.99, 35),
    ('Keyboard', 'Electronics', 129.99, 80),
    ('Mouse', 'Electronics', 79.99, 100)
ON CONFLICT DO NOTHING;

-- 插入订单数据
INSERT INTO orders (user_id, total_price, status) VALUES
    (1, 1249.98, 'completed'),
    (2, 2499.99, 'completed'),
    (3, 329.98, 'pending'),
    (1, 999.99, 'completed'),
    (4, 899.98, 'shipped')
ON CONFLICT DO NOTHING;

-- 插入订单详情数据
INSERT INTO order_items (order_id, product_id, quantity, price) VALUES
    (1, 1, 1, 999.99),
    (1, 3, 1, 249.99),
    (2, 2, 1, 2499.99),
    (3, 3, 1, 249.99),
    (3, 5, 1, 79.99),
    (4, 1, 1, 999.99),
    (5, 4, 1, 599.99),
    (5, 6, 1, 299.99)
ON CONFLICT DO NOTHING;

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

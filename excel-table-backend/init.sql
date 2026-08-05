-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role ENUM('admin','export_input','export_review','import_input','import_review') NOT NULL,
    can_edit_products BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT NOW()
);

-- 用户地区关联表（只对C/D类用户有效）
CREATE TABLE IF NOT EXISTS user_regions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    region VARCHAR(50) NOT NULL,
    created_at DATETIME DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY uk_user_region (user_id, region)
);

-- 产品基础属性表（强绑定，所有地区共用）
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    model_no VARCHAR(100) UNIQUE NOT NULL,
    brand VARCHAR(100),
    commodities VARCHAR(255),
    descriptions VARCHAR(255),
    unit VARCHAR(20),
    updated_at DATETIME DEFAULT NOW(),
    updated_by VARCHAR(50)
);

-- 产品地区价格表（弱绑定，各地区独立）
CREATE TABLE IF NOT EXISTS product_prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    model_no VARCHAR(100) NOT NULL,
    region VARCHAR(50) NOT NULL,
    export_ref_price  DECIMAL(15,4),  
    import_ref_cost  DECIMAL(15,4),
    updated_at DATETIME DEFAULT NOW(),
    updated_by VARCHAR(50),
    UNIQUE KEY uk_model_region (model_no, region),
    FOREIGN KEY (model_no) REFERENCES products(model_no)
);

-- 表单主表
CREATE TABLE IF NOT EXISTS invoices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    inv_no VARCHAR(50) NOT NULL,
    invoice_type ENUM('export','import') NOT NULL,
    region VARCHAR(50),
    invoice_date  DATE,
    contract_no VARCHAR(50),
    shipping_line VARCHAR(100),
    bl_no VARCHAR(100),
    container_no VARCHAR(100),
    port_of_loading VARCHAR(100),
    port_of_discharge VARCHAR(100),
    final_destination VARCHAR(100),
    buyer_name VARCHAR(200),
    buyer_address VARCHAR(500),
    buyer_tel VARCHAR(50),
    seller_name VARCHAR(200),
    payment_term VARCHAR(200),
    fob_total DECIMAL(15,2),
    reviewed BOOLEAN DEFAULT FALSE,
    reviewer VARCHAR(50),
    review_comment VARCHAR(500),
    created_at DATETIME DEFAULT NOW(),
    updated_at DATETIME DEFAULT NOW(),
    UNIQUE KEY uk_inv_type (inv_no, invoice_type)
);

-- 明细行表
CREATE TABLE IF NOT EXISTS invoice_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    inv_no VARCHAR(50) NOT NULL,
    invoice_type ENUM('export','import') NOT NULL,
    item_no INT,
    brand VARCHAR(100),
    commodities VARCHAR(255),
    model_no VARCHAR(100),
    carton_qty INT,
    unit_qty INT,
    unit VARCHAR(20),
    descriptions VARCHAR(255),
    export_unit_price  DECIMAL(15,4),
    total_amount DECIMAL(15,2),
    gross_weight DECIMAL(10,2),
    net_weight DECIMAL(10,2),
    volume DECIMAL(10,2),
    share_rate DECIMAL(10,8),
    sea_freight DECIMAL(15,4),
    cif_price DECIMAL(15,4),
    cif_total DECIMAL(15,2),
    cif_bs_total DECIMAL(15,2),
    import_duty_vat DECIMAL(15,2),
    transportation DECIMAL(15,2),
    others_charge DECIMAL(15,2),
    total_cost DECIMAL(15,2),
    import_unit_cost  DECIMAL(15,4),
    system_price DECIMAL(15,4),
    exchange_rate DECIMAL(10,4),
    price_alert ENUM('normal','yellow','red') DEFAULT 'normal',
    FOREIGN KEY (inv_no, invoice_type) REFERENCES invoices(inv_no, invoice_type)
);

-- 初始用户数据
INSERT IGNORE INTO users (username, password_hash, user_role, can_edit_products) VALUES 
('admin', 'changeme', 'admin', TRUE),
('dev', 'changeme', 'admin', TRUE),
('export_input', 'changeme', 'export_input', FALSE),
('export_review', 'changeme', 'export_review', FALSE),
('import_input_sa', 'changeme', 'import_input', FALSE),
('import_review_sa', 'changeme', 'import_review', FALSE);

-- 初始地区数据
INSERT IGNORE INTO user_regions (user_id, region)
SELECT id, '南美' FROM users WHERE username IN ('import_input_sa', 'import_review_sa');
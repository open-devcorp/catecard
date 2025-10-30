
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(20) NOT NULL
    );
    CREATE TABLE IF NOT EXISTS groups (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        catechist_id INT NOT NULL,
        limit_catechumens INT NOT NULL
    );

    INSERT INTO users (username, email, password, role)
    SELECT 'admin', 'admin', 'admin', '0'
    WHERE NOT EXISTS (
        SELECT 1 FROM users WHERE username = 'admin' OR email = 'admin'
    );
    --     INSERT INTO users (username, email, password, role)
    -- SELECT 'cate', 'cate', 'cate', '1'
    -- WHERE NOT EXISTS (
    --     SELECT 2 FROM users WHERE username = 'cate' OR email = 'cate'
    -- );
    -- INSERT INTO users (username, email, password, role)
    -- SELECT 'scan', 'scan', 'scan', '3'
    -- WHERE NOT EXISTS (
    --     SELECT 3 FROM users WHERE username = 'scan' OR email = 'scan'
    -- );

    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        user_json TEXT NOT NULL,
        expires_at TIMESTAMP NOT NULL
    );

    -- Nuevo esquema: qr_codes con referencia al catecúmeno
    CREATE TABLE IF NOT EXISTS qr_codes (
        id SERIAL PRIMARY KEY,
        code VARCHAR(255),
        catechumen_id INT NOT NULL,
        total_allowed INT NOT NULL,
        used_scans INT NOT NULL DEFAULT 0
    );

-- INSERT INTO groups (name, catechist_id)
-- SELECT 'Grupo A', 2
-- WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo A');

-- INSERT INTO groups (name, catechist_id)
-- SELECT 'Grupo B', 2
-- WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo B');

-- INSERT INTO groups (name, catechist_id)
-- SELECT 'Grupo C', 2
-- WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo C');


-- Datos de ejemplo opcionales para qr_codes (si existen catechumens con IDs 1..3)
-- INSERT INTO qr_codes (code, catechumen_id, total_allowed, used_scans) VALUES (NULL, 1, 1, 0);
-- INSERT INTO qr_codes (code, catechumen_id, total_allowed, used_scans) VALUES (NULL, 2, 1, 3);
-- INSERT INTO qr_codes (code, catechumen_id, total_allowed, used_scans) VALUES (NULL, 3, 2, 5);
CREATE TABLE IF NOT EXISTS catechumens (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    age VARCHAR(3) NOT NULL,
    group_id INT NOT NULL
);



CREATE TABLE IF NOT EXISTS scan_catechumens (
    id SERIAL PRIMARY KEY,
    catechumen_id INT NOT NULL,
    scan_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
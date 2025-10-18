CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
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
        catechist_id INT NOT NULL
    );

    INSERT INTO users (username, email, password, role)
    SELECT 'admin', 'admin', 'admin', '0'
    WHERE NOT EXISTS (
        SELECT 1 FROM users WHERE username = 'admin' OR email = 'admin'
    );

    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        user_json TEXT NOT NULL,
        expires_at TIMESTAMP NOT NULL
    );

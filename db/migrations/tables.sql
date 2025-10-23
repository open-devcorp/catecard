
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
        INSERT INTO users (username, email, password, role)
    SELECT 'cate', 'cate', 'cate', '1'
    WHERE NOT EXISTS (
        SELECT 2 FROM users WHERE username = 'cate' OR email = 'cate'
    );

    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        user_json TEXT NOT NULL,
        expires_at TIMESTAMP NOT NULL
    );

    CREATE TABLE IF NOT EXISTS qrs (
        id SERIAL PRIMARY KEY,
        forum INT NOT NULL,
        group_id INT NOT NULL,
        count INT DEFAULT 0

    );

INSERT INTO groups (name, catechist_id)
SELECT 'Grupo A', 2
WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo A');

INSERT INTO groups (name, catechist_id)
SELECT 'Grupo B', 2
WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo B');

INSERT INTO groups (name, catechist_id)
SELECT 'Grupo C', 2
WHERE NOT EXISTS (SELECT 1 FROM groups WHERE name = 'Grupo C');


INSERT INTO qrs (forum, group_id, count)
SELECT 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM qrs WHERE forum = 1 AND group_id = 1);

INSERT INTO qrs (forum, group_id, count)
SELECT 1, 2, 3
WHERE NOT EXISTS (SELECT 1 FROM qrs WHERE forum = 1 AND group_id = 2);

INSERT INTO qrs (forum, group_id, count)
SELECT 2, 1, 5
WHERE NOT EXISTS (SELECT 1 FROM qrs WHERE forum = 2 AND group_id = 1);

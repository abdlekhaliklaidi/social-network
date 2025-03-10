-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    gender TEXT CHECK(gender IN ('M', 'F')) NOT NULL,
    age INTEGER CHECK(age BETWEEN 1 AND 120) NOT NULL,
    password TEXT NOT NULL,
    session_token TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- func createTables() error {
-- 	_, err := DB.Exec(`
--         CREATE TABLE IF NOT EXISTS users (
--             id INTEGER PRIMARY KEY AUTOINCREMENT,
--             username TEXT UNIQUE NOT NULL ,
-- 			firstname TEXT NOT NULL,
--             lastname TEXT NOT NULL,
--             email TEXT UNIQUE NOT NULL ,
-- 			Gender TEXT CHECK(Gender IN ('M', 'F')) NOT NULL,
--             Age INTEGER CHECK(Age BETWEEN 1 AND 120) NOT NULL,
--             password TEXT NOT NULL,
-- 			session_token TEXT NOT NULL,
--             created_at DATETIME DEFAULT CURRENT_TIMESTAMP
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'users' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'users' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS users;

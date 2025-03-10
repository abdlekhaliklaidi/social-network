-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (session_token) ON DELETE CASCADE
        );

-- func createTables() error {
-- _, err = DB.Exec(`
--         CREATE TABLE IF NOT EXISTS posts (
--             id INTEGER PRIMARY KEY AUTOINCREMENT,
--             user_id INTEGER NOT NULL,
--             title TEXT NOT NULL,
--             content TEXT NOT NULL,
--             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--             FOREIGN KEY (user_id) REFERENCES users (session_token) ON DELETE CASCADE
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'posts' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'posts' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS posts;
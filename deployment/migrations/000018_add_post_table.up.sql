CREATE TABLE IF NOT EXISTS posts (
    id int NOT NULL UNIQUE AUTO_INCREMENT,
    user_id int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    body text NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES  users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (user_id, id) -- for clustered index
) CHARACTER SET utf8;
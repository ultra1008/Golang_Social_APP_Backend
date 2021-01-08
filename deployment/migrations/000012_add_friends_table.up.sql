CREATE TABLE IF NOT EXISTS friends (
    user_id int NOT NULL,
    friend_id int NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES  users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (friend_id)
        REFERENCES  users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);
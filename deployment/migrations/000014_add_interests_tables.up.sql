CREATE TABLE IF NOT EXISTS interests (
    id int NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    PRIMARY KEY (id),
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS user_interests (
    user_id int,
    interest_id int,
    FOREIGN KEY (user_id)
        REFERENCES  users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (interest_id)
        REFERENCES  interests(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);
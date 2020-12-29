CREATE TABLE IF NOT EXISTS citys (
    id int NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_by_user boolean DEFAULT 0,
    PRIMARY KEY (id)
);
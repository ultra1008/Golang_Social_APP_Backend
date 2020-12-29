CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age int NOT NULL,
    sex int NOT NULL,
    FOREIGN KEY (sex)
        REFERENCES  genders(id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY(id)
);
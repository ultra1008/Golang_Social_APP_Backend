ALTER TABLE users
DROP FOREIGN KEY users_ibfk_1,
DROP COLUMN sex,
ADD COLUMN sex_id int NOT NULL,
ADD FOREIGN KEY (sex_id)
    REFERENCES genders(id)
    ON DELETE CASCADE;
ALTER TABLE users
DROP FOREIGN KEY users_ibfk_3,
DROP COLUMN sex_id,
ADD COLUMN sex VARCHAR(15);
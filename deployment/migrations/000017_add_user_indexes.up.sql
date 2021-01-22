ALTER TABLE users
    ADD INDEX user_lastname_idx (last_name),
    ADD INDEX user_firstname_idx (first_name);
ALTER TABLE users
    ADD INDEX user_lastname_firstname_idx (last_name, first_name);
--
-- Remove everything
--
SET FOREIGN_KEY_CHECKS=0;

DELIMITER $$
DROP PROCEDURE IF EXISTS TruncateAllTables$$
CREATE PROCEDURE TruncateAllTables()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE _table_name VARCHAR(255);
    DECLARE cur CURSOR FOR SELECT table_name FROM information_schema.tables WHERE table_schema = 'binbogami' AND TABLE_TYPE = 'BASE TABLE';
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN cur;

    read_loop: LOOP
        FETCH cur INTO _table_name;
        IF done THEN
            LEAVE read_loop;
        END IF;

        SET @s = CONCAT('TRUNCATE TABLE binbogami.', _table_name);
        PREPARE stmt FROM @s;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END LOOP;

    CLOSE cur;
END$$
DELIMITER ;

CALL TruncateAllTables();
DROP PROCEDURE TruncateAllTables;

SET FOREIGN_KEY_CHECKS=1;

--
-- Add users
--
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('03a5ac74-b21b-11ee-9a7a-5ab75f0c1cab', 'george.hammond@sgc.gov', 'George', 'Hammond', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 13:52:39', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', 'jack.oneil@sgc.gov', 'Jack', 'O\'Neil', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 13:59:52', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', 'samantha.carter@sgc.gov', 'Samantha', 'Carter', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 14:00:27', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', 'daniel.jackson@sgc.gov', 'Daniel', 'Jackson', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 14:00:55', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('aff84550-b21f-11ee-8ac0-5ab75f0c1cab', 'tealc.of.chulak@sgc.gov', 'Teal\'c', 'Chulak', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 14:26:06', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('89705c9b-fe47-4c16-93f7-89cce8807f02', 'apophis@system.lords', 'Apophis', 'Goa\'uld', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 14:26:06', null, null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, created_at, updated_at, deleted_at) VALUES ('5f4dcc3b-5aa5-487d-8f6a-64cbd3665b4e', 'ra@system.lords', 'Ra', 'Goa\'uld', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', '2024-01-13 14:26:06', null, '2024-01-13 14:26:07');

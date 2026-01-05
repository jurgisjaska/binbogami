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
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('68b329da-9893-4d34-9d6b-549302554020', 'jonas.quinn@sgc.example.com', 'Jonas', 'Quinn', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', null, '2024-07-01 01:01:01', '2026-01-01 01:01:01');
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('aff84550-b21f-11ee-8ac0-5ab75f0c1cab', 'tealc.of.chulak@sgc.example.com', 'Teal\'c', 'Chulak', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', 'jack.oneil@sgc.example.com', 'Jack', 'O\'Neil', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 3, '2024-01-01 01:01:01', '2026-01-01 12:56:33', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', 'samantha.carter@sgc.example.com', 'Samantha', 'Carter', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', 'daniel.jackson@sgc.example.com', 'Daniel', 'Jackson', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', 'cameron.mitchell@sgc.example.com', 'Cameron', 'Mitchell', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', null, null, null);

--
-- Add invitations
--
INSERT INTO binbogami.invitations (id, email, role, created_by, created_at, opened_at, deleted_at, expired_at) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'random@sgc.example.com', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-05 15:46:35', null, null, '2028-01-01 01:01:01');

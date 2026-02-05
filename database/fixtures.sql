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
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('68b329da-9893-4d34-9d6b-549302554020', 'jonas.quinn@sgc.example.com', 'Jonas', 'Quinn', null, 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', null, '2024-07-01 01:01:01', '2026-01-01 01:01:01');
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('aff84550-b21f-11ee-8ac0-5ab75f0c1cab', 'tealc.of.chulak@sgc.example.com', 'Teal\'c', 'Chulak', null, 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', 'jack.oneil@sgc.example.com', 'Jack', 'O\'Neil', 'Colonel', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 3, '2024-01-01 01:01:01', '2026-02-05 13:12:16', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', 'samantha.carter@sgc.example.com', 'Samantha', 'Carter', 'Major', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-02-05 13:12:16', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', 'daniel.jackson@sgc.example.com', 'Daniel', 'Jackson', null, 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', 'cameron.mitchell@sgc.example.com', 'Cameron', 'Mitchell', null, 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', null, null, null);

--
-- Add invitations
--
INSERT INTO binbogami.invitations (id, email, role, created_by, created_at, opened_at, deleted_at, expired_at) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'random.1@sgc.example.com', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-05 15:46:35', null, null, '2028-01-01 01:01:01');
INSERT INTO binbogami.invitations (id, email, role, created_by, created_at, opened_at, deleted_at, expired_at) VALUES ('0685f091-63c3-4d24-8765-e35680621101', 'random.2@sgc.example.com', 2, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-05 15:46:35', null, null, '2028-01-01 01:01:01');

--
-- Add password resets
--
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:23', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('aa7d1712-c6d0-488a-8b0d-0172f28b05d0', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:24', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('f686a0b8-4886-42b4-a2a1-03b6bf21a045', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:27', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('c8b5bde0-6190-4ce7-84b9-0f8243ce274a', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:26', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('f729b34c-8563-4e9c-9ab8-161212b8fa1f', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:25', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('1b014ea0-fb87-469a-9c43-1da4def16164', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:56:18', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('d9cdff05-325a-46f3-a14a-223995f86334', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:56:20', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('6a8d7b6a-5de0-4961-b982-26943092884b', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:35', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('e2900b08-d6f5-4e98-9159-2f532d8c097c', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:35', null, '2027-01-01 01:01:01');
INSERT INTO binbogami.user_password_resets (id, user_id, ip, user_agent, created_at, opened_at, expire_at) VALUES ('4a60bad4-bdfd-4cb4-ba2d-30fde81eb90a', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '::1', 'PostmanRuntime/7.51.0', '2026-01-07 12:53:32', null, '2027-01-01 01:01:01');

--
-- Add books
--
INSERT INTO binbogami.books (id, name, description, created_by, created_at, updated_at, deleted_at, closed_at) VALUES ('5e6f7a8b-9c0d-4e1f-b2a3-4b5c6d7e8f9a', 'Year 2025', 'The book for year of 2025', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2025-01-12 16:46:42', null, null, '2026-01-13 16:46:52');
INSERT INTO binbogami.books (id, name, description, created_by, created_at, updated_at, deleted_at, closed_at) VALUES ('7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b', 'Year 2026', 'The book for year of 2026', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-01 16:46:32', null, null, null);
INSERT INTO binbogami.books (id, name, description, created_by, created_at, updated_at, deleted_at, closed_at) VALUES ('1a2b3c4d-5e6f-4789-a0b1-c2d3e4f5a6b7', 'Year 2025 deleted', 'Book that was created incorrect and deleted', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2025-01-12 16:46:47', null, '2026-01-13 16:46:59', null);

-- 
-- Add categories
-- 
INSERT INTO binbogami.categories (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9', 'Mission Supplies', 'P90 ammo, C4, and standard issue gear', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 08:30:00', null, null);
INSERT INTO binbogami.categories (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('f9e8d7c6-b5a4-4321-80f1-e2d3c4b5a697', 'Artifact Research', 'Tools and resources for archaeological analysis', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 09:15:00', null, null);
INSERT INTO binbogami.categories (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('1a2b3c4d-5e6f-4012-9345-6789abcdef01', 'Deep Space Telemetry', 'Cover story expenses and patent fees', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 10:00:00', null, null);
INSERT INTO binbogami.categories (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('01234567-89ab-4cde-a012-34567890abcd', 'Mess Hall Provisions', 'Candles and Jaffa cakes', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-01-15 11:45:00', null, null);
INSERT INTO binbogami.categories (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('fedcba98-7654-4321-8fed-cba987654321', 'NID Black Budget', 'Off-book acquisitions (Unauthorized)', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-01 00:00:00', '2024-01-20 12:00:00', '2024-01-20 12:00:00');

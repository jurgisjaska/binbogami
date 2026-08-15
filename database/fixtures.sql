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
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', 'jack.oneil@sgc.example.com', 'Jack', 'O\'Neil', 'Colonel', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 6, '2024-01-01 01:01:01', '2026-02-05 13:12:16', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', 'samantha.carter@sgc.example.com', 'Samantha', 'Carter', 'Major', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 4, '2024-01-01 01:01:01', '2026-02-05 13:12:16', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', 'daniel.jackson@sgc.example.com', 'Daniel', 'Jackson', null, 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', '2026-01-01 12:48:42', '2024-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', 'cameron.mitchell@sgc.example.com', 'Cameron', 'Mitchell', 'Colonel', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2024-01-01 01:01:01', null, null, null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('7c8f9a0b-1234-4567-8901-abcdef123456', 'janet.fraiser@sgc.example.com', 'Janet', 'Fraiser', 'Chief Medical Officer', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 4, '2026-01-01 01:01:01', null, '2026-01-01 01:01:01', null);
INSERT INTO binbogami.users (id, email, name, surname, position, salt, password, role, created_at, updated_at, confirmed_at, deleted_at) VALUES ('8d90ab1c-2345-5678-9012-bcdefa234567', 'vala.mal.doran@sgc.example.com', 'Vala', 'Mal Doran', 'Consultant', 'bUcdCORadqkbqHa1', '$2a$10$aIW8Elr5Q.2IXLl4ARI5hO6KGHT/DX4VGxPG0Od.CEUp7HQ.i8.Ry', 1, '2026-01-02 01:01:01', null, '2026-01-02 01:01:01', null);

--
-- Add invitations
--
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'george.hammond@sgc.example.com', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2026-01-05 15:46:35', null, null, '2028-01-01 01:01:01');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('0685f091-63c3-4d24-8765-e35680621101', 'walter.harriman@sgc.example.com', 2, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2026-01-05 15:46:35', null, null, '2028-06-01 12:00:00');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('b1fec990-1111-4ef8-bb6d-6bb9bd380a22', 'bratac@sgc.example.com', 1, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2026-01-05 15:46:35', '2026-01-06 10:00:00', null, '2028-12-31 23:59:59');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('c2afd001-2222-4ef8-bb6d-6bb9bd380a33', 'janet.fraiser@sgc.example.com', 4, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '7c8f9a0b-1234-4567-8901-abcdef123456', '2026-01-01 01:00:00', '2026-01-01 01:00:30', '2026-01-01 01:01:01', '2028-01-01 01:01:01');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('d3bae112-3333-4ef8-bb6d-6bb9bd380a44', 'vala.mal.doran@sgc.example.com', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '8d90ab1c-2345-5678-9012-bcdefa234567', '2026-01-02 01:00:00', '2026-01-02 01:00:30', '2026-01-02 01:01:01', '2028-01-01 01:01:01');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('e4cbf223-4444-4ef8-bb6d-6bb9bd380a55', 'charles.kawalsky@sgc.example.com', 1, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2026-01-05 15:46:35', null, '2026-01-07 12:00:00', '2028-12-31 23:59:59');
INSERT INTO binbogami.invitations (id, email, role, created_by, user_id, created_at, opened_at, deleted_at, expired_at) VALUES ('f5dc0334-5555-4ef8-bb6d-6bb9bd380a66', 'jacob.carter@sgc.example.com', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2026-01-05 15:46:35', '2026-01-06 09:00:00', '2026-01-07 10:00:00', '2029-01-01 00:00:00');

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
INSERT INTO binbogami.categories (id, name, description, color, created_by, created_at, updated_at, deleted_at) VALUES ('d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9', 'Mission Supplies', 'P90 ammo, C4, and standard issue gear', '#4CAF50', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 08:30:00', null, null);
INSERT INTO binbogami.categories (id, name, description, color, created_by, created_at, updated_at, deleted_at) VALUES ('f9e8d7c6-b5a4-4321-80f1-e2d3c4b5a697', 'Artifact Research', 'Tools and resources for archaeological analysis', '#2196F3', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 09:15:00', null, null);
INSERT INTO binbogami.categories (id, name, description, color, created_by, created_at, updated_at, deleted_at) VALUES ('1a2b3c4d-5e6f-4012-9345-6789abcdef01', 'Deep Space Telemetry', 'Cover story expenses and patent fees', '#9E9E9E', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 10:00:00', null, null);
INSERT INTO binbogami.categories (id, name, description, color, created_by, created_at, updated_at, deleted_at) VALUES ('01234567-89ab-4cde-a012-34567890abcd', 'Mess Hall Provisions', 'Candles and Jaffa cakes', '#FF9800', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-01-15 11:45:00', null, null);
INSERT INTO binbogami.categories (id, name, description, color, created_by, created_at, updated_at, deleted_at) VALUES ('fedcba98-7654-4321-8fed-cba987654321', 'NID Black Budget', 'Off-book acquisitions (Unauthorized)', '#000000', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-01 00:00:00', '2024-01-20 12:00:00', '2024-01-20 12:00:00');

--
-- Add locations (50 fixtures: 26 active, 24 deleted)
--
INSERT INTO binbogami.locations (id, name, description, address, created_by, created_at, updated_at, deleted_at) VALUES
('e0000000-0000-4000-8000-000000000001', 'Cheyenne Mountain Complex', 'SGC Headquarters Base', '1 Norad Rd, Colorado Springs, CO 80906, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-10 10:00:00', '2024-02-01 12:00:00', null),
('e0000000-0000-4000-8000-000000000002', 'Gate Room (Sub-Level 28)', 'Stargate Operation Room', 'Sub-Level 28, Cheyenne Mountain Complex, CO, USA', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-11 11:00:00', null, null),
('e0000000-0000-4000-8000-000000000003', 'Alpha Site', 'Primary off-world emergency evacuation site', null, '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-12 09:30:00', null, null),
('e0000000-0000-4000-8000-000000000004', 'Beta Site', 'Secondary backup off-world facility', 'P4X-639 Sector 4', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-15 14:00:00', '2024-03-10 08:15:00', null),
('e0000000-0000-4000-8000-000000000005', 'Gamma Site', 'Research facility focused on bio-weapons and alien flora', null, '7c8f9a0b-1234-4567-8901-abcdef123456', '2024-01-18 16:20:00', null, null),
('e0000000-0000-4000-8000-000000000006', 'Area 51 Groom Lake Facility', 'Top secret alien technology study lab', 'Groom Lake Rd, Rachel, NV 89001, USA', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-20 08:00:00', '2024-04-05 11:30:00', null),
('e0000000-0000-4000-8000-000000000007', 'Chulak High Council Chamber', null, 'Province of Ma\'tok, Chulak', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-01-22 13:45:00', null, null),
('e0000000-0000-4000-8000-000000000008', 'Tok\'ra Tunnels - Revanna', 'Underground crystal base of the Tok\'ra resistance', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-25 10:10:00', null, null),
('e0000000-0000-4000-8000-000000000009', 'Atlantis Control Tower', 'Central spire operations deck', 'Lantean Ocean, Pegasus Galaxy', 'c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', '2024-02-01 12:00:00', '2024-05-12 15:00:00', null),
('e0000000-0000-4000-8000-000000000010', 'Abydos Temple of Ra', 'Ancient pyramid research site', 'Great Pyramid Plaza, Abydos', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-03 09:00:00', null, null),
('e0000000-0000-4000-8000-000000000011', 'Langara Capital City Archive', 'Jonas Quinn\'s primary research vault', 'Kelowna District 1, Langara', '68b329da-9893-4d34-9d6b-549302554020', '2024-02-05 11:30:00', null, null),
('e0000000-0000-4000-8000-000000000012', 'Pentagon Deep Space Operations', 'Homeworld Command liaison office', '1400 Defense Pentagon, Washington, DC 20301, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-08 14:15:00', '2024-06-01 09:00:00', null),
('e0000000-0000-4000-8000-000000000013', 'P3X-888 Unas Village', 'Original homeworld of the Stargate symbiotes and Unas', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-10 15:45:00', null, null),
('e0000000-0000-4000-8000-000000000014', 'Tollana High Council Plaza', 'Tollan Curia administrative center', 'New Tollana Central Ring', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-12 10:20:00', null, null),
('e0000000-0000-4000-8000-000000000015', 'Dakara Temple Sanctuary', 'Ancient temple hosting the superweapon device', 'Dakara Sacred Valley', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-02-15 13:10:00', '2024-06-15 11:00:00', null),
('e0000000-0000-4000-8000-000000000016', 'SGC Infirmary (Sub-Level 21)', 'Chief Medical Officer surgical suite and quarantine', 'Sub-Level 21, Cheyenne Mountain Complex, CO, USA', '7c8f9a0b-1234-4567-8901-abcdef123456', '2024-02-18 08:30:00', null, null),
('e0000000-0000-4000-8000-000000000017', 'Lucian Alliance Trading Post', 'Black market bazaar and supply depot', 'Border Outpost 4, Sector 7', '8d90ab1c-2345-5678-9012-bcdefa234567', '2024-02-20 17:00:00', null, null),
('e0000000-0000-4000-8000-000000000018', 'Prometheus Hangar - Nevada', 'BC-303 battlecruiser assembly facility', 'Dry Lake Bed Airbase, NV, USA', 'c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', '2024-02-22 12:00:00', '2024-07-01 10:30:00', null),
('e0000000-0000-4000-8000-000000000019', 'P3X-984 Tok\'ra Alpha Base', 'Secondary subterranean outpost', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-25 14:30:00', null, null),
('e0000000-0000-4000-8000-000000000020', 'Hebridan Mining Outpost', 'Tech-exchange orbital station', 'Hebridan Outer Belt Station 3', '68b329da-9893-4d34-9d6b-549302554020', '2024-02-28 09:15:00', null, null),
('e0000000-0000-4000-8000-000000000021', 'Cimmeria Hall of Thor', 'Asgard protected world sanctuary', 'Thor\'s Valley, Cimmeria', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-03-01 11:00:00', null, null),
('e0000000-0000-4000-8000-000000000022', 'SGC Briefing Room (Sub-Level 27)', 'SG-1 mission planning conference hall', 'Sub-Level 27, Cheyenne Mountain Complex, CO, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-03-05 08:45:00', '2024-07-10 14:20:00', null),
('e0000000-0000-4000-8000-000000000023', 'P2X-555 Research Station', null, 'Polar Research Compound 1', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-03-08 16:10:00', null, null),
('e0000000-0000-4000-8000-000000000024', 'Free Jaffa Camp - Hak\'tyl', 'Ishta\'s warrior training encampment', 'Hak\'tyl Woodlands Sanctuary', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-03-10 13:00:00', null, null),
('e0000000-0000-4000-8000-000000000025', 'Vagonbrei Ruins', 'Arthurian legend search excavation site', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-03-12 10:00:00', null, null),
('e0000000-0000-4000-8000-000000000026', 'Vala\'s Warehouse Depot', 'Salvaged cargo storage hangar', 'Free Port Terminal 9, Sector 12', '8d90ab1c-2345-5678-9012-bcdefa234567', '2024-03-15 15:30:00', '2024-07-20 09:00:00', null),
('e0000000-0000-4000-8000-000000000027', 'Decommissioned Alpha Site 1', 'Compromised by Anubis forces and abandoned', 'P3X-984 Sector 2', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-05 08:00:00', '2024-01-14 12:00:00', '2024-01-14 12:05:00'),
('e0000000-0000-4000-8000-000000000028', 'Old NID Safehouse Alpha', 'Raid target, decommissioned by SGC', '404 Forest Rd, Seattle, WA 98101, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-06 09:00:00', null, '2024-01-20 18:00:00'),
('e0000000-0000-4000-8000-000000000029', 'Destroyed Tok\'ra Base - Vorash', 'Evacuated prior to star supernova expansion', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-07 10:30:00', null, '2024-02-01 10:00:00'),
('e0000000-0000-4000-8000-000000000030', 'P3X-797 Land of Light Camp', 'Temporary quarantine zone, relocated', 'Light Side Citadel Outpost', '7c8f9a0b-1234-4567-8901-abcdef123456', '2024-01-08 11:15:00', null, '2024-02-05 14:00:00'),
('e0000000-0000-4000-8000-000000000031', 'P3X-595 Drinking Club Outpost', 'Unauthorized recreation hub', null, '8d90ab1c-2345-5678-9012-bcdefa234567', '2024-01-09 16:00:00', null, '2024-02-10 11:00:00'),
('e0000000-0000-4000-8000-000000000032', 'P34-353 Naquadah Mine Site 1', 'Exhausted mineral extraction pit', 'Sector 1 Mine shaft A', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-10 13:45:00', null, '2024-02-15 16:30:00'),
('e0000000-0000-4000-8000-000000000033', 'P4X-639 Time Loop Observatory', 'Malfunctioning Ancient device site, sealed off', 'Plateau Ridge Ruins', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-12 14:00:00', null, '2024-02-20 09:10:00'),
('e0000000-0000-4000-8000-000000000034', 'Old Sub-Level 28 Armory', 'Remodeled into secondary generator room', 'Sub-Level 28 Room 14B, Cheyenne Mountain, CO, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-14 15:20:00', '2024-02-22 10:00:00', '2024-02-25 12:00:00'),
('e0000000-0000-4000-8000-000000000035', 'Ketan Research Outpost', 'Biological containment failure, destroyed', null, '7c8f9a0b-1234-4567-8901-abcdef123456', '2024-01-15 08:30:00', null, '2024-03-01 17:00:00'),
('e0000000-0000-4000-8000-000000000036', 'P3X-234 Sentinel Shrine', 'Alien planetary defense array location', 'High Mountain Summit Altar', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-16 09:50:00', null, '2024-03-05 13:40:00'),
('e0000000-0000-4000-8000-000000000037', 'Apophis Flagship Docking Bay', 'Goa\'uld vessel mothership hangar, obliterated', 'Netu Orbital Ring', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-01-18 10:10:00', null, '2024-03-10 11:20:00'),
('e0000000-0000-4000-8000-000000000038', 'P3X-562 Crystal Entity Valley', 'Restricted zone after incident', 'Canyon Sector 5', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-20 11:00:00', null, '2024-03-15 15:00:00'),
('e0000000-0000-4000-8000-000000000039', 'NID Warehouse - Chicago', 'Seized illegal arms facility', '550 W Adams St, Chicago, IL 60661, USA', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-22 12:30:00', '2024-03-18 09:00:00', '2024-03-20 10:00:00'),
('e0000000-0000-4000-8000-000000000040', 'Kelowna Trinium Refinery', 'Industrial plant collapsed in earthquake', 'Industrial Zone 4, Kelowna, Langara', '68b329da-9893-4d34-9d6b-549302554020', '2024-01-24 14:15:00', null, '2024-03-25 14:30:00'),
('e0000000-0000-4000-8000-000000000041', 'P4X-351 Icarus Base', 'Core overload explosion, evacuated', null, '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-26 15:40:00', null, '2024-04-01 08:00:00'),
('e0000000-0000-4000-8000-000000000042', 'P3X-111 Tok\'ra Temporary Camp', 'Moved after Goa\'uld patrol detection', 'Desert Oasis Ridge 2', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-28 16:50:00', null, '2024-04-05 12:15:00'),
('e0000000-0000-4000-8000-000000000043', 'P2X-338 Ziggurat Vault', 'Ancient tomb research camp, sealed', 'Babylonian Style Ruins, Sector 1', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-30 17:00:00', null, '2024-04-10 16:00:00'),
('e0000000-0000-4000-8000-000000000044', 'Trust Safehouse - Seattle', 'Rogue NID faction hideout destroyed by C4', '1200 5th Ave, Seattle, WA 98101, USA', 'c0f8c245-1b3d-4d5f-9234-8c7d6e5f4a3b', '2024-02-01 18:20:00', null, '2024-04-15 11:30:00'),
('e0000000-0000-4000-8000-000000000045', 'P3X-666 Medical Field Hospital', 'In honor of Janet Fraiser, location relocated', 'Valley Outpost B', '7c8f9a0b-1234-4567-8901-abcdef123456', '2024-02-03 09:10:00', null, '2024-04-20 14:00:00'),
('e0000000-0000-4000-8000-000000000046', 'P3X-403 Mine Headquarters', 'Unas treaty area, original camp closed', 'River Bed Valley Camp', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', '2024-02-05 10:25:00', null, '2024-04-25 10:00:00'),
('e0000000-0000-4000-8000-000000000047', 'Vala\'s Black Market Smuggling Ship', 'Cargo vessel seized and scrapped', 'Docking Slip 4, Orban Trading Station', '8d90ab1c-2345-5678-9012-bcdefa234567', '2024-02-07 11:40:00', null, '2024-05-01 09:30:00'),
('e0000000-0000-4000-8000-000000000048', 'P4S-237 Lord Mot\'s Fortress', 'Goa\'uld stronghold liberated and dismantled', 'Fortress Citadel Peak', '68b329da-9893-4d34-9d6b-549302554020', '2024-02-09 13:00:00', null, '2024-05-05 15:45:00'),
('e0000000-0000-4000-8000-000000000049', 'P3X-774 Nox Sanctuary Lodge', 'Hidden enclave entry point, sealed by Nox', null, '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-11 14:15:00', null, '2024-05-10 18:00:00'),
('e0000000-0000-4000-8000-000000000050', 'P3X-775 Talthus Colony Ship Landing Site', 'Temporary settlement site, colony relocated', 'Coastal Plains Landing Strip', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', '2024-02-13 15:30:00', null, '2024-05-15 12:00:00');

--
-- Add books_locations
--
INSERT INTO binbogami.books_locations (book_id, location_id, created_by, created_at, deleted_at) VALUES
('7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b', 'e0000000-0000-4000-8000-000000000001', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-01 17:00:00', null),
('7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b', 'e0000000-0000-4000-8000-000000000002', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-01 17:05:00', null),
('7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b', 'e0000000-0000-4000-8000-000000000003', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2026-01-01 17:10:00', null);

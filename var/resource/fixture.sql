--
-- Remove everything
--
SET FOREIGN_KEY_CHECKS=0;

TRUNCATE binbogami.members;
TRUNCATE binbogami.entries;
TRUNCATE binbogami.books_categories;
TRUNCATE binbogami.books_locations;
TRUNCATE binbogami.categories;
TRUNCATE binbogami.locations;
TRUNCATE binbogami.books;
TRUNCATE binbogami.organizations;
TRUNCATE binbogami.invitations;
TRUNCATE binbogami.users;
TRUNCATE binbogami.user_configurations;

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

--
-- Add organizations
--
INSERT INTO binbogami.organizations (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('1a530066-b21b-11ee-9a7a-5ab75f0c1cab', 'SGC', 'Stargate Command', '03a5ac74-b21b-11ee-9a7a-5ab75f0c1cab', '2024-01-13 13:53:17', null, null);
INSERT INTO binbogami.organizations (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('1d8408ef-8f56-4b95-83c4-033e9e343ead', 'SG-1', 'SG-1 Unit', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', '2024-01-13 13:53:17', null, null);
INSERT INTO binbogami.organizations (id, name, description, created_by, created_at, updated_at, deleted_at) VALUES ('7b70bf48-c6cd-4aeb-8916-4ce73e8c7511', 'System Lords', 'Goa\'uld System Lords', '89705c9b-fe47-4c16-93f7-89cce8807f02', '2024-01-13 13:53:17', null, null);

--
-- Assign memberships
--
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 4, '1a530066-b21b-11ee-9a7a-5ab75f0c1cab', '03a5ac74-b21b-11ee-9a7a-5ab75f0c1cab', '03a5ac74-b21b-11ee-9a7a-5ab75f0c1cab', '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 3, '1a530066-b21b-11ee-9a7a-5ab75f0c1cab', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1a530066-b21b-11ee-9a7a-5ab75f0c1cab', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1a530066-b21b-11ee-9a7a-5ab75f0c1cab', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1a530066-b21b-11ee-9a7a-5ab75f0c1cab', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 4, '1d8408ef-8f56-4b95-83c4-033e9e343ead', '05e7257a-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1d8408ef-8f56-4b95-83c4-033e9e343ead', '1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1d8408ef-8f56-4b95-83c4-033e9e343ead', '2b63b228-b21c-11ee-9a7a-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 1, '1d8408ef-8f56-4b95-83c4-033e9e343ead', 'aff84550-b21f-11ee-8ac0-5ab75f0c1cab', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 4, '7b70bf48-c6cd-4aeb-8916-4ce73e8c7511', '89705c9b-fe47-4c16-93f7-89cce8807f02', null, '2024-03-31 17:42:43', null, null);
INSERT INTO binbogami.members (id, role, organization_id, user_id, created_by, created_at, updated_at, deleted_at) VALUES (NULL, 3, '7b70bf48-c6cd-4aeb-8916-4ce73e8c7511', '5f4dcc3b-5aa5-487d-8f6a-64cbd3665b4e', null, '2024-03-31 17:42:43', null, null);



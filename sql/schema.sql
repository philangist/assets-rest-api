BEGIN;
CREATE TABLE IF NOT EXISTS projects (
       id serial NOT NULL,
       name varchar(128) NOT NULL,
       created_at timestamp NOT NULL DEFAULT NOW(),
       CONSTRAINT projects_pkey PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS assets (
       id serial NOT NULL,
       name varchar(128) NOT NULL,
       parent_id INT REFERENCES assets ON DELETE CASCADE,
       media_url varchar(64),
       category int NOT NULL,
       project_id int NOT NULL REFERENCES projects ON DELETE CASCADE,
       created_at timestamp NOT NULL DEFAULT NOW(),
       CONSTRAINT assets_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_assets_category_project_id ON assets (category, project_id);
CREATE INDEX idx_assets_category_parent_id ON assets (category, parent_id);

INSERT INTO projects(id, name) VALUES
       (1, 'Project 1'),
       (2, 'Project 2'),
       (3, 'Project 3');

INSERT INTO assets(id, name, parent_id, media_url, category, project_id) VALUES
       (1,  'Project 1 Folder 1', NULL, NULL, 1, 1),
       (2,  'P1F1/File 1', 1, 'long_sha256_hash_here', 2, 1),
       (3,  'P1F1/File 2', 1, 'long_sha256_hash_here', 2, 1),
       (4,  'P1F1/Subfolder 1',  1, NULL, 1, 1),
       (5,  'P1F1/S1/File 0', 4, 'long_sha256_hash_here', 2, 1),
       (6,  'Project 2 Folder 1', NULL, NULL, 1, 2),
       (7,  'P2F1/File 1',  6, 'long_sha256_hash_here', 2, 2),
       (8,  'P2F1/File 2',  6, 'long_sha256_hash_here', 2, 2),
       (9,  'P2F1/File 3',  6, 'long_sha256_hash_here', 2, 2),
       (10, 'P2F1/File 4',  6, 'long_sha256_hash_here', 2, 2),
       (11, 'P2F1/File 5',  6, 'long_sha256_hash_here', 2, 2),
       (12, 'P2F1/File 6',  6, 'long_sha256_hash_here', 2, 2),
       (13, 'P2F1/File 7',  6, 'long_sha256_hash_here', 2, 2),
       (14, 'P2F1/File 8',  6, 'long_sha256_hash_here', 2, 2),
       (15, 'P2F1/File 9',  6, 'long_sha256_hash_here', 2, 2),
       (16, 'P2F1/File 10', 6, 'long_sha256_hash_here', 2, 2),
       (17, 'P2F1/File 11', 6, 'long_sha256_hash_here', 2, 2),
       (18, 'P2F1/File 12', 6, 'long_sha256_hash_here', 2, 2),
       (19, 'P2F1/File 13', 6, 'long_sha256_hash_here', 2, 2),
       (20, 'P2F1/File 14', 6, 'long_sha256_hash_here', 2, 2),
       (21, 'P2F1/File 15', 6, 'long_sha256_hash_here', 2, 2),
       (22, 'P1F1/S1/File 1', 4, 'long_sha256_hash_here', 2, 1),
       (23, 'P1F1/S1/File 2', 4, 'long_sha256_hash_here', 2, 1),
       (24, 'P1F1/S1/File 3', 4, 'long_sha256_hash_here', 2, 1),
       (25, 'P1F1/S1/File 4', 4, 'long_sha256_hash_here', 2, 1),
       (26, 'P1F1/S1/File 5', 4, 'long_sha256_hash_here', 2, 1),
       (27, 'P1F1/S1/Subfolder 2', 4, 'long_sha256_hash_here', 1, 1),
       (28, 'P1F1/S1/S2/File 0', 27, 'long_sha256_hash_here', 2, 1),
       (29, 'P1F1/S1/S2/File 1', 27, 'long_sha256_hash_here', 2, 1);
COMMIT;

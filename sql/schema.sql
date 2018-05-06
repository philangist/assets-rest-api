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
       (1, 'Project 1 Folder 1', NULL, NULL, 1, 1),
       (2, 'P1F1/File 1', 1, 'long_sha256_hash_here', 2, 1),
       (3, 'P1F1/File 2', 1, 'long_sha256_hash_here', 2, 1),
       (4, 'Project 1 Folder 1/Subfolder 1',  1, NULL, 1, 1),
       (5, 'P1F1/F1/S1/File 1', 5, 'long_sha256_hash_here', 2, 1),
       (6, 'Project 2 Folder 1', NULL, NULL, 1, 2);
COMMIT;

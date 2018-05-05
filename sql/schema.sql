BEGIN;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS projects;
CREATE TABLE projects (
       id serial NOT NULL,
       name varchar(128) NOT NULL,
       created_at timestamp NOT NULL DEFAULT NOW(),
       CONSTRAINT projects_pkey PRIMARY KEY (id)
);
CREATE TABLE assets (
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
       (1, 'Avengers'),
       (2, 'Avengers: Age of Ultron'),
       (3, 'Avengers: Infinity War');

INSERT INTO assets(id, name, parent_id, media_url, category, project_id) VALUES
       (1, 'Tony Stark Bloopers', NULL, NULL, 1, 1),
       (2, 'Tony_falls.mkv', 1, 'long_sha256_hash_here', 2, 1),
       (3, 'Loki_Tony.mkv',  1, 'long_sha256_hash_here', 2, 1),
       (4, 'Johannesburg Footage', NULL, NULL, 1, 2);
COMMIT;

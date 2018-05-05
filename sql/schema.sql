BEGIN;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS projects;
CREATE TABLE projects (
       id serial NOT NULL,
       name varchar(128) NOT NULL,
       created_at timestamp NOT NULL,
       CONSTRAINT projects_pkey PRIMARY KEY (id)
);
CREATE TABLE assets (
       id serial NOT NULL,
       name varchar(128) NOT NULL,
       parent_id INT REFERENCES assets ON DELETE CASCADE,
       media_url varchar(64),
       category int NOT NULL,
       project_id int NOT NULL REFERENCES projects ON DELETE CASCADE,
       created_at timestamp NOT NULL,
       CONSTRAINT assets_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_assets_category_project_id ON assets (category, project_id);
CREATE INDEX idx_assets_category_parent_id ON assets (category, parent_id);
COMMIT;

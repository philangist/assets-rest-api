API

Introduction: A suite of resource representations and their corresponding
endpoints for traversing the frame.io asset storage backend

Assumptions: There is an invariant we assume holds here, that every project
             has exactly one root folder associated with it. This should be
             enforced in the api layer code through the POST, PUT, and DELETE
             verbs for both projects/ and assets/, and also possibly as a DB
             level constraint on the assets table.
             We also assume that all assets with type != 1 are leaves in the
             project tree.

Meta/Pagination: TODO: FILL THIS WHEN I DECIDE IF/HOW TO IMPLEMENT PAGINATION

Projects

- Projects are the high-level structures users to store their video and other
  media assets. All of their content is organized within a "project", and within
  each project files can be infinitely nested within folders.

Endpoints:
GET projects/ - list all projects in the database

response:
status code: 200 OK
```json
{
    data: [
        {
            id: int,
            name: string (128),
            root_folder_id: int, // references assets with type=folder and null
            // parent invariant is enforced on PUT and POST
            created_at: timestamp, // project creation date
        },
        ..
    ]
    meta: {
        total: 10
    }
}
```
supported query parameters:
type        int - filter assets by type identifier
project\_id int - filter assets by project identifier
parent\_id  int - filter assets by parent_identifier
expand     bool - returns all children of any assets that are returned by
                  any combination of the above query parameters

GET projects/:id - retrieve a specific project by id

response:
status code: 200 OK
```json
{
    id: int,
    name: string (128),
    root_folder_id: int,
    created_at: timestamp,
}
```

Assets

- Assets are the mechanism used to both refer to the media objects we're
  actually storing as well as express their heirarchical relationship with
  each other. Take note that for every project there must exist exactly one
  asset with type=folder, parent_id=null, and project_id=project.id.
  Also take note that it is a logical contradiction for an asset to have a
  parent asset of type > 1.

Endpoints:
GET assets/ - list all assets in the database

response:
status code: 200 OK
```json
{
    data: [
        {
            id: int
            name: string (128)
            parent_id: int // references a parent asset with type=folder.
            // nullable only for type=folder objects
            media_url: string (variable length, typically between 84
            - 100 characters),  // physical location of the
            // media object associated
            // with this asset. format is a http url of type
            // "http://<env>.frame.io/asset_sha256_hash". possible values of
            // env are dev, qa-<cluster_id> and cdn
            type: int media object type // current possible values are int(1)
            // for folders and (2) for video files
            project_id: int // references the encapsulating project which
            // this asset belongs to. there must exist exactly one asset of
            // type=folder and parent_id=null for every project in the database
            created_at: timestamp // asset creation date
        },
        ..
    ]
    meta: {
        total: 10
    }
}
```

GET assets/:id retrieve a specific asset by id

response:
status code: 200 OK
```json
{
    id: int
    name: string (128)
    parent_id: int
    media_url: string (100)
    type: int media object type
    project_id: int
    created_at: timestamp
}
```

Asset Picker API
====================
### Sections:
- Introduction  
- Projects  
- Assets  
- Running a local copy  
- Improvements  

### Introduction  
A suite of endpoints for traversing the frame.io asset storage backend

__Assumptions:__  
We assume that every project has exactly one root folder associated with it. This should be enforced in the api layer code through the POST, PUT, and DELETE verbs for both projects/ and assets/, and also with DB level constraints on the assets table.  

We also assume that all assets with type == 1 are folders and type != 1 are leaves/files in the project tree.  

__Pagination:__  
Standard `offset` & `limit`  query parameters are supported for paginating through large amounts of items.  

### Projects  

Projects are the namespaces users store their video and other media assets within. Within each project files can be infinitely nested within folders.  

__Endpoints:__  
`GET /projects/` - list all projects in the database  

status code: 200 OK  
```javascript
{
    data: [
        {
            id: int,
            name: string (128),
            root_folder_id: int,   // References an asset with type=folder, parent_id=NULL and project_id = project.id. Null parent invariant is enforced on PUT and POST
            created_at: timestamp, // Project creation date
        },
        ..
    ]
    page: {
        total:  100,
        limit:  10,
        offset: 0
    }
}
```


`GET /projects/:id` - retrieve a specific project by id  

status code: 200 OK  
```javascript
{
    id: int,
    name: string (128),
    root_folder_id: int,
    created_at: timestamp,
}
```

### Assets  

Assets are the mechanism used to both refer to the media objects we're actually storing as well as express their heirarchical relationship with each other. Take note that for every project there must exist exactly one asset with `type='folder'`, `parent_id=null`, and `project_id=project.id`. Also take note that it is a logical contradiction for an asset to have a parent asset of type > 1.  

__Endpoints:__  
`GET /assets/` - list all assets in the database  

status code: 200 OK  
```javascript
{
    data: [
        {
            id: int,
            name: string (128),
            parent_id: int|null, // References a parent asset with type=folder. Nullable only for type=folder objects
            media_url: string|null (variable length, typically between 84 - 100 characters),  // physical location of the  media object associated with this asset. format is a http url of type "http://<env>.frame.io/asset_sha256_hash". possible values of env are dev, qa-<cluster_id> and cdn
            type: int,  //Media object type. Current possible values are int(1) for folders and (2) for video files
            project_id: int, // References the encapsulating project which this asset belongs to. there must exist exactly one asset of type=folder and parent_id=null for every project in the database
            created_at: timestamp // Asset creation date
        },
        ..
    ]
    page: {
        total:  100,
        limit:  10,
        offset: 0
    }
}
```

supported query parameters:  
`type`        int - query assets by type identifier  
`project\_id` int - query assets by project identifier  
`parent\_id`  int - query assets by parent identifier  
`descendants` bool - returns all immediate children (1 level down the tree) of assets that are returned by  
                  any combination of the above query parameters  

`GET assets/:id` - retrieve a specific asset by id  

status code: 200 OK  
```javascript
{
    id: int
    name: string (128),
    parent_id: int,
    media_url: string (100),
    type: int,
    project_id: int,
    created_at: timestamp
}
```

### Running a local copy  

This repo requires docker to run locally. All other golang dependencies are bundled in at vendors/.
Run the following in a shell to spin up the server at localhost:8080  
```bash
$ export POSTGRES_DB=db
$ export POSTGRES_USER=postgres-dev
$ export POSTGRES_PASSWORD=Z3R0C00L
$ export POSTGRES_HOSTNAME=db
$ export POSTGRES_PORT=5432
$ docker-compose up
``` 

###  Improvements  

1. Tests - The most obvious improvement to be made here is to add comprehensive unit tests, I started down that road early on but [decided against it](https://github.com/philangist/frameio-assets/commit/bfea26ffcbc01ec71574a02821b0f90aa07e78ac#diff-b84c7556427bbbc195ca3c5d3bd5bee3) since that effort would interfere with building out the core logic.  

2. models/ is playing two roles as both an implementation of the data access layer and the api representation engine. This confusion violates the single responsibility principle and can lead to unnecessary complexity when trying to hide mismatches between the database and server api.  

3. The models/ query execution flow and object lifecycle is too rigid, and can only perform one sql query per request cycle which is really limiting  

4. The http handlers in controllers/ are all very similar to each other and their commonalities can be extracted out in a reusable way to maximize DRYness.  


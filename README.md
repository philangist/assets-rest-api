### Asset Picker API

__Introduction:__  
A suite of endpoints for traversing the frame.io asset storage backend

Assumptions: We assume that every project has exactly one root folder associated with it. This should be enforced in the api layer code through the POST, PUT, and DELETE verbs for both projects/ and assets/, and also with DB level constraints on the assets table.  

We also assume that all assets with type == 1 are folders and type != 1 are leaves/files in the project tree.  

__Pagination:__  
Standard `offset` & `limit`  query parameters are supported for paginating through large amounts of items.  

### Projects:  

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
            parent_id: int, // References a parent asset with type=folder. Nullable only for type=folder objects
            media_url: string (variable length, typically between 84 - 100 characters),  // physical location of the  media object associated with this asset. format is a http url of type "http://<env>.frame.io/asset_sha256_hash". possible values of env are dev, qa-<cluster_id> and cdn
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
`expand`      bool - returns all immediate children (1 level down the tree) of assets that are returned by  
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
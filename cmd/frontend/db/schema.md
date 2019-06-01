# Table "public.repo"
```
         Column          |           Type           |                     Modifiers                     
-------------------------+--------------------------+---------------------------------------------------
 id                      | integer                  | not null default nextval('repo_id_seq'::regclass)
 uri                     | citext                   | 
 owner                   | citext                   | 
 name                    | citext                   | 
 description             | text                     | 
 vcs                     | text                     | not null
 http_clone_url          | text                     | 
 ssh_clone_url           | text                     | 
 homepage_url            | text                     | 
 default_branch          | text                     | not null
 language                | text                     | 
 blocked                 | boolean                  | 
 deprecated              | boolean                  | 
 fork                    | boolean                  | 
 mirror                  | boolean                  | 
 private                 | boolean                  | 
 created_at              | timestamp with time zone | 
 updated_at              | timestamp with time zone | 
 pushed_at               | timestamp with time zone | 
 vcs_synced_at           | timestamp with time zone | 
 indexed_revision        | text                     | 
 freeze_indexed_revision | boolean                  | 
 origin_repo_id          | text                     | 
 origin_service          | integer                  | 
 origin_api_base_url     | text                     | 
Indexes:
    "repo_pkey" PRIMARY KEY, btree (id)
    "repo_uri_unique" UNIQUE, btree (uri)
    "repo_name" btree (name text_pattern_ops)
    "repo_name_ci" btree (name)
    "repo_owner_ci" btree (owner)
    "repo_uri_trgm" gin (lower(uri::text) gin_trgm_ops)

```

# Table "public.schema_migrations"
```
 Column  |  Type   | Modifiers 
---------+---------+-----------
 version | bigint  | not null
 dirty   | boolean | not null
Indexes:
    "schema_migrations_pkey" PRIMARY KEY, btree (version)

```

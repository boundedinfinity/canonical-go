# SQL generation


# SQLite implementation

This file contains details for generation SQL associated with the https://www.sqlite.org/ database.

When generated code for use with the SQLite database create tables that following this format:

```sql
CREATE TABLE users (
    rowid INTEGER PRIMARY KEY,  -- Fast, ordered internal clustered index
    id TEXT UNIQUE NOT NULL,    -- Handed to clients/APIs
    ... other fields
);
```

For a model with a [Natural Key](natural-key.md), the `id` column holds the natural key's value instead of a generated ID.

We want to use the internal rowid for use in INSERTs, JOINs, etc.. and use the ID field for external clients.

Queries should be writting to select on the `id` and use the rowid for any require JOINs, or other non-specificed sorting queries. This file will be the string serialized view of the [ID](../id/id.md) filed from all **Top Level  Models** models. 

Any domain model with that is marked with a **Top Level Model** annotation should have it's own SQL DDL to create a separate table for this model.  Any models that are not marked as **Top Level Model** are used in join tables or other tables related to the top level models.
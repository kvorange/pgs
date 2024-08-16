## Models

In order to transform a structure into a model, it is necessary to use `pgs.Model`. 
Use the `table` tag to specify the name of the linked table.
**For example:**
```go
type User struct {
    pgs.Model `table:"user"` // Binding the structure to the user table
    
    Id       pgs.Field[pgtype.Int8]        `json:"id"`
    Login    pgs.Field[pgtype.Text]        `json:"login"`
    Name     pgs.Field[pgtype.Text]        `json:"name"`
    CreateAt pgs.Field[pgtype.Timestamptz] `json:"create_at"`
    IsAdmin  pgs.Field[pgtype.Bool]        `json:"is_admin"`
}
```

However, in order for the model to be able to query the database, it must be **initialized**.
It is assumed that you will do this exactly once for each model, and then use the prepared models in your program.
Initialization requires a `pgs.DbClient` instance (see [Database]((./docs/dialect.md)) for more information).
You can use the following code:
```go
var userModel User
err = userModel.Init(&dbClient, &userModel)
if err != nil {
    log.Fatalf("Failed to init model User: %v", err)
}
```

## Model fields
For the correct operation of the model, it is very important to initialize the fields correctly. 
The current version of the library allows working with regular fields and nested structures (foreign keys).

### Simple fields
To define a model field that is used within a table, use `pgs.Field[T]` where T is some `pgtype`.
The behavior of the field changes depending on its tags:
* `db:"name"` To explicitly specify the field name in the table, use the `db` tag. If this tag is not present, 
the field name will be converted to **snake_case by default**.
* `db:"-"` For the model, you can create empty fields that are not associated with the table 
(for internal functionality, or, for example, working with m2m, which is not yet available in this version). 
To do this, use `db` with `-` value.
* `json` used to specify output in JSON format.

**Example:**
```go
type User struct {
    pgs.Model `table:"user"` // Binding the structure to the user table

    Id      pgs.Field[pgtype.Int8] `json:"id"`                    // id in table
    Login   pgs.Field[pgtype.Text] `json:"login" db:"user_login"` // in table this field named user_login
    MyField int                    `db:"-"`                       // the field will be skipped in queries.
}
```

### Struct fields (Foreign keys)
In cases where you have a foreign key in your table, 
you can define a nested structure with `pgs.Model` field that will describe the related table.
To do this, you **must define** the following set of tags:
* `db:"table_name"` defines an associative name for binding values from a query to a structure (not a table name, but it is usually worth specifying it. This logic is based on the logic of the `pgxscan` library)
* `fk:"from,to"` defines the names of the table fields through which the binding occurs.

**Example:**
```go
type JobTitle struct {
    pgs.Model `table:"job_title"`
    
    Id      pgs.Field[pgtype.Int8] `json:"id"`
    Name    pgs.Field[pgtype.Text] `json:"name"`
}

type User struct {
    pgs.Model `table:"user"`

    Id       pgs.Field[pgtype.Int8]        `json:"id"`
    Login    pgs.Field[pgtype.Text]        `json:"login"`
    Name     pgs.Field[pgtype.Text]        `json:"name"`
    CreateAt pgs.Field[pgtype.Timestamptz] `json:"create_at"`
    IsAdmin  pgs.Field[pgtype.Bool]        `json:"is_admin"`
	
    JobTitle JobTitle `db:"job_title" fk:"job_title_id,id" json:"job_title"` //define fk table
}
```

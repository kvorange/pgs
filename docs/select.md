## Select

To perform a **select** operation, call the appropriate **Select** method of the model. 
This method can accept a number of Selectable interface parameters to limit the selection by fields.
The Selestable interface is implemented by all structures that are initialized as a model, fields, and some special expressions.
* If you use a structure as a parameter, then all its fields will be used in the selection, without taking into account nested structures.
* If you do not specify any parameters, then the full selection will be used.
* If you use a special expression, then the selection will use the result of this expression (for example, **Count** )

### **Examples (with these models):**
```go
type JobTitle struct {
    pgs.Model `table:"job_title"`
    
    Id   pgs.Field[pgtype.Int8] `json:"id"`
    Name pgs.Field[pgtype.Text] `json:"name"`
}

type User struct {
    pgs.Model `table:"user"`
    
    Id    pgs.Field[pgtype.Int8] `json:"id"`
    Login pgs.Field[pgtype.Text] `json:"login"`
    Name  pgs.Field[pgtype.Text] `json:"name"`
    
    JobTitle JobTitle `db:"job_title" fk:"job_title_id,id" json:"job_title"`
}
```
### Full select:
```go
query := user.Select().Query() // full select
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id", "user"."login" AS "login", "user"."name" AS "name", "user__job_title"."id" AS "job_title.id", "user__job_title"."name" AS "job_title.name" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id")
```

### Select with pgs.Field
```go
query := user.Select().Query() // full select
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id", "user"."login" AS "login" FROM "user"
```

### Select with model fields
```go
query := user.Select(&user).Query() // only user fields without job_title
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id", "user"."login" AS "login" FROM "user"
```

## Count
To calculate Count use the `pgs.Count(Field[T])`. 
This function will calculate the Count expression based on the field in the parameter.

### Example:

```go
query := user.Select(pgs.Count(&user.Id)).Query()
fmt.Println(query)
```

#### Output:
```
SELECT COUNT("user"."id") FROM "user"
```

## Limit, Offset

Methods `Limit(limit uint)`, `Offset(offset uint)` are defined for the dataset.

### Example:

```go
query := user.Select().Limit(10).Offset(5).Query()
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id", "user"."login" AS "login", "user"."name" AS "name", "user__job_title"."id" AS "job_title.id", "user__job_title"."name" AS "job_title.name" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id") LIMIT 10 OFFSET 5
```


## Scanning 
To scan the query result, use the methods `Scan` and `ScanOne`. Typically, you want to scan the results into a model, so a typical example would look like this:
```go
var users []User
err := user.Select().Scan(&users) // scanning all data to user slice
// handle err
```

However, you can define your own structures. To do this, you need to define them in such a way that the names
of the structure fields correspond to the results of your selection. This behavior is inherited from the library `pgxscan`.
By default, structure fields are converted to snake_case, but you can use `db` tag to specify name explicitly.
Additionally, for `pgs.Field` you can use the `As` function to specify an explicit name.

### Example:
```go
query := user.Select(&user.Id, user.JobTitle.Id.As("job_title_id")).Query()
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id", "user__job_title"."id" AS "job_title_id" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id")
```

```go
type myUser struct {
    Id         pgtype.Int8 `json:"id"`
    JobTitleId pgtype.Int8 `json:"job_title_id" db:"job_title_id"`
}
var myUsers []myUser
err := user.Select(&user.Id, user.JobTitle.Id.As("job_title_id")).Scan(&myUsers)
// handle err
```
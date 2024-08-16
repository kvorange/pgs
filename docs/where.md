# !ATTENTION!
**In the current version of the library, the problem of conditions when updating and deleting using fk is not solved in the best way. At the moment, any fields of nested models in such situations are interpreted by the `fk` tag from value and redefined, even if you do not specify a field of the primary key. Be careful with this functionality**

## Where

For datasets select, update and delete, the functionality is defined where. 
You can build your conditions using special methods from `pgs.Field`.
The following methods are currently defined:
* `Eq(value interface{})`
* `NotEq(value interface{})`
* `Like(value interface{})`
* `NotLike(value interface{})`
* `Regex(value interface{})`
* `NotRegex(value interface{})`
* `NotRegexI(value interface{})`
* `Lt(value interface{})`
* `Lte(value interface{})`
* `Gt(value interface{})`
* `Gte(value interface{})`
* `IsNotNull()`
* `IsNull()`

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

```go
query := user.Select(&user.Id).Where(
    user.Id.Lt(10),
    user.JobTitle.Id.Eq(1),
).Query()
fmt.Println(query)
```

#### Output:
```
SELECT "user"."id" AS "id" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id") WHERE (("user"."id" < 10) AND ("user__job_title"."id" = 1))
```

### Or

You can use `pgs.Or()` function to define or condition.

#### Example:
```go
query := user.Select(&user.Id).Where(
    pgs.Or(
        user.Id.Lt(10),
        user.JobTitle.Id.Eq(1),
    ),
).Query()
fmt.Println(query)
```
#### Output:
```
SELECT "user"."id" AS "id" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id") WHERE (("user"."id" < 10) OR ("user__job_title"."id" = 1))
```

### Subquery
The selection dataset can be used as the condition value.

#### Example:
```go
var user User
    err = user.Init(&dbClient, &user)
    if err != nil {
    panic(err)
}

var jobTitle JobTitle
    err = jobTitle.Init(&dbClient, &jobTitle)
    if err != nil {
    panic(err)
}

query := user.Select(&user.Id).Where(
        user.JobTitle.Id.In(
            jobTitle.Select(&jobTitle.Id).Where(jobTitle.Id.Lt(10)),
        ),
    ).Query()
fmt.Println(query)
```
#### Output:
```
SELECT "user"."id" AS "id" FROM "user" LEFT JOIN "job_title" AS "user__job_title" ON ("user"."job_title_id" = "user__job_title"."id") WHERE (("user"."id" < 10) OR ("user__job_title"."id" = 1))
```

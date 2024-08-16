## Update

To perform **update** operation, call the appropriate **Update** method of the model.
To determine which fields need to be updated, use `pgs.Record{}`.

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
query := user.Update(pgs.Record{&user.Name: "new_name"}).Where(user.Id.Eq(1)).Query()
fmt.Println(query)
```

#### Output:
```
UPDATE "user" SET "name"='new_name' WHERE ("user"."id" = 1)
```

## Returning

To get the returning values, use the Returning() method and pass some number of fields. 
Then use the methods `Scan()` or `ScanOne()`.

### Example:
```go
query := user.Update(pgs.Record{&user.Name: "new_name"}).Where(user.Id.Eq(1)).Returning(&user.Id).Query()
fmt.Println(query)
```

#### Output:
```
UPDATE "user" SET "name"='new_name' WHERE ("user"."id" = 1) RETURNING "user"."id"
```
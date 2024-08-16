## Update

To perform **insert** operation, call the appropriate **Insert** method of the model.
To determine which fields need to be inserted, use `pgs.Record{}`.

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
query := user.Insert(pgs.Record{
    &user.Login:       "login",
    &user.Name:        "name",
    &user.JobTitle.Id: 1,
}).Query()
fmt.Println(query)
```

#### Output:
```
INSERT INTO "user" ("job_title_id", "login", "name") VALUES (1, 'login', 'name')
```
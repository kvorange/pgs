## Update

To perform **delete** operation, call the appropriate **Delete** method of the model.

### **Examples:**
```go
type User struct {
    pgs.Model `table:"user"`
    
    Id    pgs.Field[pgtype.Int8] `json:"id"`
    Login pgs.Field[pgtype.Text] `json:"login"`
    Name  pgs.Field[pgtype.Text] `json:"name"`
    
    JobTitle JobTitle `db:"job_title" fk:"job_title_id,id" json:"job_title"`
}
```

```go
query := user.Delete().Where(user.Id.Eq(1)).Query()
fmt.Println(query)
```

#### Output:
```
DELETE FROM "user" WHERE ("user"."id" = 1)
```
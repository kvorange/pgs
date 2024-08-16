## Database connection

In order to work with models, you first need to initialize the connection to the PostgreSQL database. To do this, use the following code:

```go

package main

import (
	"context"
	"fmt"
	"github.com/kvorange/pgs"
)


func main() {
	// Db config settings
	dbConfig := pgs.DbConfig{
		Host:      "localhost",
		Port:      5432,
		User:      "user",
		Password:  "password",
		Name:      "dbname",
		PollCount: 10,
	}

	// Init db connection
	dbClient := pgs.DbClient{}
	err := dbClient.Connect(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Use connection in your program
}
```

Additionally, you can use this connection as a normal connection in cases where you need to execute your pure sql queries.
To do this, use the fields `Ctx` and `Pool` and the [`pgxscan`](https://github.com/georgysavva/scany) library functions to perform queries.
```go
var result []interface{}
query := ```SELECT * FROM "some_table"```
err := pgxscan.Select(dbClient.Ctx, dbClient.Pool, &result, query)
// handle error
```
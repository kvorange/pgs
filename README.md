# PGS - PostgreSQL Go Library

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

`pgs` is a powerful and easy-to-use library for interacting with PostgreSQL databases in Go,
built on top of the robust foundations of `goqu`, `pgx`, and `pgtype`.
It provides a high-level abstraction for building SQL queries, handling transactions,
and managing database connections with ease.
## Features

- **Simple API**: Intuitive and easy-to-use API for common database operations.
- **Transaction Support**: Built-in support for managing database transactions.
- **Query Building**: Flexible query building with support for complex queries and subqueries.
- **Connection Pooling**: Efficient connection pooling using `pgxpool`.
- **Type Safety**: Strong type safety for database interactions.

## Installation

To install PGS, use the following command:

```sh
go get github.com/kvorange/pgs
```

## Quick start
Quick Start
Here's a quick example to get you started with `pgs`:

```go
package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kvorange/pgs"
	"log"
)

// User simple table for example
type User struct {
	pgs.Model `table:"user"`

	Id       pgs.Field[pgtype.Int8]        `json:"id"`
	Login    pgs.Field[pgtype.Text]        `json:"login"`
	Name     pgs.Field[pgtype.Text]        `json:"name"`
	CreateAt pgs.Field[pgtype.Timestamptz] `json:"iin"`
	IsAdmin  pgs.Field[pgtype.Bool]        `json:"is_fired"`
}

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
    
	// Init user model. Use this model for all you operations with user table
	var userModel User
	err = userModel.Init(&dbClient, &userModel)
	if err != nil {
		log.Fatalf("Failed to init model User: %v", err)
	}

	// Select all users for example
	var users []User
	err = userModel.Select().Scan(&users)
	if err != nil {
		log.Fatalf("Failed to select users: %v", err)
	}

	// Your database operations here
}
```

## Acknowledgments
Thanks to the Go community for their support and contributions.

Inspired by the need for a simple and powerful PostgreSQL library in Go, leveraging the strengths of `goqu`, `pgx`, and `pgtype`.
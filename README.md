# go-schemaless

A basic implementation of schemaless data on mysql inspired by [Uber](https://eng.uber.com/schemaless-part-one/)

## Installing
```bash
go get github.com/thelenilson/go-schemaless
```

## Usage

```go
package main

import (
	"log"

	"github.com/thelenilson/go-schemaless"
)

// you can define your own data models
type Person struct {
	FirstName string `msgpack:"first_name"` // changing the field name stored in database by msgpack
	LastName  string
}

func main() {

	// Define a connection object
	conn := schemaless.NewConnection(&schemaless.ConnectionOptions{
		Host:       "localhost",
		Port:       "3306",
		Username:   "root",
		Password:   "mysecretpass",
		Parameters: map[string]string{"parseTime": "true"}, // Custom database parameters
	})

	// Define a DataSource, aka: shard
	// the package will take care of creating the database and the table
	ds := schemaless.NewDataSource("people", conn)

	p := Person{
		FirstName: "John",
		LastName:  "Doe",
	}

	// Creating a new data cell
	cell, err := ds.CreateCell(
		RowID: "3c4eecc2-84df-466b-a023-9c044c289934",
		ColumnName: "PERSONAL_INFO",
		Body: &p,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(cell.AddedID)
}
```

## API's

### DataSource
CreateCell will create and return a new cell
```go
CreateCell(RowID string, ColumnName string, Body interface{}) (*DataCell, error)
```
FindCell will return a cell by given arguments
```go
FindCell(RowID string, ColumnName string, RefKey int) (*DataCell, error)
```
FindLatestCell will search for the latest RefKey of a cell an return it
```go
FindLatestCell(RowID string, ColumnName string) (*DataCell, error)
```
FindAllLatest will search for the latest RefKey of all cells by a given ColumnName
```go
FindAllLatest(ColumnName string) (*[]DataCell, error)
```

### DataCell
MarshalBody will receive any data type and convert to a msgpack binary, then save it to DataCell.Body
```go
MarshalBody(v interface{}) (err error)
```
UnmarshalBody will convert the msgpack binary in DataCell.Body back to a data type
```go
UnmarshalBody(v interface{}) (err error)
```

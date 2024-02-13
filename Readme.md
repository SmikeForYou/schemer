# schemer

schemer is a Go package designed to interact with databases using different dialects. 
It provides a unified interface to query metadata from databases and manipulate it in a structured way.


## Installation

To install schemer, you need to have Go installed on your machine. You can then use the `go get` command:

```bash
go get github.com/SmikeForYou/schemer
```

## Usage

The package provides different dialects to interact with databases. Each dialect implements the `Dialect` interface, which includes methods to query metadata from the database.

Here is an example of how to use the `PostgresDialect`:

```go
import (
    _ "github.com/lib/pq"
)

func main() {
    engine, err := schemer.NewEngine("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
        log.Fatal(err)
    }
	defer engine.Close()
	metadata = engine.GetMetadata()
	db_tree := metadata.BuildDBTree()
    // Do something with the metadata
}
```


## Contributing

Contributions are welcome. Please submit a pull request or create an issue to discuss the changes you want to make.

## License

This project is licensed under the MIT License.
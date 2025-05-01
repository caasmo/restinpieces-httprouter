# restinpieces-httprouter

This package provides an implementation of the `router.Router` interface for the [restinpieces](https://github.com/caasmo/restinpieces) framework, using the popular [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter).

## Usage

To use `httprouter` as the router in your `restinpieces` application, import the package and use the `httprouter.WithRouterHttprouter()` option when initializing the application with `restinpieces.New`.

```go
import (
	"github.com/caasmo/restinpieces"
	"github.com/caasmo/restinpieces-httprouter" // Import the httprouter package
	// ... other necessary imports
)

func main() {
	// ... setup database, age key, etc.

	_, srv, err := restinpieces.New(
		// ... other options like WithDbZombiezen, WithAgeKeyPath, WithCacheRistretto, WithTextLogger
		httprouter.WithRouterHttprouter(), // Add this option to use httprouter
	)
	if err != nil {
		// handle error
	}

	srv.Run()
}

```

## Compatibility

This router implements the `router.Router` interface defined by `restinpieces`.

Route patterns must follow the standard `http.ServeMux` format: `"METHOD /path"`. For example, `"GET /users"` or `"POST /items/:id"`. If the method is omitted, `GET` is assumed.

## Example

An example application demonstrating the usage of this router can be found in the `cmd/example` directory. You can run it using:

```bash
go run ./cmd/example -dbpath <path-to-db> -age-key <path-to-key>
```

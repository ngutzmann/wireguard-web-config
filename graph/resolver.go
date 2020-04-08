package graph

import (
	"database/sql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver - contains all the shared objects in the app
type Resolver struct {
	DB *sql.DB
}

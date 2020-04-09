package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	// This is required for the correct DB drivers
	_ "github.com/lib/pq"

	"github.com/ngutzmann/wireguard-web-config/graph"
	"github.com/ngutzmann/wireguard-web-config/graph/generated"
)

func createDBConnection() (*sql.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Successfully created db connection")
	return db, nil
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func cORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Defining the Graphql handler
func graphqlHandler(db *sql.DB) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func initDB() *sql.DB {
	var createPeersTable string = `CREATE TABLE IF NOT EXISTS peers (
		id UUID PRIMARY KEY,
		name VARCHAR(256) NOT NULL,
		public_key VARCHAR(256) NOT NULL,
		allowed_ip inet NOT NULL,
		created_on TIMESTAMP NOT NULL default current_timestamp
	)`

	tables := []string{createPeersTable}

	db, err := createDBConnection()
	if err != nil {
		log.Fatalln("Could not create DB connection", err)
	}

	for _, create := range tables {
		_, err = db.Exec(create)

		if err != nil {
			log.Fatalln("Could not create table:", err)
		}
	}
	return db
}

// Server - the root server for the Wireguard Web Config web server
func Server() {
	db := initDB()

	if mode := os.Getenv("GIN_MODE"); mode == RELEASE {
		gin.SetMode(mode)
	}
	r := gin.Default()
	r.Use(cORSMiddleware())
	r.GET("/ping", health)
	r.POST("/query", graphqlHandler(db))
	r.GET("/playground", playgroundHandler())
	r.Run()
}

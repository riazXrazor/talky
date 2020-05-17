package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iamsayantan/talky"
	"github.com/iamsayantan/talky/server"
	"github.com/iamsayantan/talky/store/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	defaultDBHost     = getFromEnv("MYSQL_HOST", "localhost")
	defaultDBPort     = getFromEnv("MYSQL_PORT", "3306")
	defaultDBUsername = getFromEnv("MYSQL_USERNAME", "root")
	defaultDBPassword = getFromEnv("MYSQL_PASSWORD", "12345")
	defaultDBName     = getFromEnv("DATABASE_NAME", "talky")

	defaultServerPort = "9050"
)

func main() {
	dbHost := flag.String("db.host", defaultDBHost, "Database host url")
	dbPort := flag.String("db.port", defaultDBPort, "Database port")
	dbUsername := flag.String("db.username", defaultDBUsername, "Database username")
	dbPassword := flag.String("db.password", defaultDBPassword, "Database password")
	serverPort := flag.String("server.port", defaultServerPort, "Server port where the server runs")

	flag.Parse()

	// connect to the database
	// format: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	dbCred := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", *dbUsername, *dbPassword, *dbHost, *dbPort, defaultDBName)
	log.Printf("Database Credential: %s", dbCred)

	db, err := gorm.Open("mysql", dbCred)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.AutoMigrate(talky.User{})

	userRepo := mysql.NewUserRepository(db)
	srv := server.NewServer(userRepo)

	log.Printf("Server starting on port %s", *serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), srv))
}

func getFromEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}
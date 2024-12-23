package pkg

    import (
       "database/sql"
       "fmt"
       "log"
       "os"

       "github.com/joho/godotenv"
       _ "github.com/lib/pq"
    )

    // DB  is a database connection variable.
    var DB *sql.DB

    // InitDB establishes a connection to the database
    func InitDB() {
       err := godotenv.Load()
       if err != nil {
          log.Fatal("Error loading .env file")
       }
       connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
          os.Getenv("DB_HOST"),
          os.Getenv("DB_PORT"),
          os.Getenv("DB_USER"),
          os.Getenv("DB_PASSWORD"),
          os.Getenv("DB_NAME"))

       db, err := sql.Open("postgres", connStr)

       if err != nil {
          log.Println("Failed to open DB Connection")
          log.Fatal(err)
       }

       err = db.Ping()
       if err != nil {
          log.Fatal(err)
       }

       DB = db
       log.Println("Successfully connected to the database!")

    }

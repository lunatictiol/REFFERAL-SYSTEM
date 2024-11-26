package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"referalsystem/internal/types"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	FindUserByEmail(email string) (*types.User, error)
	FindUserByID(id string) (*types.User, error)
	CreateUser(user types.User) (string, error)
	EditPoints(id string, points int) error
	CheckReferralCodeExists(code string) (bool, error)
	CheckReferralData(code string) (types.ReferalData, error)
	InsertReferralCode(code string, userId string) error
	ChangeReferalCodeUseStatus(code string) error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) CheckReferralCodeExists(code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "referals" WHERE "referalcode" = $1)`
	err := s.db.QueryRow(query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}
	return exists, nil
}
func (s *service) CheckReferralData(code string) (types.ReferalData, error) {
	var referal types.ReferalData
	fmt.Println(code)
	query := `SELECT * FROM referals WHERE referalcode = $1`
	err := s.db.QueryRow(query, code).Scan(
		&referal.Id,
		&referal.ReferalCode,
		&referal.ReferedBy,
		&referal.IsUsed,
	)
	if err != nil {
		return types.ReferalData{}, fmt.Errorf("database query failed: %w", err)
	}

	return referal, nil
}
func (s *service) InsertReferralCode(code string, userId string) error {

	_, err := s.db.Exec(`INSERT INTO referals (referalcode, referedby) VALUES ($1, $2)`, code, userId)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}
	return nil
}

func (s *service) ChangeReferalCodeUseStatus(code string) error {

	_, err := s.db.Exec(`UPDATE referals SET isUsed = true WHERE referalcode = $1`, code)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}

func (s *service) FindUserByEmail(email string) (*types.User, error) {
	user := new(types.User)
	err := s.db.QueryRow("SELECT * FROM refusers WHERE email =$1 ", email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Points,
	)

	if err != nil {
		fmt.Println("here i fucked s")
		fmt.Println(err)
		return nil, err
	}
	if user.ID == "" {
		fmt.Printf("hereeeeeee")
		return nil, fmt.Errorf("user not found")
	}
	println(user.Email)
	return user, nil
}
func (s *service) FindUserByID(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM refusers WHERE id= $1", id)
	if err != nil {
		fmt.Println("i fuxcked up")
		return nil, err
	}
	us := new(types.User)
	for rows.Next() {
		us, err = scanUsersFromRows(rows)
		if err != nil {
			fmt.Println("i fuxcked up")
			return nil, err
		}

	}
	if us.ID == "" {
		return nil, fmt.Errorf("user not found")
	}
	return us, nil
}
func (s *service) CreateUser(user types.User) (string, error) {
	_, err := s.db.Exec("INSERT INTO refusers (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	u, err := s.FindUserByEmail(user.Email)
	if err != nil {
		fmt.Println("Heree toooo")
		return "", err
	}
	println(u.ID)
	return u.ID, nil
}

func (s *service) EditPoints(id string, points int) error {
	u, err := s.FindUserByID(id)
	if err != nil {
		println("hereeeeeeeee")
		return err
	}
	p := u.Points
	newPoints := p + points

	_, err = s.db.Exec("UPDATE refusers SET points = $1 WHERE id = $2 ", newPoints, id)
	if err != nil {
		println("hereeeeeeeeezzzzzzz")
		return err
	}
	return nil

}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func scanUsersFromRows(row *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Points,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

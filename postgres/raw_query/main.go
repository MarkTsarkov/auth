package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=users user=admin password=admin"
)
var (
	name = gofakeit.Name()
	email = gofakeit.Email()
	password = gofakeit.City()
	role = "user"
)

func main() {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil{
		log.Fatalf("failed to connect to db: #{err}")
	}
	defer con.Close(ctx)

	create, err := con.Exec(ctx, "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)", name, email, password, role)
	if err != nil {
		log.Fatalf("failed to insert user: #{err}")
	}
	log.Printf("inserted %d rows", create.RowsAffected())


	rows, err := con.Query(ctx, "SELECT id, name, email, password, role, created_at, updated_at FROM users")
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email, password, role string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan")
		}
		log.Printf("id: %v, \nname: %v, \nemail: %v, \npassword: %v, \nrole: %v, \ncreated at: %v, \nupdated at: %v\n\n", id, name, email, password, role, createdAt, updatedAt)
	}
	

	update, err := con.Exec(ctx, "UPDATE users SET name = 'ivan', updated_at = now() WHERE id = 3")
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
		
	}
	log.Printf("Updated %v", update)

	delete, err := con.Exec(ctx, "DELETE FROM users WHERE id = 5")
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
		
	}
	log.Printf("Deleted %v", delete)
}
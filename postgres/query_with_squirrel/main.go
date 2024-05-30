package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
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
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer con.Close(ctx)

	//create
	builderInsert := sq.Insert("users").
	PlaceholderFormat(sq.Dollar).
	Columns("name", "email", "password", "role", "created_at").
	Values(name, email, password, role, time.Now()).
	Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int 
	err = con.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	//read
	builderRead := sq.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
	From("users").
	PlaceholderFormat(sq.Dollar).
	OrderBy("id ASC").
	Limit(10)

	query, args, err = builderRead.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := con.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to read users: %v", err)
	}

	var id int
	var name, email, password, role string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next(){
		err = rows.Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}
		log.Printf(
			"id: %v, \nname: %v, \nemail: %v, \npassword: %v, \nrole: %v, \ncreated at: %v, \nupdated at: %v\n\n", id, name, email, password, role, createdAt, updatedAt)
	}

	//update
	builderUpdate := sq.Update("users").
	PlaceholderFormat(sq.Dollar).
	Set("name", "maria").
	Set("email", "maria@db.conn").
	Set("updated_at", time.Now()).
	Where("id = ?", 37)
	
	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := con.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update: %v", err)
	}
	log.Printf("updated %v rows", res.RowsAffected())

	//read updated
	builderReadUpdated := sq.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
	From("users").
	Where("id = ?", 37).
	PlaceholderFormat(sq.Dollar).
	OrderBy("id ASC").
	Limit(10)

	query, args, err = builderReadUpdated.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = con.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to scan user: %v", err)
	}
	log.Printf(
		"id: %v, \nname: %v, \nemail: %v, \npassword: %v, \nrole: %v, \ncreated at: %v, \nupdated at: %v\n\n", id, name, email, password, role, createdAt, updatedAt)

	//delete
	builderDelete := sq.Delete("users").
	Where("id < 0").
	PlaceholderFormat(sq.Dollar)

	query, args, err = builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err = con.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete: %v", err)
	}

	log.Printf("deleted %v rows", res.RowsAffected())
}
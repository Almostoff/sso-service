package storage

import (
	"AuthService/config"
	"fmt"

	"github.com/guregu/null"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Comment struct {
	Id              int64    `db:"id"`
	UserId          int      `db:"user_id"`
	CommentId       null.Int `db:"comment_id"`
	Content         string   `db:"content"`
	Level           int
	VoteCount       int     `db:"voteCount"`
	UpdatedAt       float64 `db:"updated_at"`
	UpdatedAtNormal string
}

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	//connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	//	c.Postgres.Host,
	//	c.Postgres.Port,
	//	c.Postgres.User,
	//	c.Postgres.Password,
	//	c.Postgres.DBName,
	//	c.Postgres.SSLMode)
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"localhost",
		"54123",
		"admin",
		"root",
		"tS",
		"disable")
	database, err := sqlx.Connect("pgx", connectionUrl)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(5)
	database.SetMaxIdleConns(5)

	return database, nil
}

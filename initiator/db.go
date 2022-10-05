package initiator

import (
	"2f-authorization/platform/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(url string, log logger.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	config.ConnConfig.Logger = log.Named("pgx")
	config.MaxConns = 1000 // Not tested yet
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Statment := Statment{
	// 	Effect:   "asjd",
	// 	Resource: "jghds",
	// 	Action:   "fjgshfj",
	// }
	// s, _ := Statment.Value()

	// c.CreateUser(context.Background(), db.CreateUserParams{
	// 	Name:        sql.NullString{String: "asd", Valid: true},
	// 	Description: sql.NullString{String: "kjda", Valid: true},
	// 	Statment: pgtype.JSON{
	// 		Bytes:  s,
	// 		Status: pgtype.Present,
	// 	},
	// })
	return conn
}

type Statment struct {
	Effect   string `json:"effect"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

func (a Statment) Value() ([]byte, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Statment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

package sink

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
)

type Sink struct {
	client *gorm.DB
}

func NewSinkClient(host, port, user, dbname, password string) (*Sink, error) {
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		port,
		user,
		password,
		sslMode)
	db, err := gorm.Open("postgres", connectionStr)

	return &Sink{client: db}, err
}

func (s *Sink) GetLastUpdatedTime() (string, error) {
	q := fmt.Sprintf("SELECT MAX(timestamp) FROM %s;", os.Getenv("LOGGING_TABLE"))

	clientQuery := s.configure(q)

	ctx, cancel := cancelIn(10)
	defer cancel()
	j, err := clientQuery.Read(ctx)

	if err != nil {
		return "", err
	}

	var defaultTime time.Time
	var updatedAt time.Time

	for {
		var row []bq.Value
		err := j.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		if len(row) > 0 {
			updatedAt = row[0].(time.Time)
			break
		}
	}

	if updatedAt == defaultTime {
		return "", nil
	}

	return updatedAt.Format(time.RFC3339), nil
}

func cancelIn(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(seconds)*time.Second))
}

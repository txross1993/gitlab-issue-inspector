package sink

import (
	"context"
	"fmt"
	"os"
	"time"

	bq "cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type Sink struct {
	client *bq.Client
}

func NewSinkClient() (*Sink, error) {
	ctx, cancel := cancelIn(10)
	client, err := bq.NewClient(ctx, os.Getenv("GCP_PROJECT"))

	if err != nil {
		cancel()
		return nil, err
	}

	return &Sink{
		client: client,
	}, nil
}

func (s *Sink) GetLastUpdatedTime() (string, error) {
	q := fmt.Sprintf("SELECT MAX(timestamp) FROM %s;", os.Getenv("LOGGING_TABLE"))

	clientQuery := s.configure(q)

	ctx, cancel := cancelIn(10)
	j, err := clientQuery.Read(ctx)

	if err != nil {
		cancel()
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

func (s *Sink) configure(query string) *bq.Query {
	q := s.client.Query(query)
	q.QueryConfig.DefaultProjectID = os.Getenv("GCP_PROJECT")
	q.QueryConfig.DefaultDatasetID = os.Getenv("DATASET_ID")

	return q
}

func (s *Sink) Save(table string) {}

func cancelIn(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(seconds)*time.Second))
}

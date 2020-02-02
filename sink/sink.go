package sink

import (
	"context"
	"fmt"
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

func (s *Sink) Initialize() error {
	if err := s.autoMigrate(); err != nil {
		return err
	}

	if err := s.setForeignKeys(); err != nil {
		return err
	}

	return nil

}

func (s *Sink) GetLastUpdatedTime() (string, error) {
	rows, err := db.Table("issues").Select("max(updated_at) as max")
	defer rows.Close()

	if err != nil {
		return "", err
	}

	var max time.Time
	for rows.Next() {
		rows.Scan(&max)
	}

	return max.Format(time.RFC3339)
}

func (s *Sink) Read()   {}
func (s *Sink) Update() {}
func (s *Sink) Create() {}

// Relation interface provides the abstraction of fetching a data model's foreign key relationships
type Relation interface {
	// GetForeignKeyMapping() returns a map of the table's key to the reference key of the foreign table
	// i.e. For Note table,
	// Notes(issue_id) is mapped to Issues(id)
	// issue_id:issues(id)
	GetForeignKeyMapping() map[string]string
}

func (s *Sink) autoMigrate() error {
	models := data.GetModels()

	for _, m := range models {
		errs := s.DB.AutoMigrate(m).GetErrors()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Error(err)
			}
			return fmt.Errorf("Unable to migrate schemas")
		}
	}

	return nil
}

func setForeignKeys(models []interface{}) error {
	for _, m := range models {
		if m.(Relation) {
			fkMaps := m.GetForeignKeyMapping()
			// Create FK relationships
			for fk, mappedReference := range fkMaps {
				if err := s.DB.Model(m).AddForeignKey(fk, mappedReference, "CASCADE", "CASCADE").Err; err != nil {
					return err
				}

			}
		}

	}

}

func cancelIn(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(seconds)*time.Second))
}

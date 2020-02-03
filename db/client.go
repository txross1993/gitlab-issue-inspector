package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	data "github.com/txross1993/gitlab-issue-inspector/data"
)

type DB struct {
	*gorm.DB
}

func NewDBClient(host, port, user, dbname, password, sslMode string) (*DB, error) {
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		port,
		user,
		dbname,
		password,
		sslMode)

	db, err := gorm.Open("postgres", connectionStr)

	return &DB{db}, err
}

func (d *DB) Initialize() error {
	if err := d.autoMigrate(); err != nil {
		return err
	}

	return nil

}

func (d *DB) GetLastUpdatedTime() (string, error) {
	//TODO: Query DB for max
	//rows := d.DB.Table("issues").Select("max(updated_at) as max")
	max := "2020-01-20T19:57:54.153Z"
	//max.Format(time.RFC3339), nil
	return max, nil
}

// Relation interface provides the abstraction of fetching a data model'd foreign key relationships
type Relation interface {
	// GetForeignKeyMapping() returns a map of the table'd key to the reference key of the foreign table
	// i.e. For Note table,
	// Notes(issue_id) is mapped to Issues(id)
	// issue_id:issues(id)
	GetForeignKeyMapping() map[string]string
}

type Model interface {
	ToString() string
}

func (d *DB) autoMigrate() error {
	models := data.GetModels()

	for _, m := range models {
		modelName := m.(Model).ToString()
		fmt.Printf("Migrating model: %s\n", modelName)
		errs := d.DB.AutoMigrate(m).GetErrors()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error migrating schema %s: %v", modelName, err)
			}
			return fmt.Errorf("Unable to migrate schemas")
		}
	}

	d.setForeignKeys(models)

	return nil
}

func (d *DB) setForeignKeys(models []interface{}) error {
	for _, m := range models {
		if _, ok := m.(Relation); ok {
			fkMaps := m.(Relation).GetForeignKeyMapping()
			// Create FK relationships
			for fk, mappedReference := range fkMaps {
				if err := d.DB.Model(m).AddForeignKey(fk, mappedReference, "CASCADE", "CASCADE").Error; err != nil {
					return err
				}

			}
		}

	}

	return nil
}

func cancelIn(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(seconds)*time.Second))
}

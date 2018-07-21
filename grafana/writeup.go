package grafana

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Writeup struct {
	Sections []Section
}

type Section struct {
	Title string
	Content string
}

type WriteupClient interface {
	GetWriteup() (Writeup, error)
}

type writeupClient struct {
	username, password, host, port, database string
	ids []interface{}
	queryStr string
}

func NewWriteupClient(host, port, username, password, database string, ids []interface{}, queryStr string) WriteupClient {
	return &writeupClient{username, password, host, port, database, ids, queryStr}
}

func (c *writeupClient) GetWriteup() (Writeup, error) {
	if len(c.ids) == 0 {
		return Writeup{}, nil
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.username, c.password, c.host, c.port, c.database))
	if err != nil {
		return Writeup{}, err
	}
	defer db.Close()

	results, err := db.Query(fmt.Sprintf("%s", c.queryStr), c.ids...)
	if err != nil {
		return Writeup{}, err
	}
	defer results.Close()

	var sections []Section
	for results.Next() {
		var section Section
		err = results.Scan(&section.Title, &section.Content)
		if err != nil {
			return Writeup{}, err
		}
		section.Title = sanitizeLaTexInput(section.Title)
		section.Content = sanitizeLaTexInput(section.Content)
		sections = append(sections, section)
	}
	if err != nil {
		return Writeup{}, err
	}
	return Writeup{sections}, nil
}

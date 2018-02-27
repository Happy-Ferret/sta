/*
Package db manages database access for github.com/ribacq/sta.
*/
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ribacq/sta/context"
)

func conn() (db *sql.DB, err error) {
	return sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", username, password, hostname, dbname))
}

// GetContext fetches a context from the database, given its ID.
func GetContext(id int) (c *context.Context, err error) {
	// open connection
	db, err := conn()
	if err != nil {
		return
	}

	// let’s fetch data from the various tables
	// name and description
	row := db.QueryRow("select (name, description) from sta.contexts where id = $1", id)
	err = row.Scan(&(c.Name), &(c.Description))
	if err != nil {
		return
	}

	// container
	row = db.QueryRow("select (container) from sta.contexts_container where context = $1", id)
	err = row.Scan(&(c.Container))
	if err != nil {
		return
	}

	// contents
	rows, err := db.Query("select (context) from sta.contexts_container where container = $1", id)
	if err != nil {
		return
	}
	for rows.Next() {
		var ctxID int
		err = rows.Scan(&ctxID)
		if err != nil {
			return
		}
		c.Contents = append(c.Contents, ctxID)
	}

	// links
	rows, err = db.Query("select (name, key, locked, target) from sta.contexts_links where context = $1", id)
	for rows.Next() {
		var name, key string
		var locked bool
		var target int
		err = rows.Scan(&name, &key, &locked, &target)
		if err != nil {
			return
		}
		c.AddLink(name, key, locked, target)
	}

	// commands
	c.Commands = map[string]context.CommandFunc{}
	rows, err = db.Query("select (command) from sta.contexts_commands where context = $1", id)
	if err != nil {
		return
	}
	for rows.Next() {
		var cmd string
		err = rows.Scan(&cmd)
		if err != nil {
			return
		}
		if f, ok := context.GetCommandFunc(cmd); ok {
			c.Commands[cmd] = f
		}
	}

	// properties
	c.Properties = map[string]string{}
	rows, err = db.Query("select (name, value) from sta.contexts_properties where context = $1", id)
	if err != nil {
		return
	}
	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)
		if err != nil {
			return
		}
		c.Properties[name] = value
	}

	return
}

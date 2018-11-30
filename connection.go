package schemaless

import (
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/jmoiron/sqlx"
)

// ConnectionOptions sets the options for Connection
type ConnectionOptions struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Protocol   string
	Parameters map[string]string
}

// Connection connects to a database
type Connection struct {
	Options *ConnectionOptions
	DB      *sqlx.DB
}

// GetParameters parses the connection parameters as an encoded URL query string
func (c *Connection) GetParameters() string {
	params := "?"
	for key, value := range c.Options.Parameters {
		params = params + fmt.Sprintf("%s=%s&", url.QueryEscape(key), url.QueryEscape(value))
	}

	params = strings.TrimSuffix(params, "&")
	return params
}

// GetString returns the database connection string
func (c *Connection) GetString(withoutDatabase bool) string {
	connStr := fmt.Sprintf("%s:%s@%s(%s:%s)/",
		c.Options.Username,
		c.Options.Password,
		c.Options.Protocol,
		c.Options.Host,
		c.Options.Port,
	)
	params := c.GetParameters()
	if withoutDatabase {
		return connStr + params
	}
	return connStr + c.Options.Database + params
}

// SelectOrCreate will select or create a database
func (c *Connection) SelectOrCreate() error {
	connStr := c.GetString(true)
	db, err := sqlx.Connect("mysql", connStr)

	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("USE %s;", c.Options.Database))

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unknown database") {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", c.Options.Database))

			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

// Connect will try to connect to a database
func (c *Connection) Connect() error {
	connStr := c.GetString(false)

	err := c.SelectOrCreate()

	if err != nil {
		return err
	}

	db, err := sqlx.Connect("mysql", connStr)

	if err != nil {
		return err
	}

	c.DB = db

	return nil
}

// Disconnect will try to disconnect from a database
func (c *Connection) Disconnect() {
	defer c.DB.Close()
}

// NewConnection creates a Connection
func NewConnection(options ConnectionOptions) *Connection {
	return &Connection{Options: &options}
}

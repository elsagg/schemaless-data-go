package schemaless

import (
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v4"
)

// DataSource is a data source
type DataSource struct {
	Name       string
	Connection *Connection
}

func (d *DataSource) findTable() error {
	err := d.Connection.Connect()

	if err != nil {
		return err
	}

	_, err = d.Connection.DB.Exec(fmt.Sprintf(CreateTableQuery, d.Name))

	if err != nil {
		return err
	}

	d.Connection.Disconnect()

	return nil
}

// GetCell will search for a cell in the DataSource
func (d *DataSource) GetCell(RowKey string, ColumnKey string, RefKey int64) (*DataCell, error) {
	var cell DataCell

	err := d.findTable()

	if err != nil {
		return nil, err
	}

	err = d.Connection.Connect()

	if err != nil {
		return nil, err
	}

	err = d.Connection.DB.Get(
		&cell,
		fmt.Sprintf(GetCellQuery, d.Name),
		RowKey,
		ColumnKey,
		RefKey,
	)

	d.Connection.Disconnect()

	return &cell, err
}

// GetCellLatest will search for the lastest ref of cell in the DataSource
func (d *DataSource) GetCellLatest(RowKey string, ColumnKey string) (*DataCell, error) {
	var cell DataCell

	err := d.findTable()

	if err != nil {
		return nil, err
	}

	err = d.Connection.Connect()

	if err != nil {
		return nil, err
	}

	err = d.Connection.DB.Get(
		&cell,
		fmt.Sprintf(GetCellLatestQuery, d.Name),
		RowKey,
		ColumnKey,
	)

	d.Connection.Disconnect()

	return &cell, err
}

// GetAllLatest will search for all latest cells in the DataSource
func (d *DataSource) GetAllLatest(ColumnKey string) (*[]DataCell, error) {
	var cells []DataCell

	err := d.findTable()

	if err != nil {
		return nil, err
	}

	err = d.Connection.Connect()

	if err != nil {
		return nil, err
	}

	err = d.Connection.DB.Select(
		&cells,
		fmt.Sprintf(GetAllLatestQuery, d.Name),
		ColumnKey,
	)

	d.Connection.Disconnect()

	return &cells, err
}

// PutCell will create a new cell in the DataSource
func (d *DataSource) PutCell(RowKey string, ColumnKey string, Body interface{}) (*DataCell, error) {
	cell, err := d.GetCellLatest(RowKey, ColumnKey)

	if err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "no rows in result set") {
			return nil, err
		}
	}

	body, err := msgpack.Marshal(&Body)

	if err != nil {
		return nil, err
	}

	err = d.Connection.Connect()

	if err != nil {
		return nil, err
	}

	_, err = d.Connection.DB.Exec(
		fmt.Sprintf(PutCellQuery, d.Name),
		RowKey,
		ColumnKey,
		body,
		cell.RefKey+1,
	)

	if err != nil {
		return nil, err
	}

	return d.GetCellLatest(RowKey, ColumnKey)
}

// NewDataSource returns a new DataSource instance
func NewDataSource(name string, connection *Connection) *DataSource {
	connection.Options.Database = name
	return &DataSource{Name: name, Connection: connection}
}

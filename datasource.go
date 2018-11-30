package schemaless

import (
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack"
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

// FindCell will search for a cell in the DataSource
func (d *DataSource) FindCell(RowID string, ColumnName string, RefKey int) (*DataCell, error) {
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
		fmt.Sprintf(FindCellQuery, d.Name),
		RowID,
		ColumnName,
		RefKey,
	)

	d.Connection.Disconnect()

	return &cell, err
}

// FindLatestCell will search for the lastest ref of cell in the DataSource
func (d *DataSource) FindLatestCell(RowID string, ColumnName string) (*DataCell, error) {
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
		fmt.Sprintf(FindLatestCellQuery, d.Name),
		RowID,
		ColumnName,
	)

	d.Connection.Disconnect()

	return &cell, err
}

// FindAllLatest will search for all latest cells in the DataSource
func (d *DataSource) FindAllLatest(ColumnName string) (*[]DataCell, error) {
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
		fmt.Sprintf(FindAllLatestQuery, d.Name, d.Name),
		ColumnName,
	)

	d.Connection.Disconnect()

	return &cells, err
}

// CreateCell will create a new cell in the DataSource
func (d *DataSource) CreateCell(RowID string, ColumnName string, Body interface{}) (*DataCell, error) {
	cell, err := d.FindLatestCell(RowID, ColumnName)

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
		fmt.Sprintf(CreateCellQuery, d.Name),
		RowID,
		ColumnName,
		body,
		cell.RefKey+1,
	)

	if err != nil {
		return nil, err
	}

	return d.FindLatestCell(RowID, ColumnName)
}

// NewDataSource returns a new DataSource instance
func NewDataSource(name string, connection *Connection) *DataSource {
	connection.Options.Database = name
	return &DataSource{Name: name, Connection: connection}
}

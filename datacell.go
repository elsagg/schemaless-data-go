package schemaless

import (
	"time"

	"github.com/vmihailenco/msgpack/v4"
)

// DataCell represents a data cell in the database
type DataCell struct {
	AddedID   string    `db:"added_id"`
	RowKey    string    `db:"row_key"`
	ColumnKey string    `db:"column_key"`
	RefKey    int64     `db:"ref_key"`
	Body      []byte    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

// MarshalBody transforms a body struct into a MessagePack JSON
func (d *DataCell) MarshalBody(v interface{}) (err error) {
	b, err := msgpack.Marshal(&v)

	if err != nil {
		return err
	}

	d.Body = b
	return nil
}

// UnmarshalBody transforms a MessagePack JSON into a body struct
func (d *DataCell) UnmarshalBody(v interface{}) (err error) {
	err = msgpack.Unmarshal(d.Body, &v)
	return
}

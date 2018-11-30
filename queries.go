package schemaless

// Query constants
const (
	FindLatestCellQuery = "SELECT * FROM %s WHERE row_id=? AND column_name=? ORDER BY ref_key DESC LIMIT 1"
	FindCellQuery       = "SELECT * FROM %s WHERE row_id=? AND column_name=? AND ref_key=?"
	CreateCellQuery     = "INSERT INTO %s (row_id, column_name, body, ref_key) VALUES (?,?,?,?)"
	FindAllLatestQuery  = `SELECT a.*
	FROM %s a
	LEFT OUTER JOIN %s b
		ON a.row_id = b.row_id AND a.column_name = b.column_name AND a.ref_key < b.ref_key
	WHERE b.row_id IS NULL AND a.column_name=?`
	CreateTableQuery = `CREATE TABLE IF NOT EXISTS %s (
		added_id int(11) NOT NULL AUTO_INCREMENT,
		row_id varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		column_name varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		body blob,
		ref_key int(11) DEFAULT NULL,
		created_at datetime DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (added_id)
	  )`
)

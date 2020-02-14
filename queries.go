package schemaless

// Query constants
const (
	GetCellLatestQuery = "SELECT * FROM %s WHERE row_key=? AND column_key=? ORDER BY ref_key DESC LIMIT 1"
	GetCellQuery       = "SELECT * FROM %s WHERE row_key=? AND column_key=? AND ref_key=?"
	PutCellQuery       = "INSERT INTO %s (row_key, column_key, body, ref_key) VALUES (?,?,?,?)"
	GetAllLatestQuery  = `SELECT a.*
	FROM %[1]s a
	LEFT OUTER JOIN %[1]s b
		ON a.row_key = b.row_key AND a.column_key = b.column_key AND a.ref_key < b.ref_key
	WHERE b.row_key IS NULL AND a.column_key=?`
	CreateTableQuery = `CREATE TABLE IF NOT EXISTS %[1]s (
		added_id bigint(20) NOT NULL AUTO_INCREMENT,
		row_key varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		column_key varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		body blob,
		ref_key bigint(20) DEFAULT NULL,
		created_at datetime DEFAULT CURRENT_TIMESTAMP,
		INDEX %[1]s_cell_index (row_key, column_key, ref_key),
		PRIMARY KEY (added_id)
	  )`
)

package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSupernovaLocks, downCreateTableSupernovaLocks)
}

func upCreateTableSupernovaLocks(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table supernova_locks
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		logged_by text not null,
		user_address text not null,
		locked_until bigint,
		locked_at bigint,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	create index user_locked_until_idx
		on supernova_locks (user_address asc, included_in_block desc, log_index desc);

	`)
	return err
}

func downCreateTableSupernovaLocks(tx *sql.Tx) error {
	_, err := tx.Exec("drop table supernova_locks")
	return err
}

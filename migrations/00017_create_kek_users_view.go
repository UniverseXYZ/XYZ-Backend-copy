package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(UpKekUsersView, DownKekUsersView)
}

func UpKekUsersView(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create or replace view kek_users_with_balance as
		with users as ( select distinct sender as address
						from kek_transfers
						union
						select distinct receiver
						from kek_transfers ),
			 value_out as ( select sum(value) as val, sender as addr from kek_transfers group by sender ),
			 value_in as ( select sum(value) as val, receiver as addr from kek_transfers group by receiver )
		select address,
			   coalesce(( select val from value_in where addr = u.address ), 0) -
			   coalesce(( select val from value_out where addr = u.address ), 0) as balance
		from users u
		order by balance desc;
	`)
	return err
}

func DownKekUsersView(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop view kek_users_with_balance;
	`)
	return err
}

package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upReplaceViewKekUsersWithBalance, downReplaceViewKekUsersWithBalance)
}

func upReplaceViewKekUsersWithBalance(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create or replace view kek_users_with_balance as
		with transfers as ( select sender as address, -value as amount
							from kek_transfers
							union all
							select receiver as address, value as amount
							from kek_transfers )
		select address, sum(amount) as balance
		from transfers
		group by address;

		create or replace view kek_users_with_balance_no_staking as
		with transfers as ( select sender as address, -value as amount
							from kek_transfers
							where sender not in ('0xb0fa2beee3cf36a7ac7e99b885b48538ab364853')
							  and receiver not in ('0xb0fa2beee3cf36a7ac7e99b885b48538ab364853')
							union all
							select receiver as addrress, value as amount
							from kek_transfers
							where sender not in ('0xb0fa2beee3cf36a7ac7e99b885b48538ab364853')
							  and receiver not in ('0xb0fa2beee3cf36a7ac7e99b885b48538ab364853') )
		select address, sum(amount) as balance
		from transfers
		group by address;

	`)
	return err
}

func downReplaceViewKekUsersWithBalance(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop view if exists kek_users_with_balance_no_staking;
	`)
	return err
}

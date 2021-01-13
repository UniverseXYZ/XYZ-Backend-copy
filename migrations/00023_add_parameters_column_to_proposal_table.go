package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddParametersColumnToProposalTable, downAddParametersColumnToProposalTable)
}

func upAddParametersColumnToProposalTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	alter table governance_proposals 
			add warm_up_duration bigint,
	    	add active_duration bigint,
	    	add queue_duration bigint,
	    	add grace_period_duration bigint,
	    	add acceptance_threshold bigint,
	    	add min_quorum bigint
	`)
	return err
}

func downAddParametersColumnToProposalTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table governance_proposals 
		    add warm_up_duration bigint,
	    	add active_duration bigint,
	    	add queue_duration bigint,
	    	add grace_period_duration bigint,
	    	add acceptance_threshold bigint,
	    	add min_quorum bigint
	    	`)
	return err
}

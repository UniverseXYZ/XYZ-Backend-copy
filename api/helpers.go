package api

import (
	"database/sql"
	"strconv"

	"github.com/kekDAO/kekBackend/api/types"
	"github.com/pkg/errors"
)

func calculateOffset(limit string, page string) (string, error) {
	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return "", err
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return "", err
	}

	offset := (p - 1) * l

	return strconv.FormatInt(offset, 10), nil
}

func (a *API) getProposalEvents(id uint64) ([]types.Event, error) {
	rows, err := a.db.Query(`
		select proposal_id,
		       caller,
		       event_type,
		       event_data,
		       block_timestamp,
		       tx_hash
		from governance_events 
		where proposal_id = $1`, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "could not query proposal events")
	}

	var eventsList []types.Event

	for rows.Next() {
		var event types.Event
		err := rows.Scan(&event.ProposalID, &event.Caller, &event.EventType, &event.Eta, &event.CreateTime, &event.TxHash)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan proposal event")
		}

		eventsList = append(eventsList, event)
	}

	return eventsList, nil
}

func (a *API) getHighestBlock() (*int64, error) {
	var number int64

	err := a.db.QueryRow(`select number from blocks order by number desc limit 1;`).Scan(&number)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "could not get highest block")
	}

	return &number, nil
}

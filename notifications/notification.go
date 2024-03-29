package notifications

import (
	"context"
	"database/sql"

	"github.com/kekDAO/kekBackend/types"
	"github.com/pkg/errors"
)

type Notification struct {
	Target           string
	NotificationType string
	StartsOn         int64
	ExpiresOn        int64
	Message          string
	Metadata         types.JSONObject
	IncludedInBlock  int64
}

func (n *Notification) ToDBWithTx(ctx context.Context, tx *sql.Tx) error {
	ins := `
		INSERT INTO
			"notifications" ("target", "type", "starts_on", "expires_on", "message", "metadata", "included_in_block")
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		;
	`
	_, err := tx.ExecContext(ctx, ins, n.Target, n.NotificationType, n.StartsOn, n.ExpiresOn, n.Message, n.Metadata, n.IncludedInBlock)
	if err != nil {
		return errors.Wrap(err, "could not exec statement")
	}

	return nil
}

func NewNotification(target string, typ string, starts int64, expires int64, msg string, metadata map[string]interface{}, block int64) Notification {
	return Notification{
		Target:           target,
		NotificationType: typ,
		StartsOn:         starts,
		ExpiresOn:        expires,
		Message:          msg,
		Metadata:         metadata,
		IncludedInBlock:  block,
	}
}

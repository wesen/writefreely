/*
 * Copyright © 2019-2021 A Bunch Tell LLC.
 *
 * This file is part of WriteFreely.
 *
 * WriteFreely is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, included
 * in the LICENSE file in this source code package.
 */

package migrations

import (
	"context"
	"database/sql"

	wf_db "github.com/writefreely/writefreely/db"
)

func oauth(db *datastore) error {
	dialect := wf_db.DialectMySQL
	if db.driverName == driverSQLite {
		dialect = wf_db.DialectSQLite
	}
	return wf_db.RunTransactionWithOptions(context.Background(), db.DB, &sql.TxOptions{}, func(ctx context.Context, tx *sql.Tx) error {
		createTableUsersOauth, err := dialect.
			Table("oauth_users").
			SetIfNotExists(false).
			Column(dialect.Column("user_id", wf_db.ColumnTypeInteger, wf_db.UnsetSize)).
			Column(dialect.Column("remote_user_id", wf_db.ColumnTypeInteger, wf_db.UnsetSize)).
			PrimaryKeyConstraint("user_id", "remote_user_id").
			ToSQL()
		if err != nil {
			return err
		}
		createTableOauthClientState, err := dialect.
			Table("oauth_client_states").
			SetIfNotExists(false).
			Column(dialect.Column("state", wf_db.ColumnTypeVarChar, wf_db.OptionalInt{Set: true, Value: 255})).
			Column(dialect.Column("used", wf_db.ColumnTypeBool, wf_db.UnsetSize)).
			Column(dialect.Column("created_at", wf_db.ColumnTypeDateTime, wf_db.UnsetSize).SetDefaultCurrentTimestamp()).
			PrimaryKeyConstraint("state").
			UniqueConstraint("state").
			ToSQL()
		if err != nil {
			return err
		}

		for _, table := range []string{createTableUsersOauth, createTableOauthClientState} {
			if _, err := tx.ExecContext(ctx, table); err != nil {
				return err
			}
		}
		return nil
	})
}

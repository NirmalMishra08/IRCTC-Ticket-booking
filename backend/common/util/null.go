package util

import "database/sql"

// ToNullString converts a string to sql.NullString.

func ToNullString(s string) sql.NullString {

	if s == "" {

		return sql.NullString{Valid: false}

	}

	return sql.NullString{String: s, Valid: true}

}

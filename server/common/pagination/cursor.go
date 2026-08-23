package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

func Encode(cursor Cursor) string {
	data, _ := json.Marshal(cursor)
	return base64.RawURLEncoding.EncodeToString(data)
}

func Decode(value string) (Cursor, error) {
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, errors.New("invalid cursor")
	}
	var cursor Cursor
	if json.Unmarshal(data, &cursor) != nil || cursor.ID <= 0 || cursor.CreatedAt.IsZero() {
		return Cursor{}, errors.New("invalid cursor")
	}
	return cursor, nil
}

func DecodeOptional(value string) (*Cursor, error) {
	if value == "" {
		return nil, nil
	}
	cursor, err := Decode(value)
	if err != nil {
		return nil, err
	}
	return &cursor, nil
}

func Limit(value uint32) int {
	if value == 0 {
		return 25
	}
	if value > 100 {
		return 100
	}
	return int(value)
}

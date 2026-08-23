package pagination

import (
	"testing"
	"time"
)

func TestCursorRoundTrip(t *testing.T) {
	want := Cursor{CreatedAt: time.Date(2026, 8, 24, 10, 20, 30, 123, time.UTC), ID: 42}
	encoded := Encode(want)
	got, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !got.CreatedAt.Equal(want.CreatedAt) || got.ID != want.ID {
		t.Fatalf("Decode(Encode(cursor)) = %#v, want %#v", got, want)
	}
}

func TestDecodeRejectsMalformedCursor(t *testing.T) {
	for _, value := range []string{"%%%", "e30", Encode(Cursor{})} {
		if _, err := Decode(value); err == nil {
			t.Fatalf("Decode(%q) succeeded", value)
		}
	}
}

func TestLimit(t *testing.T) {
	for _, tc := range []struct {
		input uint32
		want  int
	}{{0, 25}, {1, 1}, {100, 100}, {101, 100}} {
		if got := Limit(tc.input); got != tc.want {
			t.Fatalf("Limit(%d) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

package models

import "testing"

func TestInitialMessageStatus(t *testing.T) {
	if got := InitialMessageStatus("visitor"); got != MessageStatusUnread {
		t.Fatalf("visitor message status = %q, want %q", got, MessageStatusUnread)
	}
	if got := InitialMessageStatus("kefu"); got != MessageStatusRead {
		t.Fatalf("kefu message status = %q, want %q", got, MessageStatusRead)
	}
}

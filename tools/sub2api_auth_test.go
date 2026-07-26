package tools

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestCheckSub2apiAuthConfigReportsMissingURL(t *testing.T) {
	withServerConfEnv(t, "")

	err := CheckSub2apiAuthConfig()
	if err == nil {
		t.Fatal("expected missing sub2api API URL error")
	}
	if got := Sub2apiAuthReason(err); got != Sub2apiAuthReasonConfigMissing {
		t.Fatalf("reason = %q, want %q", got, Sub2apiAuthReasonConfigMissing)
	}
}

func TestCheckSub2apiAuthConfigRejectsInvalidURL(t *testing.T) {
	withServerConfEnv(t, "://bad-url")

	err := CheckSub2apiAuthConfig()
	if err == nil {
		t.Fatal("expected invalid sub2api API URL error")
	}
	if got := Sub2apiAuthReason(err); got != Sub2apiAuthReasonInvalidBaseURL {
		t.Fatalf("reason = %q, want %q", got, Sub2apiAuthReasonInvalidBaseURL)
	}
}

func TestSub2apiAuthReasonFallsBackForPlainErrors(t *testing.T) {
	if got := Sub2apiAuthReason(errors.New("plain error")); got != Sub2apiAuthReasonUnknown {
		t.Fatalf("reason = %q, want %q", got, Sub2apiAuthReasonUnknown)
	}
}

func withServerConfEnv(t *testing.T, apiURL string) {
	t.Helper()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "config"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOFLY_SUB2API_API_URL", apiURL)

	oldConf := serverConf
	oldOnce := serverConfOnce
	serverConf = nil
	serverConfOnce = sync.Once{}

	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
		serverConf = oldConf
		serverConfOnce = oldOnce
	})
}

package telegramtest

import (
	"bytes"
	"regexp"
	"testing"
)

func TestGenerateCodeReturnsSixDigits(t *testing.T) {
	t.Parallel()

	code, err := GenerateCode(bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x2a}))
	if err != nil {
		t.Fatalf("GenerateCode() error = %v", err)
	}

	if ok, _ := regexp.MatchString(`^[0-9]{6}$`, code); !ok {
		t.Fatalf("GenerateCode() = %q, want six digits", code)
	}
}

func TestBuildChallengeMessage(t *testing.T) {
	t.Parallel()

	got := BuildChallengeMessage("123456")
	want := "这是一个测试，请回复 \"[123456]\""
	if got != want {
		t.Fatalf("BuildChallengeMessage() = %q, want %q", got, want)
	}
}

func TestIsMatchingReply(t *testing.T) {
	t.Parallel()

	if !IsMatchingReply("123456", "123456") {
		t.Fatal("IsMatchingReply() = false, want true")
	}
	if !IsMatchingReply(" \"123456\" ", "123456") {
		t.Fatal("IsMatchingReply() with quoted text = false, want true")
	}
	if !IsMatchingReply("[123456]", "123456") {
		t.Fatal("IsMatchingReply() with brackets = false, want true")
	}
	if IsMatchingReply("123457", "123456") {
		t.Fatal("IsMatchingReply() = true, want false")
	}
}

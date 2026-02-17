package telegramtest

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

func GenerateCode(r io.Reader) (string, error) {
	var b [4]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	n := binary.BigEndian.Uint32(b[:]) % 1000000
	return fmt.Sprintf("%06d", n), nil
}

func BuildChallengeMessage(code string) string {
	return fmt.Sprintf("这是一个测试，请回复 \"[%s]\"", code)
}

func IsMatchingReply(reply string, code string) bool {
	trimmed := strings.TrimSpace(reply)
	trimmed = strings.Trim(trimmed, "\"'“”‘’[]")
	return trimmed == code
}

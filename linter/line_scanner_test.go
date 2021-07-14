package linter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var linterLineOutput string = "main.go:16: File is not `gofmt`-ed with `-s` (gofmt)\n" +
	"vmconf.go:4:2: struct field `Id` should be `ID` (golint)\n"

func TestIssueLineScanner_Next_ScanTwoLines(t *testing.T) {
	scanner := NewLineScanner(bytes.NewBuffer([]byte(linterLineOutput)))

	for i := 0; i < 2; i++ {
		ok := scanner.Next()

		require.True(t, ok, i)
		require.NoError(t, scanner.Err())

		_ = scanner.Get()
	}
}

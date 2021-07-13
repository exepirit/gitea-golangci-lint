package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var linterLineOutput string = "main.go:16: File is not `gofmt`-ed with `-s` (gofmt)\n" +
	"vmconf.go:4:2: struct field `Id` should be `ID` (golint)"

func TestIssueLineScanner_Next_ScanTwoLines(t *testing.T) {
	scanner := NewLineScanner(bytes.NewBuffer([]byte(linterLineOutput)))

	for i := 0; i < 2; i++ {
		ok := scanner.Next()

		if i == 2-1 {
			require.False(t, ok)
		} else {
			require.True(t, ok)
		}
		require.NoError(t, scanner.Err())

		_ = scanner.Get()
	}
}

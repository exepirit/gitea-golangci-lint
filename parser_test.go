package main

import (
	"reflect"
	"testing"
)

func TestParseLinterRegexp(t *testing.T) {
	input := []string{
		"main.go:16: File is not `gofmt`-ed with `-s` (gofmt)\n",
		"vmconf.go:4:2: struct field `Id` should be `ID` (golint)",
	}

	for _, inp := range input {
		if !linterOutputRegex.MatchString(inp) {
			t.Fatalf("regexp not workng on %q", inp)
		}
	}
}

func Test_parseIssue(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Issue
	}{
		{
			name: "Single line",
			args: []string{"", "models/devices.go", "4", "", "msg", "golint"},
			want: Issue{
				Filename:   "models/devices.go",
				LineNum:    4,
				Message:    "msg",
				LinterName: "golint",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIssue(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIssue() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

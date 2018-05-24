package cli

import (
	"testing"
)

func TestParseBuildInfo(t *testing.T) {
	input := `time:"Sat May 13 19:53:08 UTC 2017" branch:master commit:320279c patches:1234`

	buildinfo := ParseBuildInfo(input)

	if buildinfo.Timestamp != "Sat May 13 19:53:08 UTC 2017" {
		t.Error("parsed time is wrong")
	}
	if buildinfo.GitBranch != "master" {
		t.Error("parsed branch is wrong")
	}
	if buildinfo.GitCommit != "320279c" {
		t.Error("parsed commit is wrong")
	}
	if buildinfo.GitRevCount != "1234" {
		t.Error("parsed patches is wrong")
	}
}

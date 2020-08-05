package main

import (
	"testing"
	"time"

	"github.com/m-lab/go/prometheusx"

	"github.com/m-lab/go/flagx"
)

func Test_main(t *testing.T) {
	// Run once with a cancelled main context to return immediately.
	mainCancel()
	projects = flagx.StringArray{"fake-project-id"}
	exit := 0
	logFatal = func(...interface{}) {
		exit = 1
	}

	main()

	if exit != 1 {
		t.Fatal("Expected exit")
	}
}

func Test_main_success(t *testing.T) {
	// Run once with a cancelled main context to return immediately.
	mainCancel()
	*prometheusx.ListenAddress = ":0"
	projects = flagx.StringArray{"fake-project-id"}
	collectTimes = flagx.DurationArray{time.Second}

	main()
}

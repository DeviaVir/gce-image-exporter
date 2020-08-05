package main

import (
	"testing"
	"time"

	"github.com/m-lab/go/prometheusx"

	"github.com/m-lab/go/flagx"
	"google.golang.org/api/option"
)

func Test_main(t *testing.T) {
	// Run once with a cancelled main context to return immediately.
	mainCancel()
	opts = []option.ClientOption{option.WithoutAuthentication()}
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
	opts = []option.ClientOption{option.WithoutAuthentication()}
	*prometheusx.ListenAddress = ":0"
	projects = flagx.StringArray{"fake-project-id"}
	collectTimes = flagx.DurationArray{time.Second}

	main()
}

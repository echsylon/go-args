package args_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/echsylon/go-args"
)

func Test_EndToEndSmokeTest_Success(t *testing.T) {
	// Setup mock call arguments
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()
	os.Args = []string{"xmpl", "-m", "5", "--verbose", "file 1.txt", "2000", "file2.txt"}

	// Setup test configuration
	args.Reset()
	args.SetApplicationDescription("A beautiful example app.")
	args.DefineOption("m", "Max lines to read.") // simple short-name option with no filter
	args.DefineOptionStrict("v", "verbose", "Print detailed output.", `^(true|false)$`)
	args.DefineArgumentStrict("FILES", "Files to read from.", 1, 2, `^*\.txt$`)
	args.DefineArgument("TIMEOUT", "Read timeout milliseconds.")
	args.Parse()

	// Extract state
	var maxLines int64 = args.GetOptionIntValue("m", 2)
	var printMuch bool = args.GetOptionBoolValue("verbose", false)
	var files []string = args.GetArgumentValues("FILES")
	var time []int64 = args.GetArgumentIntValues("TIMEOUT")

	// Validate
	fmt.Printf("verbose:  %t\n", printMuch)
	fmt.Printf("lines:    %d\n", maxLines)
	fmt.Printf("files:    %v (%d)\n", files, len(files))
	fmt.Printf("time:     %v\n", time)

	if printMuch != true {
		t.Errorf("Expected 'printMuch' to be <true>, but it's <%t>", printMuch)
	}

	if maxLines != 5 {
		t.Errorf("Expected 'maxLines' to be <5>, but it's <%d>", maxLines)
	}

	if files[0] != "file 1.txt" || files[1] != "file2.txt" || len(files) != 2 {
		t.Errorf("Expected 'files' to be <[file 1.txt file2.txt]>, but it's <%v>", files)
	}

	if time[0] != 2000 || len(time) != 1 {
		t.Errorf("Expected 'time' to be <[2000]>, but it's <%v>", time)
	}
}

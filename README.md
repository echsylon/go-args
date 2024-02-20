# A smarter command line argument parser for GO.

This library provides a smart and easy way for GO application developers to parse and validate their command line arguments.

## Key features

* Conceptual separation of "options" (optional) and "arguments" (mandatory).
* Support for short- and long name options, e.g. `-v` and `--verbose`.
* RegEx validation on user provided option and argument values.
* Range constraints on argument values (min/max number of accepted values)
* Typed value extraction (e.g "getOptionBoolValue")

## A concrete example
Consider below example code:

```go
package main

import (
    "fmt"
    "github.com/echsylon/go-args"
)

func main() {
	// Make help text nice
	args.SetApplicationDescription("A beautiful example app.")

	// Define optional user input
	args.DefineOption("m", "Max lines to read.") // simple short-name option with no filter
	args.DefineOptionStrict("v", "verbose", "Print detailed output.", `^(true|false)$`)

	// Define mandatory user input
	args.DefineArgumentStrict("FILES", "Files to read from.", 1, 2, `^*\.txt$`)
	args.DefineArgument("TIMEOUT", "Read timeout milliseconds.")

	// Read and validate provided user data
	args.Parse()

	// Get the values...
	var maxLines int64 = args.GetOptionIntValue("m", 2)
	var printMuch bool = args.GetOptionBoolValue("verbose", false)
	var files []string = args.GetArgumentValues("FILES")
	var time []int64 = args.GetArgumentIntValues("TIMEOUT")

	// ...and use them as you see fit
	fmt.Printf("verbose:  %t\n", printMuch)
	fmt.Printf("lines:    %d\n", maxLines)
	fmt.Printf("files:    %v (%d)\n", files, len(files))
	fmt.Printf("time:     %v\n", time)
}
```

Assuming that the name of the built application is `xmpl`, we can call it with this command:
```bash
$ ./xmpl -m 5 --verbose "file 1.txt" 2000 file2.txt
```

Since all configured constraints are met, the application would output:
```
verbose:  true
lines:    5
files:    [file 1.txt file2.txt] (2)
time:     [2000]
```



## Worth mentioning

The constraint patterns for both options and arguments in above example have been chosen carefully, hence the order of the FILES and TIMEOUT arguments can be mixed freely and still end up under the correct definition. The same goes for the options: since the `--verbose` flag only accepts `true` or `false` as a value, it won't capture the `"file 1.txt"` argument.

Would the command instead have been the below (with an additional `file3.txt` argument):
```bash
$ ./xmpl -m 5 --verbose "file 1.txt" 2000 file2.txt file3.txt
```

The output from the application would instead have been like so (since the `FILES` argument is configured to only accept 2 values):
```
Unexpected input: file3.txt

Usage: ./xmpl [OPTIONS...] FILES... TIMEOUT
A beautiful example app.

Arguments:
  FILES   Files to read from.
  TIMEOUT Read timeout milliseconds.

Options:
  -m            Max lines to read.
  -v, --verbose Print detailed output.
```

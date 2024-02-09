# go-args: A smarter command line argument parser for GO.

This library provides a smart and easy way for GO application developers to parse and validate their command line arguments.

## Key features

* Separation of "options" (optional) and "arguments" (mandatory).
* Short- and long name options, e.g. `-v` and `--verbose`.
* RegEx validation on user provided option and argument values.
* Range constraints on user provided argument values (min values, max values)
* Typed value extraction (e.g "getOptionBoolValue")

## An example
Consider below example code:

```go
package main

import (
    "fmt"
    "github.com/echsylon/go-args"
)

func main() {
    // Make help text nice
    args.Description = "The example app showcases how to use the go-args module."

    // Define optional user input
    args.DefineOption("m", "Max lines to read.") // simple short-name option with no filter
    args.DefineConstrainedOption("v", "verbose", "Print detailed output.", `(true|false)`)

    // Define mandatory user input
    args.DefineConstrainedArgument("FILES", "Files to read from.", `^*\.txt$`, 1, 2)
    args.DefineArgument("TIMEOUT", "Read timeout milliseconds")

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

Key observation here is that the constraint patterns have been chosen carefully, hence the argument list can be mixed and still end up under the correct definition. The same goes for the options. Since the `--verbose` flag only accepts `true` or `false` as a value, it won't capture `"file 1.txt"`.

Would the command instead have been the below (with an additional `file3.txt` argument):
```bash
$ ./xmpl -m 5 --verbose "file 1.txt" 2000 file2.txt file3.txt
```

The output from the application would instead have been like so (since the `TIMEOUT` argument only accepts 1 value and `FILES` only 2):
```
Unexpected input: file3.txt

Usage: ./xmpl [OPTIONS...] FILES... TIMEOUT
The example app showcases how to use the go-args module.

Arguments:
  FILES   Files to read from.
  TIMEOUT Read timeout milliseconds

Options:
  -m            Max lines to read.
  -v, --verbose Print detailed output.
```

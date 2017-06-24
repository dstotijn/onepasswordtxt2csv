# onepasswordtxt2csv

Tool for converting a plain text file exported from 1Password on Windows to
a CSV file for use with [Utilities for 1Password](https://github.com/agilebits/onepassword-utilities).

## Installation (Linux/macOS)

```
$ go get github.com/dstotijn/onepasswordtxt2csv
```

## Usage (Linux/macOS)

Assuming you have `bin` directory used by the `go` command in your `$PATH`, you
can run the program like this:

Example:
```
$ onepasswordtxt2csv --txtFile="/path/to/input_file.txt" output_file.csv
```

## License

[MIT](/LICENSE.md)

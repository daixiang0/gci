# GCI

GCI, a tool that controls golang package import order and makes it always deterministic.

The desired output format is highly configurable and allows for more custom formatting than `goimport` does.

GCI considers a import block based on AST as below:
```
Doc
Name Path Comment
```
All comments will keep as they were, except the independent comment blocks(line breaks before and after).

GCI splits all import blocks into different sections, now support three section type:
- standard: Golang official imports, like "fmt"
- custom: Custom section, use full and the longest match(match full string first, if multiple matches, use the longest one)
- default: All rest import blocks

The priority is standard>custom>default, all sections sort alphabetically inside.

All import blocks use one TAB(`\t`) as Indent.

**Note**:

`nolint` is hard to handle at section level, GCI will consider it as a single comment.

## Download

```shell
$ go get github.com/daixiang0/gci
```

## Usage

Now GCI provides two command line methods, mainly for backward compatibility.

### New style
GCI supports three modes of operation

```shell
$ gci print -h
Print outputs the formatted file. If you want to apply the changes to a file use write instead!

Usage:
  gci print path... [flags]

Aliases:
  print, output

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for print
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry. The Section order is the same as below, default value is [Standard,Default].
                          Std | Standard - Captures all standard packages if they do not match another section
                          Prefix(github.com/daixiang0) | pkgPrefix(github.com/daixiang0) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
                          Def | Default - Contains all imports that could not be matched to another section type
                          [DEPRECATED] Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                          [DEPRECATED] NL | NewLine - Prints an empty line
```

```shell
$ gci write -h
Write modifies the specified files in-place

Usage:
  gci write path... [flags]

Aliases:
  write, overwrite

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for write
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry. The Section order is the same as below, default value is [Standard,Default].
                          Std | Standard - Captures all standard packages if they do not match another section
                          Prefix(github.com/daixiang0) | pkgPrefix(github.com/daixiang0) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
                          Def | Default - Contains all imports that could not be matched to another section type
                          [DEPRECATED] Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                          [DEPRECATED] NL | NewLine - Prints an empty line
```

```shell
$ gci diff -h
Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.

Usage:
  gci diff path... [flags]

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for diff
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry. The Section order is the same as below, default value is [Standard,Default].
                          Std | Standard - Captures all standard packages if they do not match another section
                          Prefix(github.com/daixiang0) | pkgPrefix(github.com/daixiang0) - Groups all imports with the specified Prefix. Imports will be matched to the full and longest Prefix. All groups are in alphabetical order.
                          Def | Default - Contains all imports that could not be matched to another section type
                          [DEPRECATED] Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                          [DEPRECATED] NL | NewLine - Prints an empty line
```

### Old style

```shell
Usage:
  gci [-diff | -write] [--local localPackageURLs] path... [flags]

Flags:
  -d, --diff            display diffs instead of rewriting files
  -h, --help            help for gci
  -l, --local strings   put imports beginning with this string after 3rd-party packages, separate imports by comma
  -v, --version         version for gci
  -w, --write           write result to (source) file instead of stdout

```

**Note**::

The old style is only for local tests, will be deprecated, please uses new style, `golangci-lint` uses new style as well.

## Examples

Run `gci write -s standard -s default -s "prefix(github.com/daixiang0/gci)" main.go` and you will handle following cases:

### simple case

```go
package main
import (
  "golang.org/x/tools"
  
  "fmt"
  
  "github.com/daixiang0/gci"
)
```

to

```go
package main
import (
    "fmt"

    "github.com/daixiang0/gci"

    "golang.org/x/tools"
)
```

### with alias

```go
package main
import (
  "fmt"
  go "github.com/golang"
  "github.com/daixiang0/gci"
)
```

to

```go
package main
import (
  "fmt"

  "github.com/daixiang0/gci"

  go "github.com/golang"
)
```

## TODO

- Ensure only one blank between `Name` and `Path` in an import block
- Ensure only one blank between `Path` and `Comment` in an import block
- Format comments
- Add more testcases
- Support imports completion (please use `goimports` first then use GCI)
- Optimize comments

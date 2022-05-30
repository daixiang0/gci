# GCI

GCI, a tool that controls golang package import order and makes it always deterministic.

The desired output format is highly configurable and allows for more custom formatting than `goimport` does.

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
      --NoInlineComments           Drops inline comments while formatting
      --NoPrefixComments           Drops comment lines above an import statement while formatting
  -s, --Section strings            Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry.
                                   Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                                   Def | Default - Contains all imports that could not be matched to another section type
                                   NL | NewLine - Prints an empty line
                                   Prefix(gitlab.com/myorg) | pkgPrefix(gitlab.com/myorg) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
                                   Std | Standard - Captures all standard packages if they do not match another section
                                    (default [Standard,Default])
  -x, --SectionSeparator strings   SectionSeparators are inserted between Sections (default [NewLine])
  -h, --help                       help for print
```

```shell
$ gci write -h
Write modifies the specified files in-place

Usage:
  gci write path... [flags]

Aliases:
  write, overwrite

Flags:
      --NoInlineComments           Drops inline comments while formatting
      --NoPrefixComments           Drops comment lines above an import statement while formatting
  -s, --Section strings            Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry.
                                   Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                                   Def | Default - Contains all imports that could not be matched to another section type
                                   NL | NewLine - Prints an empty line
                                   Prefix(gitlab.com/myorg) | pkgPrefix(gitlab.com/myorg) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
                                   Std | Standard - Captures all standard packages if they do not match another section
                                    (default [Standard,Default])
  -x, --SectionSeparator strings   SectionSeparators are inserted between Sections (default [NewLine])
  -h, --help                       help for write
```

```shell
$ gci diff -h
Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.

Usage:
  gci diff path... [flags]

Flags:
      --NoInlineComments           Drops inline comments while formatting
      --NoPrefixComments           Drops comment lines above an import statement while formatting
  -s, --Section strings            Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry.
                                   Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
                                   Def | Default - Contains all imports that could not be matched to another section type
                                   NL | NewLine - Prints an empty line
                                   Prefix(gitlab.com/myorg) | pkgPrefix(gitlab.com/myorg) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
                                   Std | Standard - Captures all standard packages if they do not match another section
                                    (default [Standard,Default])
  -x, --SectionSeparator strings   SectionSeparators are inserted between Sections (default [NewLine])
  -d, --debug                      Enables debug output from the formatter
  -h, --help                       help for diff
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

## Examples

Run `gci  write --Section Standard --Section Default --Section "Prefix(github.com/daixiang0/gci)" main.go` and you will handle following cases:

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

  "golang.org/x/tools"

  "github.com/daixiang0/gci"
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

  go "github.com/golang"

  "github.com/daixiang0/gci"
)
```

## TODO

- Add more testcases

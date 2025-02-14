# arcpy2go

Convert arcpy tools into Go code.

## Description

This little program generates Go code from a documentation of a arcpy tool.
This enables the creation of Go programs that output python code that uses arcpy
code.

## Installation

To install this program in your computer, you can use the `go install`command.

```shell
go install github.com/marcelo-fm/arcpy2go@latest
```

Or if you use the Go versions 1.24 or higher, you can use the `go get -tool`
directive:

```shell
go get -tool github.com/marcelo-fm/arcpy2go@latest
```

and then call it with `go tool arcpy2go`.

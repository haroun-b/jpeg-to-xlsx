# JPEG To XLSX

A simple command-line tool to encode JPEGs as XLSX files, where each pixel's RGB values are encoded using three cells' background colors, similar to how a screen would display an image.

## Motivation

Used as teaching support asset to demonstrate the fundamental fact that the basis of computer science is data representation and manipulation, and that most of what any computer does is retrieve, store, and manipulate data.

## Installation

### Using `go install`

```sh
go install github.com/haroun-b/jpeg-to-xlsx@latest
```

**Don't forget to add `$GOPATH/bin` to your `$PATH` if you haven't already**

```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

### Using A Binary

- Download the latest binary from the [releases page](https://github.com/haroun-b/jpeg-to-xlsx/releases)
- Add the binary to your binary directory (e.g. `/usr/local/bin`)

**Make sure the binary has execute permissions**

## Usage

```sh
jpeg-to-xlsx <./dir/source-img.jpeg> [<./dir/output.xlsx>]
```

## Example Output

![Example Output](./examples/example-output.gif)

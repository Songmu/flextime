flextime
=======

[![Test Status](https://github.com/Songmu/flextime/workflows/test/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/Songmu/flextime/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/flextime?status.svg)][godoc]

[actions]: https://github.com/Songmu/flextime/actions?workflow=test
[coveralls]: https://coveralls.io/r/Songmu/flextime?branch=master
[license]: https://github.com/Songmu/flextime/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/flextime

flextime improves time testability by replacing the backend clock flexibly.

## Synopsis

```go
now := flextime.Now() // returned normal current time

func () { // Set time
    restore := flextime.Set(time.Date(2001, time.May, 1, 10, 10, 10, 0, time.UTC))
    defer restore()

    now = flextime.Now() // returned set time
}()

func () { // Fix time
    restore := flextime.Fix(time.Date(2001, time.May, 1, 10, 10, 10, 0, time.UTC))
    defer restore()

    now = flextime.Now() // returned fixed time
}()
```

## Description

The flextime improves time testability by replacing the backend clock flexibly.

It has a set of functions similar to the standard time package, making it easy to utilize.

By default, it behaves the same as the standard time package, but allows us to change or fix
the current time by using `Fix` and `Set` function.

Also, we can replace the backend clock by implementing our own `Clock` interface and combining
it with the Switch function.

## Installation

```console
% go get github.com/Songmu/flextime
```

## Migration

You can almost migrate from standard time package to Songmu/flextime with the following command.

```console
% go get github.com/Songmu/flextime
% find . -name '*.go' | xargs perl -i -pe 's/\btime\.((?:N(?:ewTi(?:ck|m)er|ow)|After(?:Func)?|Sleep|Until|Tick))/flextime.$1/g'
% goimport -w .
```

## Author

[Songmu](https://github.com/Songmu)

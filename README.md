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
import "github.com/Songmu/flextime"

now := flextime.Now() // returned normal current time by default
flextime.Sleep()
d := flextime.Until(date)
d := flextime.Since(date)
<-flextime.After(5*time.Second)
flextime.AfterFunc(5*time.Second, func() { fmt.Println("Done") })
timer := flextime.NewTimer(10*time.Second)
ticker := flextime.NewTicker(10*time.Second)
ch := flextime.Tick(3*time.Second)

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

It has a set of functions similar to the standard time package, making it easy to migrate
from standard time package.

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

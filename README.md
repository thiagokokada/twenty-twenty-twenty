# Twenty-Twenty-Twenty

20-20-20 is a program that alerts every 20 minutes to look something at 20 feet
away for 20 seconds, written in Golang. This is done to reduce eye fadigue [1].
While controversial[2][3], it is a simple rule to remember and can help long
screen usage sessions.


## How to use

```
$ ./twenty-twenty-twenty -help
Usage of ./twenty-twenty-twenty:
  -duration uint
    	how long each pause should be in seconds (default 20)
  -frequency uint
    	how often the pause should be in minutes (default 20)
$ ./twenty-twenty-twenty # the defaults are recommended
```

## How to build

Needs Go 1.18+.

```
$ go build
# or
$ make
```

Also in macOS, it needs to be built with `gogio`:

```
$ go install gioui.org/cmd/gogio@0.4.0
$ make TwentyTwentyTwenty.app
```

[1]: https://www.allaboutvision.com/conditions/refractive-errors/what-is-20-20-20-rule/
[2]: https://pubmed.ncbi.nlm.nih.gov/36473088/
[3]: https://modernod.com/articles/2023-july-aug/myth-busting-the-202020-rule?c4src=article:infinite-scroll

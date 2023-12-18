# Twenty-Twenty-Twenty

20-20-20 is a program that alerts every 20 minutes to look something at 20 feet
away for 20 seconds, written in Golang. This is done to reduce eye fadigue [1].
While controversial[2][3], it is a simple rule to remember and can help long
screen usage sessions.


## How to use

```
$ ./twenty-twenty-twenty -help
Usage of ./twenty-twenty-twenty:
  -duration int
    	how long to show the notification in seconds (does not work in macOS) (default 20)
  -frequency int
    	how often to show the notification in minutes (default 20)
$ ./twenty-twenty-twenty # the defaults are recommended
```

## How to build

Needs Go 1.18+.

```
$ go build
```

## License

The code itself is licensed in the MIT license.

[beeep](https://github.com/gen2brain/beeep) is the only third-party
dependency, and it is licensed in [BSD
2-clause](https://github.com/gen2brain/beeep/blob/master/LICENSE).

The eye icon is from Font Awesome and is licensed in [CC-BY
4.0](https://creativecommons.org/licenses/by/4.0/).

[1]: https://www.allaboutvision.com/conditions/refractive-errors/what-is-20-20-20-rule/
[2]: https://pubmed.ncbi.nlm.nih.gov/36473088/
[3]: https://modernod.com/articles/2023-july-aug/myth-busting-the-202020-rule?c4src=article:infinite-scroll

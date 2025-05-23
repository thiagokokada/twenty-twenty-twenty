# Twenty-Twenty-Twenty

![Twenty-Twenty-Twenty screenshot](./assets/screenshot.png "Screenshot")

20-20-20 is a program that alerts every 20 minutes to look something at 20 feet
away for 20 seconds, written in Golang. This is done to reduce eye fadigue [1].
While controversial [2], it is a simple rule to remember and can help long
screen usage sessions.

## How to use

If you have Go installed, you can run:

```console
# doesn't work in macOS since it needs an signed app bundle
go run github.com/thiagokokada/twenty-twenty-twenty@latest
```
If you have [`nix`](https://nixos.org/download) installed, you can run:

```console
# works in macOS too
nix run github:thiagokokada/twenty-twenty-twenty
```

You can also download the pre-build binaries in
[releases](https://github.com/thiagokokada/twenty-twenty-twenty/releases).

Below are instructions on how to get the best experience in each OS.

### macOS

> [!WARNING]
> The releases available for macOS are signed with the adhoc certificate. They
> will not work as-is unless re-signed. If you're using macOS, the easiest way
> to get it working is compiling it yourself.

In macOS, just run `TwentyTwentyTwenty_<arch>.app`. I recommend moving
`TwentyTwentyTwenty_<arch>.app` to `~/Applications`

For a better experience, configure the notification type in "System Settings"
-> "Notifications" -> "TwentyTwentyTwenty" to use "Alerts". This will result in
the notification to stay until the end of the duration period instead of
disappearing in a few seconds.

### Windows

In Windows, just run `twenty-twenty-twenty-windows-<arch>.exe`.

You may want to disable the sounds either from the application (using the
system tray) or from Windows itself (you can disable sounds for notifications
from a specific application).

The second option is probably a better idea since the application play sounds
when the notification starts and ends, to remind you that you can look at the
screen again.

### Linux/*BSDs

Run in terminal:

```console
$ ./twenty-twenty-twenty -help
Usage of ./twenty-twenty-twenty:
  -disable-sound
        disable notification sound
  -duration uint
        how long each pause should be in seconds (default 20)
  -frequency uint
        how often the pause should be in seconds (default 1200)
  -pause uint
        how long the pause (from systray) should be in seconds (default 3600)
  -version
        print program version and exit
```

## How to build

> [!IMPORTANT]
> For macOS releases, you will need to install Xcode Command-line Tools first
> to sign the application:
> ```console
> xcode-select --install
> ```

Needs Go 1.24+.

```console
make
make run
```

### With Nix

If you have `nix` installed, you can also setup the development setup with it:

```console
nix develop
```

Or if you have [`nix-direnv`](https://github.com/nix-community/nix-direnv):

```console
direnv allow
```

### Optional features

You can enable/disable optional features by passing the following build tags
(e.g.: `go build -tags=<tag>`):

- `cgo`: [CGO](https://pkg.go.dev/cmd/cgo) is necessary for some features in
some platforms. It is required for `darwin`, not used in `windows`, and
enables/disables sound support in `linux` (and other *nix systems)
- `nosound`: disable sound support.
- `nosystray`: disable system tray support.

## Credits

[Notification sound 1](https://bigsoundbank.com/sound-1111-message-1.html) and
[Notification sound 2](https://bigsoundbank.com/sound-1112-message-2.html) by
Joseph SARDIN - https://bigsoundbank.com.

Eye open extracted from Font Awesome, in SVG format from
[Wikipedia](https://en.m.wikipedia.org/wiki/File:Eye_open_font_awesome.svg).

[1]: https://www.allaboutvision.com/conditions/refractive-errors/what-is-20-20-20-rule/
[2]: https://modernod.com/articles/2023-july-aug/myth-busting-the-202020-rule

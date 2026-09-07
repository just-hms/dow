<p align="center">
    <img style="width:8em;" src="./assets/logo.png" alt="jim">
</p>

# dow

Dow is a command-line tool designed to move the most recent download to a new location.

usage:

```sh
mv the last downloaded file in the current (or the specified) destination

Usage:
  dow [flags]

Examples:
  dow
  dow /path/to/target
  dow -v | xargs -rd '\n' code

Flags:
  -x, --extract   move and extract the last downloaded archive
  -h, --help      help for dow
  -v, --verbose   show the name of the moved file
  -y, --yes       force dow to move the latest file even if it's old
```

## Install

```sh
go install github.com/just-hms/dow@latest
```

## Problems with the download folder?

Your download folder is not in `~/Downloads`? Set up the `DOW_DOWNLOAD_PATH` environment variable to explicitly define the path to your downloads folder. This ensures Dow knows exactly where to look.

```sh
export DOW_DOWNLOAD_PATH=path/to/your/download-folder
```

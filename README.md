# dow

Dow is a command-line tool designed to move the most recent download to a new location.

usage:

```shell
dow <destination path>
```

## Install

```shell
go install github.com/just-hms/dow@latest
```

## Problems with the download folder?

Your folder is not in `~/Downloads`? Set up the DOW_DOWNLOAD_PATH environment variable to explicitly define the path to your downloads folder. This ensures Dow knows exactly where to look.

```shell
export DOWNLOAD_FOLDER_PATH=path/to/your/download-folder
```

# Rclone Permissions Mapper

Tool to convert `uid` and `gid` between mac and linux defaults when syncing files via cloud storage using [rclone](https://github.com/rclone/rclone) [`sync`](https://rclone.org/commands/rclone_sync/).

## Usage
```bash
rclone sync source:path dest:path --metadata-mapper /path/to/rclone-permissions-mapper
```

or, to see input and output:
```bash
rclone sync source:path dest:path --metadata-mapper /path/to/rclone-permissions-mapper -v --dump mapper
```

## Background
The default UID of the first regular user on macOS is `501`. On Linux, it is usually `1000`. If you [`rclone sync`](https://rclone.org/commands/rclone_sync/) a file created on one to the other (by way of a cloud storage remote) and use the [`--metadata`](https://rclone.org/docs/#m-metadata) flag (without using `sudo`), by default you will probably get an error like this one:

```text
ERROR : file.txt: Failed to copy: failed to set metadata: failed to change ownership: chown /testing/file.txt.fekayen6.partial: operation not permitted
```

This is because it is trying to `chown 1000:1000 /testing/file.txt` when actually it should be `501:20` (or vice versa.)

This tool uses rclone's new `--metadata-mapper` feature to automatically detect and correct this during the sync. It does so by simply omitting the `uid` and `gid` (when necessary) in the metadata blob it passes back to rclone, so that the default values are kept.

## Installation
[Download](https://github.com/nielash/rclone-permissions-mapper/releases) and unzip (or build from source with `go build`), and then move the executable to your `$PATH`:
```bash
sudo rclone moveto /Users/yourusername/Downloads/rclone-permissions-mapper-1.0-osx-arm64/rclone-permissions-mapper /usr/local/bin/rclone-permissions-mapper -v
```

Test if it's working:

```bash
echo '{"Metadata": {"hello": "world"}}' | rclone-permissions-mapper
```
should output: `{"Metadata":{"hello":"world"}}`

You can test what it will do by giving it different `uid` and `gid` values:
``` bash
echo '{"Metadata":{"gid":"20","uid":"501"}}' | rclone-permissions-mapper
// on mac: {"Metadata":{"gid":"20","uid":"501"}}
// on linux: {"Metadata":{}}

echo '{"Metadata":{"gid":"1000","uid":"1000"}}' | rclone-permissions-mapper
// on mac: {"Metadata":{}}
// on linux: {"Metadata":{"gid":"20","uid":"501"}}
```

## Resources
* [`--metadata-mapper` docs](https://rclone.org/docs/#metadata-mapper)
* [rclone's handling of `uid` and `gid`](https://github.com/rclone/rclone/blob/c69eb84573c85206ab028eda2987180e049ef2e4/backend/local/metadata.go#L113-L128)
* [Downloads](https://github.com/nielash/rclone-permissions-mapper/releases)
* [Rclone Forum](https://forum.rclone.org/)
# blart

![](http://www.wingclips.com/system/movie-clips/paul-blart-mall-cop/police-boot-camp/images/paul-blart-mall-cop-movie-clip-screenshot-police-boot-camp_large.jpg)

Watch files/directories for changes, and report to child process.

## Installation

```bash
$ go get github.com/mattrobenolt/blart
```

## Usage

```
usage: blart [flags] [command]
  -d=3s: time to wait after change before signalling child
  -f="": files and directories to watch, split by ':'
  -s="HUP": signal to send on change
```

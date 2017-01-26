go-selfupdate
=============

This is was based on work from [sanbornm/go-selfupdate](https://github.com/sanbornm/go-selfupdate) but with updates to make suite the needs of my software

[![GoDoc](https://godoc.org/github.com/jwhitcraft/go-selfupdate/selfupdate?status.svg)](https://godoc.org/github.com/jwhitcraft/go-selfupdate/selfupdate)
[![Build Status](https://travis-ci.org/jwhitcraft/go-selfupdate.svg?branch=master)](https://travis-ci.org/jwhitcraft/go-selfupdate)

Enable your Golang applications to self update.  Inspired by Chrome based on Heroku's [hk](https://github.com/heroku/hk).

## Features

* Tested on Mac, Linux, Arm, and Windows
* Creates binary diffs with [bsdiff](http://www.daemonology.net/bsdiff/) allowing small incremental updates
  * `brew install bsdiff` on Mac
* Falls back to full binary update if diff fails to match SHA

## QuickStart

### Enable your App to Self Update
```go
	var updater = &selfupdate.Updater{
		CurrentVersion: version,
		ApiURL:         "http://updates.yourdomain.com/",
		BinURL:         "http://updates.yourdomain.com/",
		DiffURL:        "http://updates.yourdomain.com/",
		Dir:            "update/",
		CmdName:        "myapp", // app name
	}
	
	if updater != nil {
		// check if an update is available first
		new_version, err := updater.HasUpdate()
		
		if new_version != "" {
			fmt.Printf("Upgrading to v%s\n", new_version)
			go updater.BackgroundRun()
		}
	}
```
### Push Out and Update

    go-selfupdate myapp 1.2

This will create a folder in your project called, *public* you can then rsync or transfer this to your webserver or S3.

If you are cross compiling you can specify a directory:

    go-selfupdate /tmp/mybinares/ 1.2

The directory should contain files with the name, $GOOS-$ARCH. Example:

    windows-386
    darwin-amd64
    linux-arm

If you are using [goxc](https://github.com/laher/goxc) you can output the files with this naming format by specifying this config:

    "OutPath": "{{.Dest}}{{.PS}}{{.Version}}{{.PS}}{{.Os}}-{{.Arch}}",

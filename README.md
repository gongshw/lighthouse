# lighthouse 

[![Build Status](https://travis-ci.org/gongshw/lighthouse.svg)](https://travis-ci.org/gongshw/lighthouse)

An amazing web proxy application!

## A demo online

just visit [https://lighthouse.gongshw.com/](https://lighthouse.gongshw.com/) and click the `Go!` button!

## Features

What can be proxied:

 - a single html page of course
 - multiple resources(script/image/css) from the page
 - all anchors from the page
 - and, most amazing, all ajax request!
 - cookies are not suppored for now :-(
 
## Install lighthouse

1. install golang

2. run `go get github.com/gongshw/lighthouse`

3. make a conf.json file at the work dir:

```json
{
    "StaicFileDir": "static",
    "ServerBaseUrl": "http(s)://your_server:8080/",
    "ServerPort": 8080
}

```

4. run `$GOPATH/bin/lighthouse`


5. visit `http(s)://your_server:8080/` from your restricted devices and enjoy!

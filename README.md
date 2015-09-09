# lighthouse 

[![Build Status](https://travis-ci.org/gongshw/lighthouse.svg)](https://travis-ci.org/gongshw/lighthouse)

install golang

run `go get github.com/gongshw/lighthouse`

make a conf.json file at the work dir:

```json
{
    "StaicFileDir": "static",
    "ServerBaseUrl": "http(s)://your_server:8080/",
    "ServerPort": 8080
}

```

run `$GOPATH/bin/lighthouse`


visit `http(s)://your_server:8080/` from your restricted devices and enjoy!

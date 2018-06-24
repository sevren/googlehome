# Go version Google Home Kodi

This is an attempt to create a golang version of the Google Home Kodi project

## Dependencies

Ideally you should just be able to use: to get all the dependencies used

```bash
dep init
dep ensure
```

Should you not be able to get them via dep, use the following list
``` bash
go get -u -v github.com/parnurzeal/gorequest
go get -u -v github.com/sirupsen/logrus
go get -u -v github.com/spf13/pflag
go get -u -v github.com/spf13/viper
```

## Configuration

This application takes in a configuration file with the following structure

```json
{
    "kodi" : {
        "id": "kodi",
        "protocol": "http",
        "ip": "192.168.1.111",
        "port": 8080,
        "user": "kodi-user",
        "pass": "kodi-pass"
    },
    "authtoken" : "Whatever-you-want",
    "listenerport" : 9998,
    "youtubeapikey" : "XXXXXXXXXXXXXXXX"
}
```

You can override this at the start of the application by passing in an environment variable for each property you wish to override

```bash
KODI_IP="x.x.x.x" AUTHTOKEN=xxxxx go run main.go 

```

## Docker Building and running

This app will be built in a alpine/armhf container. There is heavy use of the net/http package therefore we must force the build to be static and use netgo.


```bash
docker build -t sevren/googlehomekodi . 
```
```bash
docker run -d --network=ghome -p 9998:9998 -v $PWD/cfg.json:/app/cfg.json sevren/googlehomekodi 
```

## Standalone Building and Running

To build
```bash
go build main.go
```
To run
```bash
go run main.go
```

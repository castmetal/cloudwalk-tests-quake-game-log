# Cloudwalk Tests Quake Game Log

This script creates a quake log file reader, creating two reports named **Grouped Information** or **Kills By Means**, which the first report provides grouped information for each killer Player containing the number of kills, on the other hand, the second report provides per game all kills grouped by kill mod.

![](https://github.com/castmetal/cloudwalk-tests-quake-game-log/blob/main/render1691455196169.gif)

## Starting this script (local)

You'll need a go version higher than 19.1. It's been recommended 1.20+.

	> Before starting the script install dependencies and mods
	
## Installing dependencies

Install local go version: [link here](https://go.dev/dl/)

Run:

```sh
	go mod download
```
or
```sh
	go mod tidy
```

## Executing as dev mode


```sh
	cd cmd/reader_log_script
```
Run script:
```sh
	go run main.go
```

## Executing as production mode

```sh
	cd cmd/reader_log_script
```
Run script:
```sh
	go run main.go reader_log_script --execute=true
```

## Starting using Docker

Build the image:
```sh
	docker build -t cloudwalk-tests-quake-game-log:v1 .
```
After that, run:
```sh
	docker run -it cloudwalk-tests-quake-game-log:v1
```

## Running all tests
In the root directory:

```sh
	go test ./...
```

## Licence
- MIT
- Created by [castmetal](https://github.com/castmetal)
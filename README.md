# Synerex Simulation

Simulation Services for Person and Traffic Trip using synerex-alpha

# Introduction

Synerex alpha is an alpha version of Synergic Exchange and its support systems.
This project is supported by JST MIRAI.

## Requirements

go 1.11 or later (we use go.mod files for module dependencies)
nodejs(10.13.0) / npm(6.4.1) / yarn(1.12.1) for web client development.

## How to start

Starting from SynerexEngine.
Synerex Engine is a daemon tool for controlling Synerex.

```
  cd cli/daemon
  go build
```

Then move to provider directory and build provider.

```
  cd provider/simulation/scenario
  go build

  cd provider/simulation/car
  go build

  cd provider/simulation/pedestrian
  go build

  cd provider/simulation/area
  go build

  cd provider/simulation/visualization
  go build
```

## Source Directories

### cli

#### deamon

se-daemon for cli service
It can start all providers.

```
go build se-daemon.go se-daemon_[os].go
```

#### se

command line client for Synerex Engine

```
 go build   // build se command
 se clean all   // remove all binaries
 se run all     // start all servers and providers
 se stop all    // stop all servers and providers
 se ps -l       // list current running server and providers
```

#### api

Protocl Buffer / gRPC definition of Synergex API
If you changed any protocol, you should re-generate pb.go files.
To do so, you should change directory "server", and then

```
 go generate
```

synerex-server.go contains grpc protocl compile code.

#### server

Synerex Server alpha version

#### provider

Synerex Service Providers

##### ad

##### taxi

##### multi

##### user

##### fleet

##### map

##### datastore

##### ecotan

Local community bus system. (only for regional restricted demo)

#### sxutil

Synerex Utility Package Both server and provider package will
use this.

monitor Monitoring Server

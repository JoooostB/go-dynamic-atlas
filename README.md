#

```text
   _____  ____         _______     ___   _          __  __ _____ _____          _______ _                _____ 
  / ____|/ __ \       |  __ \ \   / / \ | |   /\   |  \/  |_   _/ ____|      /\|__   __| |        /\    / ____|
 | |  __| |  | |______| |  | \ \_/ /|  \| |  /  \  | \  / | | || |   ______ /  \  | |  | |       /  \  | (___  
 | | |_ | |  | |______| |  | |\   / | . ` | / /\ \ | |\/| | | || |  |______/ /\ \ | |  | |      / /\ \  \___ \ 
 | |__| | |__| |      | |__| | | |  | |\  |/ ____ \| |  | |_| || |____    / ____ \| |  | |____ / ____ \ ____) |
  \_____|\____/       |_____/  |_|  |_| \_/_/    \_\_|  |_|_____\_____|  /_/    \_\_|  |______/_/    \_\_____/
```

[![Go Report Card][5]][6]  
[![Apache licensed][3]][4] [![Docker][1]][2]  
[![Build Status][7]][8] [![Lint Status][9]][10] [![Release Status][11]][12]

[1]: https://img.shields.io/docker/image-size/joooostb/go-dynamic-atlas/latest
[2]: https://hub.docker.com/r/joooostb/go-dynamic-atlas
[3]: https://img.shields.io/badge/license-Apache-blue.svg
[4]: LICENSE
[5]: https://goreportcard.com/badge/github.com/joooostb/go-dynamic-atlas
[6]: https://goreportcard.com/report/github.com/joooostb/go-dynamic-atlas
[7]: https://img.shields.io/github/actions/workflow/status/joooostb/go-dynamic-atlas/docker.yml?label=docker%20build
[8]: https://github.com/JoooostB/go-dynamic-atlas/actions/workflows/docker.yml
[9]: https://img.shields.io/github/actions/workflow/status/joooostb/go-dynamic-atlas/go.yml?label=ci
[10]: https://github.com/JoooostB/go-dynamic-atlas/actions/workflows/go.yml
[11]: https://img.shields.io/github/actions/workflow/status/joooostb/go-dynamic-atlas/release.yml?label=release
[12]: https://github.com/JoooostB/go-dynamic-atlas/actions/workflows/release.yml

[Go Dynamic Atlas](https://github.com/joooostb/go-dynamic-atlas) is a simple server that listens for [go-dynamic-atlas](https://github.com/joooostb/go-dynamic-atlas) webhook events in order to update MongoDB Atlas with the latest dynamic IP address of your server. This ensures that your MongoDB Atlas cluster is always accessible from your server, while keeping the attack surface as small as possible by only allowing access from your server's IP address.

## Required Access

To manage IP Access List entries, you must have Project Owner access to the project that contains the cluster you want to modify. Users with Organization Owner access must add themselves to the project as a Project Owner.

## Prerequisites

Make sure the following environment variables are set:

- `MONGODB_ATLAS_PUBLIC_KEY`: Your MongoDB Atlas public key
- `MONGODB_ATLAS_PRIVATE_KEY`: Your MongoDB Atlas private key
- `MONGODB_ATLAS_PROJECT_ID`: Your MongoDB Atlas project ID
- `GIN_MODE`(optional): The GIN environment in which the server is running (e.g. `debug`, `release`)

## Features

- [x] Listen for go-dynamic-atlas webhook events
- [x] Update MongoDB Atlas IP Access List with the latest IP address
- [x] Automatically remove old IP addresses from the IP Access List
- [x] Docker support
- [x] Kubernetes support

## Usage

Setup a webhook with POST request to `http://<server-ip>:8080/api/v1/updateIP` and the following JSON body:

```json
{
  "ip": "<ip-address>", 
  "comment": "<comment>"
}
```

### Webhook with HTTP POST request in go-dynamic-atlas

```yml
"webhook": {
  "enabled": true,
  "url": "http://<server-ip>:8080/api/v1/updateIP",
  "request_body": "{ \"ip\": \"{{.CurrentIP}}\", \"comment\": \"Updated from go-dynamic-atlas.\" }"
}
```

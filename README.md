<div align="center">
  <img src="/docs/img/chopchop_logo.png" width="180" height="150"/>
</div>

[![Build Status](https://github.com/michelin/ChopChop/workflows/Build%20ChopChop/badge.svg)](https://github.com/michelin/ChopChop/actions)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/michelin/ChopChop)](https://goreportcard.com/report/github.com/michelin/ChopChop)

# ChopChop

**ChopChop** is a command-line tool for dynamic application security testing on web applications, initially written by the Michelin CERT.

Its goal is to scan several endpoints and identify exposition of services/files/folders through the webroot.
Checks/Signatures are declared in a config file (by default: `chopchop.yml`), fully configurable, and especially by developers.

<div align="center">
  <img src="/docs/img/demo.gif?raw=true"/>
</div>

> "Chop chop" is a phrase rooted in Cantonese. "Chop chop" means "hurry" and suggests that something should be done now and **without delay**.

---

## Table of Contents

* [Building](#building)
* [Usage](#usage)
  * [Available flags](#available-flags)
  * [Advanced usage](#advanced-usage)
* [Creating a new check/signature](#creating-a-new-check)
* [External Libraries](#external-libraries)
* [Talks](#talks)
* [Licence](#licence)
* [Authors](#authors)

## Building

We tried to make the build process painless and hopefully, it should be as easy as: 

```bash
go build -o gochopchop cmd/main.go
```

There should be a resulting `gochopchop` binary in the folder.

### Using Docker

Thanks to [Github Container Registry](https://github.blog/2020-09-01-introducing-github-container-registry/), we are able to provide you some freshly-build Docker images!

```
docker run ghcr.io/michelin/gochopchop scan -v debug https://example.com
```

But if you prefer, you can also build it locally, see below: 

#### Build locally

```bash
docker build -t gochopchop .
```

## Usage

We are continuously trying to make `gochopchop` as easy as possible. Scanning a host with this utility is as simple as: 

```bash
./gochopchop scan https://example.com
```

Notice you can specify multiple URLs.

### Using Docker

```bash
docker run gochopchop scan https://example.com
```

Notice by default the Docker image has the configuration file at
`/etc/chopchop.yml`, so you may add `-c /etc/chopchop.yml` to your
command. If so, run the following command.

```bash
docker run gochopchop scan -c /etc/chopchop.yml https://example.com
```

#### Custom configuration file

Of course you can use your own configuration files, using the following.

```bash
docker run -v $(pwd):/app gochopchop scan -c /app/chopchop.yml https://example.com
```

## What's next

The Golang rewrite took place a couple of months ago but there's so much to do, still. Here are some features we are planning to integrate :
 - [ ] Improve logging
 - [ ] HTTP & SOCKS5 proxies
 - [ ] Plugin method (GET, POST, ...)
 - [ ] Plugin Cookie
 - [ ] Improve caching (Docker build & CI)
 - [ ] Implement a request gateway to avoid brute-forcing websites
 - [ ] Re-implement `query_string` for the HTTP GET method
 - [ ] Improve "severity reached" cases (currently ChopChop crashes if matches a plugin)
 - [ ] Fix default status_code (200 if not specified)

## Testing

### Unit tests

Unit tests are achieved using Go-tests. You can run them using the following.

```bash
go test ./... -cover
```

To visualize the code coverage, for developement purposes, you can also run.

```bash
go test ./... -coverprofile=cov.out -count=1 && go tool cover -html=cov.out && rm cov.out 
```

### Acceptance tests

For acceptance tests, those are achieved using RobotFramework. Notice you can
build ChopChop using another language and validate the CLI using those tests.
Run them using the following.

```bash
cd robot
./run.sh
```

## Available flags

You can find the available flags and doc for each command using `gochopchop [cmd] -h`.

Available commands are:
 - `scan` to scan for endpoints ;
 - `plugins` to parse and check the configuration file.

## Advanced usage

Here is a list of advanced usage that you might be interested in.
Note: Redirectors like `>` for post processing can be used.

- Ability to scan and disable SSL verification

```bash
./gochopchop scan --insecure https://foobar.com
```

- Ability to scan with a custom configuration file (including custom plugins)

```bash
./gochopchop scan --insecure --signature test_config.yml https://foobar.com
```

- Ability to specify number of concurrent threads (in Go those are goroutines): `--threads 4` for 4 workers

```bash
./gochopchop scan --threads 4 https://foobar.com
```

- Ability to specify specific signatures to be checked, with a debug log level

```bash
./gochopchop scan --timeout=1 --verbosity=debug --export=csv --export=json --export-filename=boo --plugin-filters=Git,Zimbra,Jenkins https://foobar.com
```

- Set a list or URLs located in a file

```bash
./gochopchop scan --url-file url_file.txt
```

- Export GoChopChop results in CSV and JSON format

```bash
./gochopchop scan https://foobar.com  --export csv --export json --export-filename results
```

- Ability to list all the plugins

```bash
./gochopchop plugins
```

- Ability to list all the plugins or by severity : `plugins` or  `plugins --severity High`

```bash
./gochopchop plugins --severity High
```

## Creating a new check

Writing a new check is as simple as : 

```yaml
  - endpoints:
      - "/.git/config"
    checks:
      - name: Git exposed
        match:
          - "[branch"
        remediation: Do not deploy .git folder on production servers
        description: Verifies that the GIT repository is accessible from the site
        severity: High
```

An endpoint (e.g. `/.git/config`) is mapped to multiple checks which avoids
sending X requests for X checks. Multiple checks are achieved through a
single HTTP request.
Each check needs those fields:

| Attribute | Type | Description | Optional ? | Example | 
|---|---|---|---|---|
| name | string | Name of the check | No | Git exposed |
| description | string | A small description for the check| No |  Ensure .git repository is not accessible from the webroot |
| remediation | string | Give a remediation for this specific "issue" | No | Do not deploy .git folder on production servers |
| severity | Enum("High", "Medium", "Low", "Informational") | Rate the criticity if it triggers in your environment| No | High |
| status_code | integer | The HTTP status code that should be returned | Yes | 200 |
| headers | List of string | List of headers there should be in the HTTP response | Yes | N/A |
| no_headers | List of string | List of headers there should NOT be in the HTTP response | Yes | N/A |
| match | List of string| List the strings there should be in the HTTP response  | Yes |  "[branch" |
| no_match | List of string | List the strings there should NOT be in the HTTP response | Yes | N/A |

## External Libraries

| Library Name | Link                                  | License              | 
|--------------|---------------------------------------|----------------------|
| go-md2man    | https://github.com/cpuguy83/go-md2man | MIT License          | 
| strfmt       | https://github.com/go-openapi/strfmt  | Apache License 2.0   |
| go-cmp       | https://github.com/google/go-cmp      | BSD-3-Clause License |
| go-pretty    | https://github.com/jedib0t/go-pretty  | MIT License          |
| go-runewidth | https://github.com/mattn/go-runewidth | MIT License          |
| logrus       | https://github.com/sirupsen/logrus    | MIT License          |
| cli          | https://github.com/urfave/cli/v2      | MIT License          |
| yaml         | https://github.com/go-yaml/yaml       | Apache License 2.0   |

## Talks

- PyCon FR 2019 (The tool was initially developed in Python) - https://docs.google.com/presentation/d/1uVXGUpt7tC7zQ1HWegoBbEg2LHamABIqfDfiD9MWsD8/edit
- DEFCON AppSec Village 2020 "Turning offsec mindset to developer's toolset" - https://drive.google.com/file/d/15P8eSarIohwCVW-tR3FN78KJPGbpAtR1/view

## License

ChopChop has been released under Apache License 2.0. 
Please, refer to the `LICENSE` file for further information.

## Authors

- Paul A. 
- David R. (For the Python version)
- Stanislas M. (For the Golang version)

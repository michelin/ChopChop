<p align="center"><img src="/docs/img/chopchop_logo.png" width="180" height="150"/></p>

[![Build Status](https://github.com/michelin/ChopChop/workflows/Build%20ChopChop/badge.svg)](https://github.com/michelin/ChopChop/actions)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/michelin/ChopChop)](https://goreportcard.com/report/github.com/michelin/ChopChop)

# ChopChop

**ChopChop** is a command-line tool for dynamic application security testing on web applications, initially written by the Michelin CERT.

Its goal is to scan several endpoints and identify exposition of services/files/folders through the webroot.
Checks/Signatures are declared in a config file (by default: `chopchop.yml`), fully configurable, and especially by developers.

<p align="center"><img src="/docs/img/demo.gif?raw=true"/></p>

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
$ go mod download
$ go build .
```

There should be a resulting `gochopchop` binary in the folder.

### Using Docker

Thanks to [Github Container Registry](https://github.blog/2020-09-01-introducing-github-container-registry/), we are able to provide you some freshly-build Docker images!

```
docker run ghcr.io/michelin/gochopchop scan https://foobar.com -v debug
```

But if you prefer, you can also build it locally, see below: 

#### Build locally

```bash
docker build -t gochopchop .
```

## Usage

We are continuously trying to make `goChopChop` as easy as possible. Scanning a host with this utility is as simple as : 

```bash
$ ./gochopchop scan https://foobar.com
```

### Using Docker

```bash
docker run gochopchop scan https://foobar.com
```

#### Custom configuration file

```bash
docker run -v ./:/app chopchop scan -c /app/chopchop.yml https://foobar.com
```

## What's next

The Golang rewrite took place a couple of months ago but there's so much to do, still. Here are some features we are planning to integrate :
[x] Threading for better performance
[x] Ability to specify the number of concurrent threads
[x] Colors and better formatting
[x] Ability to filter checks/signatures to search for
[x] Mock and unit tests
[x] Github CI
And much more!

## Testing

To launch a functionnal test, you have an example of web server available under `tests/server.go` that will launch a web server listening on http://localhost:8000.
You can then launch a scan on this dummy web server : `./gochopchop scan http://localhost:8000 --verbosity Debug`.

There are also unit test that you can launch with `go test -v ./...`.
These tests are integrated in the github CI workflow.

## Available flags

You can find the available flags available for the `scan` command :

| Flag | Full flag | Description |
|---|---|---|
| `-h` | `--help` | Help wizard |
| `-v` | `--verbosity` | Verbose level of logging |
| `-c` | `--signature` | Path of custom signature file |
| `-k` | `--insecure` | Disable SSL Verification |
| `-u` | `--url-file` | Path to a specified file containing urls to test |
| `-b` | `--max-severity` | Block the CI pipeline if severity is over or equal specified flag |
| `-e` | `--export` | Export type of the output (csv and/or json) |
|| `--export-filename` | Specify the filename for the export file(s) |
| `-t` | `--timeout` | Timeout for the HTTP requests |
|| `--severity-filter` | Filter Plugins by severity |
|| `--plugin-filter` | Filter Plugins by name of plugin |
|| `--threads` | Number of concurrent threads | 

## Advanced usage

Here is a list of advanced usage that you might be interested in.
Note: Redirectors like `>` for post processing can be used.

- Ability to scan and disable SSL verification

```bash
$ ./gochopchop scan https://foobar.com --insecure
```

- Ability to scan with a custom configuration file (including custom plugins)

```bash
$ ./gochopchop scan https://foobar.com --insecure --signature test_config.yml
```

- Ability to list all the plugins or by severity : `plugins` or  ` plugins --severity High`

```bash
$ ./gochopchop plugins --severity High
```

- Ability to specify number of concurrent threads : `--threads 4` for 4 workers

```bash
$ ./gochopchop plugins --threads 4
```

- Ability to block the CI pipeline by severity level (equal or over specified severity) : `--max-severity Medium`

```bash
$ ./gochopchop scan https://foobar.com --max-severity Medium
```

- Ability to specify specific signatures to be checked 

```bash
./gochopchop scan https://foobar.com --timeout 1 --verbosity --export=csv,json --export-filename boo --plugin-filters=Git,Zimbra,Jenkins
```

- Ability to list all the plugins

```bash
$ ./gochopchop plugins
```

- List High severity plugins

```bash
$ ./gochopchop plugins --severity High
```

- Set a list or URLs located in a file

```bash
$ ./gochopchop scan --url-file url_file.txt
```

- Export GoChopChop results in CSV and JSON format

```bash
$ ./gochopchop scan https://foobar.com  --export=csv,json --export-filename results
```

## Creating a new check

Writing a new check is as simple as : 

```yaml
  - endpoint: "/.git/config"
    checks:
      - name: Git exposed
        match:
          - "[branch"
        remediation: Do not deploy .git folder on production servers
        description: Verifies that the GIT repository is accessible from the site
        severity: "High"
```

An endpoint (eg. ```/.git/config```) is mapped to multiple checks which avoids sending X requests for X checks. Multiple checks can be done through a single HTTP request.
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
| query_string | GET parameters that have to be passed to the endpoint | String | Yes | `query_string: "id=FOO-chopchoptest"` |

## External Libraries

| Library Name | Link | License | 
|---|---|---|
| Viper | https://github.com/spf13/viper | MIT License |
| Go-pretty |  https://github.com/jedib0t/go-pretty| MIT License |
| Cobra | https://github.com/spf13/cobra| Apache License 2.0 |
| strfmt |https://github.com/go-openapi/strfmt | Apache License 2.0 |
| Go-homedir | https://github.com/mitchellh/go-homedir| MIT License |
| pkg-errors | https://github.com/pkg/errors| BSD 2 (Simplified License)|
| Go-runewidth | https://github.com/mattn/go-runewidth | MIT License |

Please, refer to the `third-party.txt` file for further information.

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

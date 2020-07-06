<p align="center"><img src="/img/chopchop_logo.png" width="180" height="110"/></p>

[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/michelin/ChopChop)](https://goreportcard.com/report/github.com/michelin/ChopChop)

# ChopChop

**ChopChop** is a command-line tool for dynamic application security testing on web applications, initially written by the Michelin CERT.

Its goal is to scan several endpoints and identify exposition of services/files/folders through the webroot.
Checks/Signatures are declared in a config file (by default: `chopchop.yml`), fully configurable, and especially by developers.

<p align="center"><img src="/img/demo.gif?raw=true"/></p>

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

## Usage

We are continuously trying to make `goChopChop` as easy as possible. Scanning a host with this utility is as simple as : 

```bash
$ ./gochopchop scan --url https://foobar.com
```

## What's next

The Golang rewrite took place a couple of months ago but there's so much to do, still. Here are some features we are planning to integrate :
* Threading for better performance
* Colors and better formatting
* Ability to filter checks/signatures to search for
* And much more!


## Available flags

You can find the available flags here :

| Flag | Full flag | Description |
|---|---|---|
| `-h` | `--help` | Help wizard |
| `-u` | `--url` | Set the target URL |
| `-i` | `--insecure` | Disable SSL Verification |
| `-c` | `--config-file` | Set a custom configuration file |
| `-f` | `--url-file` | Set a file containing a list of URLs |
| | `--csv` | Export results in CSV | 
| | `--json` | Export results in JSON | 

## Advanced usage

Here is a list of advanced usage that you might be interested in.
Note: Redirectors like `>` for post processing can be used.

- Ability to scan and disable SSL verification

```bash
$ ./gochopchop scan --url https://foobar.com --insecure
```

- Ability to scan with a custom configuration file (including custom plugins)

```bash
$ ./gochopchop scan --url https://foobar.com --insecure --config-file test_config.yml
```

- Ability to list all the plugins or by severity : `plugins` or  ` plugins --severity High`

```bash
$ ./gochopchop plugins --severity High
```

- Ability to block the CI pipeline by severity level (equal or over specified severity) : `--block Medium`

```bash
$ ./gochopchop scan --url https://foobar.com --insecure --block Medium
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
$ ./gochopchop  scan --url https://foobar.com --json --csv 
```

## Creating a new check

Writing a new check is as simple as : 

```yaml
  - uri: "/.git/config"
    checks:
      - name: Git exposed
        match:
          - "[branch"
        remediation: Do not deploy .git folder on production servers
        description: Verifies that the GIT repository is accessible from the site
        severity: "High"
```

An URI (eg. ```/.git/config```) is mapped to multiple checks which avoids sending X requests for X checks. Multiple checks can be done through a single HTTP request.
Each check needs those fields:

| Attribute | Type | Description | Optional ? | Example | 
|---|---|---|---|---|
| name | string | Name of the check | No | Git exposed |
| description | string | A small description for the check| No |  Ensure .git repository is not accessible from the webroot |
| remediation | string | Give a remediation for this specific "issue" | No | Do not deploy .git folder on production servers |
| severity | Enum("High", "Medium", "Low", "Informational") | Rate the criticity if it triggers in your environment| No | High |
| status_code | integer | The HTTP status code that should be returned |
| headers | List of string | List of headers there should be in the HTTP response | Yes | N/A |
| match | List of string| List the strings there should be in the HTTP response  | Yes |  "[branch" |
| no_match | List of string | List the strings there should NOT be in the HTTP response | Yes | N/A |

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

## License

ChopChop has been released under Apache License 2.0. 
Please, refer to the `LICENSE` file for further information.

## Authors

- Paul A. 
- David R. (For the Python version)
- Stanislas M. (For the Golang version)

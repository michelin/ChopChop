*** Settings ***
Documentation  A test suite to validate ChopChop CLI for the
...            plugins command.

Resource       ../data/common.resource

Library        ChopChop
Library        OperatingSystem

*** Variables ***
${SIGNATURES_RESULTS}=  SEPARATOR=\n
...  +-----+-------------+---------------+-------------+
...  | URL | PLUGIN NAME | SEVERITY${DS}${DS}${DS}| DESCRIPTION |
...  +-----+-------------+---------------+-------------+
...  | [/] | EXAMPLE${DS}${DS} | Informational | Description |
...  +-----+-------------+---------------+-------------+
...  |${DS}${DS} |${DS}${DS}${DS}${DS}${DS}${DS} | TOTAL CHECKS${DS}| 1${DS}${DS}${DS}${DS}${DS} |
...  +-----+-------------+---------------+-------------+
...

*** Test Cases ***
Simple Plugins
    [Documentation]    This test ensures ChopChop can display plugins.
    [Tags]    server
    [Setup]    Create File    ${SIGNATURES_FILENAME}    ${SIGNATURES}
    [Teardown]    Remove File    ${SIGNATURES_FILENAME}

    Chopchop Plugins    ${SIGNATURES_FILENAME}

    ${rc}=    Chopchop Get Returncode
    ${stdout}=    Chopchop Get Stdout

    Should Be Equal As Integers    ${rc}    0
    Should Be Equal As Strings   ${stdout}    ${SIGNATURES_RESULTS}

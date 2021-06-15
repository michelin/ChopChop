*** Settings ***
Documentation  A test suite to validate ChopChop CLI for the
...            scan command.

Resource       ../data/common.resource

Library        ChopChop
Library        MockServer
Library        OperatingSystem

*** Variables ***
${PORT}=  8080
${URL_FILE}=  url-file.txt

${SIGNATURES_RESULTS}=  SEPARATOR=\n
...  +-----------------------+----------+---------------+---------+-------------+
...  | URL${DS}${DS}${DS}${DS}${DS}${DS}${DS}${DS}${DS} | ENDPOINT | SEVERITY${DS}${DS}${DS}| PLUGIN${DS}| REMEDIATION |
...  +-----------------------+----------+---------------+---------+-------------+
...  | http://127.0.0.1:${PORT} | /${DS}${DS}${DS}${DS}| Informational | EXAMPLE | Remediation |
...  +-----------------------+----------+---------------+---------+-------------+
...

${DOUBLE_SIGNATURES_RESULTS}=  SEPARATOR=\n
...  +-----------------------+----------+---------------+-----------+-------------+
...  | URL${DS}${DS}${DS}${DS}${DS}${DS}${DS}${DS}${DS} | ENDPOINT | SEVERITY${DS}${DS}${DS}| PLUGIN${DS}${DS}| REMEDIATION |
...  +-----------------------+----------+---------------+-----------+-------------+
...  | http://127.0.0.1:${PORT} | /${DS}${DS}${DS}${DS}| Informational | EXAMPLE${DS} | Remediation |
...  | http://127.0.0.1:${PORT} | /2${DS}${DS}${DS} | Informational | EXAMPLE 2 | Remediation |
...  +-----------------------+----------+---------------+-----------+-------------+
...

*** Test Cases ***
Simple Scan
    [Documentation]    This test ensures the basic job of ChopChop: scanning endpoints.
    [Tags]    server
    [Setup]    Setup Server And Signatures File    ${PORT}    ${SIGNATURES_FILENAME}    ${SIGNATURES}
    [Teardown]    Teardown Server And Signatures File    ${SIGNATURES_FILENAME}

    Chopchop Scan    ${SIGNATURES_FILENAME}    http://127.0.0.1:${PORT}

    ${rc}=    Chopchop Get Returncode
    ${stdout}=    Chopchop Get Stdout

    Should Be Equal As Integers    ${rc}    0
    Should Be Equal As Strings   ${stdout}    ${SIGNATURES_RESULTS}

Scan With No Finding
    [Documentation]    This test ensures ChopChop fails in case there is no match.
    [Tags]    server
    [Setup]    Setup Server And Signatures File    ${PORT}    ${SIGNATURES_FILENAME}    ${FAILING_SIGNATURES}
    [Teardown]    Teardown Server And Signatures File    ${SIGNATURES_FILENAME}

    Chopchop Scan    ${SIGNATURES_FILENAME}    http://127.0.0.1:${PORT}

    ${rc}=    Chopchop Get Returncode

    Should Be Equal As Integers    ${rc}    1

Parallel Scan
    [Documentation]    This test ensures ChopChop works with multiple threads (or
    ...                concurrent routines).
    [Tags]    threads    server
    [Setup]    Setup Server And Signatures File    ${PORT}    ${SIGNATURES_FILENAME}    ${DOUBLE_SIGNATURES}
    [Teardown]    Teardown Server And Signatures File    ${SIGNATURES_FILENAME}

    Chopchop Scan    ${SIGNATURES_FILENAME}    http://127.0.0.1:${PORT}    2

    ${rc}=    Chopchop Get Returncode
    ${stdout}=    Chopchop Get Stdout

    Should Be Equal As Integers    ${rc}    0
    Should Be Equal As Strings   ${stdout}    ${DOUBLE_SIGNATURES_RESULTS}

Scan With Url From File
    [Documentation]   This test ensures ChopChop works with url from file.
    [Tags]   server
    [Setup]   Setup Server Signatures And Url Files    ${PORT}    ${SIGNATURES_FILENAME}    ${SIGNATURES}    ${URL_FILE}    http://127.0.0.1:${PORT}
    [Teardown]    Teardown Server Signatures And Url Files    ${SIGNATURES_FILENAME}    ${URL_FILE}

    Chopchop Scan Url File    ${SIGNATURES_FILENAME}    ${URL_FILE}

    ${rc}=    Chopchop Get Returncode
    ${stdout}=    Chopchop Get Stdout

    Should Be Equal As Integers    ${rc}    0
    Should Be Equal As Strings   ${stdout}    ${SIGNATURES_RESULTS}

Invalid Signatures File
    [Documentation]   This test ensures it crashes if a given signatures file is
    ...               invalid.
    [Tags]    no-server

    Chopchop Scan    ${SIGNATURES_FILENAME}    http://127.0.0.1:${PORT}

    ${rc}=    Chopchop Get Returncode

    Should Be Equal As Integers    ${rc}    1

No URL provided
    [Documentation]   This test ensures it crashes if a given signatures file is
    ...               invalid.
    [Tags]    no-server
    
    Chopchop Scan    ${SIGNATURES_FILENAME}    /

    ${rc}=    Chopchop Get Returncode

    Should Be Equal As Integers    ${rc}    1

*** Keywords ***
Setup Server And Signatures File
    [Arguments]    ${port}    ${sign_filename}    ${sign}
    Start Mock Server    ${port}
    Create File    ${sign_filename}    ${sign}

Setup Server Signatures And Url Files
    [Arguments]    ${port}    ${sign_filename}    ${sign}    ${url_filename}    ${urls}
    Setup Server And Signatures File    ${port}    ${sign_filename}    ${sign}
    Create File    ${url_filename}    ${urls}

Teardown Server And Signatures File
    [Arguments]    ${sign_filename}
    Remove File    ${sign_filename}
    Stop Mock Server

Teardown Server Signatures And Url Files
    [Arguments]    ${sign_filename}    ${url_filename}
    Remove File    ${url_filename}
    Teardown Server And Signatures File    ${sign_filename}

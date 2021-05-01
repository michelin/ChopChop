*** Variables ***
${SIGNATURES_FILENAME}=  chopchop.yml
${DS}=  ${SPACE * 2}

${SIGNATURES}=  SEPARATOR=\n
...  insecure: false
...  plugins:
...  ${DS}- endpoints:
...  ${DS}${DS}${DS}- "/"
...  ${DS}${DS}checks:
...  ${DS}${DS}${DS}- name: EXAMPLE
...  ${DS}${DS}${DS}${DS}match:
...  ${DS}${DS}${DS}${DS}${DS}- 'Example'
...  ${DS}${DS}${DS}${DS}remediation: Remediation
...  ${DS}${DS}${DS}${DS}description: Description
...  ${DS}${DS}${DS}${DS}severity: Informational
...  ${DS}${DS}${DS}${DS}status_code: 200

${DOUBLE_SIGNATURES}=  SEPARATOR=\n
...  ${SIGNATURES}
...  ${DS}- endpoints:
...  ${DS}${DS}${DS}- "/2"
...  ${DS}${DS}checks:
...  ${DS}${DS}${DS}- name: EXAMPLE 2
...  ${DS}${DS}${DS}${DS}match:
...  ${DS}${DS}${DS}${DS}${DS}- 'Example'
...  ${DS}${DS}${DS}${DS}remediation: Remediation
...  ${DS}${DS}${DS}${DS}description: Description
...  ${DS}${DS}${DS}${DS}severity: Informational
...  ${DS}${DS}${DS}${DS}status_code: 200

${FAILING_SIGNATURES}=  SEPARATOR=\n
...  insecure: false
...  plugins:
...  ${DS}- endpoints:
...  ${DS}${DS}${DS}- "/"
...  ${DS}${DS}checks:
...  ${DS}${DS}${DS}- name: EXAMPLE
...  ${DS}${DS}${DS}${DS}match:
...  ${DS}${DS}${DS}${DS}${DS}- 'Example'
...  ${DS}${DS}${DS}${DS}remediation: Remediation
...  ${DS}${DS}${DS}${DS}description: Description
...  ${DS}${DS}${DS}${DS}severity: Informational
...  ${DS}${DS}${DS}${DS}status_code: 400

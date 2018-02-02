## google-id-token

A small application to retrieve or verify google ID token

[![Build Status](https://travis-ci.org/grepplabs/google-id-token.svg?branch=master)](https://travis-ci.org/grepplabs/google-id-token)


### Requirements

- gcloud Command-Line Tool to retrieve Application Default Credentials - [Installing Cloud SDK](https://cloud.google.com/sdk/downloads)

### Install binary release

1. Download the latest release

   Linux

        curl -Lso google-id-token https://github.com/grepplabs/google-id-token/releases/download/v0.0.1/linux.amd64.google-id-token 

   macOS

        curl -Lso google-id-token https://github.com/grepplabs/google-id-token/releases/download/v0.0.1/darwin.amd64.google-id-token 

2. Make the google-id-token binary executable

    ```
    chmod +x ./google-id-token
    ```

3. Move the binary in to your PATH.

    ```
    sudo mv ./google-id-token /usr/local/bin/google-id-token
    ```

### Building

	go build

### Help output

    Retrieve or verify google ID token

    Usage:
      google-id-token [flags]
      google-id-token [command]

    Available Commands:
      get         get cached or new token
      help        Help about any command
      new         get a new token
      print       print the decoded token
      verify      verify the token

    Flags:
          --client-id string   Client ID (optional)
      -h, --help               help for google-id-token
      -t, --timeout int        Timeout in seconds (default 5)


### Acquire new user credentials and google ID token
[Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login)

    gcloud auth application-default login
    google-id-token new
    google-id-token verify --id-token="$(google-id-token get)"

## google-id-token

A small application to retrieve or verify google ID token

[![Build Status](https://travis-ci.org/grepplabs/google-id-token.svg?branch=master)](https://travis-ci.org/grepplabs/google-id-token)

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


### Acquire new user credentials and oogle ID token
[Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login)

    gcloud auth application-default login
    google-id-token new
    google-id-token verify --id-token="$(google-id-token get)"

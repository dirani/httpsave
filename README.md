# httpsave

## Build

    git clone https://github.com/udhos/httpsave
    cd httpsave
    export CGO_ENABLED=0
    go install ./httpsave

## Run

If `~/go/bin` is in the PATH env var:

    httpsave

If `~/go/bin` is NOT in the PATH env var:

    ~/go/bin/httpsave

## Test

    curl -v --data-binary @/etc/passwd localhost:8080/save
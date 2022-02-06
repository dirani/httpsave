# httpsave

## Build

    git clone https://github.com/udhos/httpsave
    cd httpsave
    export CGO_ENABLED=0
    go install ./httpsave

## Test

    curl -v --data-binary @/etc/passwd localhost:8080/save
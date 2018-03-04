# sta
sta (Simple Text-Adventure) is a text-based multiplayer game. It is an ssh server written in go.

## Installation
```bash
# download
go get github.com/ribacq/sta
cd $(go env GOPATH)/src/github/ribacq/sta

# generate ssh key for the server
ssh-keygen -C 'sta' -f ./id_rsa -P ''

# install
go install
```

## Running
Launch the server in a first terminal:
```bash
cd $(go env GOPATH)/src/github/ribacq/sta
go build
./sta
```

And connect as a client with a second terminal:
```
ssh -p 2222 username@localhost
```

## Acknowledgements

* https://github.com/gliderlabs/ssh for the wrapper of https://golang.org/x/crypto/ssh

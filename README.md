# sta
sta (Simple Text-Adventure) is a text-based multiplayer game. It is an ssh server written in go.

## Installation
```bash
# download
go get github.com/ribacq/sta

# generate ssh key for the server
ssh-keygen -C 'sta' -f $(go env GOPATH)/src/github.com/ribacq/sta/id_rsa -P ''

# install
go install github.com/ribacq/sta
```

# Running
Launch the server in a first terminal:
```bash
sta
```

And connect as a client with a second terminal:
```
ssh -p 2222 localhost
```

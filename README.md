
# golang tcp chat

Just learning and experimenting with go.

Very simple tcp chat app. Experimenting with net package and learning go.

## Building
Super flat project

```
go build talk.go
```

## Running
Usage

```
talk <local_listen_port> <remote_host:remote_port>
```

To run both sides on the same computer do the following

```
talk 5001 localhost:5002
```

```
talk 5002 localhost:5001
```


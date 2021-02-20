# TCP And UDP Server and client for file transfer
A file transfer client and server written using both TCP and UDP protocols using Go

## UDP File Transfer Time
![](./UDPFileTransferTiming.png)

## TCP File Transfer Time
![](./TCPTiming.png)


### Installation
1. Install [go](https://golang.org/doc/install)
2. Go the project directory and for compiling the server 
```
cd TCPServer
go build TCPFileServer.go
```
3. For compiling the client 
```
cd TCPClient
go build TCPFileClient.go
```
Same steps for UDP


### Running
If you don't want to compile or install go the executable files are also attached with project for both the client and the server in both the TCP and UDP implementations.For TCP the server can be started using
```
./TCPServer/TCPFileServer
```

For starting an instance of the client

```
./TCPClient/TCPFileClient
```
Multiple clients maybe started depending on the number of ports configured in the server
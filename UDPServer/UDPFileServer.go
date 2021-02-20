package main

import (
	"os"
	// "bytes"
	"net"
	// "math"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

const (
	chunkSize  int    = 20
	serverLocation string = "server.txt"
)
// InitRequest is..
type InitRequest struct {
	fileName     string
	maxChunkSize int
	fileSize     int
	totalChunks  int
}
// Request is..
type Request struct {
	// Headers
	SeqNumber int
	ChunkSize int
}
// JSONInitRequest is..
type JSONInitRequest struct {
	FileName     string `json:"fileName"`
	MaxChunkSize int    `json:"maxChunkSize"`
	FileSize     int    `json:"fileSize"`
	TotalChunks  int    `json:"totalChunks"`
}
// JSONRequest is..
type JSONRequest struct {
	SeqNumber int    `json:"seqNumber"`
	ChunkSize int    `json:"chunkSize"`
	Data      []byte `json:"data"`
}

func main() {
	port := ":8000"
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	checkError(err)
	for {
		fmt.Println("Opening UDP connection")
		conn, err := net.ListenUDP("udp", udpAddr)
		var firstReq JSONInitRequest
		var buffer = make([]byte, 200)
		fmt.Println("Reading intitial information, fileSize, fileName and chunkSize")
		n, addr, err := conn.ReadFromUDP(buffer)
		err = json.Unmarshal(buffer[0:n], &firstReq)
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = conn.WriteToUDP([]byte("ACK"), addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	
		file, err := os.Create(serverLocation)
		checkError(err)
		fmt.Println("Receiving File")
		for i := 0; i < firstReq.TotalChunks; i++ {
			receiveFilePacket(conn, file)
		}
		file.Close()
		fmt.Println("Computing and sending the md5 hash")
		clientHash, err := hashFileMd5(serverLocation)
		conn.WriteToUDP([]byte(clientHash), addr)
		conn.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func receiveFilePacket(conn *net.UDPConn, f *os.File) {

	var response JSONRequest
	var buffer = make([]byte, 200)
	n, addr, _ := conn.ReadFromUDP(buffer)
	json.Unmarshal(buffer[0:n], &response)

	f.WriteString(string(response.Data))
	conn.WriteToUDP([]byte("ACK"), addr)
}


func hashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
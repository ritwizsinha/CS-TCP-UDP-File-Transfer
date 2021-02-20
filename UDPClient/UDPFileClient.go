package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"time"
)

const (
	chunkSize   int    = 20
	serverLocation  string = "server.txt"
	resendCount int    = 2
)
// InitRequest is...
type InitRequest struct {
	FileName    string
	FileSize    int
	TotalChunks int
}
// Request is...
type Request struct {
	// Headers
	SeqNumber int
	ChunkSize int
	Data      []byte
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Enter filename")
		os.Exit(1)
	}
	filePath := os.Args[1]
	f, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("File to be transferred doesn't exist")
	}
	fmt.Println("Started a udp connection to the server")
	serverPort := ":8000"
	udpAddr, err := net.ResolveUDPAddr("udp", serverPort)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	send(conn, f)

}

func send(conn *net.UDPConn, fileInfo os.FileInfo) {
	f, err := os.Open(fileInfo.Name())
	checkError(err)
	defer f.Close()
	// Init connection
	initRequest(conn, fileInfo)
	err = receiveACKWithTimeout(conn)
	if err != nil {
		fmt.Println("No acknowledgment for init packet closing connection")
		return
	}
	var chunk [chunkSize]byte
	fmt.Println("Sending datagrams of chunk size 4 bytes")
	for sequence := 0; sequence < int(math.Ceil(float64(fileInfo.Size())/float64(chunkSize))); sequence++ {
		n, err := f.Read(chunk[0:])
		checkError(err)
		sendChunk(conn, chunk[0:n], sequence)
		err = receiveACKWithTimeout(conn)
		if err != nil {
			fmt.Printf("Missed packet with sequence number %d", sequence)
		}
	}
	fmt.Println("Computing MD5 hash")
	md5ServerCheckSum := make([]byte, 100)
	n, err := conn.Read(md5ServerCheckSum)
	if err != nil {
		fmt.Println("Invalid Checksum")
	}
	clientHash, err := hashFileMd5(fileInfo.Name())
	if err != nil {
		fmt.Println("Error in calculating hash")
	}
	if string(md5ServerCheckSum[0:n]) == clientHash {
		fmt.Println("MD5 hash match")
	} else {
		fmt.Println("MD5 hash do not match")
	}
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func sendChunk(conn *net.UDPConn, chunk []byte, sequence int) {
	request := Request{sequence, chunkSize, chunk}
	enc := json.NewEncoder(conn)
	err := enc.Encode(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func receiveACKWithTimeout(conn *net.UDPConn) error {
	conn.SetDeadline(time.Now().Add(time.Second))
	resp := make([]byte, 10)
	n, err := conn.Read(resp)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("FAIL")
	}
	if string(resp[0:n]) != "ACK" {
		return errors.New("FAIL")
	}
	return error(nil)
}

func initRequest(conn *net.UDPConn, fileInfo os.FileInfo) {
	request := InitRequest{fileInfo.Name(), int(fileInfo.Size()), int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))}
	enc := json.NewEncoder(conn)
	enc.Encode(request)
}

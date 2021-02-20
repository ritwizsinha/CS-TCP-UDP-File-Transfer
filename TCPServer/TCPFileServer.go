package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"math"
	"crypto/md5"
	"encoding/hex"
	"io"
	"bytes"
)

// var SERVER_FILE string = "server_store.txt" 
const (
	chunkSize = 1000
	server_store = "server_store.txt"
)

func main() {
	port := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		net, err := listener.Accept()
		logError(err)
		go HandleTCP(net)
	}
}

func HandleTCP(conn net.Conn) {
	defer conn.Close()
	var buf = make([]byte, 100)
	_, err := conn.Read(buf)
	logError(err)
	// The different parameters are separated by commas
	// The starting paramete should be hello followed by name of file
	connInitRequest := strings.Split(string(bytes.Trim(buf, "\x00")), ",")
	if connInitRequest[0] == "HELLO" {
		_, err := conn.Write([]byte("ACK"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		size, err := strconv.Atoi(connInitRequest[2])
		if err != nil {
			fmt.Println(err.Error())
			return;
		}
		// Get all file from connection and store on a local file
		err = storeOnServer(server_store, conn, size)
		if err != nil {
			fmt.Println(err.Error())
			return 
		}
		md5Hash, err := hash_file_md5(server_store)
		if err != nil {
			_, _ = conn.Write([]byte("FIN"))
			return;
		} 

		_, err = conn.Write([]byte("FIN,"+md5Hash))
		return;
 	} else {
		_, err := conn.Write([]byte("FIN"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// if bytes.Compare(buf[0:], []byte("HELLO")) == 0 {
			
		
	// } else {
	// 	_, err = conn.Write([]byte("INVALID REQUEST...ENDING CONNECTION\n"))
	// 	logError(err)
	// 	return
	// }
	// for {
		
	// }
}

func hash_file_md5(filePath string) (string, error) {
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

func logError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
} 

func storeOnServer(fileName string, conn net.Conn, size int) (error) {
		file, err := os.Create(fileName)
		checkError(err) 
		defer file.Close()
		var fileChunk = make([]byte, chunkSize)
		for i:=0;i<=int(math.Ceil(float64(size/chunkSize))); i++ {
			n, err := conn.Read(fileChunk)
			if err != nil {
				return err
			}
			_, err = file.Write(fileChunk[0:n])
			if err != nil {
				return err
			}
		}
		return error(nil)
}


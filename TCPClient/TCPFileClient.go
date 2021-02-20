package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv" 
	"strings"
)

const (
	chunkSize = 1000
)

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
	size := int(f.Size())
	fmt.Println(size)
	serverPort := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverPort)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()
	_, err = conn.Write([]byte("HELLO,FILE," + strconv.Itoa(size)))
	var response = make([]byte, 200)
	_, err = conn.Read(response[0:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(response[0:]))
	// Read the given file and send it over the connetion
	err = readAndSend(filePath, conn)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = conn.Read(response)
	response = clean(response)
	responseList := strings.Split(string(response), ",")
	md5hash, err := hash_file_md5(filePath)
	checkError(err)
	if len(responseList) == 1 || (len(responseList) == 2 && responseList[1] != md5hash) {
		fmt.Println("MD5 not matched")
		return
	}
	if len(responseList) == 2 && responseList[1] == md5hash {
		fmt.Println("MD5 Matched")
		return
	}
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

func readAndSend(filename string, conn net.Conn) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	i := 0
	var fileChunk = make([]byte, chunkSize)
	for {
		i++
		n, err := f.Read(fileChunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// fileChunk = clean(fileChunk)
		_, err = conn.Write(fileChunk[0:n])
		if err != nil {
			return err
		}
	}
	return error(nil)
}

func clean(s []byte) []byte {
	return bytes.Trim(s, "\x000")
}

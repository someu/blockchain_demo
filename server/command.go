package server

import (
	"blockchain_demo/config"
)

func commandToBytes(command string) []byte {
	var bytes [config.CommandLength]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return string(command)
}

func extractCommand(request []byte) []byte {
	return request[:config.CommandLength]
}

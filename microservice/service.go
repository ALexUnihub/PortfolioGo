package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
// если хочется, то для красоты можно разнести логику по разным файликам

func StartMyMicroservice(ctx context.Context, listenAddr, aclData string) error {
	aclMap, err := parseACL(aclData)
	if err != nil {
		return fmt.Errorf("package main, parseACL err: %#v", err.Error())
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("service, net.Listen err : %#v", err.Error())
	}

	server := grpc.NewServer()

	adm := NewAdm(aclMap)

	RegisterBizServer(server, NewBiz(aclMap, adm.Data, adm.Mu))
	RegisterAdminServer(server, adm)

	// server.Serve(lis)
	go func() {
		go func() {
			<-ctx.Done()
			server.Stop()
		}()

		err := server.Serve(lis)
		if err != nil {
			log.Println("pckg main, service, server.Serve err: ", err.Error())
		}
	}()

	return nil
}

func parseACL(aclData string) (map[string]string, error) {
	byteShiftIdx := 0
	aclMap := map[string]string{}
	aclBytes := []byte(aclData)

	for {
		tempIdx := bytes.Index(aclBytes[byteShiftIdx:], []byte(`"`))
		fstIdx := tempIdx + byteShiftIdx
		byteShiftIdx += (tempIdx + 1)

		tempIdx = bytes.Index(aclBytes[byteShiftIdx:], []byte(`"`))
		sndIdx := tempIdx + byteShiftIdx

		if fstIdx == sndIdx {
			break
		}

		key := string(aclBytes[fstIdx+1 : sndIdx])
		byteShiftIdx = (sndIdx + 1)

		err := parseValue(aclBytes, &byteShiftIdx, aclMap, key)
		if err != nil {
			return nil, fmt.Errorf("bad ACL data, value parsing: %#v", err.Error())
		}
	}

	if len(aclMap) == 0 {
		return nil, fmt.Errorf("bad ACL data")
	}

	return aclMap, nil
}

func parseValue(aclBytes []byte, byteShiftIdx *int, aclMap map[string]string, key string) error {
	tempIdx := bytes.Index(aclBytes[*byteShiftIdx:], []byte(`]`))
	if tempIdx == -1 {
		return fmt.Errorf("parseValue err")
	}
	maxIdx := tempIdx + *byteShiftIdx

	var value string

	for *byteShiftIdx < maxIdx {
		tempIdx := bytes.Index(aclBytes[*byteShiftIdx:], []byte(`"`))
		fstIdx := tempIdx + *byteShiftIdx
		*byteShiftIdx += (tempIdx + 1)

		sndIdx := bytes.Index(aclBytes[*byteShiftIdx:], []byte(`"`)) + *byteShiftIdx

		if fstIdx == -1 || sndIdx == -1 {
			return fmt.Errorf("bad ACL data in parse value")
		}

		value += (" " + string(aclBytes[fstIdx+1:sndIdx]))
		*byteShiftIdx = (sndIdx + 1)
	}

	aclMap[key] = value
	return nil
}

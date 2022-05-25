package main

import (
	"bytes"
	"context"
	"strings"
	sync "sync"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type BizServerHandler struct {
	UnimplementedBizServer
	ACLMap map[string]string

	Data *[]*Event
	Mu   *sync.RWMutex
}

func NewBiz(aclMap map[string]string, data *[]*Event, mut *sync.RWMutex) *BizServerHandler {
	return &BizServerHandler{
		ACLMap: aclMap,
		Data:   data,
		Mu:     mut,
	}
}

func (biz *BizServerHandler) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, val, err := checkCtx(ctx, biz)
	if err != nil {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	if !strings.Contains(val, "Biz/Check") && !strings.Contains(val, "Biz/*") {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	// logging
	_, ok := md[":authority"]
	if !ok {
		return nil, status.Errorf(2, "UNKNOWN")
	}
	host := md[":authority"][0]
	host = getHostValue(host)

	biz.Mu.Lock()
	currLen := len((*biz.Data))
	biz.Mu.Unlock()

	newEvent := &Event{
		Timestamp: int64(currLen),
		Consumer:  md["consumer"][0],
		Method:    "/main.Biz/Check",
		Host:      host,
	}

	biz.Mu.Lock()
	(*biz.Data) = append((*biz.Data), newEvent)
	biz.Mu.Unlock()

	return &Nothing{Dummy: true}, nil
}

func (biz *BizServerHandler) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, val, err := checkCtx(ctx, biz)
	if err != nil {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	if !strings.Contains(val, "Biz/Add") && !strings.Contains(val, "Biz/*") {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	// logging
	_, ok := md[":authority"]
	if !ok {
		return nil, status.Errorf(2, "UNKNOWN")
	}
	host := md[":authority"][0]
	host = getHostValue(host)

	biz.Mu.Lock()
	currLen := len((*biz.Data))
	biz.Mu.Unlock()

	newEvent := &Event{
		Timestamp: int64(currLen),
		Consumer:  md["consumer"][0],
		Method:    "/main.Biz/Add",
		Host:      host,
	}

	biz.Mu.Lock()
	(*biz.Data) = append((*biz.Data), newEvent)
	biz.Mu.Unlock()

	return &Nothing{Dummy: true}, nil
}

func (biz *BizServerHandler) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, val, err := checkCtx(ctx, biz)
	if err != nil {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	if !strings.Contains(val, "Biz/Test") && !strings.Contains(val, "Biz/*") {
		return nil, status.Errorf(16, "UNAUTHENTICATED")
	}

	// logging
	_, ok := md[":authority"]
	if !ok {
		return nil, status.Errorf(2, "UNKNOWN")
	}
	host := md[":authority"][0]
	host = getHostValue(host)

	biz.Mu.Lock()
	currLen := len((*biz.Data))
	biz.Mu.Unlock()

	newEvent := &Event{
		Timestamp: int64(currLen),
		Consumer:  md["consumer"][0],
		Method:    "/main.Biz/Test",
		Host:      host,
	}

	biz.Mu.Lock()
	(*biz.Data) = append((*biz.Data), newEvent)
	biz.Mu.Unlock()

	return &Nothing{Dummy: true}, nil
}

func getHostValue(host string) string {
	byteValue := []byte(host)

	idx := bytes.Index(byteValue, []byte(`:`))
	byteValue = byteValue[:idx+1]
	return string(byteValue)
}

func checkCtx(ctx context.Context, biz *BizServerHandler) (metadata.MD, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, "", status.Errorf(16, "UNAUTHENTICATED")
	}

	_, ok = md["consumer"]
	if !ok {
		return nil, "", status.Errorf(16, "UNAUTHENTICATED")
	}

	val, ok := biz.ACLMap[md["consumer"][0]]
	if !ok {
		return nil, "", status.Errorf(16, "UNAUTHENTICATED")
	}
	return md, val, nil
}

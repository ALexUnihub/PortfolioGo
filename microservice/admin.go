package main

import (
	"fmt"
	"strings"
	sync "sync"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AdmServerHandler struct {
	UnimplementedAdminServer
	ACLMap map[string]string

	Data      *[]*Event
	Statistic *Stat
	Mu        *sync.RWMutex

	FlagLog  bool
	FlagStat bool
}

func NewAdm(aclMap map[string]string) *AdmServerHandler {
	return &AdmServerHandler{
		ACLMap: aclMap,

		Data:      &[]*Event{},
		Statistic: &Stat{},
		Mu:        &sync.RWMutex{},

		FlagLog:  true,
		FlagStat: true,
	}
}

func (adm *AdmServerHandler) Logging(nothing *Nothing, inStream Admin_LoggingServer) error {
	// checking adm
	md, ok := metadata.FromIncomingContext(inStream.Context())

	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	_, ok = md["consumer"]
	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	val, ok := adm.ACLMap[md["consumer"][0]]
	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	if !strings.Contains(val, "Admin/Logging") && !strings.Contains(val, "Admin/*") {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	// logging
	dataLen := 0
	adm.Mu.Lock()
	flag := adm.FlagLog
	adm.Mu.Unlock()
	// fmt.Println("flag", flag)

	for {
		if flag {
			// fmt.Println("started sycle in logging")
			adm.Mu.Lock()
			adm.FlagLog = false
			adm.Mu.Unlock()

			flag = false

			_, ok = md[":authority"]
			if !ok {
				return status.Errorf(2, "UNKNOWN")
			}

			host := md[":authority"][0]
			host = getHostValue(host)

			newEvent := &Event{
				Timestamp: int64(0),
				Consumer:  md["consumer"][0],
				Method:    "/main.Admin/Logging",
				Host:      host,
			}

			err := inStream.Send(newEvent)
			if err != nil {
				fmt.Printf("pckg main, admin, Logging err : %#v", err.Error())
			}
		}

		adm.Mu.Lock()
		curLen := len((*adm.Data))
		adm.Mu.Unlock()

		if curLen > dataLen {
			// fmt.Println("logging adm + in len", dataLen)
			// fmt.Println("data to send", (*adm.Data)[dataLen])
			adm.Mu.Lock()
			// fmt.Println("data", (*adm.Data)[dataLen])
			err := inStream.Send((*adm.Data)[dataLen])
			adm.Mu.Unlock()
			if err != nil {
				fmt.Println("ERR:", err.Error())
			}
			dataLen++
		}

	}
}

func (adm *AdmServerHandler) Statistics(stat *StatInterval, inStat Admin_StatisticsServer) error {
	// checking adm
	md, ok := metadata.FromIncomingContext(inStat.Context())

	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	_, ok = md["consumer"]
	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	val, ok := adm.ACLMap[md["consumer"][0]]
	if !ok {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	if !strings.Contains(val, "Admin/Statistics") && !strings.Contains(val, "Admin/*") {
		return status.Errorf(16, "UNAUTHENTICATED")
	}

	shift := 0
	for {
		time.Sleep(time.Duration(stat.IntervalSeconds) * time.Second)

		adm.Mu.Lock()
		flag := adm.FlagLog
		adm.Statistic.ByConsumer = make(map[string]uint64)
		adm.Statistic.ByMethod = make(map[string]uint64)
		adm.Mu.Unlock()

		newStat := &Stat{
			Timestamp:  0,
			ByMethod:   map[string]uint64{},
			ByConsumer: map[string]uint64{},
		}
		// corr
		// fmt.Println("flag:", flag)
		// fmt.Println("created newStat:", newStat)

		if flag {
			adm.Mu.Lock()
			adm.FlagLog = false
			adm.Mu.Unlock()
			newStat.ByMethod["/main.Admin/Statistics"] = 1
			newStat.ByConsumer[md["consumer"][0]] = 1
		}

		adm.Mu.RLock()
		lenData := len(*adm.Data)
		adm.Mu.RUnlock()
		// fmt.Println("lenData", lenData)

		for idx := shift; idx < lenData; idx++ {
			adm.Mu.Lock()
			consum := (*adm.Data)[idx].Consumer
			method := (*adm.Data)[idx].Method
			adm.Mu.Unlock()

			// fmt.Printf("idx = %v, consumer = %v, method = %v, with shift = %v\n", idx, consum, method, shift)

			_, ok = newStat.ByConsumer[consum]

			if !ok {
				newStat.ByConsumer[consum] = 1
			} else {
				newStat.ByConsumer[consum] += 1
			}

			_, ok = newStat.ByMethod[method]

			if !ok {
				newStat.ByMethod[method] = 1
			} else {
				newStat.ByMethod[method] += 1
			}
		}

		shift = lenData

		err := inStat.Send(newStat)
		if err != nil {
			fmt.Printf("pckg main, admin, inStat.Send err: %#v", err.Error())
		}
	}
}

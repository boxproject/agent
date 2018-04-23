package comm

import (
	"strings"
	"sync"
	"github.com/prometheus/common/log"
)

type GrpcStreamChan struct {
	ModelChan chan *GrpcStreamModel
	QuitChan  chan bool
}

var GrpcStreamChanMap map[string]*GrpcStreamChan

var grpcChanMu sync.Mutex

func InitChanMap() {
	GrpcStreamChanMap = make(map[string]*GrpcStreamChan)
}

func AddChan(chanKey string, grpcStream *GrpcStreamChan) {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	GrpcStreamChanMap[chanKey] = grpcStream
}

func AddChanWithChan(chanKey string, modelChan chan *GrpcStreamModel, quitChan chan bool) {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	GrpcStreamChanMap[chanKey] = &GrpcStreamChan{modelChan,quitChan}
}

func RomoveChan(chanKey string) {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	if GrpcStreamChanMap != nil && GrpcStreamChanMap[chanKey] != nil && GrpcStreamChanMap[chanKey].QuitChan != nil {
		GrpcStreamChanMap[chanKey].QuitChan <- false
	}
	delete(GrpcStreamChanMap, chanKey)
}

func SendChanMsg(chanKey string, msgByte []byte) bool {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	log.Debug("SendChanMsg....key:%s",chanKey)
	for k, v := range GrpcStreamChanMap {
		if strings.HasPrefix(k, chanKey) {
			if v != nil && v.ModelChan != nil {
				v.ModelChan <- &GrpcStreamModel{Msg: msgByte}
			} else {
				return false
			}
		}
	}
	return true
}

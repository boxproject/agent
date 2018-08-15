package comm

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	log "github.com/alecthomas/log4go"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"sync"
	"time"
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
	printChanMap()
}

func AddChanWithChan(chanKey string, modelChan chan *GrpcStreamModel, quitChan chan bool) {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	GrpcStreamChanMap[chanKey] = &GrpcStreamChan{modelChan, quitChan}
}

func RomoveChan(chanKey string) {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	if GrpcStreamChanMap != nil && GrpcStreamChanMap[chanKey] != nil && GrpcStreamChanMap[chanKey].QuitChan != nil {
		GrpcStreamChanMap[chanKey].QuitChan <- false
	}
	delete(GrpcStreamChanMap, chanKey)
	printChanMap()
}

func SendChanMsg(chanKey string, msgByte []byte) bool {
	defer grpcChanMu.Unlock()
	grpcChanMu.Lock()
	log.Debug("SendChanMsg....key:%s", chanKey)
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

func printChanMap() {
	i := 1
	for key, _ := range GrpcStreamChanMap {
		log.Debug("current chan key[%d]:%s", i, key)
		i++
	}
}

//地址格式校验
func CheckAddress(address string, category int64) bool {
	isGoodAddr := false
	switch category {
	case CATEGORY_BTC: //
		if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err != nil {
			isGoodAddr = false
		} else {
			isGoodAddr = true
		}
		break
	default:
		isGoodAddr = common.IsHexAddress(address)
	}
	return isGoodAddr
}

const (
	//密码验证状态
	PASSS_STATUS_OTHER = "OTHER" //未连接
	PASSS_STATUS_TRUE  = "TRUE"  //正确
	PASSS_STATUS_FALSE = "FALSE" //错误
)

var passStatus = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func SetPassStatus(key string, value string) {
	passStatus.Lock()
	defer passStatus.Unlock()
	if value == PASSS_STATUS_TRUE || value == PASSS_STATUS_FALSE {
		passStatus.m[key] = value
	} else {
		passStatus.m[key] = PASSS_STATUS_OTHER
	}

}
func ClearPassStatus(key string) {
	passStatus.Lock()
	defer passStatus.Unlock()
	if _, ok := passStatus.m[key]; ok {
		delete(passStatus.m, key)
	}
}
func ClearAllPassStatus(key string) {
	passStatus.Lock()
	defer passStatus.Unlock()
	for key, _ := range passStatus.m {
		delete(passStatus.m, key)
	}
}
func GetPassStatus(key string) string {
	passStatus.Lock()
	defer passStatus.Unlock()
	if value, ok := passStatus.m[key]; ok {
		return value
	} else {
		return PASSS_STATUS_OTHER
	}
}

//验证密码正确，等待签名机30S
func DelayGetPassStatus(key string) string {
	waitVoucher := time.NewTicker(time.Nanosecond)
	for i := 0; i < 30; {
		select {
		case <-waitVoucher.C:
			waitVoucher = time.NewTicker(time.Second)
			i++
			if GetPassStatus(key) != PASSS_STATUS_OTHER {
				i = 100
			}
			break
		}
	}
	status := GetPassStatus(key)
	//clear pass status in map
	ClearPassStatus(key)

	return status
}

// 验证签名值
func VerifySign(msg, pubkey, signature string) (bool, error) {
	bSignData, err := base64.StdEncoding.DecodeString(signature)
	hashed := sha256.Sum256([]byte(msg))
	bKey, err := base64.RawStdEncoding.DecodeString(pubkey)
	pub, err := x509.ParsePKCS1PublicKey(bKey)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], bSignData)
	if err != nil {
		return false, err
	}
	return true, nil
}

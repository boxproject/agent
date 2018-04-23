package httpcli

import (
	"encoding/json"
	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RepCli struct {
	quitChannel chan int
	cfg         *config.Config
}

func NewRepCli(cfg *config.Config) *RepCli {
	return &RepCli{cfg: cfg}
}

//启动任务
func (r *RepCli) Start() {
	loop := true
	for loop {
		select {
		case <-r.quitChannel:
			logger.Info("PriEthHandler::SendMessage thread exitCh!")
			loop = false
		case data, ok := <-comm.VReqChan:
			if ok {
				switch data.ReqType {
				case comm.REQ_WITHDRAW: //withdraw
					r.withdrawReq(data)
				case comm.REQ_DEPOSIT: //deposit
					r.depositReq(data)
				case comm.REQ_WITHDRAW_TX: //deposit tx
					r.withdrawTxReq(data)
				case comm.REQ_TOKEN_CHANGE: //token change
					r.tokenChangeReq(data)
				case comm.REQ_REGIST: //token change
					r.registReq(data)
				default:
					logger.Info("unknow req:%v", data.ReqType)
				}
			} else {
				logger.Error("read from channel failed")
			}
		}
	}
}

//停止任务
func (r *RepCli) Stop() {
	r.quitChannel <- 0
}

//未处理完请求 TODO
func (r *RepCli) unFinishedReq() {

}

//充值
func (r *RepCli) depositReq(vReq *comm.VReq) {
	logger.Debug("RepCli depositReq: ", vReq)
	data := url.Values{"from": {vReq.From}, "to": {vReq.To}, "category": {strconv.Itoa(int(vReq.Category))}, "tx_id": {vReq.TxHash}, "amount": {vReq.Amount}}
	reqBody := strings.NewReader(data.Encode())
	resp, err := http.Post(r.cfg.DepositUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		logger.Error("http request error:%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cRsp := &comm.VRsp{}
	if err := json.Unmarshal(body, cRsp); err != nil {
		//TODO
		logger.Error("json marshal error:%v", err)
		return
	} else {
		logger.Info("cRsp:", cRsp)
	}
}

//提现
func (r *RepCli) withdrawReq(vReq *comm.VReq) {
	logger.Debug("RepCli withdrawReq: ", vReq)
	data := url.Values{"to": {vReq.To}, "category": {strconv.Itoa(int(vReq.Category))}, "wd_hash": {vReq.WdHash}, "tx_id": {vReq.TxHash}, "amount": {vReq.Amount}}
	reqBody := strings.NewReader(data.Encode())
	resp, err := http.Post(r.cfg.WithDrawUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		logger.Error("http request error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cRsp := &comm.VRsp{}
	if err := json.Unmarshal(body, cRsp); err != nil {
		//TODO
		logger.Error("json marshal error: %v", err)
		return
	} else {
		logger.Info("cRsp:", cRsp)
	}
}

//提现tx
func (r *RepCli) withdrawTxReq(vReq *comm.VReq) {
	logger.Debug("RepCli withdrawTxReq: ", vReq)
	data := url.Values{"wd_hash": {vReq.WdHash}, "tx_id": {vReq.TxHash}}
	reqBody := strings.NewReader(data.Encode())
	resp, err := http.Post(r.cfg.WithDrawTxUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		logger.Error("http request error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cRsp := &comm.VRsp{}
	if err := json.Unmarshal(body, cRsp); err != nil {
		logger.Error("json marshal error: %v", err)
		return
	} else {
		logger.Info("cRsp:", cRsp)
	}
}

//token 变动
func (r *RepCli) tokenChangeReq(vReq *comm.VReq) {
	logger.Debug("RepCli tokenChangeReq: ", vReq)
	data := url.Values{"type": {vReq.CurrencyType}}
	reqBody := strings.NewReader(data.Encode())
	resp, err := http.Post(r.cfg.TokenChangeUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		logger.Error("http request error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cRsp := &comm.VRsp{}
	if err := json.Unmarshal(body, cRsp); err != nil {
		logger.Error("json marshal error: %v", err)
		return
	} else {
		logger.Info("cRsp:", cRsp)
	}
}

//regist
func (r *RepCli) registReq(vReq *comm.VReq) {
	logger.Debug("RepCli registReq: ", vReq)
	data := url.Values{"regid": {vReq.RegId}, "status": {vReq.Status}, "ciphertext": {vReq.CipherText}, "pubkey": {vReq.PubKey}}
	reqBody := strings.NewReader(data.Encode())
	resp, err := http.Post(r.cfg.RegistApprovalUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		logger.Error("http request error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cRsp := &comm.VRsp{}
	if err := json.Unmarshal(body, cRsp); err != nil {
		logger.Error("json marshal error: %v", err)
		return
	} else {
		logger.Info("cRsp:", cRsp)
	}
}

// Copyright 2018. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/config"
	"github.com/boxproject/agent/db"
	"github.com/boxproject/agent/httpcli"
	"github.com/boxproject/agent/model"
	"github.com/ethereum/go-ethereum/common"
	lerrors "github.com/syndtr/goleveldb/leveldb/errors"
	"math/big"
	"strings"
	"time"
)

type VoucherController struct {
	baseController
}

type RspModel struct {
	RspNo         string
	ManagerIpPort string //管理端地址和端口
	Status        comm.VoucherStatus
	RegistInfos   []comm.RegistInfo
	ApprovalInfo  comm.ApprovalInfo
	TokenInfos    []comm.TokenInfo
	CoinStatus    []comm.CoinStatu
	ApprovalInfos []comm.ApprovalInfo
	HashOperates  []comm.HashOperate
}

//err json pkg
func (v *VoucherController) retErrJSON(errNo string) {
	v.Data["json"] = &RspModel{RspNo: errNo}
	v.ServeJSON()
}

//add regist
func (v *VoucherController) AddRegist() {
	log.Debug("addRegist....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	regId := v.GetString("regid")
	applyerId := v.GetString("applyerid")
	captainId := v.GetString("captainid")
	applyerAccount := v.GetString("applyeraccount")
	msg := v.GetString("msg")
	status := v.GetString("status")

	regist := &model.TRegist{RegId: regId, ApplyerId: applyerId, CaptainId: captainId, ApplyerAccount: applyerAccount, Msg: msg, Status: status, CreateTime: time.Now()}
	if _, err := db.GetDefaultNewOrmer().Insert(regist); err != nil {
		log.Error("land to db err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

//注册列表
func (v *VoucherController) RegistList() {
	log.Debug("queryRegist....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	captainId := v.GetString("captainid")
	var tRegists []*model.TRegist
	if i, err := db.GetDefaultNewOrmer().QueryTable(&model.TRegist{}).Filter("CaptainId", captainId).All(&tRegists); err != nil {
		log.Error("get hashlist err:%s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else if i > 0 {
		for _, registModel := range tRegists {
			rspModel.RegistInfos = append(rspModel.RegistInfos, comm.RegistInfo{RegId: registModel.RegId,
				ApplyerId:      registModel.ApplyerId,
				CaptainId:      registModel.CaptainId,
				ApplyerAccount: registModel.ApplyerAccount,
				Msg:            registModel.Msg,
				Status:         registModel.Status,
			})
		}
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) RegistAproval() {
	log.Debug("registAproval....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	regId := v.GetString("regid")
	consent := v.GetString("consent")
	ciphertext := v.GetString("ciphertext")
	status := v.GetString("status")
	pubKey := v.GetString("pubkey")
	registModel := &model.TRegist{RegId: regId}
	if err := db.GetDefaultNewOrmer().Read(registModel); err != nil {
		log.Error("regist load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else {
		registModel.Consent = consent
		registModel.CipherText = ciphertext
		registModel.Status = status
		registModel.PubKey = pubKey
		if _, err := db.GetDefaultNewOrmer().Update(registModel); err != nil {
			log.Error("regist update err: %s", err)
		} else {
			comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_REGIST, RegId: regId, Consent: consent, CipherText: ciphertext, PubKey: pubKey, Status: status}
		}
		//rspModel.ApprovalInfo = comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Flow: hashModel.Flow, Sign: hashModel.Sign, Status: hashModel.Status}
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

//add approval
func (v *VoucherController) AddApproval() {
	log.Debug("addApproval....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	hash := v.GetString("hash")
	name := v.GetString("name")
	appId := v.GetString("appid")
	captainId := v.GetString("captainid")
	flow := v.GetString("flow")
	sign := v.GetString("sign")

	hashO := &model.THash{Hash: hash, Name: name, AppId: appId, CaptainId: captainId, Sign: sign, Flow: flow, Status: comm.HASH_STATUS_0, CreateTime: time.Now()}
	if _, err := db.GetDefaultNewOrmer().Insert(hashO); err != nil {
		log.Error("land to db err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

//Invalid Approval
func (v *VoucherController) InvalidApproval() {
	log.Debug("InvalidApproval....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	hash := v.GetString("hash")
	appId := v.GetString("appid")
	sign := v.GetString("sign")
	// 获取用户公钥
	var tTregist model.TRegist
	ormer := db.GetDefaultNewOrmer()
	if err := ormer.QueryTable(&model.TRegist{}).Filter("ApplyerId", appId).Filter("Status", comm.REG_APPROVAL).One(&tTregist); err != nil {
		log.Error("*********get hashlist err:%s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else {
		if tTregist.PubKey == "" {
			v.retErrJSON(comm.Err_USER)
			return
		}
	}

	// 获取对应哈希信息
	var tHash model.THash
	if err := ormer.QueryTable(&model.THash{}).Filter("Hash", hash).One(&tHash); err != nil {
		log.Error("get hashlist err:%s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else {
		if tHash.Name == "" {
			v.retErrJSON(comm.Err_USER)
			return
		}
	}

	// 验签
	sign_pass, err := comm.VerifySign(tHash.Name, tTregist.PubKey, sign)

	if err != nil {
		log.Error("Disuse hash verify sign", err)
		v.retErrJSON(comm.Err_UNKNOW_REQ_TYPE)
		return
	}

	if sign_pass != true {
		log.Error("Disuse hash")
		v.retErrJSON(comm.Err_SIGN)
		return
	}

	// 更改hash状态
	//hashModel := &model.THash{Hash: hash, AppId: appId, Status: comm.HASH_STATUS_9, CreateTime: time.Now()}
	tHash.Status = comm.HASH_STATUS_9
	if _, err := db.GetDefaultNewOrmer().Update(&tHash, "Status"); err != nil {
		log.Error("Disuse hash err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) ApprovalList() {
	log.Debug("approvalList....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	captainId := v.GetString("captainid")
	rType := v.GetString("type")
	log.Debug("captainId:%v,type:%v", captainId, rType)
	if rType == comm.APPROVAL_TYPE_0 { //待自己审批
		var tHashs []*model.THash
		if i, err := db.GetDefaultNewOrmer().QueryTable(&model.THash{}).All(&tHashs); err != nil {
			log.Error("get hashlist err:%s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		} else if i > 0 {
			for _, hashModel := range tHashs {
				if hashModel.CaptainId == captainId && hashModel.Status == comm.HASH_STATUS_0 { //待私钥A审批
					rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
				} else if hashModel.CaptainId != captainId && hashModel.Status == comm.HASH_STATUS_3 { //待其他私钥审批
					if !v.getHashOperateOption(hashModel.Hash, captainId) {
						rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
					}
				}
			}
		}
	} else if rType == comm.APPROVAL_TYPE_1 {
		var tHashs []*model.THash
		if i, err := db.GetDefaultNewOrmer().QueryTable(&model.THash{}).All(&tHashs); err != nil {
			log.Error("get hashlist err:%s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		} else if i > 0 {
			for _, hashModel := range tHashs {
				if hashModel.CaptainId == captainId && hashModel.Status != comm.HASH_STATUS_0 {
					rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
				} else if hashModel.CaptainId != captainId && (hashModel.Status == comm.HASH_STATUS_4 || hashModel.Status == comm.HASH_STATUS_5 || hashModel.Status == comm.HASH_STATUS_6 || hashModel.Status == comm.HASH_STATUS_7 || hashModel.Status == comm.HASH_STATUS_9) {
					//if v.getHashOperateOption(hashModel.Hash, captainId) {
					rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
					//}
				} else if hashModel.CaptainId != captainId && (hashModel.Status == comm.HASH_STATUS_2 || hashModel.Status == comm.HASH_STATUS_3) {
					if v.getHashOperateOption(hashModel.Hash, captainId) {
						rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
					}
				}
			}
		}
	} else if rType == comm.APPROVAL_TYPE_2 { //查询公链已同意审批流
		var tHashs []*model.THash
		if i, err := db.GetDefaultNewOrmer().QueryTable(&model.THash{}).Filter("Status", comm.HASH_STATUS_7).All(&tHashs); err != nil {
			log.Error("get hashlist err:%s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		} else if i > 0 {
			for _, hashModel := range tHashs {
				rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
			}
		}
	} else if rType == comm.APPROVAL_TYPE_3 {
		var tHashs []*model.THash
		if i, err := db.GetDefaultNewOrmer().QueryTable(&model.THash{}).Filter("Status__in", comm.HASH_STATUS_7, comm.HASH_STATUS_2, comm.HASH_STATUS_5, comm.HASH_STATUS_9).All(&tHashs); err != nil {
			log.Error("get hashlist err:%s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		} else if i > 0 {
			for _, hashModel := range tHashs {
				rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
			}
		}
	} else {
		log.Debug("other type:%v", rType)
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) ApprovalDetail() {
	log.Debug("approvalDetail....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	hash := v.GetString("hash")
	hashModel := &model.THash{Hash: hash}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else {
		rspModel.ApprovalInfo = comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, CaptainId: hashModel.CaptainId, Name: hashModel.Name, Flow: hashModel.Flow, Sign: hashModel.Sign, Status: hashModel.Status}
	}

	var tHashOperates []*model.THashOperate
	ormer := db.GetDefaultNewOrmer()
	if _, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hash).All(&tHashOperates); err != nil {
		log.Error("query hashOpeate err:", err)
	} else {
		for _, hashOpeate := range tHashOperates {
			rspModel.HashOperates = append(rspModel.HashOperates, comm.HashOperate{CaptainId: hashOpeate.AppId, Option: hashOpeate.Option})
		}
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) KeyStore() {
	log.Debug("keyStore....")
	//operateType := v.GetString("type")
	appId := v.GetString("applyerid")
	appName := v.GetString("applyername")
	//password := v.GetString("password")
	reqIpPort := v.GetString("reqipport")
	code := v.GetString("code")
	publicKey := v.GetString("publickey")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_ADDKEY, AppId: appId, AppName: appName, ReqIpPort: reqIpPort, Code: code, PublicKey: publicKey}
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("add keystore marshal err: %s", err)
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) Operate() {
	log.Debug("opeate....")
	operateType := v.GetString("type")
	appId := v.GetString("applyerid")
	password := v.GetString("password")
	sign := v.GetString("sign")
	reqIpPort := v.GetString("reqipport")
	code := v.GetString("code")
	publicKey := v.GetString("publickey")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{
		Type:      operateType,
		AppId:     appId,
		Password:  password,
		Sign:      sign,
		ReqIpPort: reqIpPort,
		Code:      code,
		PublicKey: publicKey}
	log.Debug("[OPEATE]:%v\n[voucher operate=%v]", grpcStream, grpcStream.VoucherOperate)
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("opeate marshal err: %s", err)
		v.retErrJSON(comm.Err_JSON)
		return
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}

	//reqLog := &model.TReqLog{
	//	ReqType:      operateType,
	//	TransferType: grpcStream.Type,
	//	CreateTime:   time.Now(),
	//	ApplyTime:    time.Now()}
	//db.GetDefaultNewOrmer().Insert(reqLog)

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) Status() {
	log.Debug("status....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	rspModel.Status = *comm.RealTimeVoucherStatus
	v.Data["json"] = rspModel
	//db.GetDefaultNewOrmer().Insert(reqLog)
	v.ServeJSON()
}

func (v *VoucherController) AddHash() {
	log.Debug("addHash....")
	hashStr := v.GetString("hash")
	if !strings.HasPrefix(hashStr, comm.HASH_PRIFIX) {
		v.retErrJSON(comm.Err_UNENABLE_PREFIX)
		return
	}

	rspModel := &RspModel{RspNo: comm.Err_OK}

	hashModel := &model.THash{Hash: hashStr}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
	} else {
		grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_ADD_REQ, Hash: common.HexToHash(hashStr)}
		if msg, err := json.Marshal(grpcStream); err != nil {
			log.Error("allow marshal err: %s", err)
		} else {
			comm.SendChanMsg(comm.SERVER_COMPANION, msg)
		}
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) Wtihdraw() {
	log.Debug("wtihdraw....")
	rspModel := &RspModel{RspNo: comm.Err_OK}

	hashStr := v.GetString("hash")
	wdHashStr := v.GetString("wdhash")

	if !strings.HasPrefix(hashStr, comm.HASH_PRIFIX) || !strings.HasPrefix(wdHashStr, comm.HASH_PRIFIX) {
		v.retErrJSON(comm.Err_UNENABLE_PREFIX)
		return
	}

	recAddress := v.GetString("recaddress")
	amount := v.GetString("amount")
	flow := v.GetString("apply")       //提现原始数据
	wdFlow := v.GetString("applysign") //提现原始数据
	//sign := v.GetString("sign")

	fee := v.GetString("fee")

	category, err := v.GetInt64("category")
	if err != nil {
		log.Debug("category[%d] illegal", category)
		v.retErrJSON(comm.Err_UNENABLE_AMOUNT)
		return
	}

	if comm.CheckAddress(recAddress, category) == false {
		log.Debug("recAddress illegal:", recAddress)
		v.retErrJSON(comm.Err_UNENABLE_ADDRESS)
		return
	}

	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	feeBig := new(big.Int)
	feeBig.SetString(fee, 10)

	withdrawModel := &model.TWithdraw{Hash: hashStr, WdHash: wdHashStr, To: recAddress, Category: category, Amount: amount, Fee: fee, Flow: flow, WdFlow: wdFlow, Status: comm.WITHDRAW_STATUS_0, CreateTime: time.Now()}
	if _, err := db.GetDefaultNewOrmer().Insert(withdrawModel); err != nil {
		log.Error("land to db err: %s", err)
	}

	grpcStream := &comm.GrpcStream{Type: comm.GRPC_WITHDRAW_REQ, Hash: common.HexToHash(hashStr), WdHash: common.HexToHash(wdHashStr), To: recAddress, Category: big.NewInt(category), Amount: amountBig, Fee: feeBig}
	log.Debug("[OPEATE]:", grpcStream)
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("allow marshal err: %s", err)
	} else {
		comm.SendChanMsg(comm.SERVER_COMPANION, msg)
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

//验证密码
/*
	return(int):
	-1	->	voucher请求错误
	0	->	密码错误
	1	->	密码正确
*/
func checkPassword(grpcStream comm.GrpcStream) string {
	log.Debug("checkPassword....")
	//向voucher请求并验证密码
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("allow marshal err: %s", err)
		return comm.PASSS_STATUS_OTHER
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
		return comm.DelayGetPassStatus(grpcStream.VoucherOperate.AppId)
	}
}

//同意
func (v *VoucherController) Allow() {
	log.Debug("allow....")
	appId := v.GetString("captainid")
	sign := v.GetString("sign")
	hashStr := v.GetString("hash")
	if !strings.HasPrefix(hashStr, comm.HASH_PRIFIX) {
		v.retErrJSON(comm.Err_UNENABLE_PREFIX)
		return
	}
	password := v.GetString("password")
	pwdsign := v.GetString("pwdsign")
	reason := v.GetString("reason")
	opinion := v.GetString("opinion")

	rspModel := &RspModel{RspNo: comm.Err_OK}

	hashModel := &model.THash{Hash: hashStr}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}
	//向voucher请求并验证密码
	flow := hashModel.Flow
	grpcStream := &comm.GrpcStream{
		Type:  comm.GRPC_VOUCHER_OPR_REQ,
		Hash:  common.HexToHash(hashStr),
		Flow:  flow,
		AppId: appId,
		Sign:  sign}
	grpcStream.VoucherOperate = &comm.Operate{
		Type:     comm.VOUCHER_OPERATE_CHECK_KEY,
		AppId:    appId,
		Sign:     sign,
		Hash:     flow,
		Password: password,
		PassSign: pwdsign}
	switch checkPassword(*grpcStream) {
	case comm.PASSS_STATUS_TRUE:
		log.Info("pass check: true")
		break
	case comm.PASSS_STATUS_FALSE:
		//	密码错误
		log.Info("pass check: failed")
		v.retErrJSON(comm.Err_WRONG_PASS)
		return
	default:
		//req voucher error
		log.Info("pass check: voucher req error")
		v.retErrJSON(comm.Err_VOUCHER_REQERR)
		return
	}

	//hash operate，避免重复操作
	if count, err := db.GetDefaultNewOrmer().QueryTable(&model.THashOperate{}).Filter("Hash", hashStr).Filter("Type", comm.HASH_TYPE_ALLOW).Filter("AppId", appId).Filter("Sign", sign).Count(); err == nil && count >= 1 {
		log.Info("Repeat hash Operate [AppID:%v]", appId)
		v.retErrJSON(comm.Err_REPEAT_REQ)
		return
	}

	//hash operate
	hashOperate := &model.THashOperate{Hash: hashStr, Type: comm.HASH_TYPE_ALLOW, Option: opinion, AppId: appId, Sign: sign, Opinion: reason, CreateTime: time.Now()}
	db.GetDefaultNewOrmer().Insert(hashOperate)
	log.Debug("record hashOperate:%v", hashOperate)

	if hashModel.Status == comm.HASH_STATUS_0 { //私钥a确认
		if opinion == comm.TRUE { //同意
			hashModel.Status = comm.HASH_STATUS_1 //更新 私钥已申请提交
			grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_ADD_REQ, Hash: common.HexToHash(hashStr)}
			if msg, err := json.Marshal(grpcStream); err != nil {
				log.Error("allow marshal err: %s", err)
			} else {
				comm.SendChanMsg(comm.SERVER_COMPANION, msg)
			}
		} else {
			hashModel.Status = comm.HASH_STATUS_2 //更新 私钥已拒绝提交
		}

	} else if hashModel.Status == comm.HASH_STATUS_3 { //私钥b,c确认
		if opinion != comm.TRUE {
			hashModel.Status = comm.HASH_STATUS_5 //更新 私链已拒绝确认 私钥B、私钥C有不同意
		} else {
			var tHashOperates []*model.THashOperate
			ormer := db.GetDefaultNewOrmer()
			if i, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hashStr).Filter("Type", comm.HASH_TYPE_ALLOW).Filter("Option", comm.TRUE).All(&tHashOperates); err != nil {
				log.Error("query hashOpeate err:", err)
			} else {
				log.Debug("get allow operate count:%v", i)
				if comm.RealTimeVoucherStatus.Total <= i {
					//去重复
					var appOperateMap = make(map[string]string)
					for _, tHashOperate := range tHashOperates {
						appOperateMap[tHashOperate.AppId] = tHashOperate.Hash
					}
					if comm.RealTimeVoucherStatus.Total == int64(len(appOperateMap)) {
						hashModel.Status = comm.HASH_STATUS_4
						log.Debug("app all allow, waite companion confirm...")
						grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_ENABLE_REQ, Hash: common.HexToHash(hashStr)}
						if msg, err := json.Marshal(grpcStream); err != nil {
							log.Error("allow marshal err: %s", err)
						} else { //发往伴生程序
							comm.SendChanMsg(comm.SERVER_COMPANION, msg)
						}
					}
				}
			}
		}
	}

	if _, err := db.GetDefaultNewOrmer().Update(hashModel); err != nil {
		log.Error("update thash err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) ApprovalOperateList() {
	log.Debug("ApprovalOperateList....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	hash := v.GetString("hash")

	//get hash
	hashModel := &model.THash{Hash: hash}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}

	//get applyerAccount
	applyerAccount := ""
	var tRegists []*model.TRegist
	if _, err := db.GetDefaultNewOrmer().QueryTable(&model.TRegist{}).Filter("ApplyerId", hashModel.AppId).All(&tRegists); err != nil {
		log.Error("get hashlist err:%s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}
	if len(tRegists) == 1 {
		applyerAccount = tRegists[0].ApplyerAccount
	} else {
		log.Debug("get TRegist count:%v", len(tRegists))
	}

	// get hash create record
	newHashOperates := comm.HashOperate{
		ApplyerAccount: applyerAccount,
		CaptainId:      hashModel.AppId,
		Option:         comm.GREATE,
		Opinion:        "create",
		CreateTime:     hashModel.CreateTime.Format("2006-01-02 15:04:05")}
	rspModel.HashOperates = append(rspModel.HashOperates, newHashOperates)

	//get hash operate record
	var tHashOperates []*model.THashOperate
	ormer := db.GetDefaultNewOrmer()
	if _, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hash).OrderBy("CreateTime").All(&tHashOperates); err != nil {
		log.Error("query hashOpeate err:", err)
	} else {
		for _, hashOpeate := range tHashOperates {
			newHashOperates := comm.HashOperate{
				ApplyerAccount: applyerAccount,
				CaptainId:      hashOpeate.AppId,
				Option:         hashOpeate.Option,
				Opinion:        hashOpeate.Opinion,
				CreateTime:     hashOpeate.CreateTime.Format("2006-01-02 15:04:05")}
			rspModel.HashOperates = append(rspModel.HashOperates, newHashOperates)
		}
	}

	option := ""
	if hashModel.Status == comm.HASH_STATUS_9 {
		option = comm.INVALID
	} else {
		option = comm.OTHER
	}

	//add hash current status
	newHashOperates = comm.HashOperate{
		ApplyerAccount: applyerAccount,
		CaptainId:      hashModel.AppId,
		Option:         option,
		Opinion:        hashModel.Status,
		CreateTime:     hashModel.CreateTime.Format("2006-01-02 15:04:05")}
	rspModel.HashOperates = append(rspModel.HashOperates, newHashOperates)

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) DisAllow() {
	log.Debug("disallow....")
	appId := v.GetString("appid")
	sign := v.GetString("sign")
	hashStr := v.GetString("hash")
	password := v.GetString("password")
	pwdsign := v.GetString("pwdsign")
	reason := v.GetString("reason")

	if !strings.HasPrefix(hashStr, comm.HASH_PRIFIX) {
		v.retErrJSON(comm.Err_UNENABLE_PREFIX)
		return
	}
	opinion := v.GetString("opinion")
	rspModel := &RspModel{RspNo: comm.Err_OK}

	//密码验证
	hashModel := &model.THash{Hash: hashStr}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	}

	//向voucher请求并验证密码
	flow := hashModel.Flow
	grpcStream := &comm.GrpcStream{
		Type:  comm.GRPC_VOUCHER_OPR_REQ,
		Hash:  common.HexToHash(hashStr),
		Flow:  flow,
		AppId: appId,
		Sign:  sign}
	grpcStream.VoucherOperate = &comm.Operate{
		Type:     comm.VOUCHER_OPERATE_CHECK_KEY,
		AppId:    appId,
		Sign:     sign,
		Hash:     flow,
		Password: password,
		PassSign: pwdsign}
	switch checkPassword(*grpcStream) {
	case comm.PASSS_STATUS_TRUE:
		log.Info("pass check: true")
		break
	case comm.PASSS_STATUS_FALSE:
		//	密码错误
		log.Info("pass check: failed")
		v.retErrJSON(comm.Err_WRONG_PASS)
		return
	default:
		//req voucher error
		log.Info("pass check: voucher req error")
		v.retErrJSON(comm.Err_VOUCHER_REQERR)
		return
	}

	//hash operate，避免重复操作
	if num, err := db.GetDefaultNewOrmer().QueryTable(&model.THashOperate{}).Filter("Hash", hashStr).Filter("Type", comm.HASH_TYPE_DISALLOW).Filter("AppId", appId).Filter("Sign", sign).Count(); err == nil && num > 0 {
		log.Info("Repeat hash Operate [AppID:%v]", appId)
		v.retErrJSON(comm.Err_REPEAT_REQ)
		return
	}
	hashOperate := &model.THashOperate{Hash: hashStr, Type: comm.HASH_TYPE_DISALLOW, Option: opinion, AppId: appId, Sign: sign, Opinion: reason, CreateTime: time.Now()}
	db.GetDefaultNewOrmer().Insert(hashOperate)
	if opinion == comm.TRUE {
		var tHashOperates []*model.THashOperate
		ormer := db.GetDefaultNewOrmer()
		if _, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hashStr).Filter("Type", comm.HASH_TYPE_DISALLOW).Filter("Option", comm.FALSE).All(&tHashOperates); err != nil {
			log.Error("query hashOpeate err:", err)
		} else {
			if len(tHashOperates) > 0 {
				grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_DISABLE_REQ, Hash: common.HexToHash(hashStr)}
				log.Debug("[OPEATE]:%v", grpcStream)
				if msg, err := json.Marshal(grpcStream); err != nil {
					log.Error("allow marshal err: %s", err)
				} else { //发往伴生程序
					comm.SendChanMsg(comm.SERVER_COMPANION, msg)
				}
			}
		}
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) TokenList() {
	log.Debug("tokenlist....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	if tokenMap, err := comm.Ldb.GetPrifix([]byte(comm.TOKENLIST_PRIFIX)); err != nil {
		log.Error("get hashlist err:%s", err)
	} else {
		if len(tokenMap) != int(comm.RealTimeVoucherStatus.TokenCount) { //数量不足请求
			log.Info("token 数量不足", len(tokenMap), comm.RealTimeVoucherStatus.TokenCount)
			grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
			log.Debug("[OPEATE]:%v", grpcStream)
			grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_TOKEN_LIST}
			if msg, err := json.Marshal(grpcStream); err != nil {
				log.Error("allow marshal err: %s", err)
			} else {
				comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
			}
		}
		for _, hashBytes := range tokenMap {
			tokenInfo := &comm.TokenInfo{}
			if err = json.Unmarshal([]byte(hashBytes), tokenInfo); err != nil {
				v.retErrJSON(comm.Err_JSON)
			} else {
				rspModel.TokenInfos = append(rspModel.TokenInfos, *tokenInfo)
			}
		}
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) AddToken() {
	log.Debug("addToken....")
	appId := v.GetString("applyerid")
	tokenName := v.GetString("tokenname")
	sign := v.GetString("sign")
	decimals, err := v.GetInt64("decimals")
	if err != nil {
		v.retErrJSON(comm.Err_UNENABLE_AMOUNT)
	}
	contractAddr := v.GetString("contractaddr")
	if _, err := comm.Ldb.GetByte([]byte(comm.TOKENLIST_PRIFIX + contractAddr)); err != nil && err != lerrors.ErrNotFound {
		log.Error("get tokenlist err:%s", err)
		v.retErrJSON(comm.Err_LDB)
		return
	}

	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_TOKEN_ADD, TokenName: tokenName, Decimals: decimals, ContractAddr: contractAddr, Sign: sign, AppId: appId}
	log.Debug("[OPEATE]:%v", grpcStream)
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("token marshal err: %s", err)
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) DelToken() {
	log.Debug("delToken....")
	appId := v.GetString("applyerid")
	contractAddr := v.GetString("contractaddr")
	sign := v.GetString("sign")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_TOKEN_DEL, ContractAddr: contractAddr, Sign: sign, AppId: appId}
	log.Debug("[OPEATE]:%v", grpcStream)
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("allow marshal err: %s", err)
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) Coin() {
	log.Debug("delToken....")
	category, err := v.GetInt64("category")
	if err != nil {
		v.retErrJSON(comm.Err_UNENABLE_AMOUNT)
	}
	usedStr := v.GetString("used")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_COIN, CoinCategory: category, CoinUsed: usedStr == comm.TRUE}
	log.Debug("[OPEATE]:%v", grpcStream)
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("allow marshal err: %s", err)
		v.retErrJSON(comm.Err_JSON)
		return
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) CoinList() {
	log.Debug("coinlist....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	rspModel.CoinStatus = comm.RealTimeVoucherStatus.CoinStatus
	v.Data["json"] = rspModel
	v.ServeJSON()
}

//查询管理端信息
func (v *VoucherController) ManagerInfo() {
	log.Debug("managerInfo....")
	rspModel := &RspModel{RspNo: comm.Err_OK, ManagerIpPort: comm.MANAGER_SERVER_IPPORT}
	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) getHashOperateOption(hash string, appId string) bool {
	var tHashOperates []*model.THashOperate
	ormer := db.GetDefaultNewOrmer()
	if _, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hash).Filter("AppId", appId).All(&tHashOperates); err != nil {
		log.Error("query hashOpeate err:", err)
	} else {
		if len(tHashOperates) > 0 {
			return true
		}
	}
	return false
}

// 获取账户余额
type rRspModel struct {
	RspNo  string
	Assets []comm.Assets
}

func (v *VoucherController) GetAssets() {
	log.Debug("Assets...")
	//rspModel := &rRspModel{RspNo:comm.Err_OK}
	appid := v.GetString("appid")
	page := v.GetString("page")
	limit := v.GetString("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "20"
	}

	log.Debug("[Assets]appid=%v,page=%v,limit=%v", appid, page, limit)

	cfg := config.GConfig
	// 校验用户
	//ormer := db.GetDefaultNewOrmer()
	//var tTregist model.TRegist
	//if err := ormer.QueryTable(&model.TRegist{}).Filter("CaptainId", appid).Filter("Status", comm.REG_APPROVAL).One(&tTregist); err != nil {
	//	log.Error("get hashlist err:%s", err)
	//	v.retErrJSON(comm.Err_RDB)
	//	return
	//} else {
	//	if tTregist.PubKey == "" {
	//		log.Debug("tTregist.PubKey is nil")
	//		v.retErrJSON(comm.Err_USER)
	//		return
	//	}
	//}
	// 获取资产信息
	assetsRspModel := &rRspModel{RspNo: comm.Err_OK}
	result, err := httpcli.AssetsReq(appid, page, limit, cfg.Assets)
	if err != nil || result.RspNo != 0 {
		log.Error("assets", err)
		v.retErrJSON(comm.Err_APPSERVER)
		return
	}
	log.Debug("[Assets_result]: %v", result)
	assetsRspModel.Assets = result.Data
	v.Data["json"] = assetsRspModel

	log.Debug("GetAssets rspModel:%v", assetsRspModel)

	v.ServeJSON()
}

// 查询交易流水
type txRspModel struct {
	RspNo string
	Data  comm.TxInfoList
}

func (v *VoucherController) GetTradeHistory() {
	log.Debug("Trade History...")
	appid := v.GetString("appid")
	currency := v.GetString("currency")
	page := v.GetString("page")
	limit := v.GetString("limit")
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "20"
	}
	cfg := config.GConfig

	// 校验用户
	//ormer := db.GetDefaultNewOrmer()
	//var tTregist []*model.TRegist
	//if _, err := ormer.QueryTable(&model.TRegist{}).Filter("CaptainId", appid).Filter("Status", comm.REG_APPROVAL).All(&tTregist); err != nil {
	//	log.Error("get hashlist err:%s", err)
	//	v.retErrJSON(comm.Err_RDB)
	//	return
	//} else {
	//	if len(tTregist) == 0 {
	//		v.retErrJSON(comm.Err_USER)
	//		return
	//	}
	//}

	result, err := httpcli.TradeHistory(appid, currency, page, limit, cfg.TradeHistory)

	if err != nil {
		log.Error("trade history", err)
		v.retErrJSON(comm.Err_APPSERVER)
		return
	}
	rspModel := &txRspModel{RspNo: comm.Err_OK}

	rspModel.Data = result.Data
	v.Data["json"] = rspModel
	v.ServeJSON()
}

package controllers

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/db"
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

func (v *VoucherController) ApprovalList() {
	log.Debug("approvalList....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	captainId := v.GetString("captainid")
	rType := v.GetString("type")
	if rType == comm.APPROVAL_TYPE_0 { //待自己审批
		var tHashs []*model.THash
		if i, err := db.GetDefaultNewOrmer().QueryTable(&model.THash{}).All(&tHashs); err != nil {
			log.Error("get hashlist err:%s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		} else if i > 0 {
			for _, hashModel := range tHashs {
				if (hashModel.CaptainId == captainId && hashModel.Status == comm.HASH_STATUS_0) || (hashModel.CaptainId != captainId && hashModel.Status == comm.HASH_STATUS_3) { //待私钥审批员工发起，以及待私钥审批
					rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
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
				} else if hashModel.CaptainId != captainId && (hashModel.Status == comm.HASH_STATUS_4 || hashModel.Status == comm.HASH_STATUS_5 || hashModel.Status == comm.HASH_STATUS_6 || hashModel.Status == comm.HASH_STATUS_7) {
					rspModel.ApprovalInfos = append(rspModel.ApprovalInfos, comm.ApprovalInfo{Hash: hashModel.Hash, AppId: hashModel.AppId, Name: hashModel.Name, Status: hashModel.Status})
				}
			}
		}
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
			rspModel.HashOperates = append(rspModel.HashOperates, comm.HashOperate{AppId: hashOpeate.AppId, Option: hashOpeate.Option})
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
	role := v.GetString("role")
	publicKey := v.GetString("publickey")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_ADDKEY, AppId: appId, AppName: appName, ReqIpPort: reqIpPort, Role: role, PublicKey: publicKey}
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
	role := v.GetString("role")
	publicKey := v.GetString("publickey")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: operateType, AppId: appId, Password: password, Sign: sign, ReqIpPort: reqIpPort, Role: role, PublicKey: publicKey}
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("opeate marshal err: %s", err)
		v.retErrJSON(comm.Err_JSON)
		return
	} else {
		comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
	}
	//reqLog := &model.TReqLog{Amount: accountStr,CreateTime:time.Now(),ApplyTime:time.Now()}

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

	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	feeBig := new(big.Int)
	feeBig.SetString(fee, 10)

	withdrawModel := &model.TWithdraw{Hash: hashStr, WdHash: wdHashStr, To: recAddress, Category: category, Amount: amount, Fee: fee, Flow: flow, WdFlow: wdFlow, Status: comm.WITHDRAW_STATUS_0, CreateTime: time.Now()}
	if _, err := db.GetDefaultNewOrmer().Insert(withdrawModel); err != nil {
		log.Error("land to db err: %s", err)
	}

	grpcStream := &comm.GrpcStream{Type: comm.GRPC_WITHDRAW_REQ, Hash: common.HexToHash(hashStr), WdHash: common.HexToHash(wdHashStr), To: recAddress, Category: big.NewInt(category), Amount: amountBig, Fee: feeBig}
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("allow marshal err: %s", err)
	} else {
		comm.SendChanMsg(comm.SERVER_COMPANION, msg)
	}
	v.Data["json"] = rspModel
	v.ServeJSON()
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

	opinion := v.GetString("opinion")

	rspModel := &RspModel{RspNo: comm.Err_OK}

	//hash operate
	hashOperate := &model.THashOperate{Hash: hashStr, Type: comm.HASH_TYPE_ALLOW, Option: opinion, AppId: appId, Sign: sign, CreateTime: time.Now()}
	db.GetDefaultNewOrmer().Insert(hashOperate)

	hashModel := &model.THash{Hash: hashStr}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		log.Error("approval load err: %s", err)
		v.retErrJSON(comm.Err_RDB)
		return
	} else {
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
				hashModel.Status = comm.HASH_STATUS_5 //更新 私链已同意确认 私钥B、私钥C均同意
				log.Info("allow marshal failed: %s")
			} else {
				var tHashOperates []*model.THashOperate
				ormer := db.GetDefaultNewOrmer()
				if i, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hashStr).All(&tHashOperates); err != nil {
					log.Error("query hashOpeate err:", err)
				} else {
					if comm.RealTimeVoucherStatus.Total == i {
						allOpinion := true
						for _, tHashOperate := range tHashOperates {
							if tHashOperate.Option == comm.FALSE {
								allOpinion = false
								break
							}
						}
						if allOpinion {
							hashModel.Status = comm.HASH_STATUS_4
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
		if _, err = db.GetDefaultNewOrmer().Update(hashModel); err != nil {
			log.Error("update thash err: %s", err)
			v.retErrJSON(comm.Err_RDB)
			return
		}
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) ApprovalOperateList() {
	log.Debug("ApprovalOperateList....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	hash := v.GetString("hash")

	var tHashOperates []*model.THashOperate
	ormer := db.GetDefaultNewOrmer()
	if _, err := ormer.QueryTable(&model.THashOperate{}).Filter("Hash", hash).All(&tHashOperates); err != nil {
		log.Error("query hashOpeate err:", err)
	} else {
		for _, hashOpeate := range tHashOperates {
			rspModel.HashOperates = append(rspModel.HashOperates, comm.HashOperate{AppId: hashOpeate.AppId, Option: hashOpeate.Option})
		}
	}

	v.Data["json"] = rspModel
	v.ServeJSON()
}

func (v *VoucherController) DisAllow() {
	log.Debug("disallow....")
	appId := v.GetString("appid")
	sign := v.GetString("sign")
	hashStr := v.GetString("hash")
	if !strings.HasPrefix(hashStr, comm.HASH_PRIFIX) {
		v.retErrJSON(comm.Err_UNENABLE_PREFIX)
		return
	}
	opinion := v.GetString("opinion")
	rspModel := &RspModel{RspNo: comm.Err_OK}

	hashOperate := &model.THashOperate{Hash: hashStr, Type: comm.HASH_TYPE_DISALLOW, Option: opinion, AppId: appId, Sign: sign, CreateTime: time.Now()}
	db.GetDefaultNewOrmer().Insert(hashOperate)
	if opinion != comm.TRUE {
		log.Info("allow marshal failed: %s")
	} else {
		hashOperateCond := &model.THashOperate{Hash: hashStr, Type: comm.HASH_TYPE_DISALLOW, Option: comm.FALSE}
		var tHashOperates []*model.THashOperate
		ormer := db.GetDefaultNewOrmer()
		if _, err := ormer.QueryTable(hashOperateCond).All(&tHashOperates); err != nil {
			log.Error("query hashOpeate err:", err)
		} else {
			if len(tHashOperates) > 0 {
				grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_DISABLE_REQ, Hash: common.HexToHash(hashStr)}
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

//func (v *VoucherController) HashList() {
//	log.Debug("hashlist....")
//	rspModel := &RspModel{RspNo: comm.Err_OK}
//	if hashMap, err := comm.Ldb.GetPrifix([]byte(comm.HASHLIST_PRIFIX)); err != nil {
//		log.Error("get hashlist err:%s", err)
//	} else {
//		if len(hashMap) != int(comm.RealTimeVoucherStatus.HashCount) { //数量不足请求
//			grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
//			grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_HASH_LIST}
//			if msg, err := json.Marshal(grpcStream); err != nil {
//				log.Error("allow marshal err: %s", err)
//			} else {
//				comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
//			}
//		}
//		for _, hashBytes := range hashMap {
//			hashInfo := comm.HashInfo{}
//			if err = json.Unmarshal([]byte(hashBytes), hashInfo); err != nil {
//				v.retErrJSON(comm.Err_JSON)
//			} else {
//				rspModel.HashInfos = append(rspModel.HashInfos, hashInfo)
//			}
//		}
//	}
//	v.Data["json"] = rspModel
//	v.ServeJSON()
//}

func (v *VoucherController) TokenList() {
	log.Debug("tokenlist....")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	if tokenMap, err := comm.Ldb.GetPrifix([]byte(comm.TOKENLIST_PRIFIX)); err != nil {
		log.Error("get hashlist err:%s", err)
	} else {
		if len(tokenMap) != int(comm.RealTimeVoucherStatus.TokenCount) { //数量不足请求
			log.Info("token 数量不足", len(tokenMap), comm.RealTimeVoucherStatus.TokenCount)
			grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
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
	tokenName := v.GetString("tokenname")
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
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_TOKEN_ADD, TokenName: tokenName, Decimals: decimals, ContractAddr: contractAddr}
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
	contractAddr := v.GetString("contractaddr")
	rspModel := &RspModel{RspNo: comm.Err_OK}
	grpcStream := &comm.GrpcStream{Type: comm.GRPC_VOUCHER_OPR_REQ}
	grpcStream.VoucherOperate = &comm.Operate{Type: comm.VOUCHER_OPERATE_TOKEN_DEL, ContractAddr: contractAddr}
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

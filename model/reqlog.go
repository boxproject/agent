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

package model

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TReqLog struct {
	Id           int64  `orm:"pk"`
	ReqType      string //请求类型
	TransferType string //传递类型
	BlockNumber  int64  //区块号
	Hash         string //hash审批流
	WdHash       string //提现hash
	TxHash       string //交易hash
	Amount       string //金额
	Fee          string //手续费
	From         string //from地址
	To           string //to地址
	Category     int64  //币种分类
	Content      string //内容
	Status       string //hash状态

	ApplyTime  time.Time //申请时间
	CreateTime time.Time //创建时间
}

func (this *TReqLog) TableName() string {
	return "T_REQ_LOG"
}

type THashOperate struct {
	Id      int64 `orm:"pk"`
	AppId   string
	Type    string
	Hash    string //请求类型
	Option  string //同意拒绝
	Opinion string //操作意见[同意或拒绝说明]
	//Flow       string    //原始数据
	Sign       string    //签名
	CreateTime time.Time //创建时间
}

func (this *THashOperate) TableName() string {
	return "T_HASH_OPERATE"
}

type THash struct {
	Hash       string `orm:"pk"`
	AppId      string
	CaptainId  string
	Name       string
	Flow       string    //原始数据
	Sign       string    //签名
	Status     string    //状态
	CreateTime time.Time //创建时间
}

func (this *THash) TableName() string {
	return "T_HASH"
}

type TWithdraw struct {
	WdHash     string `orm:"pk"`
	Hash       string
	AppId      string
	To         string
	Amount     string
	Fee        string
	Category   int64
	Flow       string //原始数据
	WdFlow     string
	Sign       string    //签名
	Status     string    //状态
	CreateTime time.Time //创建时间
}

func (this *TWithdraw) TableName() string {
	return "T_WITHDRAW"
}

type TRegist struct {
	RegId          string `orm:"pk"`
	ApplyerId      string
	CaptainId      string
	ApplyerAccount string
	Msg            string
	Consent        string
	CipherText     string
	Status         string
	PubKey         string
	CreateTime     time.Time //创建时间
}

func (this *TRegist) TableName() string {
	return "T_REGIST"
}

func init() {
	orm.RegisterModel(new(TReqLog), new(THashOperate), new(THash), new(TWithdraw), new(TRegist))
}

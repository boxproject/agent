package comm

import (
	"github.com/boxproject/agent/db"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

const HASH_PRIFIX = "0x"
const (
	FALSE = "0"
	TRUE  = "1"
)

const (
	REQ_TYPE_ALL    = "all"    //全部
	REQ_TYPE_RANDOM = "random" //随机
)

const (
	ROUTER_TYPE_WEB  = "web"
	ROUTER_TYPE_GRPC = "grpc"
)

const (
	CHAN_MAX_SIZE = 100000
)

const (
	Err_OK              = "0"   //正确
	Err_UNKNOW_REQ_TYPE = "10"  //未知请求类型
	Err_JSON            = "11"  //JSON处理失败
	Err_LDB             = "12"  //leveldb处理失败
	Err_RDB             = "13"  //关系型数据库处理失败
	Err_SERVER_NOTFOUND = "100" //服务未发现
	Err_UNENABLE_PREFIX = "101" //非法hash前缀
	Err_UNENABLE_LENGTH = "102" //非法hash值长度
	Err_UNENABLE_AMOUNT = "103" //非法金额
	Err_DATE            = "104" //非法金额
)

//签名机状态
const (
	VOUCHER_STATUS_UNCONNETED = 0 //未连接
	VOUCHER_STATUS_UNCREATED  = 1 //未创建
	VOUCHER_STATUS_CREATED    = 2 //已创建
	VOUCHER_STATUS_DEPLOYED   = 3 //已发布
	VOUCHER_STATUS_STATED     = 4 //已启动
)

//上报类型
const (
	REQ_DEPOSIT      = "1" //充值上报
	REQ_WITHDRAW     = "2" //提现上报
	REQ_WITHDRAW_TX  = "3" //提现tx上报
	REQ_TOKEN_CHANGE = "4" //token新增
	REQ_REGIST       = "5" //token新增
)

//上报类型
const (
	HASH_TYPE_ALLOW    = "allow"    //同意
	HASH_TYPE_DISALLOW = "disallow" //拒绝

)

//grpc接口
const (
	GRPC_HASH_ADD_REQ     = "1"  //hash add申请
	GRPC_HASH_ADD_LOG     = "2"  //hans add 私链log
	GRPC_HASH_ENABLE_REQ  = "3"  //hash enable 申请
	GRPC_HASH_ENABLE_LOG  = "4"  //hash enable 私链log
	GRPC_HASH_DISABLE_REQ = "5"  //hash disable 申请
	GRPC_HASH_DISABLE_LOG = "6"  //hash disable 私链log
	GRPC_WITHDRAW_REQ     = "7"  //提现 申请
	GRPC_WITHDRAW_LOG     = "8"  //提现 私链log
	GRPC_DEPOSIT_WEB      = "9"  //充值上报
	GRPC_WITHDRAW_TX_WEB  = "10" //提现tx上报
	GRPC_WITHDRAW_WEB     = "11" //提现结果上报
	GRPC_VOUCHER_OPR_REQ  = "12" //签名机操作处理
	//GRPC_HASH_LIST_REQ    = "13" //审批流查询
	//GRPC_HASH_LIST_WEB    = "14" //审批流上报
	GRPC_TOKEN_LIST_WEB   = "15" //token上报
	GRPC_COIN_LIST_WEB    = "16" //coin上报
	GRPC_HASH_ENABLE_WEB  = "17" //hash enable 公链log
	GRPC_HASH_DISABLE_WEB = "18" //hash enable 公链log
)

const (
	VOUCHER_OPERATE_ADDKEY       = "0"  //添加公钥
	VOUCHER_OPERATE_CREATE       = "1"  //创建
	VOUCHER_OPERATE_DEPLOY       = "2"  //发布
	VOUCHER_OPERATE_START        = "3"  //启动
	VOUCHER_OPERATE_PAUSE        = "4"  //停止
	VOUCHER_OPERATE_HASH_ENABLE  = "5"  //hash同意
	VOUCHER_OPERATE_HASH_DISABLE = "6"  //hash拒绝
	VOUCHER_OPERATE_HASH_LIST    = "7"  //hash list 查询
	VOUCHER_OPERATE_TOKEN_ADD    = "8"  //token 添加
	VOUCHER_OPERATE_TOKEN_DEL    = "9"  //token 删除
	VOUCHER_OPERATE_TOKEN_LIST   = "10" //token list 查询
	VOUCHER_OPERATE_COIN         = "11" //coin 操作
)

const (
	HASHLIST_PRIFIX  = "hlp_"
	TOKENLIST_PRIFIX = "tlp_"

	REGIST_INFO_PRIFIX   = "rip_"
	APPROVAL_INFO_PRIFIX = "aip_"
)

const (
	HASH_STATUS_0 = "0" //待申请
	HASH_STATUS_1 = "1" //私钥已申请提交
	HASH_STATUS_2 = "2" //私钥已拒绝提交 私钥A拒绝
	HASH_STATUS_3 = "3" //私链已申请确认(日志)
	HASH_STATUS_4 = "4" //私链已同意确认 私钥B、私钥C均同意
	HASH_STATUS_5 = "5" //私链已拒绝确认 私钥B、私钥C有不同意
	HASH_STATUS_6 = "6" //私链已同意(日志)
	HASH_STATUS_7 = "7" //公链已同意
	HASH_STATUS_8 = "8" //公链已拒绝

	WITHDRAW_STATUS_0 = "0" //申请中
	WITHDRAW_STATUS_1 = "1" //私链已确认

)

const (
	CURRENCY_TYPE_BTC = "0"
	CURRENCY_TYPE_ETH = "1"
)

const (
	APPROVAL_TYPE_0 = "0" //待私钥审批
	APPROVAL_TYPE_1 = "1" //审批中及审批结束
	APPROVAL_TYPE_2 = "2" //审批完成
)

//grpc stream
type GrpcStream struct {
	Type           string
	BlockNumber    uint64 //区块号
	AppId          string //申请人
	Hash           common.Hash
	WdHash         common.Hash
	TxHash         string
	Amount         *big.Int
	Fee            *big.Int
	Account        string
	From           string
	To             string
	Category       *big.Int
	Flow           string //原始内容
	Sign           string //签名信息
	WdFlow         string //提现原始数据
	Status         string
	VoucherOperate *Operate
	ApplyTime      time.Time //申请时间
	TokenList      []*TokenInfo
	SignInfos      []*SignInfo
}

type TokenInfo struct {
	TokenName    string
	Decimals     int64
	ContractAddr string
	Category     int64
}

type SignInfo struct {
	AppId string
	Sign  string
}

type GrpcStreamModel struct {
	Msg []byte
}

//私钥-签名机操作
type Operate struct {
	Type         string
	AppId        string //appid
	AppName      string //app别名
	Hash         string
	Password     string
	ReqIpPort    string
	Role         string
	PublicKey    string
	TokenName    string
	Decimals     int64
	ContractAddr string
	CoinCategory int64  //币种分类
	CoinUsed     bool   //币种使用
	Sign         string //签名
}

type VoucherStatus struct {
	ServerStatus    int64            //系统状态
	Status          int64            //错误码状态
	Total           int64            //密钥数量
	HashCount       int64            //hash数量
	TokenCount      int64            //token数量
	Address         string           //账户地址
	ContractAddress string           //合约地址
	BtcAddress      string           //比特币地址
	D               int64            //随机数
	NodesAuthorized []NodeAuthorized //授权情况
	KeyStoreStatus  []KeyStoreStatu  //公钥添加状态
	CoinStatus      []CoinStatu      //币种状态
}

type NodeAuthorized struct {
	ApplyerId  string
	Authorized bool
}

type KeyStoreStatu struct {
	ApplyerId   string
	ApplyerName string
}

type CoinStatu struct {
	Name     string
	Category int64
	Decimals int64
	Used     bool
}

//请求数据
type VReq struct {
	ReqType      string
	Account      string
	From         string
	To           string
	Category     int64
	Amount       string
	WdHash       string
	TxHash       string
	CurrencyType string
	RegId        string
	Consent      string
	CipherText   string
	PubKey       string
	Status       string
}

type VRsp struct {
	Code    int
	Message string
}

//注册信息
type RegistInfo struct {
	RegId          string
	ApplyerId      string
	CaptainId      string
	ApplyerAccount string
	Msg            string
	Status         string
}

//注册信息列表
type RegistInfoList struct {
	RegistInfos []*RegistInfo
}

//hash信息
type ApprovalInfo struct {
	Hash      string //id
	Name      string //名称
	AppId     string //申请 appid
	CaptainId string // 私钥id
	Flow      string //原始内容
	Sign      string //签名内容
	Status    string //状态
}

//hash 审批操作
type HashOperate struct {
	AppId  string
	Option string //同意拒绝
}

//请求channel
var VReqChan chan *VReq = make(chan *VReq, CHAN_MAX_SIZE)

var RealTimeVoucherStatus *VoucherStatus = &VoucherStatus{ServerStatus: VOUCHER_STATUS_UNCONNETED}

var Ldb *db.Ldb

var SERVER_COMPANION = "companion"
var SERVER_VOUCHER = "voucher"

var MANAGER_SERVER_IPPORT string

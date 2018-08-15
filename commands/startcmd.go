package commands

import (
	"os"
	"os/signal"
	"syscall"

	logger "github.com/alecthomas/log4go"
	"github.com/astaxie/beego"
	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/config"
	"github.com/boxproject/agent/controllers"
	"github.com/boxproject/agent/db"
	"github.com/boxproject/agent/discovery"
	"github.com/boxproject/agent/httpcli"
	"github.com/boxproject/agent/server"
	"gopkg.in/urfave/cli.v1"
)

func StartCmd(c *cli.Context) error {
	logger.Debug("Starting agent service...")
	cfg, err := LoadConfig(c.String("c"), "config.json")
	if err != nil {
		logger.Error("Load config failed. cause: %v", err)
		return err
	}
	logger.Info("Load config.  %v", cfg)
	config.GConfig = cfg
	//init relational db
	if err = db.InitRDB(cfg.DbSource); err != nil {
		logger.Error("Init relational Db failed . cause: %v", err)
		return err
	}

	if comm.Ldb, err = db.InitLDB(cfg.LevelDbPath); err != nil {
		logger.Error("Init level Db failed . cause: %v", err)
		return err
	}

	//etcd discovery
	master := discovery.NewMaster(cfg)
	//go master.WatchWorkers()

	initServerName(cfg)
	//grpc
	go initGrpcServer(cfg, master)

	logger.Debug("######httpServer_cfg######", cfg)
	//提供http服务
	go httpServer(cfg)

	//上报程序
	repCli := httpcli.NewRepCli(cfg)
	go repCli.Start()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP, syscall.SIGKILL,
		syscall.SIGUSR1, syscall.SIGUSR2)
	<-signalCh

	logger.Info("companion has already been shutdown...")
	return nil
}

//http
func httpServer(cfg *config.Config) {
	beego.Router(ServiceName_KeyStore, &controllers.VoucherController{}, "get,post:KeyStore")
	beego.Router(ServiceName_Operate, &controllers.VoucherController{}, "get,post:Operate")
	beego.Router(ServiceName_Status, &controllers.VoucherController{}, "get,post:Status")
	//beego.Router(ServiceName_Hash_Add, &controllers.VoucherController{}, "get,post:AddHash")
	beego.Router(ServiceName_Allow, &controllers.VoucherController{}, "get,post:Allow")
	beego.Router(ServiceName_DisAllow, &controllers.VoucherController{}, "get,post:DisAllow")
	beego.Router(ServiceName_Regist_Add, &controllers.VoucherController{}, "get,post:AddRegist")
	beego.Router(ServiceName_Regist_Aproval, &controllers.VoucherController{}, "get,post:RegistAproval")
	beego.Router(ServiceName_Regist_List, &controllers.VoucherController{}, "get,post:RegistList")
	beego.Router(ServiceName_Approval_Add, &controllers.VoucherController{}, "get,post:AddApproval")
	beego.Router(ServiceName_Approval_Invalid, &controllers.VoucherController{}, "get,post:InvalidApproval")
	beego.Router(ServiceName_Approval_List, &controllers.VoucherController{}, "get,post:ApprovalList")
	beego.Router(ServiceName_Approval_Detail, &controllers.VoucherController{}, "get,post:ApprovalDetail")
	beego.Router(ServiceName_Approval_Operate_List, &controllers.VoucherController{}, "get,post:ApprovalOperateList")
	beego.Router(ServiceName_Token_Add, &controllers.VoucherController{}, "get,post:AddToken")
	beego.Router(ServiceName_Token_Del, &controllers.VoucherController{}, "get,post:DelToken")
	beego.Router(ServiceName_Token_List, &controllers.VoucherController{}, "get,post:TokenList")
	beego.Router(ServiceName_Coin, &controllers.VoucherController{}, "get,post:Coin")
	beego.Router(ServiceName_Coin_List, &controllers.VoucherController{}, "get,post:CoinList")
	beego.Router(ServiceName_Wtihdraw, &controllers.VoucherController{}, "get,post:Wtihdraw")
	beego.Router(ServiceName_Manager_Info, &controllers.VoucherController{}, "get,post:ManagerInfo")
	beego.Router(ServiceName_Assets, &controllers.VoucherController{}, "get,post:GetAssets")
	beego.Router(ServiceName_Trade_History, &controllers.VoucherController{}, "get,post:GetTradeHistory")
	beego.Run()
}

//init grpc
func initGrpcServer(cfg *config.Config, master *discovery.Master) error {
	return server.RpcServerStart(cfg, master)
}

func initServerName(cfg *config.Config) {
	comm.SERVER_COMPANION = cfg.CompanionServer
	comm.SERVER_VOUCHER = cfg.VoucherServer
	comm.MANAGER_SERVER_IPPORT = cfg.ManagerIpPort
}

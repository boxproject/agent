package discovery

import (
	//"time"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/config"
	//"github.com/coreos/etcd/client"
	//"golang.org/x/net/context"
)

type Master struct {
	//KeysAPI client.KeysAPI
}

func NewMaster(config *config.Config) *Master {
	//cfg := client.Config{
	//	Endpoints:               config.EtcdEndpoints,
	//	Transport:               client.DefaultTransport,
	//	HeaderTimeoutPerRequest: time.Second,
	//}
	//
	//etcdClient, err := client.New(cfg)
	//if err != nil {
	//	log.Error("Error: cannot connec to etcd:", err)
	//}

	master := &Master{
		//KeysAPI: client.NewKeysAPI(etcdClient),
	}
	return master
}

func (m *Master) Init() {
	log.Info("init master")
	comm.InitChanMap()
}

func (m *Master) AddWorker(key string, modelChan chan *comm.GrpcStreamModel, quitChan chan bool) {
	comm.AddChan(key, &comm.GrpcStreamChan{modelChan, quitChan})
	log.Debug("AddWorker...start", key)
	//m.KeysAPI.Set(context.Background(), "workers/"+key, "", &client.SetOptions{
	//	Dir: true,
	//	//TTL: time.Second * 10,
	//})
	log.Debug("AddWorker...end")
}

func (m *Master) UpdateWorker(key string) {
	log.Debug("UpdateWorker...start", key)
	//m.KeysAPI.Update(context.Background(), "workers/"+key, "")
}

func (m *Master) RemoveWorkerByName(key string) {
	m.RemoveWorkerByKey("workers/" + key)
}

func (m *Master) RemoveWorkerByKey(key string) {
	//if rsp, err := m.KeysAPI.Delete(context.Background(), key, &client.DeleteOptions{
	//	Dir: true,
	//}); err != nil {
	//	log.Error("delete err: %s", err)
	//} else {
	//	log.Debug("rsp:%v", rsp)
	//	comm.RomoveChan(key)
	//}
}

func (m *Master) RouteMsg(routerName string, msg []byte) {
	comm.SendChanMsg(routerName, msg)
}

func (m *Master) WatchWorkers() {
	//api := m.KeysAPI
	//watcher := api.Watcher("workers/", &client.WatcherOptions{
	//	Recursive: true,
	//})
	//for {
	//	res, err := watcher.Next(context.Background())
	//	if err != nil {
	//		log.Error("Error watch workers:", err)
	//		break
	//	}
	//	if res.Action == "expire" {
	//		log.Info("expire", res.Node.Key)
	//
	//		//m.RemoveWorkerByKey(res.Node.Key)
	//	} else if res.Action == "set" {
	//		log.Info("set")
	//		//comm.AddChanMap()
	//	} else if res.Action == "delete" {
	//		log.Info("delete")
	//		//m.RemoveWorker(res.Node.Key)
	//	}
	//}
}

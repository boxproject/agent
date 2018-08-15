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
	log.Debug("AddWorker...start", key)
	comm.AddChan(key, &comm.GrpcStreamChan{modelChan, quitChan})
	log.Debug("AddWorker...end")
}

func (m *Master) UpdateWorker(key string) {
	log.Debug("UpdateWorker...start", key)
	//m.KeysAPI.Update(context.Background(), "workers/"+key, "")
}

func (m *Master) RemoveWorkerByKey(key string) {
	comm.RomoveChan(key)
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

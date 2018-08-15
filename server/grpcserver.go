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

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/boxproject/agent/comm"
	"github.com/boxproject/agent/config"
	proto "github.com/boxproject/agent/pb"

	"encoding/json"

	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/db"
	"github.com/boxproject/agent/discovery"
	"github.com/boxproject/agent/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type replyServer struct {
	cfg    *config.Config
	master *discovery.Master
}

func (s *replyServer) Router(ctx context.Context, req *proto.RouterRequest) (*proto.RouterResponse, error) {
	response := &proto.RouterResponse{
		Code: comm.Err_OK,
	}
	if req.RouterType == comm.ROUTER_TYPE_GRPC { //grpc
		s.master.RouteMsg(req.RouterName, req.Msg)
	} else if req.RouterType == comm.ROUTER_TYPE_WEB { //http
		streamModel := &comm.GrpcStream{}
		if err := json.Unmarshal(req.Msg, streamModel); err != nil {
			logger.Error("router unmarshal err: %s", err)
			response.Code = comm.Err_JSON
		} else {
			switch streamModel.Type {
			case comm.GRPC_DEPOSIT_WEB: //充值上报
				comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_DEPOSIT, From: streamModel.From, To: streamModel.To, Category: streamModel.Category.Int64(), TxHash: streamModel.TxHash, Amount: streamModel.Amount.String()}
				break
			case comm.GRPC_WITHDRAW_TX_WEB:
				comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_WITHDRAW_TX, WdHash: streamModel.WdHash.Hex(), TxHash: streamModel.TxHash}
				break
			case comm.GRPC_WITHDRAW_WEB: //	提现结果上报
				comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_WITHDRAW, To: streamModel.To, Amount: streamModel.Amount.String(), WdHash: streamModel.WdHash.Hex(), TxHash: streamModel.TxHash}
				break
			//case comm.GRPC_HASH_LIST_WEB: //	审批流上报
			//	//db
			//	for _, hashInfo := range streamModel.HashList {
			//		if hashJson, err := json.Marshal(hashInfo); err != nil {
			//			logger.Error("json marshal err: %s", err)
			//		} else {
			//			if err = comm.Ldb.PutByte([]byte(comm.HASHLIST_PRIFIX+hashInfo.Hash), hashJson); err != nil {
			//				logger.Error("hashlist ldb err: %s", err)
			//			}
			//		}
			//	}
			//	//TODO 上报
			//	break
			case comm.GRPC_TOKEN_LIST_WEB: //	token上报
				//db
				if tokenMap, err := comm.Ldb.GetPrifix([]byte(comm.TOKENLIST_PRIFIX)); err != nil {
					logger.Info("get hashlist failed:%s", err)
				} else {
					for _, hashBytes := range tokenMap {
						tokenInfo := &comm.TokenInfo{}
						if err = json.Unmarshal([]byte(hashBytes), tokenInfo); err != nil {
							logger.Error("json unmarshal err: %s", err)
						} else {
							if err = comm.Ldb.DeleteByte([]byte(comm.TOKENLIST_PRIFIX + tokenInfo.ContractAddr)); err != nil {
								logger.Error("clear token err: %s", err)
							} else {
								logger.Info("clear token success.")
							}
						}
					}
				}

				logger.Debug("GRPC_TOKEN_LIST_WEB..................")
				for _, tokenInfo := range streamModel.TokenList {
					if tokenJson, err := json.Marshal(tokenInfo); err != nil {
						logger.Error("json marshal err: %s", err)
					} else {
						if err = comm.Ldb.PutByte([]byte(comm.TOKENLIST_PRIFIX+tokenInfo.ContractAddr), tokenJson); err != nil {
							logger.Error("tokenlist ldb err: %s", err)
						} else {
							logger.Info("tokenlist ldb  success")
						}
					}
				}
				comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_TOKEN_CHANGE, CurrencyType: comm.CURRENCY_TYPE_ETH}
				break
			case comm.GRPC_COIN_LIST_WEB:
				comm.VReqChan <- &comm.VReq{ReqType: comm.REQ_TOKEN_CHANGE, CurrencyType: comm.CURRENCY_TYPE_BTC}
				break
			case comm.GRPC_HASH_ADD_LOG: //私链add完成
				updateHashStatus(streamModel.Hash.Hex(), comm.HASH_STATUS_3)

				break
			case comm.GRPC_HASH_ENABLE_LOG: //私链确认完成
				if hashModel, err := updateHashStatus(streamModel.Hash.Hex(), comm.HASH_STATUS_6); err != nil {
					logger.Error("load err: %s", err)
				} else {
					if hashOperates, err := getHashOperate(streamModel.Hash.Hex(), comm.HASH_TYPE_ALLOW); err != nil {
						logger.Error("load err: %s", err)
					} else {
						grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_ENABLE_LOG, Hash: streamModel.Hash, AppId: hashModel.AppId, Flow: hashModel.Flow}
						for _, hashOperate := range hashOperates {
							grpcStream.SignInfos = append(grpcStream.SignInfos, &comm.SignInfo{AppId: hashOperate.AppId, Sign: hashOperate.Sign})
						}
						if msg, err := json.Marshal(grpcStream); err != nil {
							logger.Error("hash add marshal err: %s", err)
						} else {
							comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
						}
					}
				}

				break
			case comm.GRPC_HASH_DISABLE_LOG: //私链拒绝完成
				//if hashModel, err := updateHashStatus(streamModel.Hash.Hex(), comm.HASH_STATUS_3); err != nil {
				//	logger.Error("load err: %s", err)
				//} else {
				//	if hashOperates, err := getHashOperate(streamModel.Hash.Hex(), comm.HASH_TYPE_DISALLOW); err != nil {
				//		logger.Error("load err: %s", err)
				//	} else {
				//		grpcStream := &comm.GrpcStream{Type: comm.GRPC_HASH_DISABLE_LOG, Hash: streamModel.Hash, AppId: hashModel.AppId, Flow: hashModel.Flow}
				//		for _, hashOperate := range hashOperates {
				//			grpcStream.SignInfos = append(grpcStream.SignInfos, &comm.SignInfo{AppId: hashOperate.AppId, Sign: hashOperate.Sign})
				//		}
				//		if msg, err := json.Marshal(grpcStream); err != nil {
				//			logger.Error("hash add marshal err: %s", err)
				//		} else {
				//			comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
				//		}
				//	}
				//}
				break

			case comm.GRPC_HASH_ENABLE_WEB: //公链同意完成
				updateHashStatus(streamModel.Hash.Hex(), comm.HASH_STATUS_7)
				break

			case comm.GRPC_HASH_DISABLE_WEB: //公链拒绝完成
				//updateHashStatus(streamModel.Hash.Hex(), comm.HASH_STATUS_5)
				break

			case comm.GRPC_WITHDRAW_LOG: //私链转账申请
				logger.Debug("[router-------------->]:%v", streamModel)
				if wdModel, err := updateWithDrawStatus(streamModel.WdHash.Hex(), comm.WITHDRAW_STATUS_1); err != nil {
					logger.Error("load err: %s", err)
				} else {
					streamModel.Sign = wdModel.Sign
					streamModel.Flow = wdModel.Flow
					streamModel.WdFlow = wdModel.WdFlow
					if msg, err := json.Marshal(streamModel); err != nil {
						logger.Error("hash add marshal err: %s", err)
					} else {
						comm.SendChanMsg(comm.SERVER_VOUCHER, msg)
					}
				}
				break
			case comm.GRPC_CHECK_KEY_WEB: //密码检测
				//更新密码状态
				comm.SetPassStatus(streamModel.AppId, streamModel.Status)

				break
			default:
				response.Code = comm.Err_UNKNOW_REQ_TYPE
				logger.Info("unknow web type: %s", streamModel.Type)
			}
		}
	}
	return response, nil
}

func (s *replyServer) Heart(ctx context.Context, req *proto.HeartRequest) (*proto.HeartResponse, error) {
	//logger.Debug("Heart.....RouterType:", req.RouterType)
	response := &proto.HeartResponse{
		Code: comm.Err_OK,
	}
	s.master.UpdateWorker(req.ServerName + req.Name + req.Ip)
	if req.Msg != nil && s.cfg.VoucherName == req.ServerName { //签名机心跳信息上报
		voucherStatus := &comm.VoucherStatus{}
		if err := json.Unmarshal(req.Msg, voucherStatus); err != nil {
			logger.Error("heart unmarshal err: %s", err)
		} else {
			comm.RealTimeVoucherStatus = voucherStatus
		}
	}
	return response, nil
}

func (s *replyServer) Listen(stream proto.Synchronizer_ListenServer) error {
	defer logger.Info("grpc server listen end ......")
	logger.Info("grpc server listen start......")

	listReq, err := stream.Recv()
	logger.Debug("listReq: %s", listReq.ServerName)
	key := listReq.ServerName + listReq.Name + listReq.Ip
	grpcStreamChan := make(chan *comm.GrpcStreamModel, comm.CHAN_MAX_SIZE)
	quitCh := make(chan bool, 1)
	s.master.AddWorker(key, grpcStreamChan, quitCh)

	go func() { //监控连接情况
		_, err = stream.Recv()
		if err == io.EOF {
			logger.Debug("err EOF...", err)
			quitCh <- false
		}

		if err != nil {
			logger.Error("[LISTEN ERR] %v\n", err)
			quitCh <- false
		}
	}()

	for {
		select {
		case data, ok := <-grpcStreamChan:
			if ok {
				//logger.Debug("grpc send...", data.Msg)
				stream.Send(&proto.StreamRsp{Msg: data.Msg})
			} else {
				logger.Error("read from grpc channel failed")
			}
		case <-quitCh:
			{
				logger.Debug("recv quitch")
				if listReq.ServerName == comm.SERVER_VOUCHER {
					comm.RealTimeVoucherStatus.ServerStatus = comm.VOUCHER_STATUS_UNCONNETED
				}
				//TODO 初始化
				s.master.RemoveWorkerByKey(key)
				break
				//return nil
			}
		}
	}
	return nil
}

func RpcServerStart(cfg *config.Config, master *discovery.Master) error {
	logger.Debug("rpc server start....")

	master.Init() //init
	cred, err := loadCredential(cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	options := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: time.Minute, PermitWithoutStream: true}),
		grpc.Creds(cred),
	}

	//grpc.UseCompressor("gzip")
	server := grpc.NewServer(options...)
	proto.RegisterSynchronizerServer(server, &replyServer{cfg: cfg, master: master})
	reflection.Register(server)

	lis, err := net.Listen("tcp", cfg.GrpcSerPort)
	if err != nil {
		fmt.Printf("Can not listen to the port %v, cause: %v\n", cfg.GrpcSerPort, err)
		return err
	}
	if err = server.Serve(lis); err != nil {
		fmt.Printf("gRPC service error, cause: %v\n", err)
		return err
	}
	return nil
}

//加载服务端证书
func loadCredential(cfg *config.Config) (credentials.TransportCredentials, error) {

	cert, err := tls.LoadX509KeyPair(cfg.ServerCert, cfg.ServerKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(cfg.ClientCert)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	return credentials.NewTLS(config), nil
}

func addListerMember() {
	//etcd
	//map

}

func updateHashStatus(hash, status string) (*model.THash, error) {
	hashModel := &model.THash{Hash: hash}
	if err := db.GetDefaultNewOrmer().Read(hashModel); err != nil {
		logger.Error("approval load err: %s", err)
		return nil, err
	} else {
		hashModel.Status = status
		if _, err = db.GetDefaultNewOrmer().Update(hashModel); err != nil {
			logger.Error("approval update err: %s", err)
			return nil, err
		} else {
			logger.Info("approval update success")
		}
		return hashModel, nil
	}
}

func getHashOperate(hash, hashType string) ([]*model.THashOperate, error) {
	hashOp := &model.THashOperate{Hash: hash}
	hashOps := []*model.THashOperate{}
	if _, err := db.GetDefaultNewOrmer().QueryTable(hashOp).Filter("Hash", hash).All(&hashOps); err != nil {
		logger.Error("approval load err: %s", err)
		return nil, err
	} else {
		return hashOps, nil
	}
}

func updateWithDrawStatus(wdHash, status string) (*model.TWithdraw, error) {
	wdModel := &model.TWithdraw{WdHash: wdHash}
	if err := db.GetDefaultNewOrmer().Read(wdModel); err != nil {
		logger.Error("approval load err: %s", err)
		return nil, err
	} else {
		wdModel.Status = status
		if _, err = db.GetDefaultNewOrmer().Update(wdModel); err != nil {
			logger.Error("approval update err: %s", err)
			return nil, err
		} else {
			logger.Info("approval update success")
		}
		return wdModel, nil
	}
}

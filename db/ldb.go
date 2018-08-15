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

package db

import (
	logger "github.com/alecthomas/log4go"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Ldb struct {
	*leveldb.DB
}

//init
func InitLDB(dbFilePath string) (*Ldb, error) {
	logger.Info("initDb start... path:%v", dbFilePath)
	db, err := leveldb.OpenFile(dbFilePath, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB,
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		return nil, err
	}
	logger.Info("initDb end...")
	return &Ldb{db}, nil
}

func (this *Ldb) GetDb() *leveldb.DB {
	return this.DB
}

func (this *Ldb) PutByte(key, value []byte) error {
	return this.Put(key, value, nil)
}

func (this *Ldb) GetByte(key []byte) ([]byte, error) {
	return this.Get(key, nil)
}

func (this *Ldb) DeleteByte(key []byte) error {
	return this.Delete(key, nil)
}

func (this *Ldb) PutStrWithPrifix(keyPri, key, value string) error {
	return this.PutByte([]byte(keyPri+key), []byte(value))
}

//查询前缀
func (this *Ldb) GetPrifix(keyPrefix []byte) (map[string]string, error) {
	var resMap map[string]string = make(map[string]string)
	iter := this.NewIterator(util.BytesPrefix(keyPrefix), nil)
	if iter.Error() == leveldb.ErrNotFound {
		return nil, errors.New("no data")
	}
	if iter.Error() != nil {
		logger.Error("get prifix error")
		return nil, iter.Error()
	}
	for iter.Next() {
		resMap[string(iter.Key())] = string(iter.Value())
	}

	iter.Release()
	return resMap, nil
}

//del key
func (this *Ldb) DelKey(key []byte) error {
	if err := this.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

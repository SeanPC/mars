package storage

import (
	"encoding/json"
	"errors"

	"github.com/ouqiang/mars/internal/common"

	"github.com/ouqiang/mars/internal/common/recorder"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB 存储到levelDB
type LevelDB struct {
	db    *leveldb.DB
	queue *common.Queue
}

// NewLevelDB 创建levelDB
func NewLevelDB(db *leveldb.DB, queue *common.Queue) *LevelDB {
	l := &LevelDB{
		db:    db,
		queue: queue,
	}

	return l
}

// Get 获取transaction
func (l *LevelDB) Get(txId string) (*recorder.Transaction, error) {
	data, err := l.db.Get([]byte(txId), nil)
	if err != nil {
		return nil, err
	}
	tx := new(recorder.Transaction)
	err = json.Unmarshal(data, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Save 保存transaction
func (l *LevelDB) Save(tx *recorder.Transaction) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	var err error
	removeValue := l.queue.Add(tx.Id)
	if removeValue != nil {
		err = l.db.Delete([]byte(tx.Id), nil)
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	err = l.db.Put([]byte(tx.Id), data, nil)

	return err
}

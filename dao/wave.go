package dao

import (
	"context"
	"errors"
	"github.com/woshikedayaa/waveform-backend/dao/models"
	"gorm.io/gorm"
	"time"
)

type WaveFormDao struct {
	db *gorm.DB
}

func NewWaveFormDao() *WaveFormDao {
	wf := &WaveFormDao{db: Conn()}
	return wf
}

func (wf *WaveFormDao) Save(ctx context.Context, name string, data []byte) error {
	now := time.Now()
	w := &models.Wave{}
	if len(name) == 0 || len(data) == 0 {
		return &OpErr{
			op:    "create",
			table: w.TableName(),
			err:   errors.New("name和data 的数据不能为空"),
		}
	}
	w.Name = name
	w.CTime = now
	w.Data = now.Format("2006-01-02_15-04-05")
	if e := wf.db.WithContext(ctx).Create(w).Error; e != nil {
		ope := &OpErr{}

		ope.err = e
		ope.op = "create"
		ope.table = w.TableName()
		return ope
	}
	go w.Write(data)

	return nil
}

func (wf *WaveFormDao) GetLatest(ctx context.Context, count int) ([]models.Wave, error) {
	if count <= 0 {
		return nil, nil
	}
	res := make([]models.Wave, count)
	if e := wf.db.WithContext(ctx).Order("id desc").Limit(count).Find(res).Error; e != nil {
		return nil, e
	}
	return res, nil
}

func (wf *WaveFormDao) GetById(ctx context.Context, id int) (models.Wave, error) {
	w := models.Wave{}
	tx := wf.db.WithContext(ctx).Where("id = ?", id).First(&w)
	if tx.Error != nil {
		ope := &OpErr{}

		ope.err = tx.Error
		ope.table = w.TableName()
		ope.op = "query"
		return models.Wave{}, ope
	}
	return w, nil
}

func (wf *WaveFormDao) GetByName(ctx context.Context, name string) (models.Wave, error) {
	w := models.Wave{}
	tx := wf.db.WithContext(ctx).Where("name = ?", name).First(&w)
	if tx.Error != nil {
		ope := &OpErr{}

		ope.err = tx.Error
		ope.table = w.TableName()
		ope.op = "query"
		return models.Wave{}, ope
	}
	return w, nil
}

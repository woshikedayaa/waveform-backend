package dao

import (
	"context"
	"errors"
	"github.com/woshikedayaa/waveform-backend/dao/models"
	"github.com/woshikedayaa/waveform-backend/logf"
	"github.com/woshikedayaa/waveform-backend/pkg/wave"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type WaveFormDao struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewWaveFormDao() *WaveFormDao {
	wf := &WaveFormDao{db: Conn(), logger: logf.Open("dao/wave")}
	return wf
}

func (wf *WaveFormDao) Save(ctx context.Context, name string, data *wave.FullData) error {
	now := time.Now()
	w := &models.Wave{}
	if len(name) == 0 || len(data.Body) == 0 {
		return &OpErr{
			op:    "create",
			table: w.TableName(),
			err:   errors.New("name和data 的数据不能为空"),
		}
	}

	if e := wf.db.WithContext(ctx).Create(w).Error; e != nil {
		ope := &OpErr{}
		ope.err = e
		ope.op = "create"
		ope.table = w.TableName()
		return ope
	}
	j, err := data.Body.JSON()
	if err != nil {
		wf.logger.Sugar().Errorf("在保存 "+name+" 到数据库发送错误", zap.Error(err))
		// 只是显式声明一下不处理错误
		_ = wf.db.Update("delete_time", time.Now()).Error
		ope := &OpErr{}
		ope.err = errors.New("持久化失败 原因: data 无法转换成 json")
		ope.op = "create"
		ope.table = w.TableName()
		return ope
	}
	w.Name = name
	w.CTime = now
	w.Head = data.Header.String()
	go w.Write(j)

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

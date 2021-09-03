/**
 * @Time : 3/5/21 2:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package namespace

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type Call func() error

type Service interface {
	// 保存过更新
	Save(ctx context.Context, data *types.Namespace) (err error)
	SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error)
	FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error)
	FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error) {
	tx := s.db.Model(data).Begin()
	if err = tx.Save(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = call(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *service) FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error) {
	err = s.db.Model(&types.Namespace{}).Where("cluster_id = ? AND name = ?", clusterId, name).First(&res).Error
	return
}

func (s *service) FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error) {
	err = s.db.Model(&types.Namespace{}).Where("id IN (?)", ids).Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.Namespace) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}

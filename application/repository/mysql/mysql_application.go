package mysql

import (
	"context"
	"undina/application"
	"undina/domain"

	"gorm.io/gorm"
)

type mysqlApplicationRepository struct {
	Conn *gorm.DB
}

func NewmysqlApplicationRepository(Conn *gorm.DB) domain.ApplicationRepository {
	return &mysqlApplicationRepository{Conn}
}

func (m *mysqlApplicationRepository) Create(ctx context.Context, applicationRequest *domain.ApplicationRequest) (err error) {
	if err := m.Conn.Table("applications").Create(&applicationRequest).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlApplicationRepository) FetchAll(ctx context.Context, roomId int32) (res []domain.Application, err error) {
	var applications []domain.Application

	if err := m.Conn.Table("applications").Select("applications.created_at as application_date, applications.is_accepted, applications.application_message, users.id as user_id, users.name as user_name, users.rating as user_rating").
		Joins("JOIN users ON users.id = applications.user_id").
		Where("room_id = ?", roomId).Scan(&applications).Error; err != nil {
		return nil, err
	}

	return applications, nil
}

func (m *mysqlApplicationRepository) IsApplied(ctx context.Context, roomId int32, userId int32) (res bool, err error) {
	var application *domain.ApplicationRequest
	if err := m.Conn.Table("applications").Where("room_id = ? AND user_id = ?", roomId, userId).First(&application).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *mysqlApplicationRepository) AcceptApplication(ctx context.Context, roomId int32, userId int32) (err error) {
	res := m.Conn.Table("applications").Where("room_id = ? AND user_id = ?", roomId, userId).Update("is_accepted", true)
	if res.RowsAffected != 1 {
		return application.ErrApplicationNotFound
	}
	return nil
}

func (m *mysqlApplicationRepository) DeleteApplication(ctx context.Context, roomId int32, userId int32) (err error) {
	var application *domain.Application
	if err := m.Conn.Table("applications").Where("room_id = ? AND user_id = ?", roomId, userId).Delete(&application).Error; err != nil {
		return err
	}
	return nil
}

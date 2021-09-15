package mysql

import (
	"context"
	"undina/domain"

	"gorm.io/gorm"
)

type mysqlRoomRepository struct {
	Conn *gorm.DB
}

func NewmysqlRoomRepository(Conn *gorm.DB) domain.RoomRepository {
	return &mysqlRoomRepository{Conn}
}

func (m *mysqlRoomRepository) GetJoinedRooms(ctx context.Context, id int32) (res []domain.RoomItem, err error) {
	var rooms []domain.RoomItem
	if err := m.Conn.Table("participation").Select("service_providers.name as service_name, rooms.room_id, rooms.plan_name, rooms.room_status, rooms.is_public, participation.is_host, participation.payment_status").
		Joins("JOIN rooms ON rooms.room_id = participation.room_id").
		Joins("JOIN service_providers ON service_providers.id = rooms.service_id").
		Joins("JOIN plans ON plans.plan_name = rooms.plan_name AND plans.service_id = rooms.service_id").
		Where("participation.user_id = ?", id).Scan(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}

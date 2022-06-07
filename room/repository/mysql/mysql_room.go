package mysql

import (
	"context"
	"undina/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type mysqlRoomRepository struct {
	Conn *gorm.DB
}

func NewmysqlRoomRepository(Conn *gorm.DB) domain.RoomRepository {
	return &mysqlRoomRepository{Conn}
}

func (m *mysqlRoomRepository) Create(ctx context.Context, room *domain.Room) (roomId int32, err error) {
	logrus.Info("room matching deadline: ", room.MatchingDeadline)
	logrus.Info("room: ", room)
	if err := m.Conn.Table("rooms").Create(&room).Error; err != nil {
		return -1, err
	}

	return room.Id, nil
}

func (m *mysqlRoomRepository) GetPublicRooms(ctx context.Context) (res []domain.RoomPublic, err error) {
	var rooms []domain.RoomPublic
	if err := m.Conn.Table("rooms").Select("service_providers.name as service_name, rooms.room_id, rooms.plan_name, rooms.max_count, rooms.matching_deadline, rooms.public_message, COUNT(participation.user_id) as member_count, plans.cost, users.name as admin_name, users.rating as admin_rating").
		Joins("JOIN service_providers ON service_providers.id = rooms.service_id").
		Joins("JOIN participation ON participation.room_id = rooms.room_id").
		Joins("JOIN plans ON plans.plan_name = rooms.plan_name AND plans.service_id = rooms.service_id").
		Joins("JOIN users ON users.id = rooms.admin_id").
		Where("rooms.is_public = true").
		Group("rooms.room_id").Scan(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (m *mysqlRoomRepository) Update(ctx context.Context, roomId int32, room *domain.Room) (err error) {
	if err := m.Conn.Table("rooms").Where("room_id = ?", roomId).Updates(&room).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlRoomRepository) UpdateRoundId(ctx context.Context, roomId int32, roundId int32) (err error) {
	if err := m.Conn.Table("rooms").Where("room_id = ?", roomId).Update("round_id", roundId).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlRoomRepository) Delete(ctx context.Context, roomId int32) (err error) {
	var room *domain.Room
	if err := m.Conn.Table("rooms").Where("room_id = ?", roomId).Delete(&room).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlRoomRepository) Start(ctx context.Context, roomId int32) (err error) {
	if err := m.Conn.Table("rooms").Where("room_id = ?", roomId).Update("room_status", domain.START).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlRoomRepository) IsPublic(ctx context.Context, roomId int32) (res bool, err error) {
	if err := m.Conn.Table("rooms").Select("is_public").Where("room_id = ?", roomId).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

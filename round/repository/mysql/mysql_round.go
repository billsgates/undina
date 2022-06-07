package mysql

import (
	"context"
	"undina/domain"

	"gorm.io/gorm"
)

type mysqlRoundRepository struct {
	Conn *gorm.DB
}

func NewmysqlRoundRepository(Conn *gorm.DB) domain.RoundRepository {
	return &mysqlRoundRepository{Conn}
}

func (m *mysqlRoundRepository) AddRound(ctx context.Context, round *domain.Round) (res int32, err error) {
	if err := m.Conn.Create(&round).Error; err != nil {
		return -1, err
	}

	return round.RoundId, nil
}

func (m *mysqlRoundRepository) GetRound(ctx context.Context, roomId int32) (res *domain.RoundInfo, err error) {
	var roundInfo *domain.RoundInfo

	if err := m.Conn.Table("rooms").Select("rounds.starting_time, rounds.ending_time, rounds.round_interval, rounds.payment_deadline").Where("room_id = ?", roomId).Joins("JOIN rounds ON rounds.round_id = rooms.round_id").Take(&roundInfo).Error; err != nil {
		return roundInfo, err
	}

	return roundInfo, nil
}

func (m *mysqlRoundRepository) DeleteRound(ctx context.Context, roundId int32) (err error) {
	var round *domain.Round

	if err := m.Conn.Table("rounds").Where("round_id = ?", roundId).Delete(&round).Error; err != nil {
		return err
	}

	return nil
}

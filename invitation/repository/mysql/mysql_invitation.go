package mysql

import (
	"context"
	"undina/domain"

	"gorm.io/gorm"
)

type mysqlInvitationRepository struct {
	Conn *gorm.DB
}

func NewmysqlInvitationRepository(Conn *gorm.DB) domain.InvitationRepository {
	return &mysqlInvitationRepository{Conn}
}

func (m *mysqlInvitationRepository) GenerateInvitationCode(ctx context.Context, roomId int32, code string) (err error) {
	invitation := domain.Invitation{RoomId: roomId, InvitationCode: code, IsValid: true}

	if err := m.Conn.Table("invitation_codes").Select("room_id", "invitation_code", "is_valid").Create(&invitation).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlInvitationRepository) ConsumeInvitationCode(ctx context.Context, code string) (roomId int32, err error) {
	var invitation *domain.Invitation

	if err := m.Conn.Table("invitation_codes").Where("invitation_code = ? AND is_valid = true", code).First(&invitation).Error; err != nil {
		return -1, err
	}

	if err := m.Conn.Table("invitation_codes").Where("invitation_code = ?", code).Update("is_valid", false).Error; err != nil {
		return -1, err
	}

	return invitation.RoomId, nil
}

func (m *mysqlInvitationRepository) ResumeInvitationCode(ctx context.Context, code string) (err error) {
	if err := m.Conn.Table("invitation_codes").Where("invitation_code = ?", code).Update("is_valid", true).Error; err != nil {
		return err
	}

	return nil
}

func (m *mysqlInvitationRepository) GetInvitationCodes(ctx context.Context, roomId int32) (res []domain.InvitationCode, err error) {
	var codes []domain.InvitationCode

	if err := m.Conn.Table("invitation_codes").Select("invitation_code").Where("room_id = ? AND is_valid = true", roomId).Find(&codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

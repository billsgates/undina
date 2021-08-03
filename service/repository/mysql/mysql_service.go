package mysql

import (
	"context"
	"undina/domain"

	"gorm.io/gorm"
)

type mysqlServiceRepository struct {
	Conn *gorm.DB
}

func NewmysqlServiceRepository(Conn *gorm.DB) domain.ServiceRepository {
	return &mysqlServiceRepository{Conn}
}

func (m *mysqlServiceRepository) FetchAll(ctx context.Context) (res []domain.Service, err error) {
	var services []domain.Service

	servicesRow, err := m.Conn.Table("service_providers").Rows()
	if err != nil {
		return nil, err
	}
	defer servicesRow.Close()

	for servicesRow.Next() {
		var service domain.Service
		m.Conn.ScanRows(servicesRow, &service)

		plansRow, err := m.Conn.Table("plans").Where("plans.service_id = ?", service.Id).Rows()
		if err != nil {
			return nil, err
		}

		for plansRow.Next() {
			var plan domain.Plan
			m.Conn.ScanRows(plansRow, &plan)
			service.Plans = append(service.Plans, plan)
		}

		services = append(services, service)
	}
	return services, nil
}

func (m *mysqlServiceRepository) GetDetailByID(ctx context.Context, id string) (res []domain.ServiceDetail, err error) {
	var serviceDetails []domain.ServiceDetail
	if err := m.Conn.Table("service_providers").Select("service_providers.id, service_providers.name, plans.plan_name, plans.cost, plans.max_count").Joins("left join plans on plans.service_id = service_providers.id").Where("service_providers.id = ?", id).Scan(&serviceDetails).Error; err != nil {
		return nil, err
	}

	return serviceDetails, nil
}

func (m *mysqlServiceRepository) GetPlanByKey(ctx context.Context, planName string, serviceId string) (res *domain.Plan, err error) {
	var plan *domain.Plan
	if err := m.Conn.Table("plans").Where("plans.plan_name = ? AND plans.service_id = ?", planName, serviceId).Find(&plan).Error; err != nil {
		return nil, err
	}

	return plan, nil
}

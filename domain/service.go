package domain

import (
	"context"
)

type Service struct {
	Id    int32  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Plans []Plan `json:"plans,omitempty" gorm:"-"`
}

type ServiceDetail struct {
	Id       int32  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	Cost     int32  `json:"cost,omitempty"`
	MaxCount int32  `json:"max_count,omitempty"`
}

type Plan struct {
	PlanName string `json:"plan_name,omitempty"`
	Cost     int32  `json:"cost,omitempty"`
	MaxCount int32  `json:"max_count,omitempty"`
}

type ServiceRepository interface {
	FetchAll(ctx context.Context) ([]Service, error)
	GetDetailByID(ctx context.Context, id string) ([]ServiceDetail, error)
	GetPlanByKey(ctx context.Context, planName string, serviceId string) (*Plan, error)
}

type ServiceUsecase interface {
	FetchAll(ctx context.Context) ([]Service, error)
	GetDetailByID(ctx context.Context, id string) ([]ServiceDetail, error)
	GetPlanByKey(ctx context.Context, planName string, serviceId string) (*Plan, error)
}

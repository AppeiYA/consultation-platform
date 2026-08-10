package http

import "github.com/AppeiYA/consultation-platform/internal/consultant"

type ConsultantHandler struct {
	ConsultantModule consultant.Module
}

func NewConsultantHandler(consultantModule consultant.Module) *ConsultantHandler {
	return &ConsultantHandler{
		ConsultantModule: consultantModule,
	}
}
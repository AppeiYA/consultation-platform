package http

import "github.com/AppeiYA/consultation-platform/internal/consultationcase"

type ConsultationCaseHandler struct {
	consultationCaseModule *consultationcase.Module
}

func NewConsultationCaseHandler(consultationCaseModule *consultationcase.Module) *ConsultationCaseHandler {
	return &ConsultationCaseHandler{
		consultationCaseModule: consultationCaseModule,
	}
}
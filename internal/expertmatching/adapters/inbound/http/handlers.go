package http

import "github.com/AppeiYA/consultation-platform/internal/expertmatching"

type ExpertMatchingHandler struct {
	expertMatchingModule *expertmatching.Module
}

func NewExpertMatchingHandler(expertMatchingModule *expertmatching.Module) *ExpertMatchingHandler {
	return &ExpertMatchingHandler{
		expertMatchingModule: expertMatchingModule,
	}
}

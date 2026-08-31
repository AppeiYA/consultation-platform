package consultant_expertise

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ExpertiseModel struct {
	ID           string    `db:"id"`
	ConsultantID string    `db:"consultant_id"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
}

func (m *ExpertiseModel) ToDomain() (*domain.Expertise, error) {
	return domain.NewExpertise(m.ID, m.ConsultantID, m.Name)
}

func FromDomainToModel(expertise *domain.Expertise) *ExpertiseModel {
	return &ExpertiseModel{
		ID:           expertise.ID(),
		ConsultantID: expertise.ConsultantID(),
		Name:         expertise.Name(),
	}
}

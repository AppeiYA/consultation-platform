package domain

import "time"

const ConsultantIDPrefix = "con"

type Consultant struct {
	id                string
	userID            string
	professionID        string
	displayName       DisplayName
	bio               Bio
	yearsExperience   YearsExperience
	isAcceptingClients bool
	createdAt         time.Time
	updatedAt         time.Time
}

func NewConsultant(
	id string,
	userID string,
	professionID string,
	displayName DisplayName,
	bio Bio,
	yearsExperience YearsExperience,
	now time.Time,
) *Consultant {
	return &Consultant{
		id: id,
		userID: userID,
		professionID: professionID,
		displayName: displayName,
		bio: bio,
		yearsExperience: yearsExperience,
		isAcceptingClients: true,
		createdAt: now,
		updatedAt: now,
	}
}

func ReconstitueConsultant(
	id string,
	userID string,
	professionID string,
	displayName string,
	bio string,
	yearsExperience int,
	isAcceptingClients bool,
	createdAt time.Time,
	updatedAt time.Time,
) (*Consultant, error) {
	validDisplayName, err := NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}
	validBio, err := NewBio(bio)
	if err != nil {
		return nil, err
	}
	validYearsExperience, err := NewYearsExperience(yearsExperience)
	if err != nil {
		return nil, err
	}
	return &Consultant{
		id: id,
		userID: userID,
		professionID: professionID,
		displayName: validDisplayName,
		bio: validBio,
		yearsExperience: validYearsExperience,
		isAcceptingClients: isAcceptingClients,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (c *Consultant) ID() string {
	return c.id
}

func (c *Consultant) UserID() string {
	return c.userID
}

func (c *Consultant) ProfessionID() string {
	return c.professionID
}

func (c *Consultant) DisplayName() DisplayName {
	return c.displayName
}

func (c *Consultant) Bio() Bio {
	return c.bio
}

func (c *Consultant) YearsExperience() YearsExperience {
	return c.yearsExperience
}

func (c *Consultant) IsAcceptingClients() bool {
	return c.isAcceptingClients
}

func (c *Consultant) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Consultant) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Consultant) UpdateProfile(
	profession Profession,
	displayName DisplayName,
	bio Bio,
	yearsExperience YearsExperience,
	now time.Time,
) {
	c.professionID = profession.ID()
	c.displayName = displayName
	c.bio = bio
	c.yearsExperience = yearsExperience
	c.updatedAt = now
}

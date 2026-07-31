package domain

import "time"

type Consultant struct {
	id                string
	userID            string
	profession        Profession
	displayName       DisplayName
	bio               Bio
	yearsExperience   YearsExperience
	isVerified        bool
	isAcceptingClients bool
	createdAt         time.Time
	updatedAt         time.Time
}

func NewConsultant(
	id string,
	userID string,
	profession Profession,
	displayName DisplayName,
	bio Bio,
	yearsExperience YearsExperience,
	now time.Time,
) *Consultant {
	return &Consultant{
		id: id,
		userID: userID,
		profession: profession,
		displayName: displayName,
		bio: bio,
		yearsExperience: yearsExperience,
		isVerified: false,
		isAcceptingClients: true,
		createdAt: now,
		updatedAt: now,
	}
}

func (c *Consultant) ID() string {
	return c.id
}

func (c *Consultant) UserID() string {
	return c.userID
}

func (c *Consultant) Profession() Profession {
	return c.profession
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

func (c *Consultant) IsVerified() bool {
	return c.isVerified
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

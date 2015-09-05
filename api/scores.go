package api

import "log"

type ScoresInput struct {
	Distance    float64
	CountryCode string
}

type Scores struct {
	Distance float64 `json:"distance"`
	Country  float64 `json:"country"`
}

func NewScores(input *ScoresInput) *Scores {
	scores := &Scores{
		Distance: 0,
		Country:  0,
	}

	scores.calculateCountryScore(input.CountryCode)
	scores.calculateDistanceScore(input.Distance)

	return scores
}

func (s *Scores) calculateCountryScore(isoCountry string) {
	if Countries[isoCountry] {
		if Countries[isoCountry] == true {
			s.Country = 0.5
		} else {
			s.Country = -0.5
		}
	} else {
		s.Country = 0
	}

	log.Printf("[scores] Country: %f\n", s.Country)
}

func (s *Scores) calculateDistanceScore(distance float64) {
	s.Distance = distance / 24901.0

	log.Printf("[scores] Distance: %f\n", s.Distance)
}

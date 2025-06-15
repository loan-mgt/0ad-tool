package model

// Update this
type Unit struct {
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	Rank         string            `json:"rank,omitempty"`
	GenericName  string            `json:"genericName,omitempty"`
	SpecificName string            `json:"specificName,omitempty"`
	Icon         string            `json:"icon,omitempty"`
	Health       int               `json:"health,omitempty"`
	Attack       map[string]Attack `json:"attack,omitempty"`
	Cost         map[string]int    `json:"cost,omitempty"`
	Promotion    string            `json:"promotion,omitempty"`
	VisualActor  string            `json:"visualActor,omitempty"`
	Parent       string            `json:"parent,omitempty"`
}

type Attack struct {
	AttackName string         `json:"attackName,omitempty"`
	Damage     map[string]int `json:"damage,omitempty"`
	MaxRange   int            `json:"maxRange,omitempty"`
}

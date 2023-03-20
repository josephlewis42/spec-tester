package v1

const (
	TypeSpecification = "Specification"
)

type Specification struct {
	Base `json:",inline"`

	Sections   []SpecificationSection   `json:"sections"`
	Exclusions []SpecificationExclusion `json:"exclusions,omitempty"`
}

type SpecificationSection struct {
	Base         `json:",inline"`
	TestSelector LabelSelector          `json:"testSelector"`
	Optional     bool                   `json:"optional,omitempty"`
	Sections     []SpecificationSection `json:"sections,omitempty"`
}

type SpecificationExclusion struct {
	Base         `json:",inline"`
	TestSelector LabelSelector `json:"testSelector"`
}

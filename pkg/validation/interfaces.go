package validation

// Validator validates input parameters
type Validator interface {
	ValidateHost(host string) error
	ValidateAction(action string, validActions []string) error
}

// ActionValidator validates control actions
type ActionValidator interface {
	IsValid(action string, validActions []string) bool
	GetValidActions() []string
}

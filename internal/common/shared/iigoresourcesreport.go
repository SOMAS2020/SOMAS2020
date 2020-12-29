package shared

// ResourcesReport is a struct returned by the Client when asked to report it's resources.
// The client can choose to report the resources by setting the Reported entry to true, along with ReportedAmount.
// If client doesn't want to share the information about its resources with president, it can set Reported to false.
type ResourcesReport struct {
	ReportedAmount Resources
	Reported       bool
}

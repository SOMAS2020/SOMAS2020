package shared

// Coordinate is a floating point number that can be used represent a position on the map
type Coordinate = float64

// Magnitude defines the severity of a disaster
type Magnitude = float64

// SpatialPDFType is an enum for xy prob. density function types for disaster simulation
type SpatialPDFType int

const (
	// Uniform xy distribution: disaster peak occurs uniformally over xy bounds of env
	Uniform SpatialPDFType = iota

	// add other PDFs here post-MVP
)

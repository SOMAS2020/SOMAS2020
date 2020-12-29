package shared

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

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

	// DO NOT TOUCH THIS
	spatialPDFTypeEnd
)

func (s SpatialPDFType) String() string {
	strings := [...]string{"Uniform"}
	if s >= 0 && int(s) < len(strings) {
		return strings[s]
	}
	return fmt.Sprintf("UNKNOWN ForageType '%v'", int(s))
}

// GoString implements GoStringer
func (s SpatialPDFType) GoString() string {
	return s.String()
}

// MarshalText implements TextMarshaler
func (s SpatialPDFType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(s.String())
}

// MarshalJSON implements RawMessage
func (s SpatialPDFType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(s.String())
}

// ParseSpatialPDFType gets the SpatialPDFType based on the number
func ParseSpatialPDFType(x int) SpatialPDFType {
	if x >= 0 && SpatialPDFType(x) < spatialPDFTypeEnd {
		return SpatialPDFType(x)
	}
	log.Printf("Unknown SpatialPDFType specified: '%v'\nUse --help. Defaulting to %v", x, Uniform)
	return Uniform
}

// HelpSpatialPDFType returns a help string for SpatialPDFType
func HelpSpatialPDFType() string {
	help := "Set x,y prob. distribution of the disaster's epicentre (more post MVP)\n"

	for i := 0; i < int(spatialPDFTypeEnd); i++ {
		help += fmt.Sprintf("%v: %v\n", i, SpatialPDFType(i))
	}

	return help
}

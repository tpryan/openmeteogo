package openmateo

// TemperatureUnit defines the unit for temperature values.
type TemperatureUnit int64

const (
	// Celsius is the default temperature unit.
	Celsius TemperatureUnit = iota
	// Fahrenheit temperature unit.
	Fahrenheit
)

// String returns the string representation of the TemperatureUnit for the API.
func (t TemperatureUnit) String() string {
	switch t {
	case Celsius:
		return "celsius"
	case Fahrenheit:
		return "fahrenheit"
	}
	return "unknown"
}

// WindSpeedUnit defines the unit for wind speed values.
type WindSpeedUnit int64

const (
	// KMH is the default wind speed unit (kilometers per hour).
	KMH WindSpeedUnit = iota
	// MS is meters per second.
	MS
	// MPH is miles per hour.
	MPH
	// KN is knots.
	KN
)

// String returns the string representation of the WindSpeedUnit for the API.
func (w WindSpeedUnit) String() string {
	switch w {
	case KMH:
		return "kmh"
	case MS:
		return "ms"
	case MPH:
		return "mph"
	case KN:
		return "kn"
	}
	return "unknown"
}

// PrecipitationUnit defines the unit for precipitation values.
type PrecipitationUnit int64

const (
	// MM is the default precipitation unit (millimeters).
	MM PrecipitationUnit = iota
	// IN is inches.
	IN
)

// String returns the string representation of the PrecipitationUnit for the API.
func (p PrecipitationUnit) String() string {
	switch p {
	case MM:
		return "mm"
	case IN:
		return "in"
	}
	return "unknown"
}

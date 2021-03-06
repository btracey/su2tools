package nondimensionalize

import (
	"math"
)

// package for helping with running SU^2 non-dimensionally

const (
	GasConstantAir = 287.87
)

// Values returns the freestream pressure and density from the given Reynolds number,
// Mach number, freestream temperature and gamma
// TODO: Add tests
func Values(temperature, reynolds, Mach, gasConstant, length, gamma float64) (pressure, density float64) {
	speedOfSound := SpeedOfSound(gamma, gasConstant, temperature)
	ViscosityFreestream := 1.853E-5 * (math.Pow(temperature/300.0, 3.0/2.0) * (300.0 + 110.3) / (temperature + 110.3))
	VelocityFreestream := Mach * speedOfSound
	density = reynolds * ViscosityFreestream / (VelocityFreestream * length)
	pressure = density * gasConstant * temperature
	return
}

func TotalTemperature(temperature, mach, gamma float64) float64 {
	return temperature * (1 + ((gamma-1)/2)*math.Pow(mach, 2))
}

func TotalPressure(pressure, mach, gamma float64) float64 {
	return pressure * math.Pow((1+((gamma-1)/2)*math.Pow(mach, 2)), gamma/(gamma-1))
}

func SpeedOfSound(gamma, gasConstant, temperature float64) float64 {
	return math.Sqrt(gamma * gasConstant * temperature)
}

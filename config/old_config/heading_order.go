package config

// This file is autogenerated using github.com/btracey/su2tools/config_writer/parser/write_options_structs.go based on the output from parse_config.py
import "github.com/btracey/su2tools/config/common"

var categoryOrder map[common.ConfigCategory]int = map[common.ConfigCategory]int{
	"Problem Definition":                           0,
	"Boundary Markers":                             1,
	"Grid adaptation":                              2,
	"Time-marching":                                3,
	"Linear solver definition":                     4,
	"Dynamic mesh definition":                      5,
	"Wind Gust":                                    6,
	"Convergence":                                  7,
	"Multi-grid":                                   8,
	"Spatial Discretization":                       9,
	"Convect Option:  CONV_NUM_METHOD_FLOW":        10,
	"Convect Option:  CONV_NUM_METHOD_ADJ":         11,
	"Convect Option:  CONV_NUM_METHOD_TURB":        12,
	"Convect Option:  CONV_NUM_METHOD_ADJTURB":     13,
	"Convect Option:  CONV_NUM_METHOD_LIN":         14,
	"Convect Option:  CONV_NUM_METHOD_ADJLEVELSET": 15,
	"Convect Option:  CONV_NUM_METHOD_TNE2":        16,
	"Convect Option:  CONV_NUM_METHOD_ADJTNE2":     17,
	"Adjoint and Gradient":                         18,
	"Input/output files and formats":               19,
	"Equivalent Area":                              20,
	"Freestream Conditions":                        21,
	"Reference Conditions":                         22,
	"Reacting Flow":                                23,
	"Free surface simulation":                      24,
	"Grid deformation":                             25,
	"Rotorcraft problem":                           26,
	"FEA solver":                                   27,
	"Wave solver":                                  28,
	"Heat solver":                                  29,
	"ML Turb Options":                              30,
}
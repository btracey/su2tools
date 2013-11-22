package config

// This file is autogenerated using github.com/btracey/su2tools/config_writer/parser/write_options_structs.go based on the output from parse_config.py

// Options is a struct containing all of the possible options in SU^2
type Options struct {
	// Adjoint type
	RegimeType string
	// Write extra output
	ExtraOutput bool
	// Physical governing equations
	PhysicalProblem string
	// Mathematical problem
	MathProblem string
	// Specify turbulence model
	KindTurbModel string
	// Location of the turb model itself
	MlTurbModelFile string
	// Location of the check for the proper loading of the turbulence model
	MlTurbModelCheckFile string
	// Specify transition model
	KindTransModel string
	// Axisymmetric simulation
	Axisymmetric bool
	// Add the gravity force
	GravityForce bool
	// Perform a low fidelity simulation
	LowFidelitySimulation bool
	// Restart solution from native solution file
	RestartSol bool
	// Write a tecplot file for each partition
	VisualizePart bool
	// Marker(s) of the surface in the surface flow solution file
	MarkerPlotting string
	// Marker(s) of the surface where evaluate the non-dimensional coefficients
	MarkerMonitoring string
	// Marker(s) of the surface where objective function (design problem) will be evaluated
	MarkerDesigning string
	// Euler wall boundary marker(s)
	MarkerEuler string
	// Far-field boundary marker(s)
	MarkerFar string
	// Symmetry boundary condition
	MarkerSym string
	// Symmetry boundary condition
	MarkerPressure string
	// Near-Field boundary condition
	MarkerNearfield string
	// Zone interface boundary marker(s)
	MarkerInterface string
	// Dirichlet boundary marker(s)
	MarkerDirichlet string
	// Neumann boundary marker(s)
	MarkerNeumann string
	// poisson dirichlet boundary marker(s)
	ElecDirichlet string
	// poisson neumann boundary marker(s)
	ElecNeumann string
	// Custom boundary marker(s)
	MarkerCustom string
	// No description
	MarkerPeriodic string
	// Inlet boundary type
	InletType string
	// No description
	MarkerInlet string
	// No description
	MarkerSupersonicInlet string
	// No description
	MarkerOutlet string
	// No description
	MarkerIsothermal string
	// No description
	MarkerHeatflux string
	// No description
	MarkerNacelleInflow string
	// Engine subsonic intake region
	SubsonicNacelleInflow bool
	// No description
	MarkerNacelleExhaust string
	// Displacement boundary marker(s)
	MarkerNormalDispl string
	// Load boundary marker(s)
	MarkerNormalLoad string
	// Flow load boundary marker(s)
	MarkerFlowload string
	// Damping factor for engine inlet condition
	DampNacelleInflow float64
	// Kind of grid adaptation
	KindAdapt string
	// Percentage of new elements (% of the original number of elements)
	NewElems float64
	// Scale factor for the dual volume
	DualvolPower float64
	// Use analytical definition for surfaces
	AnalyticalSurfdef string
	// Before each computation, implicitly smooth the nodal coordinates
	SmoothGeometry bool
	// Adapt the boundary elements
	AdaptBoundary bool
	// Divide rectangles into triangles
	DivideElements bool
	// Unsteady simulation
	UnsteadySimulation string
	// Unsteady farfield boundaries
	UnsteadyFarfield bool
	// Courant-Friedrichs-Lewy condition of the finest grid
	CflNumber float64
	// No description
	CflRamp []float64
	// Reduction factor of the CFL coefficient in the adjoint problem
	AdjCflReduction float64
	// Reduction factor of the CFL coefficient in the level set problem
	TurbCflReduction float64
	// Reduction factor of the CFL coefficient in the turbulent adjoint problem
	AdjturbCflReduction float64
	// Number of total iterations
	ExtIter float64
	// Runge-Kutta alpha coefficients
	RkAlphaCoeff string
	// Time Step for dual time stepping simulations (s)
	UnstTimestep float64
	// Total Physical Time for dual time stepping simulations (s)
	UnstTime float64
	// Unsteady Courant-Friedrichs-Lewy number of the finest grid
	UnstCflNumber float64
	// Number of internal iterations (dual time method)
	UnstIntIter float64
	// Integer number of periodic time instances for Time Spectral
	TimeInstances float64
	// Iteration number to begin unsteady restarts (dual time method)
	UnstRestartIter float64
	// Starting direct solver iteration for the unsteady adjoint
	UnstAdjointIter float64
	// Time discretization
	TimeDiscreFlow string
	// Time discretization
	TimeDiscreTne2 string
	// Time discretization
	TimeDiscreAdjtne2 string
	// Time discretization
	TimeDiscreAdjlevelset string
	// Time discretization
	TimeDiscreAdj string
	// Time discretization
	TimeDiscreLin string
	// Time discretization
	TimeDiscreTurb string
	// Time discretization
	TimeDiscreAdjturb string
	// Time discretization
	TimeDiscreWave string
	// Time discretization
	TimeDiscreFea string
	// Time discretization
	TimeDiscreHeat string
	// Time discretization
	TimeDiscrePoisson string
	// Linear solver for the implicit, mesh deformation, or discrete adjoint systems
	LinearSolver string
	// Preconditioner for the Krylov linear solvers
	LinearSolverPrec string
	// Minimum error threshold for the linear solver for the implicit formulation
	LinearSolverError float64
	// Maximum number of iterations of the linear solver for the implicit formulation
	LinearSolverIter float64
	// Relaxation of the linear solver for the implicit formulation
	LinearSolverRelax float64
	// Roe-Turkel preconditioning for low Mach number flows
	RoeTurkelPrec bool
	// Time Step for dual time stepping simulations (s)
	MinRoeTurkelPrec float64
	// Time Step for dual time stepping simulations (s)
	MaxRoeTurkelPrec float64
	// Linear solver for the turbulent adjoint systems
	AdjturbLinSolver string
	// Preconditioner for the turbulent adjoint Krylov linear solvers
	AdjturbLinPrec string
	// Minimum error threshold for the turbulent adjoint linear solver for the implicit formulation
	AdjturbLinError float64
	// Maximum number of iterations of the turbulent adjoint linear solver for the implicit formulation
	AdjturbLinIter float64
	// Mesh motion for unsteady simulations
	GridMovement bool
	// Type of mesh motion
	GridMovementKind []string
	// Marker(s) of moving surfaces (MOVING_WALL or DEFORMING grid motion).
	MarkerMoving string
	// Mach number (non-dimensional, based on the mesh velocity and freestream vals.)
	MachMotion float64
	// Coordinates of the rigid motion origin
	MotionOriginX string
	// Coordinates of the rigid motion origin
	MotionOriginY string
	// Coordinates of the rigid motion origin
	MotionOriginZ string
	// Translational velocity vector (m/s) in the x, y, & z directions (RIGID_MOTION only)
	TranslationRateX string
	// Translational velocity vector (m/s) in the x, y, & z directions (RIGID_MOTION only)
	TranslationRateY string
	// Translational velocity vector (m/s) in the x, y, & z directions (RIGID_MOTION only)
	TranslationRateZ string
	// Angular velocity vector (rad/s) about x, y, & z axes (RIGID_MOTION only)
	RotationRateX string
	// Angular velocity vector (rad/s) about x, y, & z axes (RIGID_MOTION only)
	RotationRateY string
	// Angular velocity vector (rad/s) about x, y, & z axes (RIGID_MOTION only)
	RotationRateZ string
	// Pitching angular freq. (rad/s) about x, y, & z axes (RIGID_MOTION only)
	PitchingOmegaX string
	// Pitching angular freq. (rad/s) about x, y, & z axes (RIGID_MOTION only)
	PitchingOmegaY string
	// Pitching angular freq. (rad/s) about x, y, & z axes (RIGID_MOTION only)
	PitchingOmegaZ string
	// Pitching amplitude (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingAmplX string
	// Pitching amplitude (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingAmplY string
	// Pitching amplitude (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingAmplZ string
	// Pitching phase offset (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingPhaseX string
	// Pitching phase offset (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingPhaseY string
	// Pitching phase offset (degrees) about x, y, & z axes (RIGID_MOTION only)
	PitchingPhaseZ string
	// Plunging angular freq. (rad/s) in x, y, & z directions (RIGID_MOTION only)
	PlungingOmegaX string
	// Plunging angular freq. (rad/s) in x, y, & z directions (RIGID_MOTION only)
	PlungingOmegaY string
	// Plunging angular freq. (rad/s) in x, y, & z directions (RIGID_MOTION only)
	PlungingOmegaZ string
	// Plunging amplitude (m) in x, y, & z directions (RIGID_MOTION only)
	PlungingAmplX string
	// Plunging amplitude (m) in x, y, & z directions (RIGID_MOTION only)
	PlungingAmplY string
	// Plunging amplitude (m) in x, y, & z directions (RIGID_MOTION only)
	PlungingAmplZ string
	// Value to move motion origins (1 or 0)
	MoveMotionOrigin string
	//
	MotionFilename string
	// Uncoupled Aeroelastic Frequency Plunge.
	FreqPlungeAeroelastic float64
	// Uncoupled Aeroelastic Frequency Pitch.
	FreqPitchAeroelastic float64
	// Apply a wind gust
	WindGust bool
	// Type of gust
	GustType string
	// Gust wavelenght (meters)
	GustWavelength float64
	// Number of gust periods
	GustPeriods float64
	// Gust amplitude (m/s)
	GustAmpl float64
	// Time at which to begin the gust (sec)
	GustBeginTime float64
	// Location at which the gust begins (meters)
	GustBeginLoc float64
	// Direction of the gust X or Y dir
	GustDir string
	// Convergence criteria
	ConvCriteria string
	// Residual reduction (order of magnitude with respect to the initial value)
	ResidualReduction float64
	// Min value of the residual (log10 of the residual)
	ResidualMinval float64
	// Iteration number to begin convergence monitoring
	StartconvIter float64
	// Number of elements to apply the criteria
	CauchyElems float64
	// Epsilon to control the series convergence
	CauchyEps float64
	// Flow functional for the Cauchy criteria
	CauchyFuncFlow string
	// Adjoint functional for the Cauchy criteria
	CauchyFuncAdj string
	// Linearized functional for the Cauchy criteria
	CauchyFuncLin string
	// Epsilon for a full multigrid method evaluation
	FullmgCauchyEps float64
	// Full multi-grid
	Fullmg bool
	// Start up iterations using the fine grid only
	StartUpIter float64
	// Multi-grid Levels
	Mglevel float64
	// Multi-grid Cycle (0 = V cycle, 1 = W Cycle)
	Mgcycle float64
	// Multi-grid pre-smoothing level
	MgPreSmooth string
	// Multi-grid post-smoothing level
	MgPostSmooth string
	// Jacobi implicit smoothing of the correction
	MgCorrectionSmooth string
	// Damping factor for the residual restriction
	MgDampRestriction float64
	// Damping factor for the correction prolongation
	MgDampProlongation float64
	// CFL reduction factor on the coarse levels
	MgCflReduction float64
	// Maximum number of children in the agglomeration stage
	MaxChildren float64
	// Maximum length of an agglomerated element (relative to the domain)
	MaxDimension float64
	// Numerical method for spatial gradients
	NumMethodGrad string
	// Coefficient for the limiter
	LimiterCoeff float64
	// Coefficient for detecting the limit of the sharp edges
	SharpEdgesCoeff float64
	// No description
	ConvNumMethodFlow string
	// Viscous numerical method
	ViscNumMethodFlow string
	// Source term numerical method
	SourNumMethodFlow string
	// Slope limiter
	SlopeLimiterFlow string
	// 1st, 2nd and 4th order artificial dissipation coefficients
	AdCoeffFlow []float64
	// No description
	ConvNumMethodAdj string
	// Viscous numerical method
	ViscNumMethodAdj string
	// Source term numerical method
	SourNumMethodAdj string
	// Slope limiter
	SlopeLimiterAdjflow string
	// 1st, 2nd and 4th order artificial dissipation coefficients
	AdCoeffAdj []float64
	// Slope limiter
	SlopeLimiterTurb string
	// Convective numerical method
	ConvNumMethodTurb string
	// Viscous numerical method
	ViscNumMethodTurb string
	// Source term numerical method
	SourNumMethodTurb string
	// Slope limiter
	SlopeLimiterAdjturb string
	// Convective numerical method
	ConvNumMethodAdjturb string
	// Viscous numerical method
	ViscNumMethodAdjturb string
	// Source term numerical method
	SourNumMethodAdjturb string
	// Convective numerical method
	ConvNumMethodLin string
	// Viscous numerical method
	ViscNumMethodLin string
	// Source term numerical method
	SourNumMethodLin string
	// 1st, 2nd and 4th order artificial dissipation coefficients
	AdCoeffLin []float64
	// Slope limiter
	SlopeLimiterAdjlevelset string
	// Convective numerical method
	ConvNumMethodAdjlevelset string
	// Viscous numerical method
	ViscNumMethodAdjlevelset string
	// Source term numerical method
	SourNumMethodAdjlevelset string
	// Convective numerical method
	ConvNumMethodTne2 string
	// Viscous numerical method
	ViscNumMethodTne2 string
	// Source term numerical method
	SourNumMethodTne2 string
	// Slope limiter
	SlopeLimiterTne2 string
	// 1st, 2nd and 4th order artificial dissipation coefficients
	AdCoeffTne2 []float64
	// Convective numerical method
	ConvNumMethodAdjtne2 string
	// Viscous numerical method
	ViscNumMethodAdjtne2 string
	// Source term numerical method
	SourNumMethodAdjtne2 string
	// Slope limiter
	SlopeLimiterAdjtne2 string
	// 1st, 2nd and 4th order artificial dissipation coefficients
	AdCoeffAdjtne2 []float64
	// Viscous numerical method
	ViscNumMethodWave string
	// Source term numerical method
	SourNumMethodWave string
	// Viscous numerical method
	ViscNumMethodPoisson string
	// Source term numerical method
	SourNumMethodPoisson string
	// Viscous numerical method
	ViscNumMethodFea string
	// Source term numerical method
	SourNumMethodFea string
	// Viscous numerical method
	ViscNumMethodHeat string
	// Source term numerical method
	SourNumMethodHeat string
	// Source term numerical method
	SourNumMethodTemplate string
	// Limit value for the adjoint variable
	AdjLimit float64
	// Adjoint problem boundary condition
	AdjObjfunc string
	// No description
	GeoSectionLimit []float64
	// Mode of the GDC code (analysis, or gradient)
	GeoMode string
	// Drag weight in sonic boom Objective Function (from 0.0 to 1.0)
	DragInSonicboom float64
	// Sensitivity smoothing
	SensSmoothing string
	// Continuous governing equation set
	ContinuousEqns string
	// Discrete governing equation set
	DiscreteEqns string
	// Adjoint frozen viscosity
	FrozenVisc bool
	//
	CteViscousDrag float64
	// Remove sharp edges from the sensitivity evaluation
	SensRemoveSharp bool
	// I/O
	OutputFormat string
	// Mesh input file format
	MeshFormat string
	// Convert a CGNS mesh to SU2 format
	CgnsToSu2 bool
	// Mesh input file
	MeshFilename string
	// Mesh output file
	MeshOutFilename string
	// Output file convergence history (w/o extension)
	ConvFilename string
	// Restart flow input file
	SolutionFlowFilename string
	// Restart flow input file
	FarfieldFilename string
	// Restart linear flow input file
	SolutionLinFilename string
	// Restart adjoint input file
	SolutionAdjFilename string
	// Output file restart flow
	RestartFlowFilename string
	// Output file linear flow
	RestartLinFilename string
	// Output file restart adjoint
	RestartAdjFilename string
	// Output file restart wave
	RestartWaveFilename string
	// Output file flow (w/o extension) variables
	VolumeFlowFilename string
	// Output file structure (w/o extension) variables
	VolumeStructureFilename string
	// Output file structure (w/o extension) variables
	SurfaceStructureFilename string
	// Output file structure (w/o extension) variables
	SurfaceWaveFilename string
	// Output file structure (w/o extension) variables
	SurfaceHeatFilename string
	// Output file wave (w/o extension) variables
	VolumeWaveFilename string
	// Output file wave (w/o extension) variables
	VolumeHeatFilename string
	// Output file adj. wave (w/o extension) variables
	VolumeAdjwaveFilename string
	// Output file adjoint (w/o extension) variables
	VolumeAdjFilename string
	// Output file linear (w/o extension) variables
	VolumeLinFilename string
	// Output objective function gradient
	GradObjfuncFilename string
	// Output objective function
	ValueObjfuncFilename string
	// Output file surface flow coefficient (w/o extension)
	SurfaceFlowFilename string
	// Output file surface adjoint coefficient (w/o extension)
	SurfaceAdjFilename string
	// Output file surface linear coefficient (w/o extension)
	SurfaceLinFilename string
	// Writing solution file frequency
	WrtSolFreq float64
	// Writing solution file frequency
	WrtSolFreqDualtime float64
	// Writing convergence history frequency
	WrtConFreq float64
	// Writing convergence history frequency for the dual time
	WrtConFreqDualtime float64
	// Write a volume solution file
	WrtVolSol bool
	// Write a surface solution file
	WrtSrfSol bool
	// Write a surface CSV solution file
	WrtCsvSol bool
	// Write a restart solution file
	WrtRestart bool
	// Output residual info to solution/restart file
	WrtResiduals bool
	// Output the rind layers in the solution files
	WrtHalo bool
	// Output sectional forces for specified markers.
	WrtSectionalForces bool
	// Evaluate equivalent area on the Near-Field
	EquivArea bool
	// Integration limits of the equivalent area ( xmin, xmax, Dist_NearField )
	EaIntLimit []float64
	// Specific gas constant (287.87 J/kg*K (air), only for compressible flows)
	GasConstant float64
	// Ratio of specific heats (1.4 (air), only for compressible flows)
	GammaValue float64
	// Reynolds number (non-dimensional, based on the free-stream values)
	ReynoldsNumber float64
	// Reynolds length (1 m by default)
	ReynoldsLength float64
	// Laminar Prandtl number (0.72 (air), only for compressible flows)
	PrandtlLam float64
	// Turbulent Prandtl number (0.9 (air), only for compressible flows)
	PrandtlTurb float64
	// Value of the Bulk Modulus
	BulkModulus float64
	// Artifical compressibility factor
	ArtcompFactor float64
	// Mach number (non-dimensional, based on the free-stream values)
	MachNumber float64
	// No description
	MixtureMolarMass float64
	// Free-stream pressure (101325.0 N/m^2 by default)
	FreestreamPressure float64
	// Free-stream density (1.2886 Kg/m^3 (air), 998.2 Kg/m^3 (water))
	FreestreamDensity float64
	// Free-stream temperature (273.15 K by default)
	FreestreamTemperature float64
	// Free-stream vibrational-electronic temperature (273.15 K by default)
	FreestreamTemperatureVe float64
	// Free-stream velocity (m/s)
	FreestreamVelocity []float64
	// Free-stream viscosity (1.853E-5 Ns/m^2 (air), 0.798E-3 Ns/m^2 (water))
	FreestreamViscosity float64
	//
	FreestreamIntermittency float64
	//
	FreestreamTurbulenceintensity float64
	//
	FreestreamNuFactor float64
	//
	FreestreamTurb2lamviscratio float64
	// Side-slip angle (degrees, only for compressible flows)
	SideslipAngle float64
	// Angle of attack (degrees, only for compressible flows)
	Aoa float64
	// X Reference origin for moment computation
	RefOriginMomentX string
	// Y Reference origin for moment computation
	RefOriginMomentY string
	// Z Reference origin for moment computation
	RefOriginMomentZ string
	// Reference area for force coefficients (0 implies automatic calculation)
	RefArea float64
	// Reference length for pitching, rolling, and yawing non-dimensional moment
	RefLengthMoment float64
	// Reference element length for computing the slope limiter epsilon
	RefElemLength float64
	// Reference coefficient for detecting sharp edges
	RefSharpEdges float64
	// Reference pressure (1.0 N/m^2 by default, only for compressible flows)
	RefPressure float64
	// Reference temperature (1.0 K by default, only for compressible flows)
	RefTemperature float64
	// Reference density (1.0 Kg/m^3 by default, only for compressible flows)
	RefDensity float64
	// Reference velocity (incompressible only)
	RefVelocity float64
	// Reference viscosity (incompressible only)
	RefViscosity float64
	// Factor for converting the grid to meters
	ConvertToMeter float64
	// Write a new mesh converted to meters
	WriteConvertedMesh bool
	// Specify chemical model for multi-species simulations
	GasModel string
	//
	GasComposition string
	// Ratio of density for two phase problems
	RatioDensity float64
	// Ratio of viscosity for two phase problems
	RatioViscosity float64
	// Location of the freesurface (y or z coordinate)
	FreesurfaceZero float64
	// Free surface depth surface (x or y coordinate)
	FreesurfaceDepth float64
	// Thickness of the interface in a free surface problem
	FreesurfaceThickness float64
	// Free surface damping coefficient
	FreesurfaceDampingCoeff float64
	// Free surface damping length (times the baseline wave)
	FreesurfaceDampingLength float64
	// Location of the free surface outlet surface (x or y coordinate)
	FreesurfaceOutlet float64
	// Kind of deformation
	DvKind []string
	// Marker of the surface to which we are going apply the shape deformation
	DvMarker string
	// New value of the shape deformation
	DvValue string
	// No description
	DvParam string
	// Hold the grid fixed in a region
	HoldGridFixed bool
	// Coordinates of the box where the grid will be deformed (Xmin, Ymin, Zmin, Xmax, Ymax, Zmax)
	HoldGridFixedCoord []float64
	// Visualize the deformation
	VisualizeDeformation bool
	// Number of iterations for FEA mesh deformation (surface deformation increments)
	GridDeformIter float64
	// No description
	CyclicPitch float64
	// No description
	CollectivePitch float64
	// Modulus of elasticity
	ElasticityModulus float64
	// Poisson ratio
	PoissonRatio float64
	// Material density
	MaterialDensity float64
	// Constant wave speed
	WaveSpeed float64
	// Thermal diffusivity constant
	ThermalDiffusivity float64
}
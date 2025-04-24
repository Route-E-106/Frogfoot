package buildings

// const (
// 	_ = iota
// 	2000
// 	2000
// 	3000
// 	4000
// 	5000
// 	6000
// )

type BuildingCost struct {
	GasCost   int64 `json:"gasCost"`
	MetalCost int64 `json:"metalCost"`
}

var MetalBuildingCostPerLevel = map[int64]BuildingCost{
	1: BuildingCost{
		GasCost:   1000,
		MetalCost: 1000,
	},
	2: BuildingCost{
		GasCost:   20,
		MetalCost: 40,
	},
	3: BuildingCost{
		GasCost:   4000,
		MetalCost: 10000,
	},
	4: BuildingCost{
		GasCost:   5000,
		MetalCost: 15000,
	},
	5: BuildingCost{
		GasCost:   80000,
		MetalCost: 25000,
	},
	6: BuildingCost{
		GasCost:   10000,
		MetalCost: 35000,
	},
	7: BuildingCost{
		GasCost:   12000,
		MetalCost: 45000,
	},
	8: BuildingCost{
		GasCost:   15000,
		MetalCost: 100000,
	},
	9: BuildingCost{
		GasCost:   20000,
		MetalCost: 200000,
	},
	10: BuildingCost{
		GasCost:   30000,
		MetalCost: 500000,
	},
}

var GasBuildingCostPerLevel = map[int64]BuildingCost{
	1: BuildingCost{
		GasCost:   1000,
		MetalCost: 1000,
	},
	2: BuildingCost{
		GasCost:   2000,
		MetalCost: 6000,
	},
	3: BuildingCost{
		GasCost:   4000,
		MetalCost: 10000,
	},
	4: BuildingCost{
		GasCost:   5000,
		MetalCost: 15000,
	},
	5: BuildingCost{
		GasCost:   80000,
		MetalCost: 25000,
	},
	6: BuildingCost{
		GasCost:   10000,
		MetalCost: 35000,
	},
	7: BuildingCost{
		GasCost:   12000,
		MetalCost: 45000,
	},
	8: BuildingCost{
		GasCost:   15000,
		MetalCost: 100000,
	},
	9: BuildingCost{
		GasCost:   20000,
		MetalCost: 200000,
	},
	10: BuildingCost{
		GasCost:   30000,
		MetalCost: 500000,
	},
}

var MetalIncomePerLevel = map[int64]int64{
	1:  1000,
	2:  2000,
	3:  5000,
	4:  7000,
	5:  9000,
	6:  11000,
	7:  15000,
	8:  20000,
	9:  30000,
	10: 50000,
}
var GasIncomePerLevel = map[int64]int64{
	1:  1000,
	2:  2000,
	3:  4000,
	4:  5000,
	5:  7000,
	6:  9000,
	7:  11000,
	8:  15000,
	9:  20000,
	10: 30000,
}

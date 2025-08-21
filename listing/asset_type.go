package listing

import "fmt"

type AssetType int

const (
	UnknownAssetType AssetType = iota
	Stock
	ETF
	InvalidAssetType
)

func (a AssetType) String() string {
	switch a {
	case Stock:
		return "Stock"
	case ETF:
		return "ETF"
	default:
		panic("invalid value of AssetType")
	}
}

func (a AssetType) isValid() bool {
	if a <= UnknownAssetType || a >= InvalidAssetType {
		return false
	}
	return true
}

func parseAssetType(s string) (AssetType, error) {
	switch s {
	case "Stock":
		return Stock, nil
	case "ETF":
		return ETF, nil
	default:
		return UnknownAssetType, fmt.Errorf("unable to parse Asset Type from: %s", s)
	}
}

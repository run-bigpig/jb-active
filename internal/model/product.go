package model

type Ide struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PluginList struct {
	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	ID           int64  `json:"id"`
	XmlID        string `json:"xmlId"`
	Name         string `json:"name"`
	Preview      string `json:"preview"`
	PricingModel string `json:"pricingModel"`
	Icon         string `json:"icon"`
	CDate        int64  `json:"cdate"`
}

type VendorInfo struct {
	Name       string `json:"name"`
	IsVerified bool   `json:"isVerified"`
}

type PluginCode struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	PurchaseInfo PurchaseInfo `json:"purchaseInfo"`
}

type PurchaseInfo struct {
	ProductCode string `json:"productCode"`
}

package dto

type Request struct {
	Base64      string `json:"base64"`
	FileType    int    `json:"fileType"` // 1.图片 2.PDF
	NeedRawText bool   `json:"needRawText"`
}

type LLMOutput struct {
	OrderDate           string    `json:"orderDate" describe:"订单日期，格式为YYYY-MM-DD"`
	OrderNumber         string    `json:"orderNumber" describe:"订单号"`
	OrderType           string    `json:"orderType" describe:"订单类型"`
	CustomerCompany     string    `json:"customerCompany" describe:"客户公司"`
	CustomerPhone       string    `json:"customerPhone" describe:"客户电话"`
	BuyerName           string    `json:"buyerName" describe:"采购人名"`
	TotalCostWithTax    float64   `json:"totalCostWithTax" describe:"含税金额合计"`
	TotalCostWithoutTax float64   `json:"totalCostWithoutTax" describe:"不含税金额合计"`
	Products            []Product `json:"products" describe:"订购的产品列表"`
	Address             string    `json:"address" describe:"订单中标明的送货地址"`
	AdditionalPoint     string    `json:"additionalPoint" describe:"其他重要的点，例如特殊要求和其他不能从json中提出的信息。你需要用1、xxx；2、xxxx 来填入。"`
}

type Product struct {
	ProductName         string  `json:"productName" describe:"订购产品名称"`
	Model               string  `json:"model" describe:"订购产品的规格型号"`
	OrderNumber         string  `json:"orderNumber" describe:"订购产品的订单号"`
	PartNumber          string  `json:"partNumber" describe:"零件号"`
	Unit                string  `json:"unit" describe:"单位，例如：个、pcs"`
	Quantity            string  `json:"quantity" describe:"数量，例如200，不含单位"`
	PriceWithTax        float64 `json:"priceWithTax" describe:"该产品的含税单价，如果没有标明或单价不含税就填入-1"`
	PriceWithoutTax     float64 `json:"priceWithoutTax" describe:"该产品的不含税单价，如果没有标明或单价含税就填入-1"`
	LineTotalWithTax    float64 `json:"lineTotalWithTax" describe:"订购该类产品的含税金额，未标明或不符合就填-1"`
	LineTotalWithoutTax float64 `json:"lineTotalWithoutTax" describe:"订购该类产品的不含税金额，未标明或不符合就填-1"`
	TaxRate             float64 `json:"taxRate" describe:"订单中标明的税率"`
	Annotation          string  `json:"annotation" describe:"订单中注明的备注"`
	DeliveryDate        string  `json:"deliveryDate" describe:"交货日期，格式为YYYY-MM-DD"`
}

type Output struct {
	RawText   string    `json:"rawText"`
	LLMOutput LLMOutput `json:"llmOutput"`
	Message   string    `json:"message"`
}

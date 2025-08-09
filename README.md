# 识别采购订单PDF/图片中的数据

# 百度手写体识别OCR(可更换，一些OCR api通用)
 apiUrl: "https://aip.baidubce.com/rest/2.0/ocr/v1/table"

# ![百度表格识别OCR](https://cloud.baidu.com/doc/OCR/s/Al1zvpylt)
 apiUrl: "https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting"
 
# langchaingo


# 请求参数
localhost:8086/parse

```json
{
    "needRawText": true, // 相应中是否需要原始文本
    "fileType": 2, // 1.图片 2.PDF
    "base64": "" // 图片base64/PDFbase64 不需要encode
}
```

# 响应
```json
{
    "rawText": "(原始ocr文本)",
    "llmOutput": {
        "orderDate": "2025-07-23",
        "orderNumber": "xxx",
        "orderType": "xxx对接采购订单",
        "customerCompany": "xxx有限公司",
        "customerPhone": "-",
        "buyerName": "xxx",
        "totalCostWithTax": 630,
        "totalCostWithoutTax": 557.51,
        "products": [
            {
                "productName": "xxx",
                "model": "xxx",
                "orderNumber": "xxx",
                "partNumber": "xxx",
                "unit": "pcs",
                "quantity": "1",
                "priceWithTax": -1,
                "priceWithoutTax": 203.89,
                "lineTotalWithTax": 230.4,
                "lineTotalWithoutTax": 203.89,
                "taxRate": 13,
                "annotation": "-",
                "deliveryDate": "2025-07-24"
            },
            {
                "productName": "xxxx",
                "model": "xxx",
                "orderNumber": "xxx",
                "partNumber": "xxx",
                "unit": "pcs",
                "quantity": "1",
                "priceWithTax": -1,
                "priceWithoutTax": 149.73,
                "lineTotalWithTax": 169.2,
                "lineTotalWithoutTax": 149.73,
                "taxRate": 13,
                "annotation": "-",
                "deliveryDate": "2025-07-24"
            },
            {
                "productName": "xxx",
                "model": "xxx",
                "orderNumber": "xxx",
                "partNumber": "xxx",
                "unit": "pcs",
                "quantity": "2",
                "priceWithTax": -1,
                "priceWithoutTax": 101.95,
                "lineTotalWithTax": 230.4,
                "lineTotalWithoutTax": 203.89,
                "taxRate": 13,
                "annotation": "-",
                "deliveryDate": "2025-07-24"
            }
        ],
        "address": "xxx",
        "additionalPoint": "1、送货时随货附送货单,送货单上写上订购单号；2、若产品指标、性能达不到需方技术要求或不能满足需要,供方接受退货并尽快生产出合格产品交至需方；3、除不可抗力情况外,因产品质量问题或供方交货时间延误造成需方损失,若需方提出合理要求,供方应做出赔偿。如交货期每拖延一天,扣该单总金额的1%作为需方的损失赔偿；4、运输方式及到达站港和费用负担:由供方负责办理托运手续并承担托运费用；5、供方必须具有100%的按时交付能力,接单6小时内办签回,逾期本单视同默认。"
    },
    "message": "success"
}
```

# 响应参数说明：
```go
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
```
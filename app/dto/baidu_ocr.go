package dto

// OCRResponse 百度OCR API返回的结果结构
type BaiduOCRResponse struct {
	LogID          int64       `json:"log_id"`
	WordsResultNum int         `json:"words_result_num"`
	WordsResult    []WordsInfo `json:"words_result"`
	Direction      int         `json:"direction,omitempty"`
	PdfFileSize    int         `json:"pdf_file_size,omitempty"`
}

type WordsInfo struct {
	Words      string  `json:"words"`
	Probability float64 `json:"probability,omitempty"`
}
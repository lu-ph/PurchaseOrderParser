package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	OCR_TABLE_URL   = "https://aip.baidubce.com/rest/2.0/ocr/v1/table"
	OCR_URL         = "https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting"
	TOKEN_URL       = "https://aip.baidubce.com/oauth/2.0/token"
	FILE_TYPE_IMAGE = 1
	FILE_TYPE_PDF   = 2
)

type OCRService struct {
	apiKey            string
	secretKey         string
	accessToken       string
	tokenExpiry       int64
	generateTokenTime int64
}

func NewOCRService(apiKey string, secretKey string) *OCRService {
	return &OCRService{apiKey: apiKey, secretKey: secretKey}
}

func (s *OCRService) getAccessToken() (string, error) {
	if s.accessToken != "" && s.tokenExpiry > 0 {
		// 检查token是否过期，提前5分钟更新token
		currentTime := time.Now().Unix()
		// 计算token剩余有效期（秒）
		remainTime := s.generateTokenTime + s.tokenExpiry - currentTime
		// 如果剩余时间大于5分钟（300秒），继续使用当前token
		if remainTime > 300 {
			return s.accessToken, nil
		}
		// 否则重新获取token
	}

	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", s.apiKey, s.secretKey)
	resp, err := http.Post(TOKEN_URL, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取access_token响应失败: %v", err)
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析access_token响应失败: %v", err)
	}

	token, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token不存在或格式错误")
	}

	s.accessToken = token
	// 记录当前时间戳
	s.generateTokenTime = time.Now().Unix()
	// 保存过期时间
	if expiresIn, ok := tokenResp["expires_in"].(float64); ok {
		s.tokenExpiry = int64(expiresIn)
	} else {
		// 如果API没有返回过期时间，使用默认的30天(2592000秒)
		s.tokenExpiry = 2592000
	}
	return token, nil
}

func (s *OCRService) RecognizeFileWithTableAndText(base64Str string, fileType int) (string, error) {
	tableResult, err := s.RecognizeSingleFile(base64Str, fileType, true, map[string]string{"return_excel": "true"})
	if err != nil {
		tableResult = "文件中没有表格或表格识别失败。"
	}
	textResult, err := s.RecognizeSingleFile(base64Str, fileType, false, nil)
	if err != nil {
		return "", fmt.Errorf("文本识别失败: %v", err)
	}
	result := textResult + "\n上文中，表格数据的csv格式是\n" + tableResult
	println(result)
	return result, nil
}

func (s *OCRService) RecognizeSingleFile(base64Str string, fileType int, tableOnly bool, opts map[string]string) (string, error) {
    var ocrUrl string
    if tableOnly {
        ocrUrl = OCR_TABLE_URL
    } else {
        ocrUrl = OCR_URL
    }
    token, err := s.getAccessToken()
    if err != nil {
        return "", err
    }
    requestURL := fmt.Sprintf("%s?access_token=%s", ocrUrl, token)
    params := url.Values{}
    switch fileType {
    case FILE_TYPE_IMAGE:
        params.Add("image", base64Str)
        params.Add("detect_direction", "false")
        params.Add("probability", "false")
        params.Add("detect_alteration", "false")
    case FILE_TYPE_PDF:
        params.Add("pdf_file", base64Str)
        params.Add("pdf_file_num", "1")
        params.Add("detect_direction", "false")
        params.Add("probability", "false")
    default:
        return "", fmt.Errorf("不支持的文件类型")
    }
    // 添加可选参数
    for k, v := range opts {
        params.Add(k, v)
    }
    client := &http.Client{}
    req, err := http.NewRequest("POST", requestURL, strings.NewReader(params.Encode()))
    if err != nil {
        return "", fmt.Errorf("创建请求失败: %v", err)
    }
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Accept", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("发送请求失败: %v", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("读取响应失败: %v", err)
    }
    if tableOnly {
        fmt.Println("表格识别响应:", string(body))
        return s.ParseOCRTableResponseToCSV(string(body))
    }
    return s.ParseOCRResponse(body)
}

// ParseOCRResponse 解析OCR文字识别API的响应，只提取文字
func (s *OCRService) ParseOCRResponse(responseData []byte) (string, error) {
	var ocrResp struct {
		WordsResult []struct {
			Words string `json:"words"`
		} `json:"words_result"`
	}
	if err := json.Unmarshal(responseData, &ocrResp); err != nil {
		var errResp map[string]interface{}
		if jsonErr := json.Unmarshal(responseData, &errResp); jsonErr == nil {
			if errMsg, ok := errResp["error_msg"].(string); ok {
				return "", fmt.Errorf("OCR API错误: %s", errMsg)
			}
		}
		return "", fmt.Errorf("解析OCR响应失败: %v, 响应内容: %s", err, string(responseData))
	}
	var result strings.Builder
	for _, word := range ocrResp.WordsResult {
		result.WriteString(word.Words)
	}
	return result.String(), nil
}

// ParseOCRTableResponseToCSV 从OCR表格识别API的返回中解析数据为CSV格式
// 智能处理：优先使用API返回的Excel文件，否则使用表格数据结构
func (s *OCRService) ParseOCRTableResponseToCSV(jsonStr string) (string, error) {
    // 首先尝试解析Excel文件（如果API返回了）
    var excelRaw struct {
        ExcelFile string `json:"excel_file"`
    }
    if err := json.Unmarshal([]byte(jsonStr), &excelRaw); err == nil && excelRaw.ExcelFile != "" {
        excelBytes, err := base64.StdEncoding.DecodeString(excelRaw.ExcelFile)
        if err != nil {
            return "", fmt.Errorf("Excel文件Base64解码失败: %v", err)
        }
        
        // 将Excel转换为CSV
        return s.ExcelBytesToCSV(excelBytes)
    }

    // 如果没有Excel文件，使用表格数据结构解析
    var raw struct {
        TablesResult []struct {
            Header []struct {
                Words string `json:"words"`
            } `json:"header"`
            Body []struct {
                Words    string `json:"words"`
                RowStart int    `json:"row_start"`
                RowEnd   int    `json:"row_end"`
                ColStart int    `json:"col_start"`
                ColEnd   int    `json:"col_end"`
            } `json:"body"`
            Footer []struct {
                Words string `json:"words"`
            } `json:"footer"`
        } `json:"tables_result"`
    }

    if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
        return "", fmt.Errorf("json解析失败: %v", err)
    }

    var sb strings.Builder
    for _, table := range raw.TablesResult {
        // 处理表头（跳过非表格行的元数据）
        isDataHeader := false
        for _, h := range table.Header {
            // 检查是否是真正的表头行（包含列名）
            if strings.Contains(h.Words, ",") || strings.Contains(h.Words, "：") {
                isDataHeader = true
                sb.WriteString(h.Words + ",")
            }
        }
        if isDataHeader {
            sb.WriteString("\n")
        }

        // 计算表格大小
        maxRow, maxCol := 0, 0
        for _, cell := range table.Body {
            if cell.RowEnd > maxRow {
                maxRow = cell.RowEnd
            }
            if cell.ColEnd > maxCol {
                maxCol = cell.ColEnd
            }
        }

        // 初始化网格
        grid := make([][]string, maxRow+1)
        for r := range grid {
            grid[r] = make([]string, maxCol+1)
        }

        // 填充数据（处理换行符和合并单元格）
        for _, cell := range table.Body {
            // 清理换行符和多余空格
            cleanWords := strings.ReplaceAll(cell.Words, "\n", " ")
            cleanWords = strings.ReplaceAll(cleanWords, "  ", " ")
            
            for r := cell.RowStart; r <= cell.RowEnd; r++ {
                for c := cell.ColStart; c <= cell.ColEnd; c++ {
                    grid[r][c] = cleanWords
                }
            }
        }

        // 写入CSV
        for _, row := range grid {
            // 跳过空行
            isEmpty := true
            for _, col := range row {
                if col != "" {
                    isEmpty = false
                    break
                }
            }
            if !isEmpty {
                sb.WriteString(strings.Join(row, ","))
                sb.WriteString("\n")
            }
        }

        // 处理表尾（作为备注）
        if len(table.Footer) > 0 {
            sb.WriteString("\n备注：")
            for i, f := range table.Footer {
                cleanFooter := strings.ReplaceAll(f.Words, "\n", " ")
                if i > 0 {
                    sb.WriteString("; ")
                }
                sb.WriteString(cleanFooter)
            }
            sb.WriteString("\n")
        }
    }

    return sb.String(), nil
}

// ExcelBytesToCSV 将Excel文件字节数据转换为CSV字符串
func (s *OCRService) ExcelBytesToCSV(excelBytes []byte) (string, error) {
    f, err := excelize.OpenReader(bytes.NewReader(excelBytes))
    if err != nil {
        return "", fmt.Errorf("打开Excel文件失败: %v", err)
    }
    defer f.Close()

    var csvBuilder strings.Builder
    
    // 获取所有工作表
    sheets := f.GetSheetList()
    if len(sheets) == 0 {
        return "", fmt.Errorf("Excel文件中没有工作表")
    }
    
    // 处理每个工作表
    for sheetIndex, sheetName := range sheets {
        rows, err := f.GetRows(sheetName)
        if err != nil {
            return "", fmt.Errorf("读取工作表失败: %v", err)
        }
        
        // 添加工作表标题
        if len(sheets) > 1 {
            csvBuilder.WriteString(fmt.Sprintf("\n--- 工作表 %d: %s ---\n", sheetIndex+1, sheetName))
        }
        
        // 写入CSV数据
        for _, row := range rows {
            // 清理每列数据
            for i, col := range row {
                // 移除换行符和多余空格
                row[i] = strings.ReplaceAll(col, "\n", " ")
                row[i] = strings.ReplaceAll(row[i], "  ", " ")
            }
            csvBuilder.WriteString(strings.Join(row, ","))
            csvBuilder.WriteString("\n")
        }
    }
    
    return csvBuilder.String(), nil
}
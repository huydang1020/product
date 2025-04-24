package utils

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

func MakeCategoryId() string {
	return "cat" + xid.New().String()
}

func MakeProductId() string {
	return "pro" + xid.New().String()
}

func MakeProductTypeId() string {
	return "pty" + xid.New().String()
}

// Chuẩn hóa văn bản về chữ thường
func normalize(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

// Đọc danh sách từ cấm từ file, trả về slice regex pattern
func LoadBannedWords(filePath string) ([]*regexp.Regexp, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Lỗi khi mở file từ cấm:", err)
		return nil, err
	}
	defer file.Close()

	var regexList []*regexp.Regexp
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalize(scanner.Text())
		if word != "" {
			// \b để khớp từ nguyên vẹn (cả từ)
			pattern := `\b` + regexp.QuoteMeta(word) + `\b`
			reg, err := regexp.Compile(pattern)
			if err != nil {
				log.Printf("Không thể biên dịch regex: %s\n", pattern)
				continue
			}
			regexList = append(regexList, reg)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Lỗi khi đọc file từ cấm:", err)
		return nil, err
	}
	return regexList, nil
}

// Hàm kiểm tra nội dung có chứa từ cấm không
func ContainsBannedWords(text string) bool {
	bannedPatterns, err := LoadBannedWords("assets/banner_words.txt")
	if err != nil {
		log.Println("Không thể tải danh sách từ cấm:", err)
		return false
	}

	normalizedText := normalize(text)

	for _, pattern := range bannedPatterns {
		if pattern.MatchString(normalizedText) {
			log.Println("Phát hiện từ cấm:", pattern.String())
			return true
		}
	}
	return false
}

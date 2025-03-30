package goat

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/mteolis/note-goat/internal/gemini"
	"github.com/xuri/excelize/v2"
)

const (
	CLIENT_FIRST_NAME  = "CLIENT_FIRST_NAME"
	CLIENT_LAST_NAME   = "CLIENT_LAST_NAME"
	ADVISOR_FIRST_NAME = "ADVISOR_FIRST_NAME"
	ADVISOR_LAST_NAME  = "ADVISOR_LAST_NAME"

	SUMMARY      = "Summary"
	DISCUSSION   = "discussion"
	PORTFOLIO    = "portfolio"
	REVIEW       = "review"
	WITH         = "with"
	CONVERSATION = "conversation"
	LED          = "led"
	BY           = "by"
	COLON        = ":"
)

var (
	logger     *slog.Logger
	excelFile  string
	promptFile string
)

func InitGoat(slogger *slog.Logger, excelPath string, promptPath string) {
	logger = slogger
	excelFile = excelPath
	promptFile = promptPath
	gemini.InitModel(logger)
}

func AddAISummary() {
	xl, err := excelize.OpenFile(excelFile)
	if err != nil {
		logger.Error("Error opening file %s: %+v\n", "path", excelFile, "err", err)
		log.Printf("Error opening file %s: %+v\n", excelFile, err)
		return
	}
	defer xl.Close()

	sheetList := xl.GetSheetList()
	for sheetIndex, sheetName := range sheetList {
		progress(sheetIndex, sheetName, len(sheetList))
		sheetString := ""
		clientFirstName := ""
		clientLastName := ""
		advisorFirstName := ""
		advisorLastName := ""
		summaryTarget := ""
		summaryRow := -1
		summaryCol := -1

		rows, err := xl.GetRows(sheetName)
		if err != nil {
			logger.Error("Error reading rows from sheet %s: %+v\n", "sheetName", sheetName, "err", err)
			log.Printf("Error reading rows from sheet %s: %+v\n", sheetName, err)
			continue
		}

		for rowIndex, row := range rows {
			if len(row) < 1 {
				continue
			}

			rowString := strings.Join(row, " ")

			if containsClientName(rowString) {
				rowString, clientFirstName, clientLastName = redactClientName(rowString)
			}
			if containsAdvisorName(rowString) {
				rowString, advisorFirstName, advisorLastName = redactAdvisorName(rowString)
			}

			if containsSummary(rowString) {
				summaryRow = rowIndex
				for colIndex, cell := range row {
					if containsSummary(cell) {
						summaryCol = colIndex
						// coordinates are 1-based so add 1 to offset index
						summaryTarget, _ = excelize.CoordinatesToCellName(summaryCol+2, summaryRow+1)
					}
				}
			}
			sheetString += rowString + "\n"
		}

		prompt := buildPrompt(sheetString)

		response, err := gemini.Prompt(prompt)
		if err != nil {
			logger.Error("Error prompting Gemini AI: %+v\n", "err", err)
			log.Printf("Error prompting Gemini AI: %+v\n", err)
			return
		}

		summary := gemini.ExtractSummary(response)
		summary = insertNames(summary, clientFirstName, clientLastName, advisorFirstName, advisorLastName)

		xl.SetCellValue(sheetName, summaryTarget, summary)
		xl.Save()
	}
}

func buildPrompt(sheetString string) string {
	content, err := os.ReadFile(promptFile)
	if err != nil {
		logger.Error("Error reading %s file: %+v\n", "path", promptFile, "err", err)
		log.Fatalf("Error reading %s file: %+v", promptFile, err)
	}

	return string(content) + "\n" + sheetString
}

func containsClientName(str string) bool {
	substrings := []string{
		DISCUSSION, PORTFOLIO, REVIEW, WITH, COLON,
	}
	return containsAllSubstrings(str, substrings)
}

func containsAdvisorName(str string) bool {
	substrings := []string{
		CONVERSATION, LED, BY, COLON,
	}
	return containsAllSubstrings(str, substrings)
}

func containsSummary(str string) bool {
	substrings := []string{
		SUMMARY, COLON,
	}
	return containsAllSubstrings(str, substrings)
}

func containsAllSubstrings(str string, substrings []string) bool {
	for _, substring := range substrings {
		if !strings.Contains(strings.ToLower(str), strings.ToLower(substring)) {
			return false
		}
	}
	return true
}

func redactClientName(rowString string) (string, string, string) {
	firstName, lastName := extractNames(rowString)

	str := strings.Split(rowString, COLON)[0] + COLON + " " + CLIENT_FIRST_NAME + " " + CLIENT_LAST_NAME

	return str, firstName, lastName
}

func redactAdvisorName(rowString string) (string, string, string) {
	firstName, lastName := extractNames(rowString)

	str := strings.Split(rowString, COLON)[0] + COLON + " " + ADVISOR_FIRST_NAME + " " + ADVISOR_LAST_NAME

	return str, firstName, lastName
}

func extractNames(rowString string) (string, string) {
	row := strings.Split(rowString, COLON)

	fullName := strings.Split(strings.TrimSpace(row[len(row)-1]), " ")
	firstName := fullName[0]
	lastName := strings.Join(fullName[1:], " ")

	return firstName, lastName
}

func insertNames(summary string, clientFirstName string, clientLastName string, advisorFirstName string, advisorLastName string) string {
	summary = strings.Replace(summary, CLIENT_FIRST_NAME, clientFirstName, -1)
	summary = strings.Replace(summary, CLIENT_LAST_NAME, clientLastName, -1)
	summary = strings.Replace(summary, ADVISOR_FIRST_NAME, advisorFirstName, -1)
	summary = strings.Replace(summary, ADVISOR_LAST_NAME, advisorLastName, -1)

	return summary
}

func progress(sheetIndex int, sheetName string, totalSheets int) {
	percentage := (float64(sheetIndex+1) / float64(totalSheets)) * 100
	percentageStr := fmt.Sprintf("%.2f", percentage)
	for len(percentageStr) < 6 {
		percentageStr = " " + percentageStr
	}
	logger.Info("(%d/%d | %s%%) NoteGoating sheet: %s\n", "sheetIndex", sheetIndex+1, "totalSheets", totalSheets, "percentageStr", percentageStr, "sheetName", sheetName)
	log.Printf("(%d/%d | %s%%) NoteGoating sheet: %s\n", sheetIndex+1, totalSheets, percentageStr, sheetName)
}

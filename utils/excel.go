package utils

import (
	"fmt"
	"telegramBot/model"

	"github.com/xuri/excelize/v2"
)

func ExportSolitaireFile(fn string, records []model.SolitaireExported) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// header
	f.SetCellValue("Sheet1", "A1", "UserId")
	// f.SetColWidth(sheet, startCol, endCol string, width float64)
	f.SetCellValue("Sheet1", "B1", "Username")
	f.SetCellValue("Sheet1", "C1", "Nickname")
	f.SetCellValue("Sheet1", "D1", "Message")
	f.SetCellValue("Sheet1", "E1", "CreatedAt")

	index := 2
	for _, record := range records {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), record.UserId)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), "@"+record.UserName)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), record.NickName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), record.Message)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), record.CreateAt)

		index++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(fn); err != nil {
		fmt.Println("save filed:", err)
	}
}

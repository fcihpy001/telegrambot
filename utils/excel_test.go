package utils

import (
	"telegramBot/model"
	"testing"
	"time"
)

func TestExportExcel(t *testing.T) {
	records := []model.SolitaireExported{
		{UserId: 5394405541, UserName: "cencen", NickName: "ccc", Message: "哈哈1", CreateAt: time.Now()},
		{UserId: 6297349406, UserName: "Mm-hmm. Okay?", NickName: "Mm-hmm. Okay?", Message: "哈哈daf大发大发的啥饭", CreateAt: time.Now()},
		{UserId: 6450102772, UserName: "fcihpy", NickName: "fcihpy", Message: "的发射东风大厦范德萨范德萨范德萨", CreateAt: time.Now()},
		{UserId: 1091633677, UserName: "liurunq", NickName: "炒币养家", Message: "发动机快乐哦i委屈哦你发的昆仑山麻烦", CreateAt: time.Now()},
		{UserId: 6616020782, UserName: "smart_vbot", NickName: "小微", Message: "诶我去你心里是女款的", CreateAt: time.Now()},
	}
	fn := "./solitaire20230902222310.xlsx"

	ExportSolitaireFile(fn, records)
}

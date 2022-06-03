package server

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

// SettingProperties 設定資料的 struct
type SettingProperties struct {
	APPID  string // PTX API 的 ID
	APPKey string // PTX API 的 Key

	ServerHost string // 後端伺服器要監聽的 host
	ServerPort int32  // 後端伺服器要監聽的 port
}

// readSettingFile 讀取設定檔案
func readSettingFile() (setting SettingProperties, err error) {
	const settingFilePath = "setting.toml"
	// 檢查 setting.conf 是否存在
	var settingFileInfo fs.FileInfo
	if settingFileInfo, err = os.Stat(settingFilePath); err == nil {
		// 如果 setting.conf 存在
		if settingFileInfo.IsDir() {
			err = fmt.Errorf("%s 此資料夾與設定檔同名", settingFilePath)
			return
		}
		_, err = toml.DecodeFile(settingFilePath, &setting)
		if err != nil {
			err = fmt.Errorf("failed to read setting file, error=%s", err)
			return
		}
		if setting.ServerPort <= 0 || setting.ServerPort > 65535 {
			err = fmt.Errorf("listen port need to be 1~65535")
		}
	} else if os.IsNotExist(err) {
		// 如果 setting.conf 不存在
		// 自動生成設定檔
		var settingFile *os.File
		settingFile, err = os.Create(settingFilePath)
		if err != nil {
			err = fmt.Errorf("failed to create %s file, error=%s", settingFilePath, err)
			return
		}
		tEncoder := toml.NewEncoder(settingFile)
		tEncoder.Indent = "\t"
		err = tEncoder.Encode(setting)
		if err != nil {
			err = fmt.Errorf("自動生成設定檔發生錯誤, error=%s", err)
			return
		}

		err = fmt.Errorf("找不到檔案 %s， 由程式自動生成，需要進行編輯", settingFilePath)
		return

	} else {
		// 其他錯誤
		err = fmt.Errorf("failed to check %s status, error=%s", settingFilePath, err)
		return
	}

	return
}

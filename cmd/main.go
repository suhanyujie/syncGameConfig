package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/suhanyujie/syncGameConfig/pkg/filex"
	"github.com/suhanyujie/syncGameConfig/pkg/jsonx"
)

var (
	AllDirMap = map[string]map[string]string{
		"throw": {
			"D:\\tech\\repo2Company\\other\\word\\config\\wrestle": "D:\\tech\\repo2Company\\golang\\throw_gun\\gameConfig",
		},
		"farm": {
			"D:\\tech\\repo2Company\\other\\word\\config\\farm": "D:\\tech\\repo2Company\\golang\\farm_go\\config\\gameConfig",
		},
		"rainbow": {
			"D:\\tech\\repo2Company\\other\\word\\config\\paoku": "D:\\tech\\repo2Company\\golang\\rainbow-bridge\\config\\gameConfig",
		},
	}
)

func main() {
	syncOnePro("rainbow")
}

func syncOnePro(proName string) {
	dirMap1 := AllDirMap["rainbow"]
	for srcDir, dstDir := range dirMap1 {
		syncGameConfig(srcDir, dstDir)
		break
	}
}

func syncGameConfig(srcDir string, targetDir string) {
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		targetFile := path.Join(targetDir, file.Name())
		// 读取文件
		fromFile := path.Join(srcDir, file.Name())
		_, err := os.Stat(fromFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[syncGameConfig] file not exist: %v, file: %v", err, fromFile)
			} else {
				log.Printf("[syncGameConfig] os.Stat err: %v, file: %v", err, fromFile)
			}
			continue
		}
		_, err = os.Stat(targetFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[syncGameConfig] file not exist: %v, file: %v", err, targetFile)
			} else {
				log.Printf("[syncGameConfig] os.Stat err: %v, file: %v", err, targetFile)
			}
			continue
		}
		fromCont := filex.ReadFile(fromFile)
		if fromCont == "" {
			continue
		}
		// 格式化
		resStr := jsonx.JsonStrFormat(fromCont)
		// fmt.Printf("res: %v\n", resStr)
		if len(resStr) == 0 {
			continue
		}
		// 写入到目标目录
		targetFs, err := os.OpenFile(targetFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("[syncGameConfig] OpenFile err: %v", err)
			continue
		}
		targetFs.WriteString(resStr)
		targetFs.Close()
	}
}

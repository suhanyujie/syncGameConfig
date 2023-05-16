package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/suhanyujie/syncGameConfig/pkg/filex"
	"github.com/suhanyujie/syncGameConfig/pkg/jsonx"
	"github.com/urfave/cli/v2"
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

var (
	// default value
	convertProName = ""
)

const (
	UTF8_BOM = "\uFEFF"
)

func main() {
	app := &cli.App{
		Name:   "toJson",
		Usage:  "将策划放置的游戏配置文件，同步到程序所在的仓库",
		Action: DoWork,
		Flags: []cli.Flag{
			//&cli.StringFlag{
			//	Name:        "input",
			//	Value:       "./",
			//	Aliases:     []string{"i"},
			//	Usage:       "要转换的文件所在路径",
			//	Destination: &input,
			//},
			//&cli.StringFlag{
			//	Name:        "output",
			//	Value:       "./output",
			//	Aliases:     []string{"o"},
			//	Usage:       "转换成 json 存放的路径",
			//	Destination: &output,
			//},
			&cli.StringFlag{
				Name:        "project",
				Value:       "",
				Aliases:     []string{"p"},
				Usage:       "要同步的项目名，eg: throw, farm, rainbow",
				Destination: &convertProName,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func DoWork(ctx *cli.Context) error {
	params := ctx.Args().Slice()
	if len(params) != 1 {
		return errors.New("请输入要同步的项目名。eg: throw, farm, rainbow")
	} else {
		if convertProName != ctx.Args().Get(0) {
			convertProName = ctx.Args().Get(0)
		}

		if err := syncOnePro(convertProName); err != nil {
			log.Printf("[DoWork] err: %v", err)
			return err
		}
	}
	fmt.Printf("[ok] 同步完成...\n")
	return nil

}

func syncOnePro(proName string) error {
	dirMap1, ok := AllDirMap[proName]
	if !ok {
		return errors.New("没有配置该项目的目录。")
	}
	for srcDir, dstDir := range dirMap1 {
		syncGameConfig(srcDir, dstDir)
		break
	}
	return nil
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
		if strings.HasPrefix(fromCont, UTF8_BOM) {
			fromCont, _ = strings.CutPrefix(fromCont, UTF8_BOM)
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

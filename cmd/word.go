package cmd

import (
	"awesomeProject/internal/word"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

const (
	ModeUpper = iota + 1
	ModeLower
)

var str string
var mode int8

var desc = strings.Join([]string{
	"该子命令支持各种单词格式转换，模式如下：",
	"1:单词转换为大写",
	"2:单词转换为小写",
}, "\n")

//建立word的子命令，应该是flagSet的封装,所以就是flag包
//wordCmd在root.go中完成子命令的注册
var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  desc,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeLower:
			content = word.ToLower(str)
		default:
			log.Fatalf("暂不支持该模式转换，请执行help word查看帮助文档")
		}
		log.Printf("输出结果:%s", content)
	},
}

func init() {
	//wordCmd就好比是flag.Newflagset("word")的返回一样，它就是子flag
	//这里就是flag的参数注册
	wordCmd.Flags().StringVarP(&str, "str", "s", "这是？", "请输入单词内容")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词转换的模式")
}

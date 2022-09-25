package cmd

import (
	//这个需要自己下载，安装 go get -u  github.com/spf13/cobra@latest
	"github.com/spf13/cobra"
)

//空root命令集合
var rootCmd = &cobra.Command{}

func Execute() error {
	return rootCmd.Execute()
}

//把word命令加入rootCmd
//这也就是子命令
func init() {
	//子命令的注册
	rootCmd.AddCommand(wordCmd)
}

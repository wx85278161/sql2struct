package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/liu-xuewen/sql2struct/sql2struct"
)

func main() {
	parser()
}

func parser() {
	dsn := flag.String("dsn", "", "数据库dsn配置")
	file := flag.String("file", "", "保存路径")
	table := flag.String("table", "", "要迁移的表")
	realNameMethod := flag.String("realNameMethod", "", "结构体对应的表名")
	packageName := flag.String("packageName", "db", "生成的struct包名")
	tagKey := flag.String("tagKey", "gorm", "字段tag的key")
	prefix := flag.String("prefix", "", "表前缀")
	version := flag.Bool("version", false, "版本号")
	v := flag.Bool("v", false, "版本号")
	enableJsonTag := flag.Bool("enableJsonTag", true, "是否添加json的tag,默认false")
	h := flag.Bool("h", false, "帮助")
	help := flag.Bool("help", false, "帮助")

	// 开始
	flag.Parse()

	if *h || *help {
		flag.Usage()
		return
	}

	// 版本号
	if *version || *v {
		fmt.Println(fmt.Sprintf("\n version: %s\n %s\n using -h param for more help \n",
			sql2struct.VERSION, sql2struct.VERSION_TEXT))
		return
	}

	// Dsn  can set in sql2struct.localConfig.go at first
	if sql2struct.Dsn != "" {
		fmt.Println("[ql2struct.Dsn]: ", sql2struct.Dsn)
		*dsn = sql2struct.Dsn
	}

	// 初始化
	t2t := sql2struct.NewTable2Struct()
	// 个性化配置
	t2t.Config(&sql2struct.T2tConfig{
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: false,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: false,
		//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		//SeperatFile: false,
	})

	if *file == "" {

		*file = "./db/" + *table + ".go"
	}

	// 开始迁移转换
	err := t2t.
		// 指定某个表,如果不指定,则默认全部表都迁移
		Table(*table).
		// 表前缀
		Prefix(*prefix).
		// 是否添加json tag
		EnableJsonTag(*enableJsonTag).
		// 生成struct的包名(默认为空的话, 则取名为: package model)
		PackageName(*packageName).
		// tag字段的key值,默认是orm
		TagKey(*tagKey).
		// 是否添加结构体方法获取表名
		RealNameMethod(*realNameMethod).
		// 生成的结构体保存路径
		SavePath(*file).
		// 数据库dsn
		Dsn(*dsn).
		// 执行
		Run()

	if err != nil {
		log.Println(err.Error())
	}
}

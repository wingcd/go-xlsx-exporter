package main

import (
	"fmt"
	"go-xlsx-protobuf/generator"
	"io/ioutil"
	"os"
	"testing"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/serialize"
	"github.com/wingcd/go-xlsx-protobuf/settings"
	"github.com/wingcd/go-xlsx-protobuf/utils"
	"github.com/wingcd/go-xlsx-protobuf/xlsx"

	gproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"

	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// 测试表格解析
func TestParseXlsx(t *testing.T) {
	infos := xlsx.ParseDefineSheet("data/define.xlsx", "define")

	for _, info := range infos {
		fmt.Println(info.TypeName)
	}
}

// 测试值解析
func TestParseValue(t *testing.T) {
	ret, val, _ := utils.ParseValue("int", "6")
	fmt.Printf("parse int 6: %v, %v \n", ret, val)

	ret, val, arr := utils.ParseValue("int[]", "6|7|8")
	fmt.Printf("parse int[] 6|7|8: %v, %v, %v \n", ret, val, arr)

	ret, val, _ = utils.ParseValue("bool", "true")
	fmt.Printf("parse bool true: %v, %v \n", ret, val)

	ret, val, _ = utils.ParseValue("bool[]", "true|false|true")
	fmt.Printf("parse bool[] true|false|true: %v, %v \n", ret, val)

	ret, val, arr = utils.ParseValue("string[]", "ab||b||c|d|7|8")
	fmt.Printf("parse string[] ab||b||c|d|7|8: %v, %v, %v \n", ret, val, arr)
}

// 测试枚举值解析
func TestParseEnumValue(t *testing.T) {
	infos := xlsx.ParseDefineSheet("data/define.xlsx", "define")

	ret, val, _ := utils.ParseEnumValue(infos["EDataType"], "EDataType", "XML")
	fmt.Printf("parse EDataType XML: %v, %v \n", ret, val)

	ret, val, arr := utils.ParseEnumValue(infos["EDataType"], "EDataType[]", "XML|JSON|GOLANG")
	fmt.Printf("parse EDataType[] XML|JSON: %v, %v, %v \n", ret, val, arr)
}

// 测试数据表解析
func TestParseDataXlsx(t *testing.T) {
	t_user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	t_user.TypeName = "User"

	t_class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	t_class.TypeName = "PClass"
}

// 测试proto数据序列化
func TestPBSerialize(t *testing.T) {
	settings.SetDefines(xlsx.ParseDefineSheet("data/define.xlsx", "define"))

	user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	user.TypeName = "User"

	class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	class.TypeName = "PClass"

	settings.SetTables([]*model.DataTable{class, user})

	fd, _ := utils.BuildFileDesc("test")

	fooMessageDescriptor := fd.Messages().ByName("User")
	msg := dynamicpb.NewMessage(fooMessageDescriptor)
	msg.Set(fooMessageDescriptor.Fields().ByName("ID"), pref.ValueOfUint32(42))
	msg.Set(fooMessageDescriptor.Fields().ByNumber(2), pref.ValueOfString("张三"))
	userMsg := msg
	bytes, err := gproto.Marshal(msg)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	msg = dynamicpb.NewMessage(fooMessageDescriptor)
	err = gproto.Unmarshal(bytes, msg)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	v := msg.Get(fooMessageDescriptor.Fields().ByName("ID"))
	fmt.Printf("get %v \n", v)

	v = msg.Get(fooMessageDescriptor.Fields().ByName("Name"))
	fmt.Printf("get %v \n", v)

	// list
	listMD := fd.Messages().ByName("User_ARRAY")
	msg = dynamicpb.NewMessage(listMD)
	lf := listMD.Fields().ByName("Items")
	lst := msg.NewField(lf).List()
	lst.Append(pref.ValueOf(userMsg))
	lst.Append(pref.ValueOf(userMsg))
	lst.Append(pref.ValueOf(userMsg))
	lst.Append(pref.ValueOf(userMsg))
	msg.Set(lf, pref.ValueOf(lst))

	bytes, err = gproto.Marshal(msg)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	msg = dynamicpb.NewMessage(listMD)
	err = gproto.Unmarshal(bytes, msg)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	lst = msg.Get(listMD.Fields().ByName("Items")).List()
	for i := 0; i < lst.Len(); i++ {
		item := lst.Get(i).Message()

		v := item.Get(fooMessageDescriptor.Fields().ByName("ID"))
		fmt.Printf("id:%v get %v \n", i, v)

		v = item.Get(fooMessageDescriptor.Fields().ByName("Name"))
		fmt.Printf("id:%v get %v \n", i, v)
	}
}

func TestSaveSerializeDefineData(t *testing.T) {
	settings.SetDefines(xlsx.ParseDefineSheet("data/define.xlsx", "define"))

	var pbname = ""
	fd, _ := utils.BuildFileDesc(pbname)

	serialize.GenDefineTables(pbname, fd, "./gen/bytes/", settings.CONSTS)
}

func TestSaveSerializeData(t *testing.T) {
	settings.SetDefines(xlsx.ParseDefineSheet("data/define.xlsx", "define"))

	class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	class.TypeName = "PClass"

	user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	user.TypeName = "User"

	settings.SetTables([]*model.DataTable{class, user})

	var pbname = ""
	fd, _ := utils.BuildFileDesc(pbname)

	serialize.GenDataTables(pbname, fd, "./gen/bytes/", settings.TABLES)

	itemMD := fd.Messages().ByName("User")

	listMD := fd.Messages().ByName("User_ARRAY")
	msg := dynamicpb.NewMessage(listMD)

	f, err := os.Open("./gen/bytes/user.bytes")
	defer f.Close()
	if err != nil {
		fmt.Print(err.Error())
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = gproto.Unmarshal(bytes, msg)
	if err != nil {
		fmt.Print(err.Error())
	}

	lst := msg.Get(listMD.Fields().ByName("Items")).List()
	for i := 0; i < lst.Len(); i++ {
		item := lst.Get(i).Message()

		v := item.Get(itemMD.Fields().ByName("ID"))
		fmt.Printf("id:%v get %v \n", i, v)

		v = item.Get(itemMD.Fields().ByName("Name"))
		fmt.Printf("id:%v get %v \n", i, v)

		v = item.Get(itemMD.Fields().ByName("Age"))
		fmt.Printf("id:%v get %v \n", i, v)
	}
}

func TestGenProtoBytesFile(t *testing.T) {
	defines := xlsx.ParseDefineSheet("data/define.xlsx", "define")
	settings.SetDefines(defines)

	t_user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	t_user.TypeName = "User"
	t_class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	t_class.TypeName = "PClass"
	settings.SetTables([]*model.DataTable{t_user, t_class})

	generator.Build("proto_bytes", "./gen/bytes/")
}

func TestGenProtoFile(t *testing.T) {
	defines := xlsx.ParseDefineSheet("data/define.xlsx", "define")
	settings.SetDefines(defines)

	t_user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	t_user.TypeName = "User"
	t_class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	t_class.TypeName = "PClass"
	settings.SetTables([]*model.DataTable{t_user, t_class})

	generator.Build("proto", "./gen/all.proto")
}

func TestGenCSharpFile(t *testing.T) {
	defines := xlsx.ParseDefineSheet("data/define.xlsx", "define")
	settings.SetDefines(defines)

	t_user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	t_user.TypeName = "User"
	t_class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	t_class.TypeName = "PClass"
	settings.SetTables([]*model.DataTable{t_user, t_class})

	generator.Build("csharp", "./gen/DataMode.cs")
}

func TestGenGolangFile(t *testing.T) {
	defines := xlsx.ParseDefineSheet("data/define.xlsx", "define")
	settings.SetDefines(defines)

	t_user := xlsx.ParseDataSheet("data/model.xlsx", "user")
	t_user.TypeName = "User"
	t_class := xlsx.ParseDataSheet("data/model.xlsx", "class")
	t_class.TypeName = "PClass"
	settings.SetTables([]*model.DataTable{t_user, t_class})

	settings.PackageName = "gen"
	generator.Build("golang", "./gen/DataMode.pb.go")
}

func TestGenGolangFileWithComment(t *testing.T) {
	t_comment := xlsx.ParseDataSheet("data/model.xlsx", "comment")
	t_comment.TypeName = "Comment"
	settings.SetTables([]*model.DataTable{t_comment})

	settings.PackageName = "gen"
	generator.Build("golang", "./gen/Comment.pb.go")
}

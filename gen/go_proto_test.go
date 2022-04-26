package gen

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	gxe "github.com/wingcd/go-xlsx-exporter/reader/golang/go_xlsx_exporter"
	"google.golang.org/protobuf/proto"
)

func TestGoPBMode(t *testing.T) {
	var data = PClass{}
	data.ID = 1
	data.Name = "test"
	data.Type = EDataType_CSHARP
	data.Level = 1

	bts, err := proto.Marshal(&data)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var cdata = PClass{}
	err = proto.Unmarshal(bts, &cdata)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("ID:%v \n", cdata.ID)
	fmt.Printf("Name:%v \n", cdata.Name)
	fmt.Printf("Type:%v \n", cdata.Type.String())
	fmt.Printf("Level:%v \n", cdata.Level)
}

func TestGoPBListData(t *testing.T) {
	var data = PClass{}
	data.ID = 1
	data.Name = "test"
	data.Type = EDataType_CSHARP
	data.Level = 1

	var list = PClass_ARRAY{}
	list.Items = []*PClass{&data, &data}
	bts, err := proto.Marshal(&list)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	list = PClass_ARRAY{}
	err = proto.Unmarshal(bts, &list)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for _, cdata := range list.Items {
		fmt.Printf("ID:%v \n", cdata.ID)
		fmt.Printf("Name:%v \n", cdata.Name)
		fmt.Printf("Type:%v \n", cdata.Type.String())
		fmt.Printf("Level:%v \n", cdata.Level)
	}
}

func TestGoLoadPBFile(t *testing.T) {
	var classes = PClass_ARRAY{}
	var bytes, _ = ioutil.ReadFile("./data/pclass.bytes")
	err := proto.Unmarshal(bytes, &classes)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for _, cdata := range classes.Items {
		fmt.Printf("ID:%v \n", cdata.ID)
		fmt.Printf("Name:%v \n", cdata.Name)
		fmt.Printf("Type:%v \n", cdata.Type.String())
		fmt.Printf("Level:%v \n", cdata.Level)
	}
}

func TestGoLoadPBFile2(t *testing.T) {
	var users = User_ARRAY{}
	var bytes, _ = ioutil.ReadFile("./data/user.bytes")
	err := proto.Unmarshal(bytes, &users)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for _, cdata := range users.Items {
		fmt.Printf("ID:%v \n", cdata.ID)
		fmt.Printf("Name:%v \n", cdata.Name)
		fmt.Printf("Type:%v \n", cdata.Type.String())
		fmt.Printf("Level:%v \n", cdata.Head)
	}
}

func TestLoadConfigFile(t *testing.T) {
	var settings = Settings{}
	var bytes, _ = ioutil.ReadFile("./data/settings.bytes")
	err := proto.Unmarshal(bytes, &settings)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("DataType:%v \n", settings.DataType)
	fmt.Printf("Version:%v \n", settings.VERSION)
}

func TestGoLoadLanguage(t *testing.T) {
	var lans_cn = gxe.Language_ARRAY{}
	var bytes, _ = ioutil.ReadFile("./data/language.cn.bytes")
	err := proto.Unmarshal(bytes, &lans_cn)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var lans_en = gxe.Language_ARRAY{}
	bytes, _ = ioutil.ReadFile("./data/language.en.bytes")
	err = proto.Unmarshal(bytes, &lans_en)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var lans_jp = gxe.Language_ARRAY{}
	bytes, _ = ioutil.ReadFile("./data/language.jp.bytes")
	err = proto.Unmarshal(bytes, &lans_jp)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for i := 0; i < len(lans_cn.Items); i++ {
		var cn = lans_cn.Items[i]
		var en = lans_en.Items[i]
		var jp = lans_jp.Items[i]
		fmt.Printf("ID: %v\n Name: %v, %v, %v \n", cn.ID, cn.Text, en.Text, jp.Text)
	}
}
func TestReader(t *testing.T) {
	gxe.Initial("./bytes", "ID")

	// set data table key name, can different from default key name(just same with xlsx config)
	gxe.RegistDataTable("ID", reflect.TypeOf(User{}))

	// settings
	var dt = gxe.GetDataItem(reflect.TypeOf(Settings{}))
	var settings = dt.Item().(*Settings)
	fmt.Printf("settings version=%v,maxconn=%v\n\n", settings.VERSION, settings.MAX_CONNECT)

	// table
	var userTable = gxe.GetDataTable(reflect.TypeOf(User{}))
	var users = userTable.Items()
	for _, dt := range users {
		var user = dt.(*User)
		fmt.Printf("Name:%s, Age:%v, Sex: %v \n", user.Name, user.Age, user.Sex)
	}
	var user = userTable.GetItem("1").(*User)
	fmt.Printf("\nByMap ID=1 Name:%s, Age:%v, Sex: %v \n\n", user.Name, user.Age, user.Sex)

	// language
	gxe.SetLanguage(gxe.DefaultIndexKey, "cn")
	fmt.Printf("中文 cn id=1, text=%v \n", gxe.Translate("1"))

	gxe.SetLanguage(gxe.DefaultIndexKey, "en")
	fmt.Printf("English en id=1, text=%v \n", gxe.Translate("1"))
}

package gen

import (
	"fmt"
	"io/ioutil"
	"testing"

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
	var bytes, _ = ioutil.ReadFile("./bytes/pclass.bytes")
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

func TestGoLoadLanguage(t *testing.T) {
	var lans_cn = Language_ARRAY{}
	var bytes, _ = ioutil.ReadFile("./bytes/language.cn.bytes")
	err := proto.Unmarshal(bytes, &lans_cn)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var lans_en = Language_ARRAY{}
	bytes, _ = ioutil.ReadFile("./bytes/language.en.bytes")
	err = proto.Unmarshal(bytes, &lans_en)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var lans_jp = Language_ARRAY{}
	bytes, _ = ioutil.ReadFile("./bytes/language.jp.bytes")
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

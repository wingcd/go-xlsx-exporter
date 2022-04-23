package gen

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestCreateMessage(t *testing.T) {
	c2splayer := C2S_GetPlayerInfo{}
	c2splayer.Id = 1
	c2splayer.Name = "test"

	bts, err := proto.Marshal(&c2splayer)
	if err != nil {
		t.Error(err)
		return
	}

	err, dt := LoadMessage(10001, bts)
	if err != nil {
		fmt.Println(err)
	}

	msg := dt.(*C2S_GetPlayerInfo)
	fmt.Println("ID:", msg.Id)
}

func TestMessage(t *testing.T) {
	c2splayer := C2S_GetPlayerInfo{}
	c2splayer.Id = 1
	c2splayer.Name = "test"

	bts, err := proto.Marshal(&c2splayer)
	if err != nil {
		t.Error(err)
		return
	}

	c2splayer = C2S_GetPlayerInfo{}
	err = proto.Unmarshal(bts, &c2splayer)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("ID:", c2splayer.Id)
	fmt.Println("Name:", c2splayer.Name)
}

func TestMessage2(t *testing.T) {
	s2cplayer := S2C_GetPlayerInfo{}
	s2cplayer.Id = 1
	s2cplayer.Name = "test"
	s2cplayer.Type = EMsgType_XML
	s2cplayer.Items = []*Item{
		&Item{
			Id:   1,
			Name: "test",
		},
		&Item{
			Id:   2,
			Name: "test2",
		},
	}

	bts, err := proto.Marshal(&s2cplayer)
	if err != nil {
		t.Error(err)
		return
	}

	s2cplayer = S2C_GetPlayerInfo{}
	err = proto.Unmarshal(bts, &s2cplayer)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("ID:", s2cplayer.Id)
	fmt.Println("Name:", s2cplayer.Name)
	fmt.Println("Type:", s2cplayer.Type.String())
	for _, item := range s2cplayer.Items {
		fmt.Println("ID:", item.Id)
		fmt.Println("Name:", item.Name)
	}
}

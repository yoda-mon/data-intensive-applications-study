package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	genProto "github.com/yoda-mon/data-intensive-application/encoding/cmd/protobuf"
	genThrift "github.com/yoda-mon/data-intensive-application/encoding/cmd/thrift"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/golang/protobuf/proto"
	"github.com/linkedin/goavro/v2"
	"github.com/ugorji/go/codec"
)

type Person struct {
	userName       interface{}
	favoriteNumber interface{}
	interests      interface{}
}

func main() {
	/*
		MessagePack
	*/
	fmt.Println("MessagePack:")
	bs, err := ioutil.ReadFile("dist/sample.mpac")
	if err != nil {
		log.Fatal(err)
	}
	//var personMp Person
	var personMp map[string]interface{} // structにマッピングできない...
	var h codec.Handle = new(codec.MsgpackHandle)
	var dec *codec.Decoder = codec.NewDecoderBytes(bs, h)
	err = dec.Decode(&personMp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(personMp)
	fmt.Println()
	/*
		Thrift
	*/
	fmt.Println("Thrift:")
	tb, err := ioutil.ReadFile("dist/sample.tb")
	if err != nil {
		log.Fatal(err)
	}
	transp := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(tb)}
	ptc := thrift.NewTBinaryProtocol(transp, false, false)
	personTh := genThrift.NewPerson()
	personTh.Read(ptc)
	fmt.Println(personTh) // 上手くいかない（ドキュメント皆無）
	fmt.Println()
	/*
		PB
	*/
	fmt.Println("Protocol Buffers:")
	pb, err := ioutil.ReadFile("dist/sample.pb")
	if err != nil {
		log.Fatal(err)
	}
	p := &genProto.Person{}
	proto.Unmarshal(pb, p)
	fmt.Println(p)
	fmt.Println()

	/*
		Avro
	*/
	fmt.Println("Avro:")
	avsc, err := ioutil.ReadFile("config/sample.avsc")
	if err != nil {
		log.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(avsc))
	if err != nil {
		fmt.Println(err)
	}
	ab, err := ioutil.ReadFile("dist/sample.ab")
	native, _, err := codec.NativeFromBinary(ab)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(native)

}

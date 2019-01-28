package main

import (
	"bytes"
	"compress/gzip"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"io/ioutil"
	// "os"
	"fmt"
	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

var sampleTableFile = "tcaplusgatesvr.proto"
var descriptorFileName = "google/protobuf/descriptor.proto"

func tryA() {
	g := generator.New()

	g.Request.FileToGenerate = make([]string, 0)
	g.Request.FileToGenerate = append(g.Request.FileToGenerate, sampleTableFile)
	g.Request.FileToGenerate = append(g.Request.FileToGenerate, descriptorFileName)

	g.Request.ProtoFile = make([]*google_protobuf.FileDescriptorProto, 0)
	for _, protoFileName := range g.Request.FileToGenerate {
		gz := proto.FileDescriptor(protoFileName)
		if len(gz) == 0 {
			LogErr("FileDescriptor gz return 0 for:%s", protoFileName)
			continue
		}
		//本来应该使用c++ protobuf lib生成descriptor文件
		protoFileDescriptor, errExtract := ExtractDescriptorFile(gz)
		if errExtract != nil {
			LogErr("errExtract:%+v,for:%s", errExtract, protoFileName)
			continue
		}
		g.Request.ProtoFile = append(g.Request.ProtoFile, protoFileDescriptor)
	}
	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	LogInfo("before: WrapTypes")
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	LogInfo("before: GenerateAllFiles")
	g.GenerateAllFiles()
	LogInfo("after: GenerateAllFiles")

	// Send back the results.
	_, err := proto.Marshal(g.Response)
	if err != nil {
		LogErr("marshal fail:err:%+v", err)
		return
	}
	for _, rspFile := range g.Response.File {
		LogInfo("rsp,%s", rspFile.GetName())
	}
	//LogInfo("response:%s", data)
}

func ExtractDescriptorFile(gz []byte) (*google_protobuf.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %v", err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to uncompress descriptor: %v", err)
	}

	fd := new(google_protobuf.FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		return nil, fmt.Errorf("malformed FileDescriptorProto: %v", err)
	}

	return fd, nil
}

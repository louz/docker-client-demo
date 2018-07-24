package main

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os"
	"io"
	"flag"
	"bytes"
	"strings"
	"encoding/json"
)

type ImageLoadResult struct {
	// eg. {"stream":"Loaded image ID: sha256:4289425e7917e771cfbb0065616586e17789bb23008e01975b709deed0438106\n"}
	Stream string	`json:"stream"`
}

func (r *ImageLoadResult) ImageId() string {
	// eg. 要从{"stream":"Loaded image ID: sha256:4289425e7917e771cfbb0065616586e17789bb23008e01975b709deed0438106\n"}
	// 解析出 sha256:4289425e7917e771cfbb0065616586e17789bb23008e01975b709deed0438106
	firstLine := strings.Split(r.Stream, "\n")[0]
	return strings.Split(firstLine, ": ")[1]
}

func main() {
	// 解析命令行参数
	filePath, tag := parseFilePathAndTagArgs()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// 1. 读取镜像列表
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.RepoTags)
	}

	// 2. load镜像压缩文件
	dir, _ := os.Open(filePath)
	response, err := cli.ImageLoad(context.Background(), dir, false)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	// 控制台输出结果
	//io.Copy(os.Stdout, response.Body)

	// 将结果保存到字符串
	bs := bytes.NewBufferString("")
	io.Copy(bs, response.Body)
	loadResultString := bs.String()
	fmt.Println(loadResultString)

	imageId := parseImageIdFromLoadResult(loadResultString)

	// 3. 打tag
	err = cli.ImageTag(context.Background(), imageId, tag)
	if err != nil {
		panic(err)
	}

	// 4. push image
	res, err := cli.ImagePush(context.Background(), tag, types.ImagePushOptions{All: true,
		RegistryAuth: "123"}) // ImagePushOptions必须要使用类似的设置，不能使用types.ImagePushOptions{}
	if err != nil {
		panic(err)
	}

	defer res.Close()

	// 控制台输出结果
	io.Copy(os.Stdout, res)

	// 5. remove image
	rmResult, err := cli.ImageRemove(context.Background(), tag, types.ImageRemoveOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(rmResult)
}

func parseFilePathAndTagArgs() (string, string) {
	var filePath = flag.String("f", "", "镜像tar文件的绝对路径")
	var tag = flag.String("t", "", "镜像push到registry的tag")
	flag.Parse()

	if *filePath == "" || *tag == "" {
		fmt.Println("参数数目错误！")
		flag.Usage()
		os.Exit(1)
	}

	return *filePath, *tag
}

func parseImageIdFromLoadResult(loadResultString string) string {
	s := strings.Split(loadResultString, "\n")
	var loadResult ImageLoadResult
	json.Unmarshal([]byte(s[len(s)-2]), &loadResult)

	return loadResult.ImageId()
}
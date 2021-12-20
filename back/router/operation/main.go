package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
)

/**
 	现有的镜像/上传一个tar包
		上传一个tar包，
*/
func ImageLoad(input io.Reader) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("116.56.140.62"))
	if err != nil {
		fmt.Print("client NewClientWithOpts err.%v", err)
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	response, err := cli.ImageLoad(ctx, input, false)
	if err != nil {
		fmt.Println("image load err.%v", err)
		return err
	}
	defer response.Body.Close()
	return nil
}
func ImagePush(username string, passwd string, tag string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://116.56.140.62:2375"))
	if err != nil {
		fmt.Print("client NewClientWithOpts err.%v", err)
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	authConfig := types.AuthConfig{
		Username: username,
		Password: passwd,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	response, err := cli.ImagePush(ctx, tag, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Print("failed to push image %s err.%v", tag, err)
		return err
	}

	defer response.Close()
	io.Copy(os.Stdout, response)
	return nil
}
func ImagePull(img string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://116.56.140.62:2375"))
	if err != nil {
		fmt.Print("client NewClientWithOpts err.%v", err)
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ImageTag(source string, target string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Print("client NewClientWithOpts err.%v", err)
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	if err := cli.ImageTag(ctx, source, target); err != nil {
		fmt.Println("failed to tag %s as %s", source, target)
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://116.56.140.62:2375"))
	if err != nil {
		fmt.Printf("client NewClientWithOpts err.%v", err)
		return
	}
	//cli.NegotiateAPIVersion(ctx)
	authConfig := types.AuthConfig{
		Username: "root",
		Password: "Harbor12345",
	}
	encodedJSON, err := json.Marshal(authConfig)
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	a, err := cli.ImagePull(ctx, "116.56.140.62:8089/public/abc-service-provider:latest",  types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Printf("client ImagePull err.%v", err)
		return
	}
	a.Close()
	//cli, err := client.NewClientWithOpts(client.WithHost("tcp://116.56.140.62:2375"))
	//if err != nil {
	//	fmt.Printf("client NewClientWithOpts err.%v", err)
	//	return
	//}
	//ctx := context.Background()
	//cli.NegotiateAPIVersion(ctx)
	//list, err := cli.ImageList(ctx, types.ImageListOptions{})
	//if err != nil{
	//	fmt.Printf("client ImageList err.%v", err)
	//	return
	//}
	//for _, v := range list {
	//	fmt.Println(v.RepoTags)
	//}
}

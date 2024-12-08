package main

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdresty"
	"github.com/samber/lo"
	"strings"
	"time"
)

func main() {
	sep := strings.Repeat("-", 60)
	c := sdresty.New(&sdresty.Options{
		Timeout: time.Second * 20,
	})

	// CURL
	ip := lo.Must(sdresty.CURL(context.Background(), c, "https://ipinfo.io/ip"))
	fmt.Printf("ip: %s\n", ip)
	fmt.Println(sep)

	// GET
	info := lo.Must(sdresty.GET(
		context.Background(),
		c,
		"https://httpbin.org/get",
		sdresty.ForJsonValue(),
		sdresty.QueryParams(map[string]any{"k1": "v1", "k2": 22}),
	))
	fmt.Println(sdjson.MarshalStringOr(info, ""))
	fmt.Println(sep)

	// POST
	info = lo.Must(sdresty.POST(
		context.Background(),
		c,
		"https://httpbin.org/post",
		sdresty.ForJsonValue(),
		sdresty.JsonData(sdjson.Object{"body_k1": "body_v1", "body_k2": 222}),
		sdresty.QueryParams(map[string]any{"k1": "v1", "k2": 22}),
	))
	fmt.Println(sdjson.MarshalStringOr(info, ""))
	fmt.Println(sep)

	// 解析返回后的json

	type Address struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	}

	type Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	}

	type User struct {
		Id       int64   `json:"id"`
		Name     string  `json:"name"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Address  Address `json:"address"`
		Phone    string  `json:"phone"`
		Website  string  `json:"website"`
		Company  Company `json:"company"`
	}

	users := lo.Must(sdresty.GET(
		context.Background(),
		c,
		"https://jsonplaceholder.typicode.com/users",
		sdresty.ForJsonBody[[]*User](),
	))
	fmt.Println(sdjson.MarshalPretty(users[0:1]))
	fmt.Print(sep)

	// 将返回的body保存到文件
	_ = lo.Must(sdresty.GET(
		context.Background(),
		c,
		"https://httpbin.org/get",
		sdresty.SaveToFile("/Users/gaorx/Desktop/httpbin_get.json", false),
		sdresty.QueryParams(map[string]any{"k1": "v1", "k2": 22}),
	))
	fmt.Println(sep)
}

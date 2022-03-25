package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
)
var cli , _  = clientv3.New(clientv3.Config{
	Endpoints: []string{"127.0.0.1:2379"},
	DialTimeout: 5 * time.Second,
})

func main() {
	defer cli.Close()

	// get
	r := gin.Default()
	r.GET("/get",get)
	r.GET("put",put)
	r.Run(":9092")
}
func get(c *gin.Context){
	key := c.Query("key")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		c.JSON(400,err)
		return
	}
	m := make(map[string]string)

	for _, ev := range resp.Kvs {
		m[string(ev.Key)]= string(ev.Value)
	}
	c.JSON(200,m)
}
func put(c *gin.Context){
	key := c.Query("key")
	val := c.Query("value")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		c.JSON(400,err)
		return
	}
	c.JSON(200,"成功")
}

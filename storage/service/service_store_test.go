package service

import (
	"testing"
	"github.com/Jeffail/gabs"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"sync"
)

func testServiceStoreUpload(key int64,img []byte){
	json:=gabs.New()
	json.Set(img,"img")
	json.Set(key,"key")
	json.Set(123,"cookie")
	var jsonStr = []byte(json.String())
	req, err := http.NewRequest("POST", "http://192.168.10.25:10086/store/upload", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json,err=gabs.ParseJSON(body)
	if err != nil {
		fmt.Println(err)
	}

	result:=json.Path("result").String()
	if result!="0"{
		fmt.Println("result:",json.String())
	}
}

func TestServiceStoreUpload(t *testing.T) {
	fmt.Println("TestServiceStoreUpload")
	img:=make([]byte,1024*128)
	count:=int64(0)
	lock:=&sync.Mutex{}
	key:=int64(0)
	for i:=int64(0);i<16;i++{
		go func(){
			for {
				var k int64
				lock.Lock()
				key++
				k=key
				lock.Unlock()

				count+=int64(len(img))
				testServiceStoreUpload(k,img)

				if (k%100)==0{
					fmt.Println(count/1024/1024,"M")
				}
			}
		}()
	}

	for{
		time.Sleep(time.Second)
	}
}
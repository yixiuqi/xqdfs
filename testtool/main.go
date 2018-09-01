package main

import (
	"io/ioutil"
	"fmt"
	"github.com/Jeffail/gabs"
	"bytes"
	"net/http"
	"sync"
	"time"
	"os"
)

func main() {
	fmt.Println("TestServiceStoreUpload")
	img:=make([]byte,1024*64)
	count:=int64(0)
	lock:=&sync.Mutex{}
	key:=int64(0)
	for i:=int64(0);i<64;i++{
		go func(){
			for {
				var k int64
				lock.Lock()
				key++
				k=key
				lock.Unlock()

				count+=int64(len(img))
				str,ret:=testServiceStorageUpload(k,img)
				json,_:=gabs.ParseJSON([]byte(str))
				url:=json.Path("url").Data().(string)
				testServiceStorageDelete(url)
				if (k%1000)==0{
					if ret==true{
						fmt.Println(count/1024/1024,"M ",str)
					}
				}
			}
		}()
	}

	for{
		time.Sleep(time.Second)
	}
}

func testServiceStorageUpload(key int64,img []byte) (string,bool) {
	json:=gabs.New()
	json.Set(img,"img")
	var jsonStr = []byte(json.String())
	req, err := http.NewRequest("POST", "http://192.168.10.25:10087/opt/upload", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "",false
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
		os.Exit(1)
		return "",false
	}else{
		return json.String(),true
	}
}

func testServiceStorageDelete(url string) {
	json:=gabs.New()
	json.Set(url,"url")
	var jsonStr = []byte(json.String())
	req, err := http.NewRequest("POST", "http://192.168.10.25:10087/opt/delete", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
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
		os.Exit(1)
		return
	}else{
		return
	}
}
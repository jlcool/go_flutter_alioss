// File: main.go
package aliossFlutter

import (
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"os"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"bytes"
)

//  Make sure to use the same channel name as was used on the Flutter client side.
const channelName = "aliossflutter"
const (
	PARAM_ENDPOINT                  = "endpoint"
	PARAM_ACCESSKEYID                  = "accessKeyId"
	PARAM_ACCESSKEYSECRET                  = "accessKeySecret"
	PARAM_BUCKET	="bucket"
	PARAM_KEY	="key"
	PARAM_ID	="id"
	PARAM_FILE	="file"
	PARAM_BYTE	="fileByte"
	PARAM_PATH	="path"
)
type AliossFlutterPlugin struct{
	channel *plugin.MethodChannel
}
type OssProgressListener struct {
	channel *plugin.MethodChannel
	key string
	id string
}


var endpoint string
var accessKeyId string
var accessKeySecret string

var _ flutter.Plugin = &AliossFlutterPlugin{} // compile-time type check

// 定义进度变更事件处理函数。
func (l *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	m1 := make(map[string]string)
	m1["currentSize"] = strconv.FormatInt(event.ConsumedBytes,10)
	m1["totalSize"] = strconv.FormatInt(event.TotalBytes,10)
	m1["id"]=l.id
	m1["key"]=l.key
	b,err := json.Marshal(m1)

	switch event.EventType {
	case oss.TransferStartedEvent:
	case oss.TransferDataEvent:
		l.channel.InvokeMethod("onProgress",string(b))
	case oss.TransferCompletedEvent:		
	case oss.TransferFailedEvent:		
	default:
	}

	
	if err != nil {
	}
}

func (p *AliossFlutterPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	p.channel = plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	p.channel.HandleFunc("init", handleInit)
	p.channel.HandleFunc("upload", p.handleUpload)
	p.channel.HandleFunc("uploadByte", p.handleUploadByte)
	p.channel.HandleFunc("download", p.handleDownload)
	p.channel.HandleFunc("secretInit", handleSecretInit)
	p.channel.HandleFunc("signurl", handleSignurl)
	p.channel.HandleFunc("delete", handleDelete)
	p.channel.HandleFunc("doesObjectExist", handleDoesObjectExist)
	p.channel.HandleFunc("asyncHeadObject", handleAsyncHeadObject)
	p.channel.HandleFunc("listObjects", handleListObjects)
	return nil // no error
}

func handleInit(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
func (p *AliossFlutterPlugin) handleUpload(arguments interface{}) (reply interface{}, err error) {
	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}
	var param_id string
	var param_key string
	var param_bucket string 
	var param_file string

	if id1, ok := args[PARAM_ID]; ok {
		param_id = id1.(string)
	}
	if key1, ok := args[PARAM_KEY]; ok {
		param_key = key1.(string)
	}
	if bucket1, ok := args[PARAM_BUCKET]; ok {
		param_bucket = bucket1.(string)
	}
	if file1, ok := args[PARAM_FILE]; ok {
		param_file = file1.(string)
	}
	m1 := make(map[string]string)
	m1["key"] = param_key
	m1["id"] = param_id
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
    }

    // 获取存储空间。
    bucket, err := client.Bucket(param_bucket)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
    }

    // 上传文件流。
    err = bucket.PutObjectFromFile(param_key, param_file,oss.Progress(&OssProgressListener{channel:p.channel,key:param_key,id:param_id}))
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
	}
	m1["result"]="success"
	b, err := json.Marshal(m1)
	p.channel.InvokeMethod("onUpload",string(b) )
	return nil, nil
}

func (p *AliossFlutterPlugin) handleUploadByte(arguments interface{}) (reply interface{}, err error) {
	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}
	var param_id string
	var param_key string
	var param_bucket string 
	var param_byte []byte

	if id1, ok := args[PARAM_ID]; ok {
		param_id = id1.(string)
	}
	if key1, ok := args[PARAM_KEY]; ok {
		param_key = key1.(string)
	}
	if bucket1, ok := args[PARAM_BUCKET]; ok {
		param_bucket = bucket1.(string)
	}
	if file1, ok := args[PARAM_BYTE]; ok {
		param_byte = file1.([]byte)
	}
	m1 := make(map[string]string)
	m1["key"] = param_key
	m1["id"] = param_id
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
    }

    // 获取存储空间。
    bucket, err := client.Bucket(param_bucket)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
    }

    // 上传文件流。
    err = bucket.PutObject(param_key,  bytes.NewReader(param_byte),oss.Progress(&OssProgressListener{channel:p.channel,key:param_key,id:param_id}))
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onUpload", m1)
        os.Exit(-1)
	}
	m1["result"]="success"
	b, err := json.Marshal(m1)
	p.channel.InvokeMethod("onUpload",string(b) )
	return nil, nil
}
func (p *AliossFlutterPlugin) handleDownload(arguments interface{}) (reply interface{}, err error) {

	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}
	var param_key string
	var param_bucket string 
	var param_path string

	var param_id string
	if id1, ok := args[PARAM_ID]; ok {
		param_id = id1.(string)
	}
	if key1, ok := args[PARAM_KEY]; ok {
		param_key = key1.(string)
	}
	if bucket1, ok := args[PARAM_BUCKET]; ok {
		param_bucket = bucket1.(string)
	}
	if path1, ok := args[PARAM_PATH]; ok {
		param_path = path1.(string)
	}

	m1 := make(map[string]string)
	m1["key"] = param_key
	m1["id"] = param_id
	// 创建OSSClient实例。
    client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onDownload", m1)
        os.Exit(-1)
    }

    // 获取存储空间。
    bucket, err := client.Bucket(param_bucket)
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onDownload", m1)
        os.Exit(-1)
    }

    // 下载文件到本地文件。
    err = bucket.GetObjectToFile(param_key, param_path,oss.Progress(&OssProgressListener{channel:p.channel,key:param_key,id:param_id}))
    if err != nil {
		fmt.Println("Error:", err)
		m1["message"]=err.Error()
		m1["result"]="fail"
		p.channel.InvokeMethod("onDownload", m1)
        os.Exit(-1)
	}
	m1["result"]="success"
	b, err := json.Marshal(m1)
	p.channel.InvokeMethod("onDownload", string(b))
	return nil, nil
}
func handleSecretInit(arguments interface{}) (reply interface{}, err error) {
	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}

	if endpointt, ok := args[PARAM_ENDPOINT]; ok {
		endpoint = endpointt.(string)
	}
	if accessKeyIdd, ok := args[PARAM_ACCESSKEYID]; ok {
		accessKeyId = accessKeyIdd.(string)
	}
	if accessKeySecrett, ok := args[PARAM_ACCESSKEYSECRET]; ok {
		accessKeySecret = accessKeySecrett.(string)
	}
	return nil, nil
}
func handleSignurl(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
func handleDelete(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
func handleDoesObjectExist(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
func handleAsyncHeadObject(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
func handleListObjects(arguments interface{}) (reply interface{}, err error) {
	return nil, nil
}
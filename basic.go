package OceanEngineApi

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ApiVersion   = "2"
	ApiUrlPrefix = "https://ad.oceanengine.com/open_api/"
)

const (
	ContentTypeJson = "application/json"
)

type respImplement interface {
	String() string
}

type OceanEngineResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

func (r *OceanEngineResp) String() string {
	return "OceanEngineResp"
}

func (r *OceanEngineResp) doRequest(api *OceanEngineApi, req *http.Request, res respImplement) error {
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	if err := api.checkResp(req, bytes); err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.Unmarshal(bytes, res); err != nil {
		return err
	}

	return nil
}

type OceanEngineApi struct {
	accessToken string
	client      *http.Client
}

func NewOceanEngineApi(token string) *OceanEngineApi {
	return &OceanEngineApi{
		accessToken: token,
		client:      new(http.Client),
	}
}

// @function: 创建一个 request
func (api *OceanEngineApi) NewRequest(method string, url string, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Access-Token", api.accessToken)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// @function: 检查 oceanengine 返回是否正常, 不正常则设置 Debug 模式获取返回的 message 以便于调试错误
func (api *OceanEngineApi) checkResp(req *http.Request, body []byte) error {
	data := OceanEngineResp{}

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.Code == 0 {
		return nil
	}

	req.Header.Add("X-Debug-Mode", "1")

	debugResp, err := api.client.Do(req)

	if err != nil {
		return err
	}

	defer debugResp.Body.Close()

	bytes, _ := ioutil.ReadAll(debugResp.Body)
	data = OceanEngineResp{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	return errors.New(data.Message)
}

// @function: 构造文件上传表单
// params: k 表单的 key 值, v 表单的 value
// files: k 表单的 key 值, v 文件路径
func (api *OceanEngineApi) formCompose(params map[string]string, files map[string]string) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 写入文件
	for k, v := range files {
		file, err := os.Open(v)
		if err != nil {
			return nil, "", err
		}

		part, err := writer.CreateFormFile(k, filepath.Base(v))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, file)
		file.Close()
	}

	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

// @function: 压缩文件
// @remark: srcFile 可以是一个文件或者一个目录
func (api *OceanEngineApi) zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// @function: 随机字符串
func (api *OceanEngineApi) randomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

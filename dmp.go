package OceanEngineApi

import (
	"encoding/base64"
	"github.com/MrNullPoint/OceanEngineApi/pb"
	"github.com/golang/protobuf/proto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ApiDmpDataSourceFileUpload = ApiUrlPrefix + ApiVersion + "/dmp/data_source/file/upload/"
	API_DMP_DATA_SOURCE_CREATE = ApiUrlPrefix + ApiVersion + "/dmp/data_source/file/upload/"
)

// @function: 构建 DMP 所需要的上传的 zip 文件
// @params<path>: zip 文件生成的路径
// @params<name>: 可以指定 zip 文件名
func (api *OceanEngineApi) DataSourceFileCompose(data []*pb.DmpData, path string) (string, error) {
	name := api.randomString(32)

	binPath := filepath.Join(path, name)
	zipPath := filepath.Join(path, name+".zip")

	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	defer file.Close()

	for _, d := range data {
		b, err := proto.Marshal(d)
		if err != nil {
			return "", err
		}
		s := base64.StdEncoding.EncodeToString(b)
		if _, err := file.WriteString(s + "\n"); err != nil {
			return "", err
		}
	}

	if err := api.zip(binPath, zipPath); err != nil {
		return "", err
	}

	if err := os.Remove(binPath); err != nil {
		return "", err
	}

	return zipPath, nil
}

// @function: 数据源文件上传
func (api *OceanEngineApi) DataSourceFileUpload(file string, advertiserId string) (string, error) {
	params := make(map[string]string)
	params["advertiseId"] = advertiserId

	files := make(map[string]string)
	files["file"] = file

	body, contentType, err := api.formCompose(params, files)
	if err != nil {
		return "", err
	}

	req, err := api.NewRequest("POST", ApiDmpDataSourceFileUpload, contentType, body)
	if err != nil {
		return "", err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	if err := api.checkResp(req, resp); err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	return gjson.ParseBytes(bytes).Get("data.filepath").String(), nil
}

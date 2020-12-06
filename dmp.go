package OceanEngineApi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/MrNullPoint/OceanEngineApi/pb"
	"github.com/golang/protobuf/proto"
	"os"
	"path/filepath"
	"strconv"
)

const (
	ApiDmpDataSourceFileUpload = ApiUrlPrefix + ApiVersion + "/dmp/data_source/file/upload/"
	ApiDmpDataSourceCreate     = ApiUrlPrefix + ApiVersion + "/dmp/data_source/create/"
	ApiDmpDataSourceUpdate     = ApiUrlPrefix + ApiVersion + "/dmp/data_source/update/"
	ApiDmpDataSourceDetail     = ApiUrlPrefix + ApiVersion + "/dmp/data_source/read/"
	ApiDmpAudiencePublish      = ApiUrlPrefix + ApiVersion + "/dmp/custom_audience/publish/"
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

type DataSourceFileUploadResp struct {
	OceanEngineResp
	Data struct {
		FilePath string `json:"file_path"`
	} `json:"data"`
}

// @function: 数据源文件上传
// @reference: https://ad.oceanengine.com/openapi/doc/index.html?id=501
func (api *OceanEngineApi) DataSourceFileUpload(file string, advertiserId int) (*DataSourceFileUploadResp, error) {
	params := make(map[string]string)
	params["advertiseId"] = strconv.Itoa(advertiserId)

	files := make(map[string]string)
	files["file"] = file

	body, contentType, err := api.formCompose(params, files)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest("POST", ApiDmpDataSourceFileUpload, contentType, body)
	if err != nil {
		return nil, err
	}

	resp := new(DataSourceFileUploadResp)
	if err := resp.doRequest(api, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type DataSourceCreateResp struct {
	OceanEngineResp
	Data struct {
		DataSourceId string `json:"data_source_id"`
	} `json:"data"`
}

// @function: 数据源创建
// @reference: https://ad.oceanengine.com/openapi/doc/index.html?id=502
func (api *OceanEngineApi) DataSourceCreate(advertiserId int, dataSourceName string, dataSourceType string,
	desc string, format int, storageType int, paths []string) (*DataSourceCreateResp, error) {
	// 默认投放数据源类型为 UID
	if dataSourceType == "" {
		dataSourceType = "UID"
	}

	if dataSourceName == "" || len(dataSourceName) >= 30 || len(desc) >= 256 ||
		len(paths) >= 1000 || len(paths) == 0 || (dataSourceType != "UID" && dataSourceType != "DID") {
		return nil, errors.New("data source create params check failed")
	}

	params := make(map[string]interface{})

	params["advertiser_id"] = advertiserId
	params["data_source_name"] = dataSourceName
	params["description"] = desc
	params["file_paths"] = paths
	params["data_format"] = 0
	params["file_storage_type"] = 0

	if format != 0 {
		params["data_format"] = format
	}

	if storageType != 0 {
		params["file_storage_type"] = storageType
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest("POST", ApiDmpDataSourceCreate, ContentTypeJson, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp := new(DataSourceCreateResp)
	if err := resp.doRequest(api, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type DataSourceUpdateResp struct {
	OceanEngineResp
	Data struct{} `json:"data"`
}

// @function: 数据源更新
// @reference: https://ad.oceanengine.com/openapi/doc/index.html?id=504
func (api *OceanEngineApi) DataSourceUpdate(advertiserId int, dataSourceId string, operationType int,
	format int, storageType int, paths []string) (*DataSourceUpdateResp, error) {
	if dataSourceId == "" || (operationType != 1 && operationType != 2 && operationType != 3) ||
		len(paths) >= 200 || len(paths) == 0 {
		return nil, errors.New("data source update params check failed")
	}

	params := make(map[string]interface{})

	params["advertiser_id"] = advertiserId
	params["data_source_id"] = dataSourceId
	params["operation_type"] = operationType
	params["file_paths"] = paths
	params["data_format"] = 0
	params["file_storage_type"] = 0

	if format != 0 {
		params["data_format"] = format
	}

	if storageType != 0 {
		params["file_storage_type"] = storageType
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest("POST", ApiDmpDataSourceUpdate, ContentTypeJson, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp := new(DataSourceUpdateResp)
	if err := resp.doRequest(api, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type DataSourceDetailResp struct {
	OceanEngineResp
	Data struct {
		DataList []struct {
			Name                       string `json:"name"`
			DataSourceId               string `json:"data_source_id"`
			Description                string `json:"description"`
			Status                     int    `json:"status"`
			CoverNum                   int64  `json:"cover_num"`
			UploadNum                  int64  `json:"upload_num"`
			CreateTime                 int64  `json:"create_time"`
			ModifyTime                 int64  `json:"modify_time"`
			LatestPublishedChangeLogId int    `json:"latest_published_change_log_id"`
			LatestPublishedTime        int64  `json:"latest_published_time"`
			DataSourceType             string `json:"data_source_type"`
			DefaultAudience            struct {
				AdvertiserId     int    `json:"advertiser_id"`
				CustomAudienceId int    `json:"custom_audience_id"`
				Name             string `json:"name"`
				CustomType       int    `json:"custom_type"`
				Source           string `json:"source"`
				Status           int    `json:"status"`
				PushStatus       int    `json:"push_status"`
				UploadNum        int64  `json:"upload_num"`
				CoverNum         int64  `json:"cover_num"`
				ExpiryDate       string `json:"expiry_date"`
				CreateTime       int64  `json:"create_time"`
				ModifyTime       int64  `json:"modify_time"`
				Isdel            int    `json:"isdel"`
				CalculateSubType int    `json:"calculate_sub_type"`
				CalculateType    int    `json:"calculate_type"`
				DataSourceId     string `json:"data_source_id"`
				Tag              string `json:"tag"`
				ThirdPartyInfo   string `json:"third_party_info"`
				DeliveryStatus   string `json:"delivery_status"`
			} `json:"default_audience"`
		} `json:"data_list"`
	} `json:"data"`
}

// @function: 数据源详细信息
// @reference: https://ad.oceanengine.com/openapi/doc/index.html?id=505
func (api *OceanEngineApi) DataSourceDetail(advertiserId int, dataSourceIds []string) (*DataSourceDetailResp, error) {
	if len(dataSourceIds) == 0 || len(dataSourceIds) >= 400 {
		return nil, errors.New("data source detail params check failed")
	}

	params := make(map[string]interface{})

	params["advertiser_id"] = advertiserId
	params["data_source_id_list"] = dataSourceIds

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest("GET", ApiDmpDataSourceDetail, ContentTypeJson, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp := new(DataSourceDetailResp)
	if err := resp.doRequest(api, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type AudiencePublishResp struct {
	OceanEngineResp
	Data struct{} `json:"data"`
}

// @function: 发布人群包
// @reference: https://ad.oceanengine.com/openapi/doc/index.html?id=508
func (api *OceanEngineApi) AudiencePublish(advertiserId int, audienceId int) (*AudiencePublishResp, error) {
	params := make(map[string]interface{})

	params["advertiser_id"] = advertiserId
	params["custom_audience_id"] = audienceId

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest("POST", ApiDmpAudiencePublish, ContentTypeJson, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp := new(AudiencePublishResp)
	if err := resp.doRequest(api, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

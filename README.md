## 巨量引擎 Golang API

### 1. 安装 INSTALL

```shell
go get github.com/MrNullPoint/OceanEngineApi
```

### 2. 使用方式

```go
api := NewOceanEngineApi("access_token")
```

### 3. 完成功能

| 分类 | 功能                      | 接口                      | 测试   |
| ---- | ------------------------- | ------------------------- | ------ |
| DMP  | 生成 DMP 数据源文件压缩包 | api.DataSourceFileCompose | passed |
| DMP  | 上传数据源文件            | api.DataSourceFileUpload  |        |
| DMP  | 数据源创建                | api.DataSourceCreate      |        |
| DMP  | 数据源更新                | api.DataSourceUpdate      |        |
| DMP  | 数据源详细信息            | api.DataSourceDetail      |        |
| DMP  | 发布人群包                | api.AudiencePublish       |        |
| DMP  | 推送人群包                | api.AudiencePush          |        |

## 其他注意

1. 巨量引擎目前 DMP 数据格式要求是 PROTOBUF V2 的版本，我特 protoc 了一个根据头条原来 proto 文件对应的 `toutiao_dmp.pb.go` ，同时修改了头条 DMP 文件说明页面中的 python 代码几处小 bug，修改 bug 的文件对应 pb 目录下的 `toutiao_dmp_test.py` 和 `toutiao_dmp_validate.py`
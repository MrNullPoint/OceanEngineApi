package OceanEngineApi

import (
	"crypto/md5"
	"fmt"
	"github.com/MrNullPoint/OceanEngineApi/pb"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestDataSourceFileCompose(t *testing.T) {
	api := NewOceanEngineApi("test")
	ti := uint32(time.Now().Unix())
	dt := pb.IdItem_IDFA_MD5
	id1 := fmt.Sprintf("%x", md5.Sum([]byte(api.randomString(10))))
	id2 := fmt.Sprintf("%x", md5.Sum([]byte(api.randomString(10))))

	data := []*pb.DmpData{
		{
			IdList: []*pb.IdItem{
				{
					Timestamp: &ti,
					DataType:  &dt,
					Id:        &id1,
					Tags:      []string{"test1", "test2"},
				},
			},
		},
		{
			IdList: []*pb.IdItem{
				{
					Timestamp: &ti,
					DataType:  &dt,
					Id:        &id2,
					Tags:      []string{"test3", "test4"},
				},
			},
		},
	}

	zipPath, err := api.DataSourceFileCompose(data, "/tmp")
	assert.Nil(t, err)
	assert.NotEqual(t, zipPath, "")
}

func TestOceanEngineApi_DataSourceFileUpload(t *testing.T) {
	api := NewOceanEngineApi("964fafa7e7c4b6ebed752f09108429eaeac9d98a")

	paths := []string{
		"/tmp/7da65e52f9e1f17471cfe555d33487ca.zip",
		"/tmp/fdd569b1458c7fe1e5a8c6958fe7716b.zip",
		"/tmp/b88bc726cc937d2eff166fab3b5c2633.zip",
		"/tmp/01769e3ce7639c1a3a8b024b53106a35.zip",
		"/tmp/952c4be36140c651551a578abb976232.zip",
		"/tmp/1f34087616da01f89c6db9f79890b6ea.zip",
		"/tmp/7121fda0f1ecb08acf0e1e099c4a09a1.zip",
	}

	for _, p := range paths {
		resp, err := api.DataSourceFileUpload(p, 1667201949630478)
		assert.Nil(t, err)
		log.Println(resp.Data.FilePath)
	}
}

func TestOceanEngineApi_DataSourceFileCreate(t *testing.T) {
	api := NewOceanEngineApi("964fafa7e7c4b6ebed752f09108429eaeac9d98a")

	paths := []string{
		"1667201949630478-0fba4c5ebd65c89d876880059a564737",
		"1667201949630478-c0274f88cc0738da95f7d332fa66d80e",
		"1667201949630478-0f4288b5caeefba2189cb8fad484be2e",
		"1667201949630478-693182e1a952b9f44aaa8405411dcfc5",
		"1667201949630478-debb1dd33ceff27a5f6ca5468d2ebe9b",
		"1667201949630478-2f1d8044f189f0316eca99d4c541f0a9",
		"1667201949630478-c8dd6c37aeadc4f1054fa190caf97dc7",
	}

	resp, err := api.DataSourceCreate(1667201949630478, "test", "UID", "test",
		0, 0, paths)

	assert.Nil(t, err)

	log.Println(resp)
}

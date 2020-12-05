package OceanEngineApi

import (
	"crypto/md5"
	"fmt"
	"github.com/MrNullPoint/OceanEngineApi/pb"
	"github.com/stretchr/testify/assert"
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

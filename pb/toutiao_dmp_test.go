package pb

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"testing"
	"time"
)

func TestToutiaoDMPGenerate(t *testing.T) {
	dmpData1 := DmpData{}
	dmpData2 := DmpData{}

	var (
		id1 = "356145080566857"
		ts1 = uint32(time.Now().Unix())
		dt1 = IdItem_IMEI
	)

	var (
		id2 = "1E2DFA89-496A-47FD-9941-DF1FC4E6484A"
		ts2 = uint32(time.Now().Unix())
		dt2 = IdItem_IDFA
	)

	var (
		id3 = "136145080566857"
		ts3 = uint32(time.Now().Unix())
		dt3 = IdItem_IMEI
	)

	var (
		id4 = "642DFA89-496A-47FD-9941-DF1FC4E6484A"
		ts4 = uint32(time.Now().Unix())
		dt4 = IdItem_IDFA
	)

	idItem1 := IdItem{
		Timestamp: &ts1,
		DataType:  &dt1,
		Id:        &id1,
		Tags:      []string{"信用卡", "理财"},
	}

	idItem2 := IdItem{
		Timestamp: &ts2,
		DataType:  &dt2,
		Id:        &id2,
		Tags:      []string{"黄金", "理财"},
	}

	idItem3 := IdItem{
		Timestamp: &ts3,
		DataType:  &dt3,
		Id:        &id3,
		Tags:      []string{"信用卡", "股票"},
	}

	idItem4 := IdItem{
		Timestamp: &ts4,
		DataType:  &dt4,
		Id:        &id4,
		Tags:      []string{"黄金", "理财"},
	}

	dmpData1.IdList = append(dmpData1.IdList, &idItem1)
	dmpData1.IdList = append(dmpData1.IdList, &idItem2)

	dmpData2.IdList = append(dmpData2.IdList, &idItem3)
	dmpData2.IdList = append(dmpData2.IdList, &idItem4)

	binaryString1, _ := proto.Marshal(&dmpData1)
	resultString1 := base64.StdEncoding.EncodeToString(binaryString1)

	binaryString2, _ := proto.Marshal(&dmpData2)
	resultString2 := base64.StdEncoding.EncodeToString(binaryString2)

	file, err := os.OpenFile("target_pb2", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	if _, err := file.WriteString(resultString1 + "\n"); err != nil {
		log.Panic(err)
	}
	//if _, err := file.WriteString("\n"); err != nil {
	//	log.Panic(err)
	//}
	if _, err := file.WriteString(resultString2 + "\n"); err != nil {
		log.Panic(err)
	}
	//if _, err := file.WriteString("\n"); err != nil {
	//	log.Panic(err)
	//}
}

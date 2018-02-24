package main

import (
	"fmt"
	"dbOptions"
	"github.com/syndtr/goleveldb/leveldb/util"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func main()  {

	stdin := bufio.NewReader(os.Stdin)
	var dbIdStrs string
	fmt.Println("input dbIds, split by ,")
	fmt.Fscan(stdin,&dbIdStrs)
	dbIdStrArray := strings.Split(dbIdStrs, ",")
	for _, dbIdStr := range dbIdStrArray{
		dbIdS,_ := strconv.Atoi(dbIdStr)
		curDbId := uint8(dbIdS)

		imgToIndexDB := dbOptions.InitMuImgToIndexDB(curDbId)
		indexToImgDB := dbOptions.InitMuIndexToImgDB(curDbId)
		indexToClipDB := dbOptions.InitMuIndexToClipDB(curDbId)
		clipToIndexDB := dbOptions.InitMuClipToIndexDB(curDbId)

		imgToIndexDB.DBPtr.CompactRange(util.Range{nil,nil})
		indexToImgDB.DBPtr.CompactRange(util.Range{nil,nil})
		indexToClipDB.DBPtr.CompactRange(util.Range{nil,nil})
		clipToIndexDB.DBPtr.CompactRange(util.Range{nil,nil})

		imgToIndexDB.CloseDB()
		indexToImgDB.CloseDB()
		indexToClipDB.CloseDB()
		clipToIndexDB.CloseDB()
	}

}
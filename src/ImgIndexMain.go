package main

import (
	"bufio"
	"os"
	"fmt"
	"dbOptions"
)

func SaveMainImgIndexes()  {
	stdin := bufio.NewReader(os.Stdin)
	var input int
	var dbIndex uint8

	for{
		fmt.Println("select a image db to deal: ")
		fmt.Fscan(stdin,&dbIndex)
		imgDB := dbOptions.PickImgDB(dbIndex)

		imgIndexDB := dbOptions.InitIndexToImgDB()
		imgToIndexDB := dbOptions.InitImgToIndexDB()
		if nil == imgIndexDB{
			fmt.Println("open img index db failed")
			return
		}

		fmt.Println("input how many times each thread(16 in total) to deal: ")
		fmt.Fscan(stdin,&input)
		dbOptions.ImgIndexSaveRun(dbIndex, input)

		imgDB.CloseDB()
		imgIndexDB.CloseDB()
		imgToIndexDB.CloseDB()
	}
}

func main(){
	SaveMainImgIndexes()
}
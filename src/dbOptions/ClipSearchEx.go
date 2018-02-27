package dbOptions

import (
	"strconv"
	"util"
	"fmt"
	"imgCache"
	"imgIndex"
	"bytes"
)


/**
	计算哪些大图中联合出现了 imgKey 中的多个子图, imgKey 不包含在内
 */
func occInImgsEx(dbId uint8, imgKey []byte) (occedImgIndex *imgCache.MyMap, allStatBranchesIndex[] [][]byte ){
	curImgIdent := make([]byte, ImgIndex.IMG_IDENT_LENGTH)
	curImgIdent[0] = byte(dbId)
	copy(curImgIdent[1:], imgKey)
	curImgIndex := InitMuImgToIndexDB(uint8(dbId)).ReadFor(curImgIdent)

	seeker := NewMultyDBReader(GetInitedClipStatIndexToIdentDB())

	clipIndexes := QueryClipIndexesFor(dbId, imgKey)
	if nil == clipIndexes{
		fmt.Println("can't find clip indexes: ", string(ImgIndex.ParseImgKeyToPlainTxt(imgKey)))
		return
	}

	//各个子图的 imgIndexContainer 容器: 每个容器表示某些 img index 与相应的子图有关系(即子图也出现在这些 img 中)
	occedImgIndex = imgCache.NewMyMap(true)

	allStatBranchesIndex = make([] [][]byte, len(clipIndexes))

	//cmap := imgCache.NewMyMap(false)
	//var exsitsIndexes []interface{}

	curClipOccIn := imgCache.NewMyMap(false)

	for i, clipIndex := range clipIndexes{
		curStatBranches := ImgIndex.ClipStatIndexBranch(clipIndex)
		allStatBranchesIndex[i] = curStatBranches

		/*
		//当前子图可能已经处理过. 注意这里会有一个反噬的作用: 设置的 search_conf 越宽松, 越有可能认为
		//clipIndex 已经处理过. 这会丢失一些匹配.
		//实际上不需要做这个判断, 一是同一张大图中有两张相同的子图的情况很少, 二是即使有重复的子图也能正确地处理
		if cmap.ContainsAnyOneOf(curStatBranches){
			//再次确认
			skip := false
			exsitsIndexes = cmap.QueryUnion(curStatBranches)
			for _, interfaceIndex := range exsitsIndexes{
				curIndex := interfaceIndex.([]byte)
				if isSameClip(curIndex, clipIndex){
					fmt.Print("clip index has dealed: ")
					fileUtil.PrintBytes(clipIndex)
					fmt.Print("as like: ")
					fileUtil.PrintBytes(curIndex)
					skip = true
					break
				}
			}
			if skip{
				continue
			}
		}
		*/

		//计算所有与当前子图相似的子图出现在哪此大图中
		for _,branch := range curStatBranches{
		//	cmap.Put(branch, clipIndex)
			clipIdentsSet := seeker.ReadFor(branch)
			if 0 == len(clipIdentsSet){
				continue
			}
			for _,clipIdents := range clipIdentsSet{
				if 0 != len(clipIdents)%ImgIndex.IMG_CLIP_IDENT_LENGTH{
					fmt.Println("value length for stat index db is not multple of ", ImgIndex.IMG_CLIP_IDENT_LENGTH)
					continue
				}
				for l:=0;l < len(clipIdents);l += ImgIndex.IMG_CLIP_IDENT_LENGTH{
					sameClip := clipIdents[l:l+ImgIndex.IMG_CLIP_IDENT_LENGTH]
					curClipOccIn.Put(sameClip, nil)
				}
			}
		}

		clipIdents := curClipOccIn.KeySet()
		if 0 == len(clipIdents){
			continue
		}

		curClipOccedImgIndexes := getImgIndex(clipIdents)
		imgIndexes := curClipOccedImgIndexes.KeySet()
		if len(clipIdents) > 100{
	//		fmt.Println("OOPS, more than 100 same clip, is right????")
		}
	//	fmt.Print(dbId, "-", string(ImgIndex.ParseImgKeyToPlainTxt(imgKey)), "-", i, ": ")
		for _,imgIndex := range imgIndexes{
			//当前子图出现在下面的 img 中. 为了唯一性表示，使用 img index 作为键去表示
			interfaceClipIdent := curClipOccedImgIndexes.Get(imgIndex)
			if 1 != len(interfaceClipIdent){
				continue
			}
			clipIdent := interfaceClipIdent[0].([]byte)
			imgIndexDBId := clipIdent[0]
			imgIdent := clipIdent[0:5]

			//当前图跳过
			if bytes.Equal(curImgIdent, imgIdent){
				continue
			}

			//与当前大图是同一张图
			if bytes.Equal(curImgIndex, imgIndex){
				continue
			}

			//最后再使用欧拉距离验证到底是否是相似的子图
			sameClipIndex := InitMuClipToIndexDB(imgIndexDBId).ReadFor(clipIdent)
			if len(sameClipIndex) != ImgIndex.CLIP_INDEX_BYTES_LEN{
				continue
			}
			if isSameClip(sameClipIndex, clipIndex){
				continue
			}

		//	fmt.Print(imgIndexDBId,"-", string(ImgIndex.ParseImgKeyToPlainTxt(imgIdent[1:])), "-", clipIdent[5]," | ")

			occedImgIndex.Put(imgIndex, uint8(i))	//第 i 个子图出现在 imgIndex 所指示的大图中
		}
		fmt.Println()

		curClipOccIn.Clear()
	}
	return
}

func isSameClip(left, right []byte) bool {
	return ImgIndex.TheclipSearchConf.Delta_Eul_square < fileUtil.CalEulSquare(left, right)
}

func getImgIndex(clipIdents [] []byte) *imgCache.MyMap {
	cmap := imgCache.NewMyMap(false)
	for _,clipIdent := range clipIdents{
		imgIndexDBId := clipIdent[0]
		imgIdent := clipIdent[0:5]
		imgIndex := InitMuImgToIndexDB(uint8(imgIndexDBId)).ReadFor(imgIdent)
		if len(imgIndex) == 0{
			continue
		}

		cmap.Put(imgIndex, clipIdent)
	}
	return cmap
}


/*
	以 dbId 库中的 imgKey 为对象，找出 imgKey 中哪些子图共同出现在其它大图中
*/
func SearchCoordinateForClipEx(dbId uint8, imgKey []byte) (whichesGroupAndCount *imgCache.MyMap, allStatIndex [] [][]byte ) {
	imgName := strconv.Itoa(int(dbId)) + "-" + string(ImgIndex.ParseImgKeyToPlainTxt(imgKey))
	//计算哪些大图中联合出现了 imgKey 中的多个子图. 注意 imgKey 不包含在内
	occedImgIndex, allStatIndex := occInImgsEx(dbId, imgKey)

	if nil == occedImgIndex || 0 == occedImgIndex.KeyCount(){
		return
	}

	var resWhiches [][]uint8

	motherImgIndexes := occedImgIndex.KeySet()

	for _,imgIndex := range motherImgIndexes {
		interfaceWhiches := occedImgIndex.Get(imgIndex)
		if 2 > len(interfaceWhiches){
			continue
		}
		whiches := make([]uint8, len(interfaceWhiches))
		for i,which := range interfaceWhiches{
			whiches[i] = which.(uint8)
		}
		whiches = fileUtil.RemoveDupplicatedBytes(whiches)
		if 2 > len(whiches){
			continue
		}

		fileUtil.BytesSort(whiches)

		resWhiches = append(resWhiches, whiches)

	}
	if 0 == len(resWhiches){
		return
	}

	if len(resWhiches) > 1{
		fmt.Println("okay, find len(resWhiches) > 1: ", len(resWhiches), ", ",imgName)
	}

	//校准次数. 注意校验不能在 statCoordinateResult 中边统计边校准：只能最终校准.
	//原因在于: 设 1,3,4 同时出现在 A, B 图中，3,4 同时出现在 A,B,C 中则
	whichesGroupAndCount = statCoordinateResult(resWhiches)
	whichesGroups := whichesGroupAndCount.KeySet()
	for _,whiches := range whichesGroups{

		interfaceCounts := whichesGroupAndCount.Get(whiches)
		if 1 == len(interfaceCounts){
			countExclusiveCurrentImg := interfaceCounts[0].(int)
			whichesGroupAndCount.Put(whiches, countExclusiveCurrentImg + 1)
		}
	}


	//打印
	if len(whichesGroups) > 0{
		showStr := imgName + " : "
		for _,whiches := range whichesGroups{
			showStr += "["
			for _,which := range whiches{
				showStr += strconv.Itoa(int(which)) + ","
			}
			showStr += "]"
			interfaceCounts := whichesGroupAndCount.Get(whiches)
			if 1 == len(interfaceCounts){
				showStr += "-" + strconv.Itoa(interfaceCounts[0].(int)) + " | "
			}
		}
		fmt.Println(showStr)
	}

	return
}

func CalClipStatBranchIndexes(clipIdent []byte) [][]byte {
	dbId ,_,_:= ImgIndex.ParseAImgClipIdentBytes(clipIdent)
	clipIndex := InitMuClipToIndexDB(dbId).ReadFor(clipIdent)
	return ImgIndex.ClipStatIndexBranch(clipIndex)
}

func PrintClipStatBranchIndexes(clipIdent []byte)  {
	indexes := CalClipStatBranchIndexes(clipIdent)
	fmt.Println("stat indexes for ", clipIdent)
	for _,index := range indexes{
		fileUtil.PrintBytes(index)
	}
}

func TestClipStatBranchIndeses() {
	clipIdent1 := make([]byte, ImgIndex.IMG_CLIP_IDENT_LENGTH)
	clipIdent1[0] = 2
	copy(clipIdent1[1:], ImgIndex.FormatImgKey([]byte("F0000067")))
	clipIdent1[5] = 7

	clipIdent2 := make([]byte, ImgIndex.IMG_CLIP_IDENT_LENGTH)
	clipIdent2[0] = 2
	copy(clipIdent2[1:], ImgIndex.FormatImgKey([]byte("A0000000")))
	clipIdent2[5] = 3

	clipIdent3 := make([]byte, ImgIndex.IMG_CLIP_IDENT_LENGTH)
	clipIdent3[0] = 2
	copy(clipIdent3[1:], ImgIndex.FormatImgKey([]byte("E0000150")))
	clipIdent3[5] = 0

	PrintClipStatBranchIndexes(clipIdent1)
	PrintClipStatBranchIndexes(clipIdent2)
	PrintClipStatBranchIndexes(clipIdent3)
}

func TestStatIndexValue()  {

	/*
	referClipIdents := [] []byte{[]byte{2,70,0,0,67,7}, []byte{2,65,0,0,0,3}, []byte{2,69,0,0,150,0} }
	indexBytes := []byte{222, 57}
	clipIdentList := InitClipStatIndexToIdentsDB(2).ReadFor(indexBytes)

	if 0 == len(clipIdentList) || len(clipIdentList) % ImgIndex.IMG_CLIP_IDENT_LENGTH != 0{
		fmt.Println("error")
		return
	}

	clipIdents := make([][]byte, len(clipIdentList) / ImgIndex.IMG_CLIP_IDENT_LENGTH)
	ci := 0
	for i:=0;i < len(clipIdentList); i += ImgIndex.IMG_CLIP_IDENT_LENGTH{
		clipIdents[ci] = fileUtil.CopyBytesTo(clipIdentList[i: i + ImgIndex.IMG_CLIP_IDENT_LENGTH])
		ci ++
	}


	for _,clipIdent := range clipIdents{
	//	fileUtil.PrintBytes(clipIdent)
		for i,refer := range referClipIdents{
			if bytes.Equal(clipIdent, refer){
				fmt.Println("contain ", i)
			}
		}
	}

*/

	InitClipStatIndexToIdentsDB(2)
	occInImgsEx(2, []byte{65,0,0,0})

}

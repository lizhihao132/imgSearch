package ImgIndex

import (
	"config"
	"fmt"
)

var CLIP_INDEX_STAT_BYTES_LEN int = 2
var CLIP_INDEX_BYTES_LEN int = 72
var CLIP_INDEX_BRANCH_BITS int = 2
var CLIP_INDEX_BRANCH_BOUND uint8 = 10

/**
	获得一个 jpg 图像的 index
	内部实现的原理是，获得它里面每张切图的 index，将其拼接在一起
 */
func GetIndexFor(data [][][]uint8) [] byte {
	height := len(data)
	width := len(data[0])
	imgConfig := config.GetImgConfigBySize(height, width)
	clipConfig := imgConfig.TheClipConfig

	if nil == imgConfig || nil == clipConfig{
		fmt.Println("can't get img config for height: ", height, ", width: ", width)
		return nil
	}

	//自身就是大图，它不属于任何大图，所以 mainImgKey 为 nil
	clipsIndexes := GetClipsIndexOfImgEx(data,imgConfig.Id, nil,imgConfig.TheClipConfig.ClipOffsets, imgConfig.OverrideClipLength)

	return GetFlatIndexBytesFrom(clipsIndexes)
}

package bptree

/**
	查找key 在 keys[] 中的位置
	returns:
		previousKeyIndex: 前一个key的index
		nextKeyIndex: 后一个key的index
		如果 previousKeyIndex == nextKeyIndex && previousKeyIndex > 0 则key存在并且index为previousKeyIndex
 */
func binarySearchIndexKey(keys []string, key string) (previousKeyIndex int, nextKeyIndex int) {
	if keys == nil || len(keys) <= 0 {
		return -1, -1
	}
	//fmt.Println(" 二分查找。indexkey。。")
	var low int = 0
	var height int = len(keys)

	for low <= height {
		//fmt.Println(" 。。。。比较。。")
		var mid int = low + (height-low)/2
		if keys[mid] == key {
			return mid, mid
		} else if key > keys[mid] { // 如果新的key 大于中间值的key,则查找右边
			if mid == height-1 || (mid < height-1 && key < keys[mid+1]) {
				return mid, mid + 1
			}
			low = mid + 1
		} else if key < keys[mid] { // 如果新的key 小于中间值的key,则查找左边

			if mid == 0 || (mid > 0 && key > keys[mid-1]) {

				return mid - 1, mid
			}
			height = mid - 1
		}
	}
	return -1, -1
}

/**
 查找key对应的indexNode 在 children[]中的位置
 returns:
		currentIndexNode: 当前索引节点
		indexAtCurrentIndexNode： key在当前索引节点的children[]中的index
 */
func binarySearchIndexNode(indexNode *IndexNode, key string) (currentIndexNode *IndexNode, indexAtCurrentIndexNode int) {

	if indexNode == nil {
		return nil, -1
	}
	if indexNode.IsLeaf {
		return indexNode, -1
	}

	//fmt.Println(" 二分查找 index node index。。。")
	var low int = 0
	var height int = len(indexNode.Keys)

	for low <= height {
		//fmt.Println(" 。。。。。比较。。")
		var mid int = low + (height-low)/2

		if indexNode.Keys[mid] == key { // 如果存在这个key

			return indexNode.Children[mid+1], mid + 1
		} else if key > indexNode.Keys[mid] { // 如果新的key 大于中间值的key,则查找右边
			if mid == len(indexNode.Keys)-1 || (mid < len(indexNode.Keys)-1 && key < indexNode.Keys[mid+1]) {

				return indexNode.Children[mid+1], mid + 1
			}
			low = mid + 1
		} else if key < indexNode.Keys[mid] { // 如果新的key 小于中间值的key,则查找左边

			if mid == 0 || (mid > 0 && key > indexNode.Keys[mid-1]) {

				return indexNode.Children[mid], mid
			}
			height = mid - 1
		}
	}

	return indexNode, -1
}

/** 二分查找 定位key对应的leafNode在indexNode节点的DataNodes[]中的位置
  returns:
       currentIndexNode: 当前叶子结点
       previousLeafIndexAtCurrentIndexNode: 前一个leafNode的index
       nextLeafIndexAtCurrentIndexNode: 后一个leafNode的index
	如果 previousLeafIndexAtCurrentIndexNode == nextLeafIndexAtCurrentIndexNode
			&& previousLeafIndexAtCurrentIndexNode > 0 则key对应的leafNode存在，并且index为previousLeafIndexAtCurrentIndexNode
**/
func binarySearchDataNode(indexNode *IndexNode, key string) (currentIndexNode *IndexNode, previousLeafIndexAtCurrentIndexNode int, nextLeafIndexAtCurrentIndexNode int) {
	if indexNode == nil || len(indexNode.DataNodes) <= 0 {
		return indexNode, -1, -1
	}
	//fmt.Println(" 二分查找 leafkey。。。")

	DataNodes := indexNode.DataNodes
	var low int = 0
	var height int = len(indexNode.DataNodes)

	for low <= height {
		//fmt.Println(" 。。。。比较。。")
		var mid int = low + (height-low)/2

		if DataNodes[mid].KeyAndValue.Key == key { // 如果存在这个key
			// 如果是叶子结点
			return indexNode, mid, mid

		} else if key > DataNodes[mid].KeyAndValue.Key { // 如果新的key 大于中间值的key,则查找右边
			if mid == len(DataNodes)-1 || (mid < len(DataNodes)-1 && key < DataNodes[mid+1].KeyAndValue.Key) {
				return indexNode, mid, mid + 1
			}

			low = mid + 1
		} else if key < DataNodes[mid].KeyAndValue.Key { // 如果新的key 小于中间值的key,则查找左边

			if mid == 0 || (mid > 0 && key > DataNodes[mid-1].KeyAndValue.Key) {

				return indexNode, mid - 1, mid

			}
			height = mid - 1
		}
	}
	return nil, -1, -1
}

package bptree

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// B+ 树
type BPTree struct {
	// TODO  lock
	M    int        // B+ 树的阶数
	Root *IndexNode // root节点
	Head *DataNode  // 头结点
}

// 索引节点
type IndexNode struct {
	Keys       []string     // 关键字，非叶子节点才有
	Children   []*IndexNode // 叶子结点没有children
	ParentNode *IndexNode

	IsLeaf    bool        // 是否是叶子结点
	DataNodes []*DataNode // 叶子结点才有数据结点
}

// 叶子的数据节点
type DataNode struct {
	KeyAndValue      KeyAndValue
	ParentNode       *IndexNode
	PreviousDataNode *DataNode
	NextDataNode     *DataNode
}

type KeyAndValue struct {
	Key   string // key 关键字
	Value string // value是对应数据在levelDB中的id、地址或者nil，内部节点的value为nil
}

// 初始化一个B+ 树
func MallocNewBPTree(m int) *BPTree {

	node := MallocNewIndexNode(true)
	return &BPTree{
		M:    m,
		Root: node,
		Head: nil}
}

// 初始化一个node
func MallocNewIndexNode(isLeaf bool) *IndexNode {
	return &IndexNode{
		IsLeaf: isLeaf,
	}
}

// 初始化一个数据结点
func MallocNewDataNode(keyAndValue KeyAndValue) *DataNode {
	return &DataNode{
		KeyAndValue: keyAndValue}
}

// 插入一个关键字和值
func (t *BPTree) Insert(keyAndValue KeyAndValue) (*DataNode, error) {

	keyAndValue.Key = strings.TrimSpace(keyAndValue.Key)
	if keyAndValue.Key == "" {
		return nil, errors.New(" key is nil ")
	}
	fmt.Printf("--- insert %s:%s --\n", keyAndValue.Key, keyAndValue.Value)
	// 只有根节点:插入第一个数据节点
	if t.Head == nil {
		dataNode := MallocNewDataNode(keyAndValue)
		dataNode.ParentNode = t.Root
		t.Root.DataNodes = []*DataNode{dataNode}
		t.Head = dataNode
		return t.Head, nil
	} else {
		// 查询是否存在该Key
		indexNode, oldDataNode, previousDataNodeIndex, _ := t.searchDataNode(keyAndValue.Key)

		// 存在这个key,直接替换值
		if oldDataNode != nil {
			oldDataNode.KeyAndValue.Value = keyAndValue.Value
			return indexNode.DataNodes[previousDataNodeIndex], nil
		}
		// 不存在时，新增一个数据结点dataNode
		newDataNode := MallocNewDataNode(keyAndValue)
		newDataNode.ParentNode = indexNode

		// 新增节点的Key小于当前组的所有值
		if previousDataNodeIndex < 0 {

			// 将数据结点按顺序左右链接起来
			newDataNode.PreviousDataNode = indexNode.DataNodes[0].PreviousDataNode
			newDataNode.NextDataNode = indexNode.DataNodes[0]
			indexNode.DataNodes[0].PreviousDataNode = newDataNode

			// 合并到DataNodes[]中
			indexNode.DataNodes = append([]*DataNode{newDataNode}, indexNode.DataNodes[:]...)

		} else {

			// 将叶子结点按顺序左右链接起来
			newDataNode.PreviousDataNode = indexNode.DataNodes[previousDataNodeIndex]
			newDataNode.NextDataNode = indexNode.DataNodes[previousDataNodeIndex].NextDataNode

			// 合并到DataNodes[]中
			tempSlice := append([]*DataNode{}, indexNode.DataNodes[:previousDataNodeIndex+1]...)
			tempSlice = append(tempSlice, newDataNode)
			tempSlice = append(tempSlice, indexNode.DataNodes[previousDataNodeIndex+1:]...)
			indexNode.DataNodes = tempSlice
		}

		if newDataNode.PreviousDataNode != nil {
			newDataNode.PreviousDataNode.NextDataNode = newDataNode
		} else {
			t.Head = newDataNode // 新的头结点
		}

		if newDataNode.NextDataNode != nil {
			newDataNode.NextDataNode.PreviousDataNode = newDataNode
		}

		// 树分裂
		t.divide(indexNode)
		return newDataNode, nil
	}

	return nil, nil
}

// 查询关键字
func (t *BPTree) Get(key string) (*DataNode, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("key is nil")
	}

	// 此处用二分查找到数据节点
	_, dataNode, _, _ := t.searchDataNode(key)
	return dataNode, nil
}

// 修改关键字的值
func (t *BPTree) Update(keyAndValue KeyAndValue) (bool, error) {
	keyAndValue.Key = strings.TrimSpace(keyAndValue.Key)
	if keyAndValue.Key == "" {
		return false, errors.New(" key is nil")
	}
	// 找到dataNode
	_, dataNode, _, _ := t.searchDataNode(keyAndValue.Key)

	// 存在这个key,修改数据
	if dataNode != nil {
		dataNode.KeyAndValue.Value = keyAndValue.Value
		return true, nil
	} else {
		return false, nil
	}
}

/**
returns:
  currentIndexNode: dataNode所在的当前节点
  dataNode: key所在dataNode
 */
func (t *BPTree) searchDataNode(key string) (currentIndexNode *IndexNode, dataNode *DataNode, preDataNodeIndexAtCurrentIndexNode int, nextDataIndexAtCurrentNode int) {
	indexNode, _ := binarySearchIndexNode(t.Root, key)
	for !indexNode.IsLeaf {
		indexNode, _ = binarySearchIndexNode(indexNode, key)
	}
	// 此处用二分查找 找到key所在叶子节点上DataNodes中的位置
	indexNode, previousDataNodeIndex, nextDataNodeIndex := binarySearchDataNode(indexNode, key)
	// 存在这个key
	if previousDataNodeIndex == nextDataNodeIndex && previousDataNodeIndex >= 0 {
		fmt.Printf("-----Get key:%s  value:%s , key [%s] is already exist --\n", key, indexNode.DataNodes[previousDataNodeIndex].KeyAndValue.Value, key)
		return indexNode, indexNode.DataNodes[previousDataNodeIndex], previousDataNodeIndex, nextDataNodeIndex
	} else {
		fmt.Printf("----Get key [%s] is not exist -- \n", key)
		return indexNode, nil, previousDataNodeIndex, nextDataNodeIndex
	}
}

func (t *BPTree) Remove(key string) (bool, error) {
	fmt.Printf("------delete %s----- \n", key)
	key = strings.TrimSpace(key)
	if key == "" {
		return false, errors.New("key is nil")
	}

	// 找到dataNode
	indexNode, removeDataNode, indexAtParent, _ := t.searchDataNode(key)

	// 存在该key,直接删除
	if removeDataNode != nil {

		// 修改被删除dataNode的前后dataNode的左右连接关系
		if removeDataNode.PreviousDataNode != nil {
			removeDataNode.PreviousDataNode.NextDataNode = removeDataNode.NextDataNode
		} else {
			t.Head = removeDataNode.NextDataNode
		}
		if removeDataNode.NextDataNode != nil {
			removeDataNode.NextDataNode.PreviousDataNode = removeDataNode.PreviousDataNode
		}
		// 删除节点
		tempSlice := append([]*DataNode{}, indexNode.DataNodes[:indexAtParent]...)
		tempSlice = append(tempSlice, indexNode.DataNodes[indexAtParent+1:]...)
		indexNode.DataNodes = tempSlice

		// 合并
		t.merge(indexNode)
		return true, nil
	} else {
		return false, nil
	}

	return false, nil
}

// 树 合并
func (t *BPTree) merge(indexNode *IndexNode) {
	if indexNode == nil || indexNode.ParentNode == nil {
		return
	}

	if indexNode.IsLeaf {
		if len(indexNode.DataNodes) >= t.M/2 {
			return
		}
		fmt.Println("---合并叶子节点-----")
		// 找到当前indexNode所在数组中的Index
		_, indexAtParent := binarySearchIndexNode(indexNode.ParentNode, indexNode.DataNodes[0].KeyAndValue.Key)
		if indexAtParent == 0 { // 找右兄弟
			rightNode := indexNode.ParentNode.Children[indexAtParent+1]
			// 如果右兄弟的dataNode > t.M/2 ,就借一个， 右兄弟的值都比左兄弟的大
			if len(rightNode.DataNodes) > t.M/2 {
				// 借的dataNode
				borrowDataNode := rightNode.DataNodes[0]
				// 加入到indexNode
				indexNode.DataNodes = append(indexNode.DataNodes, borrowDataNode)
				borrowDataNode.ParentNode = indexNode
				// 在右兄弟中删除该节点
				rightNode.DataNodes = rightNode.DataNodes[1:]
				// 右兄弟第一个节点的Key替换父节点相应的key
				indexNode.ParentNode.Keys[0] = rightNode.DataNodes[0].KeyAndValue.Key
			} else { // 如果右兄弟的dataNode <= t.M/2  则直接合并
				// 右边合到左边
				// 修改dataNode的父节点，并append到indexNode的DataNodes中
				for _, dataNode := range rightNode.DataNodes {
					dataNode.ParentNode = indexNode
					indexNode.DataNodes = append(indexNode.DataNodes, dataNode)
				}
				// 如果父节点的keys.len >1
				if len(indexNode.ParentNode.Keys) > 1 {

					// 父节点中删除rightNode
					tempSlice := append([]*IndexNode{indexNode}, indexNode.ParentNode.Children[indexAtParent+2:]...)
					indexNode.ParentNode.Children = tempSlice
					// 父节点中删除第一个key
					indexNode.ParentNode.Keys = indexNode.ParentNode.Keys[1:]
				} else { // 说明parentNode可能是root节点
					t.Root = indexNode
					indexNode.ParentNode = nil
				}
			}

		} else { // 找左兄弟
			leftNode := indexNode.ParentNode.Children[indexAtParent-1]

			if len(leftNode.DataNodes) > t.M/2 {
				// 借的dataNode
				borrowDataNode := leftNode.DataNodes[len(leftNode.DataNodes)-1]
				// 加入到indexNode
				indexNode.DataNodes = append([]*DataNode{borrowDataNode}, indexNode.DataNodes[:]...)
				borrowDataNode.ParentNode = indexNode
				// 在左兄弟中删除该节点
				leftNode.DataNodes = leftNode.DataNodes[:len(leftNode.DataNodes)-1]
				// 借的dataNode的Key替换父节点相应的key
				indexNode.ParentNode.Keys[indexAtParent-1] = borrowDataNode.KeyAndValue.Key
			} else { // 如果左兄弟的dataNode <= t.M/2  则直接合并
				// 右边合到左边
				// 修改dataNode的父节点，并append到左兄弟的DataNodes中
				for _, dataNode := range indexNode.DataNodes {
					dataNode.ParentNode = leftNode
					leftNode.DataNodes = append(leftNode.DataNodes, dataNode)
				}

				// 父节点中删除indexNode
				tempSlice := append([]*IndexNode{}, indexNode.ParentNode.Children[:indexAtParent]...)
				tempSlice = append(tempSlice, indexNode.ParentNode.Children[indexAtParent+1:]...)
				indexNode.ParentNode.Children = tempSlice

				// 如果父节点的keys.len >1
				if len(indexNode.ParentNode.Keys) > 1 {
					// 父节点中删除相应的key
					tempKeySlice := append([]string{}, indexNode.ParentNode.Keys[:indexAtParent-1]...)
					tempKeySlice = append(tempKeySlice, indexNode.ParentNode.Keys[indexAtParent:]...)
					indexNode.ParentNode.Keys = tempKeySlice
				} else { // 说明parentNode可能是root节点

					t.Root = leftNode
					leftNode.ParentNode = nil
				}

			}
		}
	} else {
		if len(indexNode.Keys) >= t.M/2 {
			return
		}
		fmt.Println("---合并索引节点-----")
		// 找到当前indexNode所在数组中的Index
		_, indexAtParent := binarySearchIndexNode(indexNode.ParentNode, indexNode.Keys[0])

		if indexAtParent == 0 { // 找右兄弟
			rightNode := indexNode.ParentNode.Children[indexAtParent+1]
			// 右兄弟有富余key
			if len(rightNode.Keys) > t.M/2 {
				borrowKey := rightNode.Keys[0]
				// 借的key加入到indexNode的keys中
				indexNode.Keys = append(indexNode.Keys, borrowKey)
				// 借的key对应的child 加入到indexNode的children中
				rightNode.Children[0].ParentNode = indexNode
				indexNode.Children = append(indexNode.Children, rightNode.Children[0])
				// 借的key替换父节点中对应的旧的key
				indexNode.ParentNode.Keys[0] = borrowKey
				// 兄弟节点中删除被借的key
				rightNode.Keys = rightNode.Keys[1:]
				// 兄弟节点中删除被借的child
				rightNode.Children = rightNode.Children[1:]
			} else { // 将右兄弟合并到indexNode中

				// 合并keys
				indexNode.Keys = append(indexNode.Keys, indexNode.ParentNode.Keys[0])
				indexNode.Keys = append(indexNode.Keys, rightNode.Keys[:]...)
				// 合并child
				for _, rightChild := range rightNode.Children {
					rightChild.ParentNode = indexNode
					indexNode.Children = append(indexNode.Children, rightChild)
				}

				if len(indexNode.ParentNode.Keys) < 2 { // 说明parentNode是root
					indexNode.ParentNode = nil
					t.Root = indexNode
				} else {
					// 父节点中删除对应key
					indexNode.ParentNode.Keys = indexNode.ParentNode.Keys[1:]
					// 父节点中删除右兄弟节点
					indexNode.ParentNode.Children = append([]*IndexNode{indexNode}, indexNode.ParentNode.Children[2:]...)
				}
			}
		} else { // 找左兄弟
			leftNode := indexNode.ParentNode.Children[indexAtParent-1]

			if len(leftNode.Keys) > t.M/2 {
				// 借的key
				borrowKey := leftNode.Keys[len(leftNode.Keys)-1]
				// 加入到indexNode
				indexNode.Keys = append([]string{borrowKey}, indexNode.Keys[:]...)
				// 借的key对应的child加入到indexNode
				borrowChild := leftNode.Children[len(leftNode.Children)-1]
				borrowChild.ParentNode = indexNode
				indexNode.Children = append([]*IndexNode{borrowChild}, indexNode.Children[:]...)
				// 在左兄弟中删除该key
				leftNode.Keys = leftNode.Keys[:len(leftNode.Keys)-1]
				// 在左兄弟中删除对应child
				leftNode.Children = leftNode.Children[:len(leftNode.Children)-1]
				// 借的Key替换父节点相应的key
				indexNode.ParentNode.Keys[indexAtParent-1] = borrowKey

			} else { // 如果左兄弟的dataNode <= t.M/2  则直接合并
				// 右边合到左边
				// 合并keys
				leftNode.Keys = append(leftNode.Keys, indexNode.ParentNode.Keys[indexAtParent-1])
				leftNode.Keys = append(leftNode.Keys, indexNode.Keys[:]...)
				// 修改child的父节点，并append到左兄弟的children中
				for _, child := range indexNode.Children {
					child.ParentNode = leftNode
					leftNode.Children = append(leftNode.Children, child)
				}

				if len(indexNode.ParentNode.Keys) > 1 {
					// 父节点中删除indexNode
					tempSlice := append([]*IndexNode{}, indexNode.ParentNode.Children[:indexAtParent]...)
					tempSlice = append(tempSlice, indexNode.ParentNode.Children[indexAtParent+1:]...)
					indexNode.ParentNode.Children = tempSlice
					// 父节点中删除相应的key
					tempKeySlice := append([]string{}, indexNode.ParentNode.Keys[:indexAtParent-1]...)
					tempKeySlice = append(tempKeySlice, indexNode.ParentNode.Keys[indexAtParent:]...)
					indexNode.ParentNode.Keys = tempKeySlice
				} else {
					leftNode.ParentNode = nil
					t.Root = leftNode
				}

			} // end
		} // end 找左兄弟
	}
	t.merge(indexNode.ParentNode)
}

// 树分裂
func (t *BPTree) divide(indexNode *IndexNode) {
	if indexNode == nil {
		return
	}
	if indexNode.IsLeaf && len(indexNode.DataNodes) > t.M {
		fmt.Printf("-----  分裂树:叶子节点 ------- %d  \n", len(indexNode.DataNodes))
		// 子节点个数是否大于阶数M 时，分裂树

		// 新生成一个右叶子节点
		newRightIndexNode := MallocNewIndexNode(true)
		// DataNodes右边部分，成为其数据节点
		newRightIndexNode.DataNodes = indexNode.DataNodes[t.M/2:] // [2:]
		// 遍历修改newRightIndexNode.DataNodes的父节点
		for _, dataNode := range newRightIndexNode.DataNodes {
			dataNode.ParentNode = newRightIndexNode
		}

		// 左边部分，成为DataNodes原父节点的子节点
		indexNode.DataNodes = indexNode.DataNodes[:t.M/2] // [0:2]

		// 此处需要优化
		if indexNode.ParentNode == nil { // 说明此节点是当前root节点
			newRootNode := MallocNewIndexNode(false)

			newRootNode.Keys = []string{newRightIndexNode.DataNodes[0].KeyAndValue.Key}
			fmt.Println("新root keys:", newRootNode.Keys)
			indexNode.ParentNode = newRootNode
			newRightIndexNode.ParentNode = newRootNode

			newRootNode.Children = []*IndexNode{indexNode, newRightIndexNode}
			t.Root = newRootNode // root节点指针上移
		} else {

			newRightIndexNode.ParentNode = indexNode.ParentNode

			previous, _ := binarySearchIndexKey(indexNode.ParentNode.Keys, newRightIndexNode.DataNodes[0].KeyAndValue.Key)

			// 合并Keys
			if previous < 0 {
				indexNode.ParentNode.Keys = append([]string{newRightIndexNode.DataNodes[0].KeyAndValue.Key}, indexNode.ParentNode.Keys[:]...)
			} else {
				tempSlice := append([]string{}, indexNode.ParentNode.Keys[:previous+1]...)
				tempSlice = append(tempSlice, newRightIndexNode.DataNodes[0].KeyAndValue.Key)
				tempSlice = append(tempSlice, indexNode.ParentNode.Keys[previous+1:]...)
				indexNode.ParentNode.Keys = tempSlice
			}
			// 合并child
			tempChildSlice := append([]*IndexNode{}, indexNode.ParentNode.Children[:previous+2]...)
			tempChildSlice = append(tempChildSlice, newRightIndexNode)
			tempChildSlice = append(tempChildSlice, indexNode.ParentNode.Children[previous+2:]...)
			indexNode.ParentNode.Children = tempChildSlice
		}

	} else if !indexNode.IsLeaf && len(indexNode.Keys) > t.M {
		fmt.Printf("-----  分裂树：索引节点 -------keyNum: %d ,keys: %s \n", len(indexNode.Keys), indexNode.Keys)
		// 新生成一个右索引节点
		newRightIndexNode := MallocNewIndexNode(false)
		// 原keys右边部分，成为其keys
		newRightIndexNode.Keys = indexNode.Keys[t.M/2+1:]
		// 原children右边部分，成为其children
		newRightIndexNode.Children = indexNode.Children[t.M/2+1:]
		// 遍历修改newRightIndexNode.Children的父节点
		for _, childNode := range newRightIndexNode.Children {
			childNode.ParentNode = newRightIndexNode
		}

		// TODO 此处需要优化
		if indexNode.ParentNode == nil { // 说明此节点是当前root节点
			newRootNode := MallocNewIndexNode(false)

			newRootNode.Keys = []string{indexNode.Keys[t.M/2]}
			indexNode.ParentNode = newRootNode
			newRightIndexNode.ParentNode = newRootNode

			newRootNode.Children = []*IndexNode{indexNode, newRightIndexNode}
			t.Root = newRootNode // root节点指针上移
		} else {

			newRightIndexNode.ParentNode = indexNode.ParentNode
			previous, _ := binarySearchIndexKey(indexNode.ParentNode.Keys, indexNode.Keys[0])

			// 合并Keys
			if previous < 0 {
				indexNode.ParentNode.Keys = append([]string{indexNode.Keys[t.M/2]}, indexNode.ParentNode.Keys[:]...)

			} else {
				// slice  左闭右开 [ )
				tempSlice := append([]string{}, indexNode.ParentNode.Keys[:previous+1]...)
				tempSlice = append(tempSlice, indexNode.Keys[t.M/2])
				tempSlice = append(tempSlice, indexNode.ParentNode.Keys[previous+1:]...)
				indexNode.ParentNode.Keys = tempSlice
			}
			// 合并child
			tempChildSlice := append([]*IndexNode{}, indexNode.ParentNode.Children[:previous+2]...)
			tempChildSlice = append(tempChildSlice, newRightIndexNode)
			tempChildSlice = append(tempChildSlice, indexNode.ParentNode.Children[previous+2:]...)
			indexNode.ParentNode.Children = tempChildSlice
		}
		// 左边部分，成为Children原父节点的子节点
		indexNode.Children = indexNode.Children[:t.M/2+1]
		indexNode.Keys = indexNode.Keys[:t.M/2]

	}
	t.divide(indexNode.ParentNode)
}

func (t *BPTree) Traversal() {
	p := t.Head
	// 遍历
	for p != nil {
		if p.ParentNode.ParentNode == nil {
			fmt.Printf("key %s: value %s  \n", p.KeyAndValue.Key, p.KeyAndValue.Value)
		} else {

			fmt.Printf("key %s: value %s , parent keys:%s \n", p.KeyAndValue.Key, p.KeyAndValue.Value, p.ParentNode.ParentNode.Keys)
		}
		p = p.NextDataNode
	}
	fmt.Println()
}

func (t *BPTree) UpToDownPrint() {
	p := t.Root

	if p != nil {
		if p.IsLeaf {
			// 打印 DataNode
			for _, dataNode := range p.DataNodes {
				fmt.Printf("%s ", dataNode.KeyAndValue.Key)
			}
			fmt.Println()
		} else {
			fmt.Println(p.Keys)
			var tempArray []*IndexNode
			// 打印child的
			for _, child := range p.Children {
				if child.IsLeaf {
					for _, dataNode := range child.DataNodes {
						fmt.Printf("%s ", dataNode.KeyAndValue.Key)
					}
					fmt.Print("|")
				} else {
					tempArray = append(tempArray, child)
				}
			}

			for len(tempArray) > 0 {
				var newTempArray []*IndexNode
				for _, node := range tempArray {
					if node.IsLeaf {
						for _, dataNode := range node.DataNodes {
							fmt.Printf("%s ", dataNode.KeyAndValue.Key)
						}
						fmt.Print("|")
					} else {
						fmt.Print(node.Keys)
						for _, newChild := range node.Children {
							newTempArray = append(newTempArray, newChild)
						}

					}
				}
				tempArray = newTempArray
				fmt.Println()
			}

		}
	}
}
func (t *BPTree) FindLeft() *IndexNode {
	if len(t.Root.Children) > 0 {
		p := t.Root.Children[0]
		for len(p.Children) > 0 {
			p = p.Children[0]
		}
		return p
	}
	return t.Root
}

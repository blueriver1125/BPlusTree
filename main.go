package main

import (
	"strconv"
	"fmt"
	"github.com/mycode/BPlusTree/bptree"
	"math/rand"
	"time"
)

const M = 4

func main() {

	t := bptree.MallocNewBPTree(M)

	keyArray := []int{55, 34, 15, 95, 99, 98, 81, 16, 99, 14, 36, 13, 77, 57, 37, 2, 39, 3, 89, 76}
	//for _, key := range keyArray {

	for n := 0; n < 100; n++ {
		rand.Seed(time.Now().UnixNano())
		key := rand.Intn(100)
		keystr := strconv.Itoa(key)
		keyAndValue := bptree.KeyAndValue{
			"k" + keystr, "v" + keystr}
		//fmt.Printf("开始插入： key:%s  \n\n", keyAndValue.Key)
		//
		t.Insert(keyAndValue)
		t.UpToDownPrint()
		//bpTree.Traversal()
		fmt.Println()
		fmt.Println()
	}
	// 修改
	fmt.Println("---------update----------")
	updateKV1 := bptree.KeyAndValue{
		" ", "v36修改值"}
	updateKV2 := bptree.KeyAndValue{
		"k13", "v13修改值"}

	updateKV3 := bptree.KeyAndValue{
		"k15", "v15修改值"}

	updateKV4 := bptree.KeyAndValue{
		"k39", "v39修改值"}
	updateKV5 := bptree.KeyAndValue{
		"k81", "v81修改值"}

	updateKV6 := bptree.KeyAndValue{
		"k95", "v95修改值"}

	updateKV7 := bptree.KeyAndValue{
		"k55", "v55修改值"}
	updateKV8 := bptree.KeyAndValue{
		"k99", "v99修改值"}

	_, err := t.Update(updateKV1)
	if err != nil {
		fmt.Println(err)
	}
	t.Update(updateKV2)
	t.Update(updateKV3)
	t.Update(updateKV4)
	t.Update(updateKV5)
	t.Update(updateKV6)
	t.Update(updateKV7)
	t.Update(updateKV8)
	t.UpToDownPrint()

	fmt.Println()
	t.Get("k99")
	t.Get("k13")
	t.Get("k15")
	t.Get("k39")
	t.Get("k81")
	t.Get("k95")
	t.Get("k55")
	t.Get("k99")

	t.Remove("k16")
	t.UpToDownPrint()

	t.Remove("k2")
	t.UpToDownPrint()

	t.Remove("k14")
	t.UpToDownPrint()

	t.Remove("k89")
	t.UpToDownPrint()

	t.Get("k99")
	t.Get("k13")
	t.Get("k15")
	t.Get("k39")
	t.Get("k81")
	t.Get("k95")
	t.Get("k55")
	t.Get("k99")

	t.Remove("k55")
	t.UpToDownPrint()

	t.Remove("k77")
	t.UpToDownPrint()

	t.Remove("k99")
	t.UpToDownPrint()

	t.Remove("k98")
	t.UpToDownPrint()

	t.Remove("k95")
	t.UpToDownPrint()

	t.Remove("k99")
	t.UpToDownPrint()

	t.Get("k99")
	t.Get("k13")
	t.Get("k15")
	t.Get("k39")
	t.Get("k81")
	t.Get("k95")
	t.Get("k55")
	t.Get("k99")

	t.Remove("k34")
	t.UpToDownPrint()

	t.Remove("k39")
	t.UpToDownPrint()

	t.Remove("k81")
	t.UpToDownPrint()

	t.Remove("k3")
	t.UpToDownPrint()

	t.Remove("k13")
	t.UpToDownPrint()

	t.Remove("k36")
	t.UpToDownPrint()

	t.Remove("k57")
	t.UpToDownPrint()

	t.Remove("k37")
	t.UpToDownPrint()

	t.Get("k15")
	t.Get("k76")
	t.Get("k99")

	t.Remove("k15")
	t.UpToDownPrint()

	t.Remove("k76")
	t.UpToDownPrint()

	for _, key := range keyArray {

		//for n := 0; n < 100; n++ {
		//	rand.Seed(time.Now().UnixNano())
		//	key := rand.Intn(100)
		keystr := strconv.Itoa(key)
		keyAndValue := bptree.KeyAndValue{
			"k" + keystr, "v" + keystr}
		//fmt.Printf("开始插入： key:%s  \n\n", keyAndValue.Key)
		//
		t.Insert(keyAndValue)
		t.UpToDownPrint()
		//bpTree.Traversal()
		fmt.Println()
		fmt.Println()
	}
}

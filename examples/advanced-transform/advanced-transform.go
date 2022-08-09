package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

func AdvancedTransformUrlEncode(key string) *diskv.PathKey {
	return &diskv.PathKey{
		FileName: url.PathEscape(key),
	}
}

func InverseTransformUrlDecode(pathKey *diskv.PathKey) (key string) {
	k, _ := url.PathUnescape(pathKey.FileName)
	return k
}

var d *diskv.Diskv

const CachePath = "my-data-dir"

func init() {
	d = diskv.New(diskv.Options{
		BasePath:          CachePath,
		AdvancedTransform: AdvancedTransformUrlEncode,
		InverseTransform:  InverseTransformUrlDecode,
		CacheSizeMax:      1024 * 1024,
		FileTTLMax:        1,
	})
}

func main() {

	// Write some text to the key "alpha/beta/gamma".
	key := "2766848225710557_/data/user/0/com.linzihy.mgd"
	err := d.WriteString(key, "¡Hola!") // will be stored in "<basedir>/alpha/beta/gamma.txt"
	print(url.PathEscape(key))
	if err != nil {
		panic(err)
	}
	fmt.Println(d.ReadString(key))

	log.Println(time.Now())
	for i := 0; i < 5e1; i++ {
		k := strconv.Itoa(i)
		v := randStr(1024 * 1024)
		err := WriteRead1M(k, v)
		if err != nil {
			panic(err)
		}
	}
	log.Println(time.Now())
	time.Sleep(time.Second * 2)

	//d.EraseAll()
	var c <-chan struct{}
	d.Erase("*")
	cdata := d.KeysPrefix("1", c)
	for cdatum := range cdata {
		log.Println(time.Now(), cdatum, len(d.ReadString(cdatum)))
		//log.Println(len(d.ReadString(cdatum)), d.ReadString(cdatum))
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			sb.WriteByte(letters[idx])
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return sb.String()
}

var s = randStr(1024 * 1024)

func WriteRead1M(k, v string) error {
	//s := randStr(1024 * 1024)
	if err := d.WriteString(k, v); err != nil {
		return err
	}
	d.ReadString(k)
	return nil
}

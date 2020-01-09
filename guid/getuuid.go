package guid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"crypto/aes"
	"crypto/cipher"
	"bbTool/n_log"
	"bytes"
	"strings"
	"fmt"
	"time"
)

func GetId() string {
	len := 16
	data := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, data)
	if err != nil {
		return ""
	}
	return func() string {
		s1 := base64.URLEncoding.EncodeToString(data)
		i := md5.New()
		i.Write([]byte(s1))
		return hex.EncodeToString(i.Sum(nil))
	}()
	//return GetMd5(string(data))
}

func GetId_timer() string  {
	curTime := time.Now()
	t := "20060102150405"
	//n_log.Info("len  %v",len(t))
	//
	//
	str := curTime.Format(t)+fmt.Sprintf("%v%v",curTime.Nanosecond(),GetRandInt32())
	//u1 := fmt.Sprintf("%v",curTime.UnixNano())
	//u2 := fmt.Sprintf("%v",GetRandInt32())
	//n_log.Info("u1   %v   %v",u1,len(u1))
	//n_log.Info("u2   %v   %v",u2,len(u2))
	//
	//
	//str := fmt.Sprintf("%v%v",u1,u2)
	return str
}


func GetMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)

}


var (
	Commonkey = "826A3D0A97F62D9E"
	//syncMutex sync.Mutex
)


func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	res, err := base64.URLEncoding.DecodeString(data)
	fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)
}

func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

func AesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err is:", err)
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	fmt.Println("source is :", origData, string(origData))
	return origData
}

func AesEncrypt(src, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		n_log.Info("key error1", err)
	}
	if src == "" {
		n_log.Info("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// 普通base64编码加密 区别于urlsafe base64
	ret := base64.StdEncoding.EncodeToString(crypted)
//	n_log.Info("base64 result:", ret)

	//fmt.Println("base64UrlSafe result:", Base64UrlSafeEncode(crypted))
	//return crypted
	return ret
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		n_log.Panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		n_log.Panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		n_log.Panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		n_log.Panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
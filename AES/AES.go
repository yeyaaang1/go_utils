package AES

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"gitee.com/super_step/go_utils/rand"
	"github.com/kataras/golog"
)

/*CBC加密 按照golang标准库的例子代码
不过里面没有填充的部分,所以补上
*/

type AES interface {
	DesaltDecrypt(encrypt string) string
	SaltyEncrypt(password string) string
}

type myAes struct {
	key string
}

func NewAES(key string) AES {
	return &myAes{
		key: key,
	}
}

// 使用PKCS7进行填充，IOS也是7
func (AES *myAes) pCks7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (AES *myAes) pCks7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func (AES *myAes) myDecrypt(rawData string) (result []byte, err error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return
	}
	// 处理传入的key
	AESKey, err := base64.StdEncoding.DecodeString(AES.key + "=")
	if err != nil {
		return
	}
	// 开始解密过程
	return AES.myAesCBCDecrypt(data, AESKey)
}

func (AES *myAes) DesaltDecrypt(Encrypt string) string {
	defer func() {
		if err := recover(); err != nil {
			golog.Default.Error("去盐解码出错", err)
		}
	}()
	result, err := AES.myDecrypt(Encrypt)
	if err != nil {
		return ""
	}
	resultStr := string(result)
	if len(resultStr) > 16 {
		return resultStr[16:]
	}
	return ""
}

func (AES *myAes) SaltyEncrypt(word string) string {
	defer func() {
		if err := recover(); err != nil {
			golog.Default.Error("加盐编码出错", err)
		}
	}()
	word = rand.RandRunes(16) + word
	result, err := AES.myEncrypt([]byte(word))
	if err != nil {
		return ""
	}
	return result
}

func (AES *myAes) myEncrypt(rawData []byte) (string, error) {
	// 处理传入的key
	AESKey, err := base64.StdEncoding.DecodeString(AES.key + "=")
	if err != nil {
		return "", err
	}
	data, err := AES.myAesCBCEncrypt(rawData, AESKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (AES *myAes) myAesCBCDecrypt(encryptData, key []byte) (encryptResult []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		err = errors.New("ciphertext too short")
		return
	}

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}
	// 设置iv初始向量为key的前16字节
	iv := key[:blockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptResult = AES.pCks7UnPadding(encryptData)
	return
}

// aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func (AES *myAes) myAesCBCEncrypt(rawData, key []byte) (result []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	// 填充原文
	blockSize := block.BlockSize()
	rawData = AES.pCks7Padding(rawData, blockSize)
	// 初始向量IV必须是唯一，但不需要保密
	result = make([]byte, len(rawData))
	// block大小 16
	iv := key[:blockSize]

	// block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(result, rawData)
	return
}

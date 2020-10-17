package RSA

import (
	"crypto/rand"
	RSA "crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"gitee.com/super_step/go_utils/myError"
	"io/ioutil"
)

type Service interface {
	// 加密
	Encrypt(origData []byte) (encrypt []byte, err error)
	// 解密
	Decrypt(ciphertext []byte) (decrypt []byte, err error)
	// 加密文本
	EncryptStr(origData string) (encrypt string, err error)
	// 解密文本
	DecryptStr(ciphertext string) (decrypt string, err error)
}

func NewRsa(public, private string) (srv Service, err error) {
	var (
		privateBytes, publicBytes []byte
	)
	if private == "" && public == "" {
		err = myError.New("公钥和私钥至少传入一个")
		return
	}
	if private != "" {
		privateBytes, err = ioutil.ReadFile(private)
		if err != nil {
			err = myError.Warp(err, "读取pem文件出错")
			return
		}
	}
	if public != "" {
		publicBytes, err = ioutil.ReadFile(public)
		if err != nil {
			err = myError.Warp(err, "读取pem文件出错")
			return
		}
	}
	srv = &rsa{
		public:  publicBytes,
		private: privateBytes,
	}
	return
}

type rsa struct {
	public  []byte
	private []byte
}

// 加密
func (r *rsa) Encrypt(origData []byte) (encrypt []byte, err error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(r.public)
	if block == nil {
		err = myError.New("公钥解密出错")
		return
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		err = myError.Warp(err, "公钥解析出错")
		return
	}
	// 类型断言
	pub, ok := pubInterface.(*RSA.PublicKey)
	if !ok {
		err = myError.New("公钥类型异常")
		return
	}
	// 加密
	encrypt, err = RSA.EncryptPKCS1v15(rand.Reader, pub, origData)
	if err != nil {
		err = myError.Warp(err, "加密出错")
		return
	}
	return
}

// 解密
func (r *rsa) Decrypt(ciphertext []byte) (decrypt []byte, err error) {
	// 解密pem格式的私钥
	block, _ := pem.Decode(r.private)
	if block == nil {
		err = myError.New("私钥解密出错")
		return
	}
	// 解析PKCS1格式的私钥
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = myError.Warp(err, "私钥解析出错")
		return
	}
	// 解密
	decrypt, err = RSA.DecryptPKCS1v15(rand.Reader, private, ciphertext)
	if err != nil {
		err = myError.Warp(err, "解密出错")
		return
	}
	return
}

func (r *rsa) DecryptStr(ciphertext string) (decrypt string, err error) {
	var decryptByte, base64Bytes []byte
	base64Bytes, err = base64.StdEncoding.DecodeString(ciphertext)
	decryptByte, err = r.Decrypt(base64Bytes)
	if err != nil {
		return
	}
	decrypt = string(decryptByte)
	return
}

func (r *rsa) EncryptStr(origData string) (encrypt string, err error) {
	var encryptByte []byte
	encryptByte, err = r.Encrypt([]byte(origData))
	if err != nil {
		return
	}
	encrypt = base64.StdEncoding.EncodeToString(encryptByte)
	return
}

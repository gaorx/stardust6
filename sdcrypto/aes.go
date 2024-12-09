package sdcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/gaorx/stardust6/sderr"
)

var (
	// AES AES加密解密器
	AES Encrypter = &EncrypterFunc{
		Encrypter: AESEncrypt,
		Decrypter: AESDecrypt,
	}
	// AESCRC32 AES加密解密器，带CRC32校验
	AESCRC32 Encrypter = &CRC32Encrypter{AES}
)

// AESEncrypt 加密数据
func AESEncrypt(key, data []byte) ([]byte, error) {
	return AESEncryptPadding(key, data, Pkcs5)
}

// AESDecrypt 解密数据
func AESDecrypt(key, encrypted []byte) ([]byte, error) {
	return AESDecryptPadding(key, encrypted, UnPkcs5)
}

// AESEncryptPadding 加密数据，带有数据填充
func AESEncryptPadding(key, data []byte, p Padding) ([]byte, error) {
	if p == nil {
		return nil, sderr.Newf("AES nil padding")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, sderr.Wrapf(err, "AES create cipher error")
	}
	data, err = p(data, block.BlockSize())
	if err != nil {
		return nil, sderr.Wrapf(err, "AES padding error")
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	encrypted := make([]byte, len(data))
	encrypter.CryptBlocks(encrypted, data)
	return encrypted, nil
}

// AESDecryptPadding 解密数据，带有数据填充
func AESDecryptPadding(key, encrypted []byte, p Unpadding) ([]byte, error) {
	if p == nil {
		return nil, sderr.Newf("DeAES nil unpadding")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, sderr.Wrapf(err, "DeAES create cipher error")
	}
	decrypter := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	data := make([]byte, len(encrypted))
	decrypter.CryptBlocks(data, encrypted)
	r, err := p(data, block.BlockSize())
	if err != nil {
		return nil, sderr.Wrapf(err, "DeAES unpadding error")
	}
	return r, nil
}

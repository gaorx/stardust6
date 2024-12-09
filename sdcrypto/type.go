package sdcrypto

// Encrypter 加密解密接口
type Encrypter interface {
	// Encrypt 加密
	Encrypt(key, data []byte) ([]byte, error)
	// Decrypt 解密
	Decrypt(key, crypted []byte) ([]byte, error)
}

// EncrypterFunc 函数式加密解密
type EncrypterFunc struct {
	Encrypter func(key, data []byte) ([]byte, error)
	Decrypter func(key, crypted []byte) ([]byte, error)
}

// Encrypt 实现 Encrypter.Encrypt
func (e *EncrypterFunc) Encrypt(key, data []byte) ([]byte, error) {
	return e.Encrypter(key, data)
}

// Decrypt 实现 Encrypter.Decrypt
func (e *EncrypterFunc) Decrypt(key, crypted []byte) ([]byte, error) {
	return e.Decrypter(key, crypted)
}

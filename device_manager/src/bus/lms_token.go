package bus

import(
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

//SecretKey 作为秘钥，使用 HmacSHA1 方法对上面的待签名字符串做签名计算得到一个二进制数组，最后对该二进制数组做 Base64 编码得到最终的 password 签名字符串，即 “eqweq+adwe23fssf”。


func GenToken(input string, key string) string {
	/*
		签名采用HmacSHA1算法 + Base64，编码采用：UTF-8
		MAC算法结合了MD5和SHA算法的优势，并加入密钥的支持，是一种更为安全的消息摘要算法。
	*/

	key_for_sign := []byte(key)
	h := hmac.New(sha1.New, key_for_sign)
	h.Write([]byte(input))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
package colossus

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func generateBase64MD5ForContent(content []byte) (string, error) {
	if len(content) == 0 {
		return "", nil
	}
	hash := md5.New()
	_, err := hash.Write(content)
	if err != nil {
		return "", err
	}
	base64Hash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return base64Hash, nil
}

func generateCanonicalString(reqMethod string, contentType string, contentMD5 string, reqURI string, datetime string) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", reqMethod, contentType, contentMD5, reqURI, datetime)
}

func generateSignature(key []byte, secret []byte, canonicalString []byte) (string, error) {
	hmacSHA1 := hmac.New(sha256.New, secret)
	_, err := hmacSHA1.Write(canonicalString)
	if err != nil {
		return "", err
	}
	encodedHash := base64.StdEncoding.EncodeToString(hmacSHA1.Sum(nil))
	return fmt.Sprintf("APIAuth-HMAC-SHA256 %s:%s", key, encodedHash), nil
}

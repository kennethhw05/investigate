package colossus

import (
	"net/http"
	"testing"
)

func TestGenerateBase64MD5ForNoContent(t *testing.T) {
	t.Parallel()

	hash, err := generateBase64MD5ForContent([]byte(""))
	if err != nil {
		t.Fatalf("Unexpected error during test %s", t.Name())
	}
	expectedHash := ""
	if hash != expectedHash {
		t.Errorf("Expected content hash %s but got %s", expectedHash, hash)
	}
}

func TestGenerateBase64MD5ForContent(t *testing.T) {
	t.Parallel()

	hash, err := generateBase64MD5ForContent([]byte("Give me my md5"))
	if err != nil {
		t.Fatalf("Unexpected error during test %s", t.Name())
	}
	expectedHash := "zmvtKAnOCEdWkr2sb8SCIQ=="
	if hash != expectedHash {
		t.Errorf("Expected content hash %s but got %s", expectedHash, hash)
	}
}

func TestGenerateCanonicalString(t *testing.T) {
	t.Parallel()

	md5 := "rY5GyingAbnOReyXXtSqyA=="
	canonicalStr := generateCanonicalString(http.MethodPost, "application/json", md5, "/tickets", "Tue, 11 Apr 2017 16:05:42 GMT")
	expectedCanonicalStr := "POST,application/json,rY5GyingAbnOReyXXtSqyA==,/tickets,Tue, 11 Apr 2017 16:05:42 GMT"
	if canonicalStr != expectedCanonicalStr {
		t.Errorf("Expected signature %s but got %s", expectedCanonicalStr, canonicalStr)
	}
}

func TestGenerateCanonicalStringNoContent(t *testing.T) {
	t.Parallel()

	canonicalStr := generateCanonicalString(http.MethodGet, "application/json", "", "/test_auth", "Mon, 17 Sep 2018 18:30:17 GMT")
	expectedCanonicalStr := "GET,application/json,,/test_auth,Mon, 17 Sep 2018 18:30:17 GMT"
	if canonicalStr != expectedCanonicalStr {
		t.Errorf("Expected signature %s but got %s", expectedCanonicalStr, canonicalStr)
	}
}

func TestGenerateSignature(t *testing.T) {
	t.Parallel()

	actualSig, err := generateSignature([]byte("api_key"), []byte("api_secret"), []byte("POST,application/json,rY5GyingAbnOReyXXtSqyA==,/tickets,Tue, 11 Apr 2017 16:05:42 GMT"))
	if err != nil {
		t.Fatalf("Unexpected error during test %s", t.Name())
	}
	expectedSig := "APIAuth-HMAC-SHA256 api_key:J1zElRwTpRzsxttFtJTz+7YPd5OiZ17oVk5i3Vrmxk0="
	if actualSig != expectedSig {
		t.Errorf("Expected signature %s but got %s", expectedSig, actualSig)
	}
}

func TestGenerateSignatureNoContent(t *testing.T) {
	t.Parallel()

	actualSig, err := generateSignature([]byte("api_key"), []byte("api_secret"), []byte("GET,,1B2M2Y8AsgTpgAmY7PhCfg==,test_auth,Mon, 17 Sep 2018 18:30:17 GMT"))
	if err != nil {
		t.Fatalf("Unexpected error during test %s", t.Name())
	}
	expectedSig := "APIAuth-HMAC-SHA256 api_key:L7SU0BLV0a30hAdMdbtliVMrCCqa5OVG4QvsTEH4Uq4="
	if actualSig != expectedSig {
		t.Errorf("Expected signature %s but got %s", expectedSig, actualSig)
	}
}

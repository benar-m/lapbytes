package api

import "testing"

func TestGenerateRandomString(t *testing.T) {
	str, err := generateRandomString()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if str == "" {
		t.Fatalf("Got an Empty string")
	}
}

func TestHashPassword(t *testing.T) {
	testPassword := "strongpassword"
	hash, err := hashPassword(testPassword)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatalf("The hash is empty")
	}

	if testPassword == hash {
		t.Fatal("Hash cant be equal to original password")
	}

	//same passwords shouldnt equal hashes because of salt
	hash1, _ := hashPassword("strongpassword")
	if hash1 == hash {
		t.Fatal("Same hashes from the same password!")
	}

}

func TestVerifyHash(t *testing.T) {
	testPassword := "superstrongpasswordtrustme"
	hash, _ := hashPassword(testPassword)
	err := verifyPasswordHash(testPassword, hash)
	if !err {

		t.Fatal("Hash verification not working")
	}
}

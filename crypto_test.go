package main

import "testing"

func TestDecrypt(t *testing.T) {
	e, _ := encrypt([]byte("test"), []byte("1234"))
	a, _ := decrypt(e, []byte("1234"))
	if string(a) != "test" {
		t.Errorf(string(a))
	}
}

func TestCompression(t *testing.T) {
	test := `package main

  import "testing"

  func TestDecrypt(t *testing.T) {
  	e, _ := encrypt([]byte("test"), []byte("1234"))
  	a, _ := decrypt(e, []byte("1234"))
  	if string(a) != "test" {
  		t.Errorf(string(a))
  	}
  }
`
	ENCRYPTION_COMPRESSION = true
	s1, _ := encryptString(test, "1234")
	ENCRYPTION_COMPRESSION = false
	s2, _ := encryptString(test, "1234")
	if float64(len(s1))/float64(len(s2)) > 0.75 {
		t.Errorf("Weird compression ratio: %2.2f", float64(len(s1))/float64(len(s2)))
	}
}

func TestEncryptWriting(t *testing.T) {
	encryptAndWrite("test.test.test", "test", "1234")
	s, _ := openAndDecrypt("test.test.test", "1234")
	Shred("test.test.test")
	if exists("test.test.test") {
		t.Errorf("Problem with shredding")
	}
	if s != "test" {
		t.Errorf(s)
	}
}

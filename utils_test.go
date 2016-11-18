package edvcfs

import (
	"testing"
	"time"
)

func TestExists(t *testing.T) {
	if !exists("README.md") {
		t.Errorf("Doesn't exist?")
	}
}

func TestHashAndHex(t *testing.T) {
	s := hashAndHex("test")
	if s != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf(s)
	}
}

func TestStrExtract(t *testing.T) {
	text1 := `!!!!!some text<<<<<<`
	extracted := StrExtract(text1, "!!!!!", "<<<<<<", 1)

	if extracted != "some text" {
		t.Errorf("Incorrect extracted: %s", extracted)
	}
}

func TestGetRandomMD5Hash(t *testing.T) {
	a := GetRandomMD5Hash()
	time.Sleep(1 * time.Millisecond)
	b := GetRandomMD5Hash()
	if a == b {
		t.Errorf("MD5 hashes not random: %s, %s", a, b)
	}
	if len(GetRandomMD5Hash()) < 8 {
		t.Errorf("MD5 hashes are longer than expected")
	}
}

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	a := RandStringBytesMaskImprSrc(10, 1)
	b := RandStringBytesMaskImprSrc(10, 1)
	if a != b {
		t.Errorf("Seeding not working")
	}
	if len(b) != 10 {
		t.Errorf("Wrong size")
	}
}

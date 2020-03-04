package indyverify

import "testing"

func TestEmptyInput(t *testing.T) {
	var a, b, c []byte

	emptyResult, err := Indyverify(a, b, c)
	if emptyResult != false {
		t.Errorf("failed, expected false, but got %v and error is %v", emptyResult, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", emptyResult, err)
	}
}

func TestInvalidDid(t *testing.T) {
	var pb, did, signature []byte
	pb = []byte("input value")
	did = []byte("123456789012345") //modify the input
	signature = []byte("test input")

	result, err := Indyverify(pb, did, signature)
	if result != false {
		t.Errorf("failed, expected false, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", result, err)
	}
}

func TestInvalidSignature(t *testing.T) {
	var pb, did, signature []byte
	pb = []byte("input value")
	did = []byte("8KHMLmGrxuy1yJ2r7eM3xW")
	signature = []byte("wqBOesOfW8OjN8KiaBh8wrvDmsKvwrsQw5NNwpfCl8OVwo3DtMODfMOuw53CjhjCqQJYxaERQ1fDkMOSDcOZw4Nsw4HDrWR7w5jDoMKLwq3CocKpVMKpw4TDsDHCtRB/w4bDoMKtAQ==") //modify the input

	result, err := Indyverify(pb, did, signature)
	if result != false {
		t.Errorf("failed, expected false, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", result, err)
	}
}

func TestInvalidPB(t *testing.T) {
	var pb, did, signature []byte
	pb = []byte("input value") //modify the input
	did = []byte("8KHMLmGrxuy1yJ2r7eM3xW")
	signature = []byte("test input")

	result, err := Indyverify(pb, did, signature)
	if result != false {
		t.Errorf("failed, expected false, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", result, err)
	}
}
func TestSuccess(t *testing.T) {
	var pb, did, signature []byte
	pb = []byte("test") //modify the input
	did = []byte("8KHMLmGrxuy1yJ2r7eM3xW")
	signature = []byte("xb0Be2DCmUDDmxkxMAHDtsOOwoc1xaBBwpvChMOgwrXCsUY9ccOdYH07MhkiayErLMOiwolTH8Onw4bDm8OEw6AUcUoIwqnDrcO5M8Ohw5YGxZLCjQZSwp1ewrYI")

	result, err := Indyverify(pb, did, signature)
	if result != true {
		t.Errorf("failed, expected true, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected true and got %true and error is %v", result, err)
	}
}

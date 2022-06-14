package link

import "testing"

func TestLinkHandler(t *testing.T) {
	url := "https://sendsafely.dremio.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE#keyCode=MYKEYCODE"
	linkParks, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParks.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParks.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParks.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParks.PackageCode)
	}

	expectedThread := "MYTHREAD"
	if linkParks.Thread != expectedThread {
		t.Errorf("expected thread '%v' but got '%v'", expectedThread, linkParks.Thread)
	}
}

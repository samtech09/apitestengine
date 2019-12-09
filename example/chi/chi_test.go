package chirest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	apitest "github.com/samtech09/apitestengine"
)

var chiTest *apitest.APITest
var chiTestCases map[string]apitest.TestCase
var chits *httptest.Server

func initChiTest() error {
	// prepare test routes and server
	r := SetupRouter()
	chits = httptest.NewServer(r)

	// prepare test cases
	chiTestCases = make(map[string]apitest.TestCase)
	chiTestCases["TestChiIndex"] = apitest.NewTestCase("TestChiIndex", "GET", "/", "Bare minimum API server in go with chi router", "", nil)
	chiTestCases["TestChiParam1"] = apitest.NewTestCase("TestChiParam", "GET", "/param/chitest", `"key = chitest"`+"\n", "", nil)
	chiTestCases["TestChiParam2"] = apitest.NewTestCase("TestChiParam2", "GET", "/param/chi-chi-test", `"key = chi-chi-test"`+"\n", "", nil)

	expected := `{"Age":16,"Name":"Mohan"}` + "\n"
	chiTestCases["TestChiJSON"] = apitest.NewTestCase("TestChiJSON", "GET", "/getjson", expected, "", nil)

	// post payload
	p := Person{}
	p.Age = 6
	p.Name = "Krishna"

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Failed marshling payload")
	}
	chiTestCases["TestChiPostJSON"] = apitest.NewTestCase("TestChiPostJSON", "POST", "/postjson", `"OK"`+"\n", "", bytes.NewBuffer(payload))

	// initialize test engine
	chiTest = apitest.NewAPITest(chits)
	return nil
}

func TestChiInit(t *testing.T) {
	err := initChiTest()
	if err != nil {
		t.Fatal(err)
	}
}

func TestChiIndex(t *testing.T) {
	ret := chiTest.DoTest(chiTestCases["TestChiIndex"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestChiParam(t *testing.T) {
	tcs := []apitest.TestCase{}
	tcs = append(tcs, chiTestCases["TestChiParam1"])
	tcs = append(tcs, chiTestCases["TestChiParam2"])

	tresult := chiTest.DoTests(tcs)
	for _, ret := range tresult {
		if ret.Err != nil {
			fmt.Println(ret.ID)
			t.Error(ret.Err)
		} else {
			fmt.Println("PASS: ", ret.ID)
		}
	}
}

func TestChiJSON(t *testing.T) {
	ret := chiTest.DoTest(chiTestCases["TestChiJSON"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestChiPostJSON(t *testing.T) {
	ret := chiTest.DoTest(chiTestCases["TestChiPostJSON"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

//----
//
//
//----

func BenchmarkChiInit(b *testing.B) {
	err := initChiTest()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkChiIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := chiTest.DoTest(chiTestCases["TestChiParam1"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkParam(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := chiTest.DoTest(chiTestCases["TestChiIndex"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkChiJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := chiTest.DoTest(chiTestCases["TestChiJSON"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkChiPostJSON(b *testing.B) {
	// post payload
	p := Person{}
	p.Age = 6
	p.Name = "Krishna"

	payload, err := json.Marshal(p)
	if err != nil {
		b.Fatal("Failed marshling payload")
	}

	for n := 0; n < b.N; n++ {
		ret := chiTest.DoTest(apitest.NewTestCase("TestChiPostJSON", "POST", "/postjson", `"OK"`+"\n", "", bytes.NewBuffer(payload)))
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func TestChiCloser(t *testing.T) {
	chits.Close()
}

func BenchmarkChiCloser(b *testing.B) {
	chits.Close()
}

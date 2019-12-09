package ginrest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	apitest "github.com/samtech09/apitestengine"
)

var ginTest *apitest.APITest
var ginTestCases map[string]apitest.TestCase
var gints *httptest.Server

func initGinTest() error {
	// prepare test routes and server
	r := SetupRouter()
	gints = httptest.NewServer(r)

	// prepare test cases
	ginTestCases = make(map[string]apitest.TestCase)
	ginTestCases["TestGinIndex"] = apitest.NewTestCase("TestGinIndex", "GET", "/", "Bare minimum API server in go with gin router", "", nil)
	ginTestCases["TestGinParam1"] = apitest.NewTestCase("TestGinParam", "GET", "/param/ginTest", `"key = ginTest"`+"\n", "", nil)
	ginTestCases["TestGinParam2"] = apitest.NewTestCase("TestGinParam2", "GET", "/param/gin-gin-test", `"key = gin-gin-test"`+"\n", "", nil)

	expected := `{"Age":16,"Name":"Mohan"}` + "\n"
	ginTestCases["TestGinJSON"] = apitest.NewTestCase("TestGinJSON", "GET", "/getjson", expected, "", nil)

	// post payload
	p := Person{}
	p.Age = 6
	p.Name = "Krishna"

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Failed marshling payload")
	}
	ginTestCases["TestGinPostJSON"] = apitest.NewTestCase("TestGinPostJSON", "POST", "/postjson", `"OK"`+"\n", "", bytes.NewBuffer(payload))

	// initialize test engine
	ginTest = apitest.NewAPITest(gints)
	return nil
}

func TestGinInit(t *testing.T) {
	err := initGinTest()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGinIndex(t *testing.T) {
	ret := ginTest.DoTest(ginTestCases["TestGinIndex"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestGinParam(t *testing.T) {
	tcs := []apitest.TestCase{}
	tcs = append(tcs, ginTestCases["TestGinParam1"])
	tcs = append(tcs, ginTestCases["TestGinParam2"])

	tresult := ginTest.DoTests(tcs)
	for _, ret := range tresult {
		if ret.Err != nil {
			fmt.Println(ret.ID)
			t.Error(ret.Err)
		} else {
			fmt.Println("PASS: ", ret.ID)
		}
	}
}

func TestGinJSON(t *testing.T) {
	ret := ginTest.DoTest(ginTestCases["TestGinJSON"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestGinPostJSON(t *testing.T) {
	ret := ginTest.DoTest(ginTestCases["TestGinPostJSON"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

//----
//
//
//----

func BenchmarkGinInit(b *testing.B) {
	err := initGinTest()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkGinIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := ginTest.DoTest(ginTestCases["TestGinParam1"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkGinParam(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := ginTest.DoTest(ginTestCases["TestGinIndex"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkGinJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ret := ginTest.DoTest(ginTestCases["TestGinJSON"])
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func BenchmarkGinPostJSON(b *testing.B) {
	// post payload
	p := Person{}
	p.Age = 6
	p.Name = "Krishna"

	payload, err := json.Marshal(p)
	if err != nil {
		b.Fatal("Failed marshling payload")
	}

	for n := 0; n < b.N; n++ {
		ret := ginTest.DoTest(apitest.NewTestCase("TestGinPostJSON", "POST", "/postjson", `"OK"`+"\n", "", bytes.NewBuffer(payload)))
		if ret.Err != nil {
			b.Fatal(ret.Err)
		}
	}
}

func TestGinCloser(t *testing.T) {
	gints.Close()
}

func BenchmarkGinCloser(b *testing.B) {
	gints.Close()
}

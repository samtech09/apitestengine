## ApiTestEngine
It is wrapper to make http calles to given APIs using http GET, POST, PUT, DELETE methods against the given httptest.Server


### Usage
1. First create a httptest.Server, pass it to create new APITest engine.

```
// import package
import apitest "github.com/samtech09/apitestengine"
...
...
var chiTest *apitest.APITest


//create router/mux
r := chi.NewRouter()
r.Get("/", Index)
...


// create httptest server
s := httptest.NewServer(r)

// create and initialize test engine
chiTest = apitest.NewAPITest(s)
```

2. Create test-endpoints and make testcases

```
// single case
func TestChiIndex(t *testing.T) {
	ret := chiTest.DoTest(apitest.NewTestCase("TestChiIndex", "GET", "/", "[string to match with response]", apitest.MatchTypeXXX, "", nil))
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

//multiple cases
func TestChiParam(t *testing.T) {
	tcs := []apitest.TestCase{}
	tcs = append(tcs, apitest.NewTestCase("TestChiParam", "GET", "/param/1", "key = 1, apitest.MatchExact, "", nil))
	tcs = append(tcs, apitest.NewTestCase("TestChiParam2", "GET", "/param/123", "key = 123", apitest.MatchContains, "", nil))

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
```

View [Examples](https://github.com/samtech09/apitestengine/tree/master/example) for more details.

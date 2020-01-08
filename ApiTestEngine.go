package apitestengine

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

// TestCase allow to set URL calling parameters
type TestCase struct {
	ID                 string
	Method             string
	URL                string
	Expected           string
	Match              MatchType
	PostPutContentType string
	PostPutPayload     io.Reader
}

//TestCaseResult containd error (if any) of given test case
type TestCaseResult struct {
	ID  string
	Err error
}

//APITest struct for out test framework
type APITest struct {
	Ts *httptest.Server
}

// NewAPITest creates and return new ApiTest initialized with given httptest.Server
func NewAPITest(ts *httptest.Server) *APITest {
	return &APITest{ts}
}

//NewTestCase creates and return new-test-case with given parameters
func NewTestCase(id, method, url, expected string, match MatchType, postPutContentType string, postPutPayload io.Reader) TestCase {
	t := TestCase{}
	t.ID = id
	t.Method = method
	t.URL = url
	t.Expected = expected
	t.Match = match
	t.PostPutPayload = postPutPayload
	t.PostPutContentType = postPutContentType
	return t
}

func (t *APITest) executeTest(tcase TestCase) error {
	url := t.Ts.URL + tcase.URL
	//var resp *http.Response
	var err error
	var req *http.Request

	// Create client
	client := &http.Client{}

	switch tcase.Method {
	case "GET", "DELETE":
		// Create request
		req, err = http.NewRequest(tcase.Method, url, nil)
		if err != nil {
			return err
		}
		break

	case "POST", "PUT":
		// Create request
		req, err = http.NewRequest(tcase.Method, url, tcase.PostPutPayload)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", tcase.PostPutContentType)
		break

	default:
		return fmt.Errorf("unsupported http method")
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// check http status
	if status := resp.StatusCode; status != http.StatusOK {
		return fmt.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rsp := string(body)
	if tcase.Match == MatchContains {
		if !strings.Contains(rsp, tcase.Expected) {
			return fmt.Errorf("handler returned unexpected response: \n\tgot  %v \n\twant [Contains] %v", rsp, tcase.Expected)
		}
	} else if tcase.Match == MatchStartsWith {
		if !strings.HasPrefix(rsp, tcase.Expected) {
			return fmt.Errorf("handler returned unexpected response: \n\tgot  %v \n\twant [Start with] %v", rsp, tcase.Expected)
		}
	} else if tcase.Match == MatchEndsWith {
		if !strings.HasSuffix(rsp, tcase.Expected) {
			return fmt.Errorf("handler returned unexpected response: \n\tgot  %v \n\twant [Ends with] %v", rsp, tcase.Expected)
		}
	} else {
		// exact match
		if rsp != tcase.Expected {
			return fmt.Errorf("handler returned unexpected response: \n\tgot  %v \n\twant [Exact] %v", rsp, tcase.Expected)
		}
	}
	return nil
}

//DoTest run a single given test case
func (t *APITest) DoTest(cs TestCase) TestCaseResult {
	r := TestCaseResult{}
	r.ID = cs.ID
	r.Err = t.executeTest(cs)
	return r
}

//DoTests run all given test cases
func (t *APITest) DoTests(cases []TestCase) []TestCaseResult {
	res := []TestCaseResult{}
	for _, cs := range cases {
		res = append(res, t.DoTest(cs))
	}
	return res
}

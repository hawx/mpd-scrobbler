package lastfm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	ApiResponseStatusFailed = "failed"
)

type Err struct {
	Code    int
	Message string
}

func (e *Err) Error() string {
	return fmt.Sprintf("lastfm[%d]: %s", e.Code, e.Message)
}

func constructUrl(base string, params url.Values) string {
	return base + "?" + params.Encode()
}

func parseResponse(body io.Reader, result interface{}) error {
	var base Base
	if err := xml.NewDecoder(body).Decode(&base); err != nil {
		return err
	}
	if base.Status == ApiResponseStatusFailed {
		var errorDetail ApiError
		if err := xml.Unmarshal(base.Inner, &errorDetail); err != nil {
			return err
		}

		return &Err{errorDetail.Code, strings.TrimSpace(errorDetail.Message)}
	}

	if result != nil {
		return xml.Unmarshal(base.Inner, result)
	}
	return nil
}

func getSignature(params map[string]string, secret string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sigPlain string
	for _, k := range keys {
		sigPlain += k + params[k]
	}
	sigPlain += secret

	hasher := md5.New()
	hasher.Write([]byte(sigPlain))
	return hex.EncodeToString(hasher.Sum(nil))
}

//////////////
// POST API //
//////////////
func (api *Api) callPost(apiMethod string, args Args, result interface{}, withSession bool) error {
	urlParams := url.Values{}
	uri := constructUrl(api.uriBase, urlParams)

	//post data
	postData := url.Values{}
	postData.Add("method", apiMethod)
	postData.Add("api_key", api.params.apikey)
	if withSession {
		postData.Add("sk", api.params.sk)
	}

	tmp := make(map[string]string)
	tmp["method"] = apiMethod
	tmp["api_key"] = api.params.apikey
	if withSession {
		tmp["sk"] = api.params.sk
	}

	formated := args.Format()
	for k, v := range formated {
		tmp[k] = v
		postData.Add(k, v)
	}

	sig := getSignature(tmp, api.params.secret)
	postData.Add("api_sig", sig)

	res, err := http.PostForm(uri, postData)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return parseResponse(res.Body, result)
}

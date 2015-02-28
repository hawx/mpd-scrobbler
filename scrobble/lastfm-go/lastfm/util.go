package lastfm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func requireAuth(params *apiParams) (err error) {
	if params.sk == "" {
		err = newLibError(
			ErrorAuthRequired,
			Messages[ErrorAuthRequired],
		)
	}
	return
}

func constructUrl(base string, params url.Values) (uri string) {
	p := params.Encode()
	uri = base + "?" + p
	return
}

func toString(val interface{}) (str string, err error) {
	switch val.(type) {
	case string:
		str = val.(string)
	case int:
		str = strconv.Itoa(val.(int))
	case []string:
		ss := val.([]string)
		if len(ss) > 10 {
			ss = ss[:10]
		}
		str = strings.Join(ss, ",")
	default:
		err = newLibError(
			ErrorInvalidTypeOfArgument,
			Messages[ErrorInvalidTypeOfArgument],
		)
	}
	return
}

func parseResponse(body []byte, result interface{}) (err error) {
	var base Base
	err = xml.Unmarshal(body, &base)
	if err != nil {
		return
	}
	if base.Status == ApiResponseStatusFailed {
		var errorDetail ApiError
		err = xml.Unmarshal(base.Inner, &errorDetail)
		if err != nil {
			return
		}
		err = newApiError(&errorDetail)
		return
	} else if result == nil {
		return
	}
	err = xml.Unmarshal(base.Inner, result)
	return
}

func getSignature(params map[string]string, secret string) (sig string) {
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
	sig = hex.EncodeToString(hasher.Sum(nil))
	return
}

type Args interface {
	Format() (map[string]string, error)
}

//////////////
// POST API //
//////////////
func callPost(apiMethod string, baseUri string, params *apiParams, args Args, result interface{}, withSession bool) (err error) {
	if withSession {
		if err = requireAuth(params); err != nil {
			return
		}
	}

	urlParams := url.Values{}
	urlParams.Add("method", apiMethod)
	uri := constructUrl(baseUri, urlParams)

	//post data
	postData := url.Values{}
	postData.Add("method", apiMethod)
	postData.Add("api_key", params.apikey)
	if withSession {
		postData.Add("sk", params.sk)
	}

	tmp := make(map[string]string)
	tmp["method"] = apiMethod
	tmp["api_key"] = params.apikey
	if withSession {
		tmp["sk"] = params.sk
	}

	formated, err := args.Format()
	for k, v := range formated {
		tmp[k] = v
		postData.Add(k, v)
	}

	sig := getSignature(tmp, params.secret)
	postData.Add("api_sig", sig)

	//call API
	res, err := http.PostForm(uri, postData)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = parseResponse(body, result)
	return
}

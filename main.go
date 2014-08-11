package main
 
import (
	"errors"
	"fmt"
	"net/http"
	"net/http/fcgi"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"net/url"
	"log"
	"math"
	"flag"
)

type Counts struct {
	Url string `json:"url"`
	Count int64 `json:"count,string"`
}

var gplus_pat = regexp.MustCompile(`window\.__SSR\s=\s\{c:\s([0-9]+)\.0`)

func GetGPlus(pageUrl string) (int64, error) {
	safeUrl := url.QueryEscape(pageUrl)
	reqUrl := fmt.Sprintf("https://plusone.google.com/u/0/_/+1/fastbutton?url=%s&count=true", safeUrl)

	// log.Printf("calling G+ at %s", reqUrl)

	resp, err := http.Get(reqUrl)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	var count int64 = 0
	if submatches := gplus_pat.FindSubmatch(body); len(submatches) >= 2 {
		// log.Printf("pattern %s found in G+ response for %s", gplus_pat, reqUrl)

		countAsBytes := submatches[1]
		countAsStr := string(countAsBytes[:])

		// log.Printf("count is %s", countAsStr)
		
		count, err = strconv.ParseInt(countAsStr, 10, 64)

		if err != nil {
			return 0, nil
		}
	} else {
		log.Printf("warning: pattern %s not found in G+ response for %s", gplus_pat, reqUrl)
	}

	return count, nil
}

func GetJson(urlFormat string, pageUrl string, skipFront int, skipEnd int) (map[string]interface{}, error) {
	safeUrl := url.QueryEscape(pageUrl)
	reqUrl := fmt.Sprintf(urlFormat, safeUrl)
	
	resp, err := http.Get(reqUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var ret interface{}

	err = json.Unmarshal(body[skipFront:len(body)-skipEnd], &ret)

	return ret.(map[string]interface{}), err
}

func GetStumbleUpon(pageUrl string) (int64, error) {
	ret, err := GetJson("https://www.stumbleupon.com/services/1.01/badge.getinfo?url=%s", pageUrl, 0, 0)

	if err != nil {
		return 0, err
	}

	if ret["result"] == nil {
		log.Printf("StumbleUpon: bad response: result is empty")
		return 0, errors.New("bad response")
	}

	result := ret["result"].(map[string]interface{})

	if result["views"] == nil {
		return 0, nil
	}

	count := result["views"].(float64)

	return int64(math.Floor(count)), nil
	
}

func GetPinterest(pageUrl string) (int64, error) {
	ret, err := GetJson("https://api.pinterest.com/v1/urls/count.json?url=%s&callback=a", pageUrl, 2, 1)

	if err != nil {
		return 0, err
	}

	if ret["count"] == nil {
		log.Printf("Pinterest: bad response: count is empty")
		return 0, errors.New("bad response")
	}

	count := ret["count"].(float64)

	return int64(math.Floor(count)), nil
	
}

func GetReddit(pageUrl string) (int64, error) {
	// ret, err := GetJson("http://www.reddit.com/api/info.json?url=%s", pageUrl, 0, 0)

	// if err != nil {
	// 	return 0, err
	// }

	// if ret["children"] == nil {
	// 	log.Printf("Pinterest: bad response: count is empty")
	// 	return 0, errors.New("bad response")
	// }

	// children := ret["children"].([]interface{})

	// if len(children) == 0 {
	// 	return 0, nil
	// }

	// return int64(math.Floor(count)), nil

	return 0, nil
}

var vk_pat = regexp.MustCompile(`VK\.Share\.count\(([0-9]+), ([0-9]+)\);`)

func GetVK(pageUrl string) (int64, error) {
	safeUrl := url.QueryEscape(pageUrl)
	reqUrl := fmt.Sprintf("https://vk.com/share.php?act=count&index=0&url=%s", safeUrl)
	
	resp, err := http.Get(reqUrl)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	if submatches := vk_pat.FindSubmatch(body); len(submatches) >= 3 {
		// log.Printf("pattern %s found in VK response for %s", gplus_pat, reqUrl)

		countAsBytes := submatches[2]
		countAsStr := string(countAsBytes[:])

		// log.Printf("count is %s", countAsStr)
		
		return strconv.ParseInt(countAsStr, 10, 64)
	} else {
		log.Printf("warning: pattern %s not found in VK response for %s", vk_pat, reqUrl)
		return 0, errors.New("bad response")
	}

	return 0, errors.New("unreachable")
}

var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")


func Handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "max-age=300, public")
	header.Set("Pragma", "public")

	params := r.URL.Query()
	source := params.Get("type")
	pageUrl := params.Get("url")

	log.Printf("url %s type %s", pageUrl, source)
	
	counts := Counts{
		Url: pageUrl,
		Count: 0,
	}

	var count int64 = 0
	var err error = nil
	
	switch source {
	case "googlePlus":
		count, err = GetGPlus(pageUrl)
	case "stumbleupon":
		count, err = GetStumbleUpon(pageUrl)
	case "pinterest":
		count, err = GetPinterest(pageUrl)
	case "vkontakte":
		count, err = GetVK(pageUrl)
	}

	if err != nil {
		fmt.Fprintf(w, "ERROR %+v", err)
		w.WriteHeader(500)
	} else {
		counts.Count = count
	}
	
	bytes, err := json.Marshal(counts)

	log.Printf("output: %s", bytes)

	if err != nil {
		fmt.Fprintf(w, "ERROR %+v", err)
		w.WriteHeader(500)
	} else {
		w.Write(bytes)
	}
}

func main() {
	flag.Parse()

	handler := http.HandlerFunc(Handler)
	var err error
	
	if *local != "" { // Run as a local web server
        err = http.ListenAndServe(*local, handler)
    } else { // Run as FCGI via standard I/O
        err = fcgi.Serve(nil, handler)
    }
	
	if err != nil {
		fmt.Println(err)
	}
}

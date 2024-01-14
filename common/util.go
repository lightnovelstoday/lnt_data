package common

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lightnovelstoday/lnt_data/data"
)

var (
	BaseURL  = "https://www.example.com"
	Debug    = false
	Wait     = 5 * time.Second
	ImgWait  = time.Second
	dayCount = 2
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) ChromeRequestOrCache(httpPath string) string {
	fmt.Println("RequestOrCache: ", httpPath)
	if !strings.HasPrefix(httpPath, "http") && !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	if strings.HasPrefix(httpPath, "/") {
		httpPath = c.BaseURL + httpPath
	}
	if strings.Contains(httpPath, "\n") {
		httpPath = strings.ReplaceAll(httpPath, "\n", "")
	}

	u, err := url.Parse(httpPath)
	if err != nil {
		log.Fatal("Url Parse: ", err)
	}
	filename := fmt.Sprintf("cache/%x.html", md5.Sum([]byte(u.Path+u.RawQuery)))
	if Debug {
		fmt.Println(httpPath, filename)
	}
	stat, err := os.Stat(filename)
	if err == nil && stat.ModTime().After(time.Now().AddDate(0, 0, -dayCount)) {
		b, err := os.ReadFile(filename)
		if err == nil {
			if !strings.Contains("Article Not Found", string(b)) {
				return string(b)
			}
		}
	}
	if Debug {
		fmt.Println("Requesting: ", httpPath)
	}
	req, err := http.NewRequest("GET", httpPath, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Get: ", httpPath, err)
	}
	if resp.StatusCode == 404 {
		return ""
	}
	if resp.StatusCode != 200 {
		log.Fatal("Status: ", resp.StatusCode, httpPath)
	}
	f, _ := os.Create(filename)
	io.Copy(f, resp.Body)
	f.Close()
	time.Sleep(Wait)

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("ReadFile: ", filename, err)
	}
	if !strings.Contains("Article Not Found", string(b)) {
		return string(b)
	}
	log.Fatal("404: ", httpPath)
	return ""
}
func (c *Client) RequestOrCache2(httpPath string) (string, error) {
	st, _ := os.Stat("cache")
	if st == nil {
		os.Mkdir("cache", 0755)
	}
	if Debug {
		fmt.Println("RequestOrCache: ", httpPath)
	}
	if !strings.HasPrefix(httpPath, "http") && !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	if strings.HasPrefix(httpPath, "/") {
		httpPath = c.BaseURL + httpPath
	}
	if strings.Contains(httpPath, "\n") {
		httpPath = strings.ReplaceAll(httpPath, "\n", "")
	}

	u, err := url.Parse(httpPath)
	if err != nil {
		log.Fatal("Url Parse: ", err)
	}
	filename := fmt.Sprintf("cache/%x.html", md5.Sum([]byte(u.Path+u.RawQuery)))
	if Debug {
		fmt.Println(httpPath, filename)
	}
	stat, err := os.Stat(filename)
	if err == nil && stat.ModTime().After(time.Now().AddDate(0, 0, -dayCount)) {
		b, err := os.ReadFile(filename)
		if err == nil {
			if !strings.Contains("Article Not Found", string(b)) {
				return string(b), nil
			}
		}
	}
	if Debug {
		fmt.Println("Requesting: ", httpPath)
	}
	req, _ := http.NewRequest("GET", httpPath, nil)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status: %d %s", resp.StatusCode, httpPath)
	}
	f, _ := os.Create(filename)
	io.Copy(f, resp.Body)
	f.Close()
	time.Sleep(Wait)

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if !strings.Contains("Article Not Found", string(b)) {
		return string(b), nil
	}
	return "", fmt.Errorf("404: %s", httpPath)
}

func (c *Client) RequestOrCache(httpPath string) string {
	a, err := c.RequestOrCache2(httpPath)
	if err != nil {
		fmt.Println(a)
		panic(err)
	}
	return a
}
func (c *Client) RequestOrCacheBytes(httpPath string) []byte {
	if !strings.HasPrefix(httpPath, "http") && !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	if strings.HasPrefix(httpPath, "/") {
		httpPath = c.BaseURL + httpPath
	}
	if strings.Contains(httpPath, "\n") {
		httpPath = strings.ReplaceAll(httpPath, "\n", "")
	}

	u, err := url.Parse(httpPath)
	if err != nil {
		log.Fatal("Url Parse: ", err)
	}
	filename := fmt.Sprintf("cache/%x.html", md5.Sum([]byte(u.Path+u.RawQuery)))
	if Debug {
		fmt.Println(httpPath, filename)
	}
	stat, err := os.Stat(filename)
	if err == nil && stat.ModTime().After(time.Now().AddDate(0, 0, -dayCount)) {
		b, err := os.ReadFile(filename)
		if err == nil {
			return b
		}
	}
	if Debug {
		fmt.Println("Requesting: ", httpPath)
	}
	req, _ := http.NewRequest("GET", httpPath, nil)
	req.Header.Set("User-Agent", "LightNovels.today v1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Get: ", httpPath, err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Status: ", resp.StatusCode, httpPath)
	}
	f, _ := os.Create(filename)
	io.Copy(f, resp.Body)
	f.Close()
	time.Sleep(Wait)

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("ReadFile: ", filename, err)
	}
	return b
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "’", "")
	s = strings.ReplaceAll(s, "/", "-")

	return s
}
func ChromeRequestOrCache(httpPath string) string {
	return NewClient(BaseURL).ChromeRequestOrCache(httpPath)
}
func RequestOrCache2(httpPath string) (string, error) {
	return NewClient(BaseURL).RequestOrCache2(httpPath)
}

func RequestOrCache(httpPath string) string {
	return NewClient(BaseURL).RequestOrCache(httpPath)
}
func RequestOrCacheBytes(httpPath string) []byte {
	return NewClient(BaseURL).RequestOrCacheBytes(httpPath)
}

func SaveImage(href string) string {
	ext := filepath.Ext(href)
	if len(ext) > 5 {
		ext = strings.Split(ext, "?")[0]
	}
	if len(ext) > 5 {
		log.Fatal(ext)
	}
	filename := fmt.Sprintf("%x%s", sha1.Sum([]byte(href)), ext)

	filepath := "../../ln_images/img/" + filename
	cwd, _ := os.Getwd()
	if strings.HasSuffix(cwd, "ln_data") {
		filepath = "../ln_images/img/" + filename
	}

	if _, err := os.Stat(filepath); err == nil {
		return "/img/" + filename
	}
	resp, err := http.Get(href)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("GetImage: ", href, err)
	}
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal("Create: ", filename, err)
	}
	io.Copy(f, resp.Body)
	f.Close()
	resp.Body.Close()
	time.Sleep(ImgWait)
	return "/img/" + filename
}
func AttemptSaveImage(href string) (string, error) {
	ext := filepath.Ext(href)
	filename := fmt.Sprintf("%x%s", sha1.Sum([]byte(href)), ext)
	filepath := "../../ln_images/img/" + filename
	if _, err := os.Stat(filepath); err == nil {
		return "/img/" + filename, nil
	}
	resp, err := http.Get(href)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	io.Copy(f, resp.Body)
	f.Close()
	resp.Body.Close()
	time.Sleep(ImgWait)
	return "/img/" + filename, nil
}
func CleanStuff(s string) string {
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s
}
func AssignSeriesID(s data.Series) string {
	if s.Publisher == "" || s.Type == "" || s.Slug == "" {
		log.Fatal("Missing Publisher/Type/Slug: ", s)
	}
	fp := "known.json"
	k, err := os.Open(fp)
	if err != nil {
		fp = "../known.json"
		k, err = os.Open(fp)
		if err != nil {
			log.Fatal("Can't find known.json: ", err)
		}
	}
	m := make(map[string]string)
	json.NewDecoder(k).Decode(&m)
	k.Close()

	seriesLookup := strings.ToLower(
		fmt.Sprintf("%s-%s-%s", s.Publisher, s.Type, s.Slug),
	)
	if v, ok := m[seriesLookup]; ok {
		return v
	}

	uid, _ := uuid.NewRandom()
	m[seriesLookup] = uid.String()

	k, err = os.Create(fp)
	if err != nil {
		log.Fatal("Can't create known.json: ", err)
	}
	json.NewEncoder(k).Encode(m)
	k.Close()
	return uid.String()
}

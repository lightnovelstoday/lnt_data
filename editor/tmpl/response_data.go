package tmpl

import (
	"encoding/json"
	"html/template"
	"sort"
	"strconv"
	"strings"

	"github.com/acsellers/ln_shared/data"
)

type DataFile struct {
	Name     string
	Key      string
	Filename string
	Series   []data.Series
}
type ResponseData struct {
	Files      map[string]*DataFile
	Key        string
	SeriesList []data.Series
	Series     *data.Series
	Volume     *data.Volume
	Page       int
	Total      int
}

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}
func evenodd(i int) string {
	if i%2 == 0 {
		return "even"
	}
	return "odd"
}
func join(s []string) string {
	return strings.Join(s, ",")
}
func tojson(v any) template.HTML {
	b, _ := json.Marshal(v)
	return template.HTML(string(b))
}
func linkdata(links []data.PurchaseLink) string {
	s := []string{}
	for _, l := range links {
		s = append(s, l.Link)
	}
	return strings.Join(s, "\r\n")
}
func latest(volumes []data.Volume) string {
	dates := []string{}
	for _, v := range volumes {
		if v.Release != "" {
			dates = append(dates, v.Release)
		}
		if v.PrintRelease != "" {
			dates = append(dates, v.PrintRelease)
		}
		if v.DigitalRelease != "" {
			dates = append(dates, v.DigitalRelease)
		}
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i] > dates[j]
	})
	if len(dates) > 0 {
		return dates[0]
	}
	return ""
}
func itoa(i int) string {
	return strconv.Itoa(i)
}
func length[T any](Ts []T) string {
	return strconv.Itoa(len(Ts))
}

var (
	originalLangs = []string{
		"Japanese",
		"Chinese",
		"Korean",
		"Other",
	}
	publishers = []string{
		"Cross Infinite World",
		"Hanashi Media",
		"J-Novel Club",
		"Kaiten Books",
		"Kodansha USA",
		"One Peace Books",
		"Seven Seas",
		"Seven Seas Entertainment",
		"Square Enix",
		"Square Enix Manga & Books",
		"Tentai Books",
		"Viz Media",
		"Yen Press",
	}
	statuses = []string{
		"ongoing",
		"hiatus",
		"completed",
		"cancelled",
	}
	types = []string{
		"light_novel",
		"manga",
		"artbook",
	}
)

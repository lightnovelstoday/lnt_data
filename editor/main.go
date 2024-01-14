package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/lightnovelstoday/lnt_data/common"
	"github.com/lightnovelstoday/lnt_data/data"
	"github.com/lightnovelstoday/lnt_data/editor/tmpl"
	"github.com/lightnovelstoday/lnt_data/publishers"
)

var (
	port      = flag.Int("port", 8087, "Port to listen on")
	dataFiles = make(map[string]*tmpl.DataFile)
)

func main() {
	flag.Parse()

	for _, pubFile := range publishers.Files {
		if !strings.Contains(pubFile, "one-peace") {
			pubFile = strings.ReplaceAll(pubFile, "-", "_")
		}
		series := []data.Series{}
		f, err := os.Open(pubFile)
		if err != nil {
			log.Fatal("Open", pubFile, err)
		}
		err = json.NewDecoder(f).Decode(&series)
		if err != nil {
			log.Fatal("Decode", pubFile, err)
		}
		f.Close()

		key := strings.TrimSuffix(pubFile, "/output.json")
		dataFiles[key] = &tmpl.DataFile{
			Name:     series[0].Publisher,
			Key:      key,
			Filename: pubFile,
			Series:   series,
		}
		ConsistencyCheck(key)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", DashboardHandler)
	r.Get("/{key}/series", HomeHandler)
	r.Get("/{key}/series/new", NewSeriesHandler)
	r.Get("/{key}/series/{series}", SeriesHandler)
	r.Post("/{key}/series/", CreateSeriesHandler)
	r.Post("/{key}/series/{series}", UpdateSeriesHandler)
	r.Delete("/{key}/series/{series}", DeleteSeriesHandler)

	r.Get("/{key}/orphans", OrphansHandler)
	r.Post("/{key}/orphans", UpdateOrphanHandler)
	r.Get("/{key}/low-quality", LowQualityHandler)

	r.Get("/{key}/seriesimg/{page}", SeriesImgHandler)

	r.Get("/{key}/series/{series}/volumes/new", NewVolumeHandler)
	r.Get("/{key}/series/{series}/volumes/{volume}", VolumeHandler)
	r.Post("/{key}/series/{series}/volumes/", CreateVolumeHandler)
	r.Post("/{key}/series/{series}/volumes/{volume}", UpdateVolumeHandler)
	r.Delete("/{key}/series/{series}/volumes/{volume}", DeleteVolumeHandler)

	r.Get("/img/{image}", ImageHandler)

	fmt.Println("Listening on port", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Dashboard(tmpl.ResponseData{Files: dataFiles}).Render(r.Context(), w)
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	rd := tmpl.ResponseData{
		Files:      dataFiles,
		Key:        k,
		SeriesList: dataFiles[k].Series,
	}
	rd.Orphans = GetOrphans(rd.Key)
	tmpl.SeriesList(rd).Render(r.Context(), w)
}
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "image")
	filename := "../lnt_images/build/img/" + i
	st, _ := os.Stat(filename)
	if st == nil {
		filename = "../ln_images/img/" + i
		st, _ = os.Stat(filename)
		if st == nil {
			fmt.Println("Missing Image: ", i)
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
			return
		}
	}
	http.ServeFile(w, r, filename)
}
func OrphansHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	rd := tmpl.ResponseData{
		Files:      dataFiles,
		Key:        k,
		SeriesList: dataFiles[k].Series,
	}
	rd.Orphans = GetOrphans(rd.Key)
	tmpl.OrphanList(rd).Render(r.Context(), w)
}

type updateOrphanPayload struct {
	Title  string `json:"title"`
	Series string `json:"series"`
}

func UpdateOrphanHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	payload := updateOrphanPayload{}
	json.NewDecoder(r.Body).Decode(&payload)
	SaveOrphan(key, payload.Title, payload.Series)
	http.Redirect(w, r, fmt.Sprintf("/%s/orphans", key), 302)
}
func LowQualityHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	rd := tmpl.ResponseData{
		Files:      dataFiles,
		Key:        k,
		SeriesList: dataFiles[k].Series,
		Title:      "Low Quality Volumes",
		SubList:    filterLowQualityVolumes(dataFiles[k].Series),
	}
	tmpl.VolumeList(rd).Render(r.Context(), w)
}
func filterLowQualityVolumes(series []data.Series) []data.Series {
	filtered := []data.Series{}
	for _, s := range series {
		volumes := []data.Volume{}
		seen := map[string]bool{}
		for _, v := range s.Volumes {
			v := v
			issues := []string{}
			if seen[v.Title] {
				issues = append(issues, "Duplicate Title")
			}
			seen[v.Title] = true
			if strings.Contains(v.Title, "(Manga)") {
				issues = append(issues, "Format in Title")
			}
			if strings.Contains(v.Title, "(Light Novel)") {
				issues = append(issues, "Format in Title")
			}
			if v.CoverImage == "" {
				issues = append(issues, "No Cover Image")
			}
			if v.Description == "" {
				issues = append(issues, "No Description")
			}
			if strings.Contains(v.Description, "\u003c") && strings.Contains(v.Description, "\u003e") {
				issues = append(issues, "HTML in Description")
			}
			if v.PrintRelease == "" && v.DigitalRelease == "" {
				fmt.Println(v.Title, v.PrintRelease, v.DigitalRelease)
				issues = append(issues, "No Release Date")
			}
			if len(issues) > 0 {
				v.Title += " (" + strings.Join(issues, ", ") + ")"
				volumes = append(volumes, v)
			}
		}
		if len(volumes) > 0 {
			s.Volumes = volumes
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func SeriesImgHandler(w http.ResponseWriter, r *http.Request) {
	rd := tmpl.ResponseData{
		Files: dataFiles,
		Key:   chi.URLParam(r, "key"),
	}
	rd.SeriesList = dataFiles[rd.Key].Series

	rd.Page, _ = strconv.Atoi(chi.URLParam(r, "page"))
	if rd.Page < 1 {
		rd.Page = 1
	}
	total := len(rd.SeriesList) / 50
	if len(rd.SeriesList)%50 > 0 {
		total++
	}
	start := (rd.Page - 1) * 50
	end := rd.Page * 50
	if end > len(rd.SeriesList) {
		end = len(rd.SeriesList)
	}
	rd.SeriesList = rd.SeriesList[start:end]
	tmpl.SeriesImg(rd).Render(r.Context(), w)
}
func SeriesHandler(w http.ResponseWriter, r *http.Request) {
	rd := tmpl.ResponseData{
		Files: dataFiles,
		Key:   chi.URLParam(r, "key"),
	}
	rd.SeriesList = dataFiles[rd.Key].Series

	s := chi.URLParam(r, "series")
	rd.Series = GetSeries(rd.Key, s)
	if rd.Series == nil {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	tmpl.SeriesEdit(rd).Render(r.Context(), w)
}
func NewSeriesHandler(w http.ResponseWriter, r *http.Request) {
	rd := tmpl.ResponseData{
		Files:  dataFiles,
		Key:    chi.URLParam(r, "key"),
		Series: &data.Series{},
	}
	rd.SeriesList = dataFiles[rd.Key].Series
	tmpl.SeriesEdit(rd).Render(r.Context(), w)
}
func CreateSeriesHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	s := data.Series{}
	DecodeSeries(r, &s)
	if s.Image != "" && strings.HasPrefix(s.Image, "http") {
		s.LocalImage = common.SaveImage(s.Image)
	}
	id := InsertSeries(k, s)
	http.Redirect(w, r, fmt.Sprintf("/%s/series/%s", k, id), 302)
}
func UpdateSeriesHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	id := chi.URLParam(r, "series")
	s := GetSeries(k, id)

	if s.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}

	img := s.Image
	DecodeSeries(r, s)
	if s.Image == "" {
		s.LocalImage = ""
		s.Thumbnail = ""
		s.WebImage = ""
	}
	if strings.HasPrefix(s.Image, "http") && s.Image != img {
		s.LocalImage = common.SaveImage(s.Image)
		s.Thumbnail = ""
		s.WebImage = ""
	}

	UpdateSeries(k, s)
	http.Redirect(w, r, fmt.Sprintf("/%s/series/%s", k, s.ID), 302)
}
func DeleteSeriesHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	s := chi.URLParam(r, "series")
	series := GetSeries(k, s)
	if series.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	good := []data.Series{}
	for _, series := range dataFiles[k].Series {
		if series.ID != s {
			good = append(good, series)
		}
	}
	dataFiles[k].Series = good
	SaveData(k)
	ReadData(k)
	http.Redirect(w, r, fmt.Sprintf("/%s/series", k), 302)
}
func NewVolumeHandler(w http.ResponseWriter, r *http.Request) {
	rd := tmpl.ResponseData{
		Files:  dataFiles,
		Key:    chi.URLParam(r, "key"),
		Volume: &data.Volume{},
	}
	rd.SeriesList = dataFiles[rd.Key].Series

	id := chi.URLParam(r, "series")
	rd.Series = GetSeries(rd.Key, id)
	if rd.Series == nil {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	tmpl.VolumeEdit(rd).Render(r.Context(), w)
}
func VolumeHandler(w http.ResponseWriter, r *http.Request) {
	rd := tmpl.ResponseData{
		Files: dataFiles,
		Key:   chi.URLParam(r, "key"),
	}
	rd.SeriesList = dataFiles[rd.Key].Series

	s := chi.URLParam(r, "series")
	rd.Series = GetSeries(rd.Key, s)
	if rd.Series == nil {
		w.WriteHeader(404)
		w.Write([]byte("Series Not Found"))
		return
	}
	v := chi.URLParam(r, "volume")
	for _, vol := range rd.Series.Volumes {
		vol := vol
		if vol.ID == v {
			rd.Volume = &vol
		}
	}
	if rd.Volume == nil {
		w.WriteHeader(404)
		w.Write([]byte("Volume Not Found"))
		return
	}
	tmpl.VolumeEdit(rd).Render(r.Context(), w)
}
func CreateVolumeHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	s := chi.URLParam(r, "series")
	series := GetSeries(k, s)
	if series.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	v := data.Volume{}
	DecodeVolume(r, &v)

	uid, _ := uuid.NewRandom()
	v.ID = uid.String()
	v.SeriesID = series.ID
	v.Series = series.Slug
	if strings.HasPrefix(v.CoverImage, "http") {
		//v.LocalImage = common.SaveImage(v.CoverImage)
	}

	if v.Order == 0 {
		v.Order = len(series.Volumes) + 1
	}
	series.Volumes = append(series.Volumes, v)
	sort.Slice(series.Volumes, func(i, j int) bool {
		return series.Volumes[i].Order < series.Volumes[j].Order
	})

	UpdateSeries(k, series)
	http.Redirect(w, r, fmt.Sprintf("/%s/series/%s/volumes/%s", k, series.ID, v.ID), 302)
}
func UpdateVolumeHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	s := chi.URLParam(r, "series")
	v := chi.URLParam(r, "volume")
	series := GetSeries(k, s)
	if series.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	volume := data.Volume{}
	for _, vol := range series.Volumes {
		if vol.ID == v {
			volume = vol
		}
	}
	if volume.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}

	img := volume.CoverImage
	DecodeVolume(r, &volume)
	if volume.CoverImage == "" {
		volume.LocalImage = ""
		volume.Thumbnail = ""
		volume.WebImage = ""
	}
	if strings.HasPrefix(volume.CoverImage, "http") && volume.CoverImage != img {
		volume.LocalImage = common.SaveImage(volume.CoverImage)
		volume.Thumbnail = ""
		volume.WebImage = ""
	}
	for i, vol := range series.Volumes {
		if vol.ID == v {
			series.Volumes[i] = volume
		}
	}
	UpdateSeries(k, series)
	http.Redirect(w, r, fmt.Sprintf("/%s/series/%s/volumes/%s", k, series.ID, volume.ID), 302)
}
func DeleteVolumeHandler(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "key")
	s := chi.URLParam(r, "series")
	v := chi.URLParam(r, "volume")
	series := GetSeries(k, s)
	if series.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	volume := data.Volume{}
	for _, vol := range series.Volumes {
		if vol.ID == v {
			volume = vol
		}
	}
	if volume.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}

	good := []data.Volume{}
	for _, vol := range series.Volumes {
		if vol.ID != v {
			vol.Order = len(good) + 1
			good = append(good, vol)
		}
	}
	series.Volumes = good
	UpdateSeries(k, series)
	http.Redirect(w, r, fmt.Sprintf("/%s/series/%s", k, series.ID), http.StatusSeeOther)
}
func GetOrphans(key string) []data.Volume {
	orphans := []data.Volume{}
	st, _ := os.Stat(key + "/orphans.json")
	if st != nil {
		f, _ := os.Open(key + "/orphans.json")
		json.NewDecoder(f).Decode(&orphans)
		f.Close()
	}
	return orphans
}
func SaveOrphan(key string, title string, series string) {
	orphans := GetOrphans(key)
	index := -1
	for i := range orphans {
		if orphans[i].Title == title {
			index = i
		}
	}
	if index >= 0 {
		orphan := orphans[index]
		goodMatch := false
		for i := range dataFiles[key].Series {
			if dataFiles[key].Series[i].ID == series {
				orphan.Series = dataFiles[key].Series[i].Slug
				orphan.SeriesID = dataFiles[key].Series[i].ID
				dataFiles[key].Series[i].Volumes = append(dataFiles[key].Series[i].Volumes, orphan)
				goodMatch = true
			}
		}
		if goodMatch {
			SaveData(key)
			ReadData(key)
			DropOrphan(key, index)
		}
	}
}
func DropOrphan(key string, index int) {
	orphans := GetOrphans(key)
	good := []data.Volume{}
	for i := range orphans {
		if i != index {
			good = append(good, orphans[i])
		}
	}
	f, _ := os.Create(key + "/orphans.json")
	json.NewEncoder(f).Encode(good)
	f.Close()
}
func GetSeries(key, s string) *data.Series {
	for _, series := range dataFiles[key].Series {
		if series.ID == s {
			return &series
		}
	}
	return nil
}
func UpdateSeries(key string, s *data.Series) {
	for i, series := range dataFiles[key].Series {
		if series.ID == s.ID {
			dataFiles[key].Series[i] = *s
		}
	}
	SaveData(key)
	ReadData(key)
}
func InsertSeries(key string, s data.Series) string {
	id := common.AssignSeriesID(s)
	s.ID = id
	s.VersionLanguage = "English"
	dataFiles[key].Series = append(dataFiles[key].Series, s)
	SaveData(key)
	ReadData(key)
	return id
}
func ConsistencyCheck(key string) {
	for i, series := range dataFiles[key].Series {
		for j, volume := range series.Volumes {
			if volume.Order == 0 {
				dataFiles[key].Series[i].Volumes[j].Order = j + 1
			}
			if volume.SeriesID == "" {
				dataFiles[key].Series[i].Volumes[j].SeriesID = series.ID
				dataFiles[key].Series[i].Volumes[j].Series = series.Slug
			}
			if volume.ID == "" {
				uid, _ := uuid.NewRandom()
				dataFiles[key].Series[i].Volumes[j].ID = uid.String()
			}
		}
	}
}
func SaveData(key string) {
	ConsistencyCheck(key)
	dataFile := dataFiles[key].Filename
	data.OutputData(dataFiles[key].Series, dataFile)
}
func ReadData(key string) {
	f, _ := os.Open(dataFiles[key].Filename)
	json.NewDecoder(f).Decode(&dataFiles[key].Series)
	f.Close()
}
func DecodeSeries(r *http.Request, series *data.Series) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	series.Type = r.FormValue("type")
	series.Title = r.FormValue("title")
	series.Slug = common.Slugify(series.Title)

	series.Authors = parseStrings(r.FormValue("authors"))
	series.Illustrators = parseStrings(r.FormValue("illustrators"))
	series.Translators = parseStrings(r.FormValue("translators"))

	series.OriginalLanguage = r.FormValue("original_language")
	series.Publisher = r.FormValue("publisher")
	series.Status = r.FormValue("status")

	series.Description = r.FormValue("description")
	series.Image = r.FormValue("image")
	series.Website = r.FormValue("website")

	series.AutoGenres = parseStrings(r.FormValue("auto_genres"))
	series.MainGenres = parseStrings(r.FormValue("main_genres"))
	series.Setting = parseStrings(r.FormValue("setting"))
	series.Themes = parseStrings(r.FormValue("themes"))
	series.AgeLevel = parseStrings(r.FormValue("maturity_level"))
	series.Tags = parseStrings(r.FormValue("tags"))

	series.NULink = r.FormValue("nu_link")
	series.MDLink = r.FormValue("md_link")

	return nil
}
func DecodeVolume(r *http.Request, volume *data.Volume) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	volume.Title = r.FormValue("title")

	volume.ISBN = r.FormValue("isbn")
	volume.DigitalISBN = r.FormValue("digital_isbn")
	volume.PrintRelease = r.FormValue("print_release")
	volume.DigitalRelease = r.FormValue("digital_release")

	volume.Website = r.FormValue("website")
	volume.AltWebsite = r.FormValue("alt_website")
	volume.Order = parseInt(r.FormValue("order"))

	volume.Amazon.PaperbackASIN = r.FormValue("paperback_asin")
	volume.Amazon.HardcoverASIN = r.FormValue("hardcover_asin")
	volume.Amazon.DigitalASIN = r.FormValue("digital_asin")
	volume.Amazon.AudiobookASIN = r.FormValue("audiobook_asin")

	volume.CoverImage = r.FormValue("cover_image")
	volume.Description = r.FormValue("description")

	volume.PrintLinks = parseLinks(r.FormValue("print_links"))
	volume.DigitalLinks = parseLinks(r.FormValue("digital_links"))

	return nil
}
func parseLinks(s string) []data.PurchaseLink {
	links := []data.PurchaseLink{}
	if s == "" {
		return links
	}
	for _, link := range strings.Split(s, "\r\n") {
		link = strings.TrimSpace(link)
		links = append(links, data.NewPurchaseLink(link))
	}
	sort.Slice(links, func(i, j int) bool {
		return links[i].Vendor < links[j].Vendor
	})
	return links
}
func parseStrings(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

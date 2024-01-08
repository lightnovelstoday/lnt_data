package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/acsellers/ln_shared/amazon"
)

type Series struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"` // light-novel manga
	Slug             string                 `json:"slug"`
	Title            string                 `json:"title"`
	OtherTitles      []string               `json:"other_titles"`
	Authors          []string               `json:"author"`
	Translators      []string               `json:"translators"`
	Illustrators     []string               `json:"illustrators"`
	Roles            map[string][]string    `json:"roles"`
	AutoGenres       []string               `json:"auto_genres"`
	PrimaryGenres    []string               `json:"primary_genres"`
	MainGenres       []string               `json:"main_genres"`
	Setting          []string               `json:"setting"`
	Themes           []string               `json:"themes"`
	AgeLevel         []string               `json:"age_level"`
	OtherGenres      []string               `json:"other_genres"`
	Tags             []string               `json:"tags"`
	Publisher        string                 `json:"publisher"`
	Website          string                 `json:"website"`
	Image            string                 `json:"image"`
	WebImage         string                 `json:"web_image"`
	LocalImage       string                 `json:"local_image"`
	Thumbnail        string                 `json:"thumbnail"`
	Description      string                 `json:"description"`
	Universe         string                 `json:"universe"`
	ParentSeries     string                 `json:"parent_series"`
	ChildSeries      []string               `json:"child_series"`
	Status           string                 `json:"status"` // Ongoing, Complete, Hiatus, Cancelled
	OriginalLanguage string                 `json:"original_language"`
	VersionLanguage  string                 `json:"version_language"`
	AnnounceDate     string                 `json:"announce_date"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
	Formats          []string               `json:"formats"`
	Volumes          []Volume               `json:"volumes"`
	NULink           string                 `json:"nu_link,omitempty"`
	MDLink           string                 `json:"md_link,omitempty"`
}

type Volume struct {
	ID           string              `json:"id"`
	SeriesInt    int                 `json:"series_int,omitempty"`
	SeriesID     string              `json:"series_id"`
	Series       string              `json:"series"` // Use slug
	Title        string              `json:"title" form:"title"`
	Order        int                 `json:"order" form:"order"`
	Authors      []string            `json:"authors,omitempty"`
	Translators  []string            `json:"translators,omitempty"`
	Illustrators []string            `json:"illustrators,omitempty"`
	Roles        map[string][]string `json:"roles,omitempty"`
	SideVolume   bool                `json:"side_story"`
	CoverImage   string              `json:"cover_image" form:"cover_image"`
	WebImage     string              `json:"web_image"`
	LocalImage   string              `json:"local_image"`
	Thumbnail    string              `json:"thumbnail"`
	Description  string              `json:"description" form:"description"`
	Website      string              `json:"website" form:"website"`
	AltWebsite   string              `json:"alt_website,omitempty" form:"alt_website"`
	// Deprecated
	Release string `json:"release,omitempty"`
	// When we can get a release date for each medium
	DigitalRelease string `json:"digital_release"` // YYYY-MM-DD
	PrintRelease   string `json:"print_release"`
	// We can only get a single list of purchase links
	PurchaseLinks []PurchaseLink `json:"purchase_links"`

	// When we can get a list of purchase links for each medium
	DigitalLinks []PurchaseLink `json:"digital_links" form:"digital_links"`
	PrintLinks   []PurchaseLink `json:"print_links" form:"print_links"`
	Formats      []string       `json:"formats"`
	DigitalISBN  string         `json:"digital_isbn" form:"digital_isbn"`
	ISBN         string         `json:"isbn" form:"isbn"`
	Amazon       AmazonData     `json:"amazon"`

	Extra map[string]interface{} `json:"extra,omitempty"`
}
type AmazonData struct {
	PaperbackASIN  string  `json:"paperback_asin" form:"paperback_asin"`
	PaperbackPrice float32 `json:"paperback_price"`
	DigitalASIN    string  `json:"digital_asin" form:"digital_asin"`
	DigitalPrice   float32 `json:"digital_price"`
	HardcoverASIN  string  `json:"hardcover_asin" form:"hardcover_asin"`
	HardcoverPrice float32 `json:"hardcover_price"`
	AudiobookASIN  string  `json:"audiobook_asin" form:"audiobook_asin"`
	AudiobookPrice float32 `json:"audiobook_price"`
	BookRank       int     `json:"book_rank"`
	PhysicalRank   int     `json:"physical_rank"`
	DigitalRank    int     `json:"digital_rank"`
}

type SeriesPopularity struct {
	SeriesID          string             `json:"series_id"`
	SeriesName        string             `json:"series_name"`
	AveragePopularity float64            `json:"average_popularity"`
	Ranking           int                `json:"ranking"`
	LNRanking         int                `json:"ln_ranking"`
	Volumes           []VolumePopularity `json:"volumes"`
}
type VolumePopularity struct {
	SeriesID   string  `json:"series_id"`
	VolumeID   string  `json:"volume_id"`
	VolumeName string  `json:"volume_name"`
	Popularity float64 `json:"popularity"`
	Ranking    int     `json:"ranking"`
	LNRanking  int     `json:"ln_ranking"`
}

var (
	currentFolder  = fmt.Sprintf("amazon/%s/", time.Now().Format("2006-01"))
	previousFolder = fmt.Sprintf("amazon/%s/", time.Now().AddDate(0, 0, -1).Format("2006-01"))
)

func (ad AmazonData) GetProductData() []amazon.ProductData {
	asins := []string{}
	if ad.PaperbackASIN != "" {
		asins = append(asins, ad.PaperbackASIN)
	}
	if ad.DigitalASIN != "" {
		asins = append(asins, ad.DigitalASIN)
	}
	if ad.HardcoverASIN != "" {
		asins = append(asins, ad.HardcoverASIN)
	}

	ret := []amazon.ProductData{}
	for _, asin := range asins {

		data, err := os.ReadFile(currentFolder + asin + ".json")
		if err != nil {
			data, err = os.ReadFile(previousFolder + asin + ".json")
			if err != nil {
				log.Println("Error loading data for: ", asin)
				log.Println(currentFolder + asin + ".json")
				log.Println(err)
				continue
			}
		}
		var pd amazon.ProductData
		err = json.Unmarshal(data, &pd)
		if err != nil {
			log.Println("Error unmarshalling data for: ", asin)
			continue
		}
		ret = append(ret, pd)
	}
	return ret
}

func (v Volume) HasReleased() bool {
	dates := v.ReleaseDates()
	if len(dates) == 0 {
		return false
	}
	today := time.Now().Format("2006-01-02")
	for _, d := range dates {
		if d <= today {
			return true
		}
	}
	return false
}
func (v *Volume) ReleaseDates() []string {
	ret := []string{}
	if len(v.DigitalRelease) == 4 {
		fmt.Println("Bad digital release date", v.DigitalRelease, v.Series, v.Title)
	}
	if len(v.PrintRelease) == 4 {
		fmt.Println("Bad print release date", v.PrintRelease, v.Series, v.Title)
	}
	if len(v.Release) == 4 {
		fmt.Println("Bad release date", v.Release, v.Series, v.Title)
	}
	if v.DigitalRelease != "" {
		ret = append(ret, standardDate(v.DigitalRelease))
	}

	if v.PrintRelease != "" {
		ret = append(ret, standardDate(v.PrintRelease))
	}
	if v.Release != "" {
		ret = append(ret, standardDate(v.Release))
	}
	if len(ret) == 2 && ret[0] == ret[1] {
		return []string{ret[0]}
	}
	if len(ret) == 0 {
		return []string{"2099-12-31"}
	}
	return ret
}
func standardDate(d string) string {
	_, err := time.Parse("2006-01-02", d)
	if err == nil {
		return d
	}
	t, err := time.Parse("2006-1-2", d)
	if err == nil {
		return t.Format("2006-01-02")
	}
	return "2099-12-31"
}

type PurchaseLink struct {
	Link   string `json:"link"`
	Vendor string `json:"vendor"`
}

func OutputData(data []Series, filename string) error {
	sortable := SeriesSlice(data)
	sort.Sort(sortable)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(sortable)
	if err != nil {
		return err
	}
	return f.Close()
}
func OutputOrphans(data []Volume, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	return f.Close()
}

type MergeConfig struct {
	SeriesOverride map[string]bool
	VolumeOverride map[string]bool
}

func MergeData(data []Series, filename string) error {
	return MergeDataConfig(data, filename, MergeConfig{})
}
func MergeDataConfig(data []Series, filename string, config MergeConfig) error {

	if config.SeriesOverride == nil {
		config.SeriesOverride = map[string]bool{}
	}
	if config.VolumeOverride == nil {
		config.VolumeOverride = map[string]bool{}
	}
	var existing []Series
	f, _ := os.Open(filename)
	dec := json.NewDecoder(f)
	err := dec.Decode(&existing)
	if err != nil {
		return err
	}
	known := map[string]*Series{}
	for _, s := range existing {
		ls := s
		known[s.Type+"/"+s.Slug] = &ls
	}
	for _, s := range data {
		key := s.Type + "/" + s.Slug
		// new series
		if _, ok := known[key]; !ok {
			ls := s
			known[key] = &ls
			continue
		} else {
			mergeSeries(known[key], s, config)
		}

		for i := range s.Volumes {
			if i == len(known[key].Volumes) {
				known[key].Volumes = append(known[key].Volumes, s.Volumes[i])
				continue
			} else {
				mergeVolume(&known[key].Volumes[i], s.Volumes[i], config)
			}
		}
	}

	merged := []Series{}
	for _, s := range known {
		merged = append(merged, *s)
	}
	return OutputData(merged, filename)
}

var (
	missingCovers = map[string]bool{
		"https://yenpress-us.imgix.net/missing-cover.jpg":                                     true,
		"https://sevenseasentertainment.com/wp-content/uploads/2019/11/ss_nocover_header.jpg": true,
		"https://sevenseasentertainment.com/wp-content/uploads/2016/05/ss_nocover.jpg":        true,
	}
)

func mergeSeries(existing *Series, updated Series, config MergeConfig) {
	if existing.Title == "" || config.SeriesOverride["Title"] {
		existing.Title = updated.Title
	}
	if existing.Description == "" || config.SeriesOverride["Description"] {
		existing.Description = updated.Description
	}
	if existing.Website == "" || config.SeriesOverride["Website"] {
		existing.Website = updated.Website
	}
	if existing.Status == "" || config.SeriesOverride["Status"] {
		existing.Status = updated.Status
	}
	if existing.OriginalLanguage == "" || config.SeriesOverride["OriginalLanguage"] {
		existing.OriginalLanguage = updated.OriginalLanguage
	}
	if existing.VersionLanguage == "" || config.SeriesOverride["VersionLanguage"] {
		existing.VersionLanguage = updated.VersionLanguage
	}
	if existing.Authors == nil || len(existing.Authors) == 0 || config.SeriesOverride["Authors"] {
		existing.Authors = updated.Authors
	}
	if existing.Translators == nil || len(existing.Translators) == 0 || config.SeriesOverride["Translators"] {
		existing.Translators = updated.Translators
	}
	if existing.Illustrators == nil || len(existing.Illustrators) == 0 || config.SeriesOverride["Illustrators"] {
		existing.Illustrators = updated.Illustrators
	}
	if existing.AutoGenres == nil || len(existing.AutoGenres) == 0 || config.SeriesOverride["AutoGenres"] {
		existing.AutoGenres = updated.AutoGenres
	}
	if existing.PrimaryGenres == nil || len(existing.PrimaryGenres) == 0 || config.SeriesOverride["PrimaryGenres"] {
		existing.PrimaryGenres = updated.PrimaryGenres
	}
	if existing.MainGenres == nil || len(existing.MainGenres) == 0 || config.SeriesOverride["MainGenres"] {
		existing.MainGenres = updated.MainGenres
	}
	if existing.Setting == nil || len(existing.Setting) == 0 || config.SeriesOverride["Setting"] {
		existing.Setting = updated.Setting
	}
	if existing.Themes == nil || len(existing.Themes) == 0 || config.SeriesOverride["Themes"] {
		existing.Themes = updated.Themes
	}
	if len(existing.AgeLevel) == 0 || config.SeriesOverride["AgeLevel"] {
		existing.AgeLevel = updated.AgeLevel
	}
	if existing.Tags == nil || len(existing.Tags) == 0 || config.SeriesOverride["Tags"] {
		existing.Tags = updated.Tags
	}
	if existing.Image == "" || missingCovers[existing.Image] || config.SeriesOverride["Image"] {
		if existing.Image != updated.Image {
			existing.Image = updated.Image
			existing.LocalImage = updated.LocalImage
			existing.Thumbnail = updated.Thumbnail
		}
	}
	if config.SeriesOverride["LocalImage"] {
		existing.LocalImage = updated.LocalImage
	}
}
func mergeVolume(existing *Volume, updated Volume, config MergeConfig) {
	if existing.Title == "" || config.VolumeOverride["Title"] {
		existing.Title = updated.Title
	}
	if existing.Description == "" || config.VolumeOverride["Description"] {
		existing.Description = updated.Description
	}
	if existing.Authors == nil || len(existing.Authors) == 0 || config.VolumeOverride["Authors"] {
		existing.Authors = updated.Authors
	}
	if existing.Translators == nil || len(existing.Translators) == 0 || config.VolumeOverride["Translators"] {
		existing.Translators = updated.Translators
	}
	if existing.Illustrators == nil || len(existing.Illustrators) == 0 || config.VolumeOverride["Illustrators"] {
		existing.Illustrators = updated.Illustrators
	}
	if existing.CoverImage == "" || missingCovers[existing.CoverImage] || config.VolumeOverride["CoverImage"] {
		if existing.CoverImage != updated.CoverImage {
			existing.CoverImage = updated.CoverImage
			existing.LocalImage = updated.LocalImage
			existing.Thumbnail = updated.Thumbnail
		}
	}
	if config.VolumeOverride["LocalImage"] {
		existing.LocalImage = updated.LocalImage
	}
	if existing.Release == "" || config.VolumeOverride["Release"] {
		existing.Release = updated.Release
	}
	if existing.DigitalRelease == "" || config.VolumeOverride["DigitalRelease"] {
		existing.DigitalRelease = updated.DigitalRelease
	}
	if existing.PrintRelease == "" || config.VolumeOverride["PrintRelease"] {
		existing.PrintRelease = updated.PrintRelease
	}
	if existing.PurchaseLinks == nil || len(existing.PurchaseLinks) == 0 || config.VolumeOverride["PurchaseLinks"] {
		existing.PurchaseLinks = updated.PurchaseLinks
	} else {
		existing.PurchaseLinks = mergeLinks(existing.PurchaseLinks, updated.PurchaseLinks)
	}
	if existing.DigitalLinks == nil || len(existing.DigitalLinks) == 0 || config.VolumeOverride["DigitalLinks"] {
		existing.DigitalLinks = updated.DigitalLinks
	} else {
		existing.DigitalLinks = mergeLinks(existing.DigitalLinks, updated.DigitalLinks)
	}
	if existing.PrintLinks == nil || len(existing.PrintLinks) == 0 || config.VolumeOverride["PrintLinks"] {
		existing.PrintLinks = updated.PrintLinks
	} else {
		existing.PrintLinks = mergeLinks(existing.PrintLinks, updated.PrintLinks)
	}
	if existing.ISBN == "" || config.VolumeOverride["ISBN"] {
		existing.ISBN = updated.ISBN
	}
	if existing.DigitalISBN == "" || config.VolumeOverride["DigitalISBN"] {
		existing.DigitalISBN = updated.DigitalISBN
	}
}
func mergeLinks(a, b []PurchaseLink) []PurchaseLink {
	ret := append([]PurchaseLink{}, a...)
	seen := map[string]bool{}
	for _, link := range a {
		seen[link.Vendor] = true
	}
	for _, link := range b {
		if !seen[link.Vendor] {
			ret = append(ret, link)
		}
	}
	return ret
}

type SeriesSlice []Series

func (s SeriesSlice) Len() int {
	return len(s)
}
func (s SeriesSlice) Less(i, j int) bool {
	if s[i].Slug != s[j].Slug {
		return s[i].Slug < s[j].Slug
	}
	if s[i].Type != s[j].Type {
		return s[i].Type < s[j].Type
	}
	if s[i].Publisher != s[j].Publisher {
		return s[i].Publisher < s[j].Publisher
	}
	return s[i].ID < s[j].ID
}
func (s SeriesSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewPurchaseLink(link string) PurchaseLink {
	pl := PurchaseLink{
		Link: link,
	}
	switch {
	case strings.Contains(link, "amazon.com"):
		pl.Vendor = "Amazon US"
	case strings.Contains(link, "amazon.ca"):
		pl.Vendor = "Amazon Canada"
	case strings.Contains(link, "amazon.co.uk"):
		pl.Vendor = "Amazon UK"
	case strings.Contains(link, "amazon.de"):
		pl.Vendor = "Amazon Germany"
	case strings.Contains(link, "amazon.co.jp"):
		pl.Vendor = "Amazon Japan"
	case strings.Contains(link, "kobo.com"), strings.Contains(link, "kobobooks.com"):
		pl.Vendor = "Kobo"
	case strings.Contains(link, "barnesandnoble.com"):
		pl.Vendor = "Barnes & Noble"
	case strings.Contains(link, "booksamillion.com"):
		pl.Vendor = "Books-A-Million"
	case strings.Contains(link, "itunes.apple.com"):
		pl.Vendor = "Apple Books"
	case strings.Contains(link, "books.apple.com"):
		pl.Vendor = "Apple Books"
	case strings.Contains(link, "google.com"):
		pl.Vendor = "Google Play Books"
	case strings.Contains(link, "bookwalker"):
		pl.Vendor = "Bookwalker"
	case strings.Contains(link, "indigo.ca"):
		pl.Vendor = "Indigo"
	case strings.Contains(link, "gum.co"), strings.Contains(link, "gumroad.com"):
		pl.Vendor = "Gumroad"
	case strings.Contains(link, "rightstufanime.com"):
		pl.Vendor = "Right Stuf Anime"
	case strings.Contains(link, "bookdepository.com"):
		pl.Vendor = "Book Depository"
	case strings.Contains(link, "walmart.com"):
		pl.Vendor = "Walmart"
	case strings.Contains(link, "bookshop.org"):
		pl.Vendor = "Bookshop"
	case strings.Contains(link, "indiebound.org"):
		pl.Vendor = "IndieBound"
	case strings.Contains(link, "powells.com"):
		pl.Vendor = "Powell's"
	case strings.Contains(link, "comixology.com"):
		pl.Vendor = "Comixology"
	case strings.Contains(link, "penguinrandomhouse.com"):
		pl.Vendor = "Penguin Random House"
	case strings.Contains(link, "kinokuniya.com"):
		pl.Vendor = "Kinokuniya"
	case strings.Contains(link, "comicshoplocator.com"):
		pl.Vendor = "Comic Shop Locator"
	case strings.Contains(link, "kentai.com"):
		pl.Vendor = "Kentai Comics"
	case strings.Contains(link, "waterstones.com"):
		pl.Vendor = "Waterstones"
	case strings.Contains(link, "gomanga.com"):
		pl.Vendor = "GoManga"
	default:
		pl.Vendor = "Unknown"
	}
	return pl
}
func ToLinks(links []string) []PurchaseLink {
	var plinks []PurchaseLink
	for _, link := range links {
		plinks = append(plinks, NewPurchaseLink(link))
	}
	return plinks
}
func GenEditData(series []Series) []EditSeries {
	es := []EditSeries{}
	for _, s := range series {
		s2 := EditSeries{
			Title:     s.Title,
			Type:      s.Type,
			Publisher: s.Publisher,
			ID:        s.ID,
			Website:   s.Website,
		}
		for _, v := range s.Volumes {
			s2.Volumes = append(s2.Volumes, EditVolume{
				Title:          v.Title,
				ID:             v.ID,
				Website:        v.Website,
				DigitalRelease: v.DigitalRelease,
				PrintRelease:   v.PrintRelease,
			})
		}
		es = append(es, s2)
	}
	return es
}

type EditSeries struct {
	Title     string       `json:"title"`
	Type      string       `json:"type"`
	Publisher string       `json:"publisher"`
	ID        string       `json:"id"`
	Volumes   []EditVolume `json:"volumes"`
	Website   string       `json:"website"`
	Website2  string       `json:"website2"`
	GetID     string       `json:"get_id"`
}
type EditVolume struct {
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	ID             string `json:"id"`
	Website        string `json:"website"`
	Website2       string `json:"website2"`
	GetID          string `json:"get_id"`
	PrintRelease   string `json:"print_release,omitempty"`
	DigitalRelease string `json:"digital_release,omitempty"`
}

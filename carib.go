package opendmm

import (
  "fmt"
  "net/url"
  "regexp"
  "strings"
  "sync"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
)

func caribParse(keyword string, urlstr string, metach chan MovieMeta) {
  glog.Info("[CARIB] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[CARIB] Error: ", err)
    return
  }

  var meta MovieMeta
  meta.Code = fmt.Sprintf("Carib %s", keyword)
  meta.Page = urlstr

  var urlbase *url.URL
  urlbase, err = url.Parse(urlstr)
  if err != nil {
    return
  }
  var urlcover *url.URL
  urlcover, err = urlbase.Parse("./images/l_l.jpg")
  if err == nil {
    meta.CoverImage = urlcover.String()
  }
  var urlthumbnail *url.URL
  urlthumbnail, err = urlbase.Parse("./images/main_b.jpg")
  if err == nil {
    meta.ThumbnailImage = urlthumbnail.String()
  }

  meta.Title = doc.Find("#main-content > div.main-content-movieinfo > div.video-detail").Text()
  meta.Description = doc.Find("#main-content > div.main-content-movieinfo > div.movie-comment").Text()
  doc.Find("#main-content > div.detail-content.detail-content-gallery > ul > li > div > a").Each(
    func(i int, a *goquery.Selection) {
      href, ok := a.Attr("href")
      if ok {
        if !strings.Contains(href, "/member/") {
          meta.SampleImages = append(meta.SampleImages, href)
        }
      }
    })

  doc.Find("#main-content > div.main-content-movieinfo > div.movie-info > dl").Each(
    func(i int, dl *goquery.Selection) {
      dt := dl.Find("dt")
      if strings.Contains(dt.Text(), "出演") {
        meta.Actresses = dl.Find("dd a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(dt.Text(), "カテゴリー") {
        meta.Categories = dl.Find("dd a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(dt.Text(), "販売日") {
        meta.ReleaseDate = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "再生時間") {
        meta.MovieLength = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "スタジオ") {
        meta.Maker = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "シリーズ") {
        meta.Series = dl.Find("dd").Text()
      }
    })

  metach <- meta
}

func caribSearchKeyword(keyword string, metach chan MovieMeta) {
  glog.Info("[CARIB] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.caribbeancom.com/moviepages/%s/index.html",
    url.QueryEscape(keyword),
  )
  caribParse(keyword, urlstr, metach)
}

func caribSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
  glog.Info("[CARIB] Query: ", query)
  wg := new(sync.WaitGroup)
  re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
  matches := re.FindAllStringSubmatch(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%s-%s", match[1], match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      caribSearchKeyword(keyword, metach)
    }()
  }
  return wg
}

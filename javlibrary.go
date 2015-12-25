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

func javParse(urlstr string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[JAV] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[JAV] Error: ", err)
    return
  }

  var meta MovieMeta
  var ok bool
  meta.Page, ok = doc.Find("link[rel=shortlink]").Attr("href")
  if ok {
    meta.Code = doc.Find("#video_id .text").Text()
    meta.Title = strings.Replace(doc.Find("#video_title > h3").Text(), meta.Code, "", -1)
    meta.CoverImage, _ = doc.Find("#video_jacket > img").Attr("src")
    meta.ReleaseDate = doc.Find("#video_date .text").Text()
    meta.MovieLength = doc.Find("#video_length .text").Text()
    meta.Directors = doc.Find("#video_director .text span.director").Map(
      func(i int, span *goquery.Selection) string {
        return span.Text()
      })
    meta.Maker = doc.Find("#video_maker .text").Text()
    meta.Label = doc.Find("#video_label .text").Text()
    meta.Genres = doc.Find("#video_genres .text span.genre").Map(
      func(i int, span *goquery.Selection) string {
        return span.Text()
      })
    meta.Actresses = doc.Find("#video_cast .text span.cast span.star").Map(
      func(i int, span *goquery.Selection) string {
        return span.Text()
      })
    metach <- meta
  } else {
    base, err := url.Parse(urlstr)
    if err != nil {
      return
    }
    doc.Find("div.videothumblist > div.videos > div.video > a").Each(
      func(i int, s *goquery.Selection) {
        href, ok := s.Attr("href")
        if !ok {
          return
        }
        abshref, err := base.Parse(href)
        if err != nil {
          return
        }
        wg.Add(1)
        go func() {
          defer wg.Done()
          javParse(abshref.String(), metach, wg)
        }()
      })
    return
  }

}

func javSearchKeyword(kw string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[JAV] Keyword: ", kw)
  urlstr := fmt.Sprintf(
    "http://www.javlibrary.com/ja/vl_searchbyid.php?keyword=%s",
    url.QueryEscape(kw),
  )
  javParse(urlstr, metach, wg)
}

func javSearch(q string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[JAV] Query: ", q)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  matches := re.FindAllStringSubmatch(q, -1)
  for _, match := range matches {
    kw := fmt.Sprintf("%s-%s", match[1], match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      javSearchKeyword(kw, metach, wg)
    }()
  }
}

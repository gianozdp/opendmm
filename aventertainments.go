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

func aveParse(murl string, metach chan MovieMeta) {
  glog.Info("[AVE] Parse: ", murl)
  doc, err := newUtf8Document(murl)
  if err != nil {
    glog.Error("[AVE] Error: ", err)
    return
  }

  var meta MovieMeta
  var ok bool
  meta.Page = murl
  meta.Title = doc.Find("#mini-tabet > h2").Text()
  meta.CoverImage, ok = doc.Find("#titlebox > div.list-cover > img").Attr("src")
  if ok {
    meta.CoverImage = strings.Replace(meta.CoverImage, "jacket_images", "bigcover", -1)
  }
  meta.Code = strings.Replace(doc.Find("#mini-tabet > div").Text(), "商品番号:", "", -1)
  doc.Find("#titlebox > ul > li").Each(
    func(i int, li *goquery.Selection) {
      k := li.Find("span").Text()
      if strings.Contains(k, "主演女優") {
        meta.Actresses = li.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "スタジオ") {
        meta.Maker = li.Find("a").Text()
      } else if strings.Contains(k, "シリーズ") {
        meta.Series = li.Find("a").Text()
      } else if strings.Contains(k, "発売日") {
        meta.ReleaseDate = li.Text()
      } else if strings.Contains(k, "収録時間") {
        meta.MovieLength = li.Text()
      }
    })
  metach <- meta
}

func aveSearchKeyword(keyword string, metach chan MovieMeta) {
  glog.Info("[AVE] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.aventertainments.com/search_Products.aspx?keyword=%s",
    url.QueryEscape(keyword),
  )
  glog.Info("[AVE] Search: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if (err != nil) {
    glog.Error("[AVE] Error: ", err)
    return
  }

  href, ok := doc.Find("div.main-unit2 > table a").First().Attr("href")
  if ok {
    aveParse(href, metach)
  }
}

func aveSearch(query string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[AVE] Query: ", query)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  matches := re.FindAllStringSubmatch(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%s-%s", match[1], match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      aveSearchKeyword(keyword, metach)
    }()
  }
}
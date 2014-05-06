package kotoba

import (
  "fmt"
  "math"
  "strconv"
  "time"
  "github.com/yunyun/db"
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)

/*
 * Constants/Variables definitions
 */
const (
  PAGE_ERROR string = "/error/"
  FUNC_HOME string = "/"
  HOME_TEMPLATE = "home"
  REVIEW_TEMPLATE = "review"
  PAGE_REVIEW string = "/doReviews"
  FUNC_KOTOBA string = "/kotoba/"
  TEMPLATE_SHOW = "kotoba"
  TEMPLATE_EDIT = "edit"
  TEMPLATE_STATS = "stats"
)
/*
 * WEB FUNCTIONS
 */
type HomeData struct {
  Count int
  More int
  First []db.MongoKotoba
}
type ReviewData struct {
  Valid int
  ReviewUntil string
  Kotoba db.MongoKotoba
}
type StatsData struct {
  ReviewsNow int
  ReviewsHour int
  ReviewsDay int
  Beginner int
  Elementary int
  Intermediate int
  Master int
  Known int
}

func TimeUntil(t time.Time) string {
  until := t.Sub(time.Now().Local())
  if until.Hours() > 24 {
    return fmt.Sprintf("%v days %v hours %v minutes", math.Floor(until.Hours() / 24), 
                       math.Floor(math.Mod(until.Hours(), 24)), 
                       math.Ceil(math.Mod(until.Minutes(), 60)) )
  } else if until.Minutes() > 60 {
    return fmt.Sprintf("%v hours %v minutes", math.Floor(until.Hours()), math.Ceil(math.Mod(until.Minutes(), 60)))
  } else if until.Minutes() < 0 {
    return "available now"
  } else if until.Minutes() > 1 {
    return fmt.Sprintf("%v minutes", math.Ceil(until.Minutes()))
  } else {
    return fmt.Sprintf("%v minute", math.Ceil(until.Minutes()))
  }
                       
}
// WEB
func Home(rw http.ResponseWriter, req *http.Request, r render.Render, u *db.MongoUser, c *mgo.Collection) {
  k := db.GetAllKotoba(u.Id, c)
  
  data := HomeData{}
  data.Count = len(*k)
  if data.Count > 0 {
    data.First = (*k)[0:20]
    data.More = 1 // print X more
  } else {
    data.First = (*k)
  }
  r.HTML(200, HOME_TEMPLATE, data)
}
func Search(rw http.ResponseWriter, req *http.Request, r render.Render, u *db.MongoUser, c *mgo.Collection) {
  find := req.FormValue("search")
  k := db.FindKotoba(u.Id, find, c)
  var data HomeData
  data.Count = len(*k)
  if data.Count > 20 {
    data.First = (*k)[0:20]
    data.More = 1 // print X more
  } else {
    data.First = (*k)
  }
  fmt.Println(data.Count)
  r.HTML(200, HOME_TEMPLATE, data)
}
func ShowKotoba(params martini.Params, r render.Render, u *db.MongoUser, c *mgo.Collection) {
  id := params["id"]
  k := db.GetKotoba(bson.ObjectIdHex(id), c)
  r.HTML(200, TEMPLATE_SHOW, k)
}
func EditKotoba(params martini.Params, r render.Render, u *db.MongoUser, c *mgo.Collection) {
  id := params["id"]
  k := db.GetKotoba(bson.ObjectIdHex(id), c)
  r.HTML(200, TEMPLATE_EDIT, k)
}
func DoReviews(r render.Render, u *db.MongoUser, c *mgo.Collection) {
  k := db.GetRandomKotoba(u.Id, c)
  var data ReviewData
  if k != nil {
    data.Valid = 0
    data.Kotoba = *k
  } else {
    data.Valid = 1
    k = db.GetNextKotoba(u.Id, c)
    data.ReviewUntil = TimeUntil(k.Review)
  }
  r.HTML(200, REVIEW_TEMPLATE, data)
}
func AddKotoba(rw http.ResponseWriter, r *http.Request, u *db.MongoUser, c *mgo.Collection) {
  word, hatsuon, hatsuon_, imi, imi_ := r.FormValue("word"), r.FormValue("hatsuon"), r.FormValue("hatsuon_"), r.FormValue("imi"), r.FormValue("imi_")
  label := r.FormValue("label")
  diff, _ := strconv.Atoi(r.FormValue("diff"))
  _, err := db.SaveNewKotoba(u.Id, word, label, diff, hatsuon, hatsuon_ , imi, imi_, c)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, r, PAGE_ERROR, http.StatusFound)
  }
  http.Redirect(rw, r, FUNC_HOME, http.StatusFound)
}
func SaveEditKotoba(rw http.ResponseWriter, r *http.Request, params martini.Params, u *db.MongoUser, c *mgo.Collection) {
  word, hatsuon, hatsuon_, imi, imi_ := r.FormValue("word"), r.FormValue("hatsuon"), r.FormValue("hatsuon_"), r.FormValue("imi"), r.FormValue("imi_")
  label := r.FormValue("label")
  diff, _ := strconv.Atoi(r.FormValue("diff"))
  id := r.FormValue("word_id")
  _, err := db.SaveEditKotoba(id, word, label, diff, hatsuon, hatsuon_ , imi, imi_, c)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, r, PAGE_ERROR, http.StatusFound)
  }
  http.Redirect(rw, r, FUNC_KOTOBA + id, http.StatusFound)
}

func CheckReview(params martini.Params, rw http.ResponseWriter, r *http.Request, u *db.MongoUser, c *mgo.Collection) {
  checked := r.FormValue("checked")
  id := params["id"]
  if checked == "true" {
    db.ReviewUpdateKotoba(id, true, c)
  } else {
    db.ReviewUpdateKotoba(id, false, c)
  }
  http.Redirect(rw, r, PAGE_REVIEW, http.StatusFound)
}
func ShowStats(r render.Render, u *db.MongoUser, 
                c *mgo.Collection) {
  stats := StatsData{}
  stats.ReviewsNow = db.GetNumberReviewsNow(u.Id, c)
  stats.ReviewsHour = db.GetNumberReviewsHour(u.Id, c)
  stats.ReviewsDay = db.GetNumberReviewsDay(u.Id, c)
  stats.Beginner = db.GetNumberBeginner(u.Id, c)
  stats.Elementary = db.GetNumberElementary(u.Id, c)
  stats.Intermediate = db.GetNumberIntermediate(u.Id, c)
  stats.Master = db.GetNumberMaster(u.Id, c)
  stats.Known = db.GetNumberKnown(u.Id, c)
  r.HTML(200, TEMPLATE_STATS, stats)
}

// remove everything behind this point later!
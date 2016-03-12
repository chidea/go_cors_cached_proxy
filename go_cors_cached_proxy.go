package main

import (
	"container/ring"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	_ "net"
	"net/http"
	"regexp"
	"runtime"
	"time"
)

type RSSResult struct {
	XMLName xml.Name `xml:"rss"`
	Titles  []string `xml:"channel>item>title"`
}
type CacheItem struct {
	topic string
	ring  *ring.Ring
}

var cachemap = map[string]*CacheItem{
	"top":           &CacheItem{"", ring.New(10)},
	"world":         &CacheItem{"w", ring.New(10)},
	"US":            &CacheItem{"n", ring.New(10)},
	"business":      &CacheItem{"b", ring.New(10)},
	"technology":    &CacheItem{"t", ring.New(10)},
	"entertainment": &CacheItem{"e", ring.New(10)},
	"sports":        &CacheItem{"s", ring.New(10)},
	"science":       &CacheItem{"snc", ring.New(10)},
	"health":        &CacheItem{"m", ring.New(10)},
}
var cache = make([]byte, 1) // final json string output cache

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var title_company_filter = regexp.MustCompile("(.+) - [^-]+$") // regex to remove '- company' trail

func get_news(key string, topic string) bool {
	rssurl := "https://news.google.com/news?pz=1&cf=all&ned=us&hl=en&output=rss&topic=" + topic
	log.Printf("Getting news to cachemap[%s] from %s\n", key, rssurl)
	r, _ := http.Get(rssurl)
	if r == nil {
		log.Println("Unable to get news") // Probably the network is down
		return false
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	//ioutil.WriteFile("news.xml", b, 0666) // write rss to disk
	//d, _ := xml.NewDecoder(r.Body)
	v := RSSResult{}
	_ = xml.Unmarshal(b, &v)     // Parse and filter to get titles only
	for i, _ := range v.Titles { // reverse
		t := v.Titles[len(v.Titles)-1-i]
		st := title_company_filter.FindAllStringSubmatch(t, -1) // filter out '- company' in title
		var s string
		if st == nil { // not filterable
			s = t
		} else { // filterable
			s = st[0][1]
		}
		//fmt.Println(s)
		dup := false
		for _, v := range cachemap {
			v.ring.Do(func(o interface{}) { // check if already exists
				if dup || o == nil {
					return
				}
				if o == s {
					dup = true
				}
			})
		}
		if dup {
			fmt.Printf("%s [duplicated and removed from buffer]\n", s)
			continue
		}
		cachemap[key].addRing(s)
		fmt.Println(s)
	}
	return true
}

func (cacheItem *CacheItem) addRing(value interface{}) {
	cacheItem.ring = cacheItem.ring.Prev()
	cacheItem.ring.Value = value
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use all physical cores

	http.HandleFunc("/news", newshandler)
	go http.ListenAndServe(":81", nil)

	for {
		log.Println("Refreshing cachemap")
		for k, v := range cachemap {
			get_news(k, v.topic)
		}
		update_cache()
		/* //debug purpose code
		for k, v := range cachemap {
			v.ring.Do(func(o interface{}) {
				fmt.Println("cache", k, ":", o)
			})
		} */
		time.Sleep(15 * time.Minute)
	}
}

func update_cache() {
	m := make(map[string][]string)
	for k, v := range cachemap {
		v.ring.Do(func(o interface{}) {
			if o != nil {
				m[k] = append(m[k], o.(string))
			}
		})
	}
	cache, _ = json.Marshal(m)
}

func newshandler(w http.ResponseWriter, r *http.Request) {
	ori := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", ori)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	//json.NewEncoder(w).Encode(m)
	log.Printf("%s is retrieving cache\n", r.RemoteAddr)
	w.Write(cache)
}

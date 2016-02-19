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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_news(key string, topic string) {
	rssurl := "https://news.google.co.kr/news?pz=1&cf=all&ned=kr&hl=ko&output=rss&topic=" + topic
	log.Printf("Getting news to cachemap[%s] from %s\n", key, rssurl)
	r, _ := http.Get(rssurl)
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	//ioutil.WriteFile("news.xml", b, 0666) // write rss to disk
	c := regexp.MustCompile("(.+) - [^-]+$")
	//d, _ := xml.NewDecoder(r.Body)
	v := RSSResult{}
	_ = xml.Unmarshal(b, &v)
	for i, _ := range v.Titles { // reverse
		t := v.Titles[len(v.Titles)-1-i]
		s := c.FindAllStringSubmatch(t, -1)[0][1]
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
			fmt.Printf("%s is removed because of duplication\n", s)
			continue
		}
		cachemap[key].addRing(s)
		fmt.Println(s)
	}
}

type RSSResult struct {
	XMLName xml.Name `xml:"rss"`
	Titles  []string `xml:"channel>item>title"`
}
type CacheItem struct {
	topic string
	ring  *ring.Ring
}

func (cacheItem *CacheItem) addRing(value interface{}) {
	cacheItem.ring = cacheItem.ring.Prev()
	cacheItem.ring.Value = value
}

var cachemap = map[string]*CacheItem{
	"prim": &CacheItem{"", ring.New(10)},
	"soc":  &CacheItem{"y", ring.New(10)},
	"pol":  &CacheItem{"p", ring.New(10)},
	"eco":  &CacheItem{"b", ring.New(10)},
	"int":  &CacheItem{"w", ring.New(10)},
	"cul":  &CacheItem{"l", ring.New(10)},
	"cel":  &CacheItem{"e", ring.New(10)},
	"sci":  &CacheItem{"t", ring.New(10)},
}

var cacheready bool = false

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use all physical cores

	http.HandleFunc("/", handler)
	http.HandleFunc("/news", newshandler)
	http.HandleFunc("/weather", newshandler)
	go http.ListenAndServe(":81", nil)

	for {
		log.Println("Refreshing cachemap")
		for k, v := range cachemap {
			get_news(k, v.topic)
		}

		for k, v := range cachemap {
			v.ring.Do(func(o interface{}) {
				fmt.Println("cache", k, ":", o)
			})
		}
		time.Sleep(3 * time.Minute)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Got tired of life? Want some fun? Contact me instead of hacking hospitals : sbw228@gmail.com")
}

type JSONnews struct {
	Section string
	News    string
}

func newshandler(w http.ResponseWriter, r *http.Request) {
	ori := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", ori)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	m := make(map[string][]string)
	for k, v := range cachemap {
		v.ring.Do(func(o interface{}) {
			if o != nil {
				m[k] = append(m[k], o.(string))
			}
		})
	}
	json.NewEncoder(w).Encode(m)
	log.Printf("%s is retrieving cache\n", r.Host)
	/*for _, v := range cachemap {
		v.ring.Do(func(o interface{}) {

		})
	}*/
}

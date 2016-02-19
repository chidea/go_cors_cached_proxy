# go_cors_cached_proxy
CORS proxy that caches google news feeds with duplication checks into JSON

### Features
  - Gets news feed from [Google News](news.google.com) as RSS(XML) and converts it to JSON.
  - RSS filtering discards any description and images but headlines.
  - Additional string filtering to discard source company (in form of `Headline - Company`) if needed.
  - Interval based buffering which doesn't make burst of GETs to Google server to show news.
  - Ring buffer based updates. By default, google news feeds sends up to 10 newses at a time, and this server does so too.
  - CORS header is added to let client request from Cross-Origin. Useful when you use local webpage in some kiosk.

## How to use
  Run with `go run go_cors_cached_proxy.go` or `go install && go_cors_cached_proxy` (when path contains $GOPATH/bin directory).

  Try AJAX call on `localhost:81/news`. To try it with [jQuery](jquery.com), `$.getJSON('http://localhost:81/news', function(d){console.log(window.d=d);})`
  
  Note that it listens on `0.0.0.0` which means you instantly go public with running it.

## How to change language
  Open news.google.com with browser and check your final URL on browser with different language/country settings.
  Lanuage is specified with `hl` and country is with `ned`.
  Set this to `rssurl` variable in first line of `func get_news`.

  You may need to change `topic` settings of `cachemap` map to have right topic letters.
  Back on browser, check several sections on left and their link URLs.
  Top section will not contain any topic but all others will.

  Bellow sets AJAX response to have "top" as list of titles without topic letter, "world" with topic "w"
  ```
    "top":           &CacheItem{"", ring.New(10)},
    "world":         &CacheItem{"w", ring.New(10)},
  ```

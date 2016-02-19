# go_cors_cached_proxy
CORS proxy that caches google news feeds with duplication checks into JSON

## How to use
  Run with `go run go_cors_cached_proxy` or `go install && go_cors_cached_proxy` (when path contains $GOPATH/bin directory).

  Try AJAX call on `localhost:81`. To try it with jQuery, `$.getJSON('http://localhost:81/news', function(d){window.d=d;})`

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

### News APIs used

* https://developer.nytimes.com/


### Structure of config.json

```
[
  {
    "title": "articleSearch",
    "description": "With the Article Search API, you can search New York Times articles from Sept. 18, 1851 to today, retrieving headlines, abstracts, lead paragraphs, links to associated multimedia and other article...",
    "value": "YOUR-API-KEY-GOES-HERE",
    "target": "https://api.nytimes.com/svc/search/v2/articlesearch.json"
  }
]
```

### Golang notes

* In Go, string is a primitive type, it's readonly, every manipulation to it will create a new string
* In Go there is no such thing as a constant map
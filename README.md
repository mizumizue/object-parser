# object-parser

Struct convert to map[string]interface{} parser

## How to use

```
// definition struct
type Date time.Time
type CreateArticleRequest struct {
	ArticleNumber     int    `query:"articleNumber" search:"article_number,omitempty"`
	Title             string `query:"title" search:"title,omitempty"`
	AuthorID          string `query:"authorId" search:"author_id,omitempty"`
	PublishedDateFrom Date   `query:"publishedDateFrom" search:"published_date_from,omitempty"`
	PublishedDateTo   Date   `query:"publishedDateTo" search:"published_date_to,omitempty"`
}

now := time.Now()
req := &CreateArticleRequest{
	ArticleNumber       : 1, 
	Title               : "test",
	AuthorID            : "test-author",
	PublishedDateFrom   : Date(now),
	PublishedDateTo     : Date(now),
}

objectParser := NewObjectParser()
tagValueMap := objectParser.TagValueMap("search")
// tagValueMap
map[string]interface{
    "article_number": 1,
    "title": "test",
    "author_id": "test-author",
    "published_date_from": Date(now),
    "published_date_to": Date(now),
}

// can specified fields convert to other type
now := time.Now()
castMap := map[reflect.Type]reflect.Type{
    reflect.TypeOf(Date{}): reflect.TypeOf(time.Time{}),
}
req := &CreateArticleRequest{
	ArticleNumber       : 1, 
	Title               : "test",
	AuthorID            : "test-author",
	PublishedDateFrom   : Date(now),
	PublishedDateTo     : Date(now),
}

objectParser := NewObjectParser()
tagValueMap := objectParser.TagValueMap("search", castMap)

// tagValueMap
map[string]interface{
    "article_number": 1,
    "title": "test",
    "author_id": "test-author",
    "published_date_from": now,
    "published_date_to": now,
}

```

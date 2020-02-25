package object_parser

import (
	"reflect"
	"testing"
	"time"
)

func allocateStr(value string) *string {
	return &value
}

func allocateInt(value int) *int {
	return &value
}

func allocateDate(value Date) *Date {
	return &value
}

type Date time.Time

func (date Date) Convert() interface{} {
	return time.Time(date)
}

type CreateArticleRequest struct {
	ArticleNumber     int    `query:"articleNumber" search:"article_number,omitempty"`
	Title             string `query:"title" search:"title,omitempty"`
	AuthorID          string `query:"authorId" search:"author_id,omitempty"`
	PublishedDateFrom Date   `query:"publishedDateFrom" search:"published_date_from,omitempty"`
	PublishedDateTo   Date   `query:"publishedDateTo" search:"published_date_to,omitempty"`
}

type CreateArticleRequestNotOmitEmpty struct {
	ArticleNumber     int    `query:"articleNumber" search:"article_number"`
	Title             string `query:"title" search:"title"`
	AuthorID          string `query:"authorId" search:"author_id"`
	PublishedDateFrom Date   `query:"publishedDateFrom" search:"published_date_from"`
	PublishedDateTo   Date   `query:"publishedDateTo" search:"published_date_to"`
}

type CreateArticleRequestFieldPtr struct {
	ArticleNumber     *int    `query:"articleNumber" search:"article_number"`
	Title             *string `query:"title" search:"title"`
	AuthorID          *string `query:"authorId" search:"author_id"`
	PublishedDateFrom *Date   `query:"publishedDateFrom" search:"published_date_from"`
	PublishedDateTo   *Date   `query:"publishedDateTo" search:"published_date_to"`
}

func TestNewObjectParser(t *testing.T) {
	type Hoge struct {
		hoge string
		foo  string
		bar  string
	}
	type args struct {
		object interface{}
	}
	tests := []struct {
		name string
		args args
		want *ObjectParser
	}{
		{
			name: "",
			args: args{object: Hoge{}},
			want: &ObjectParser{
				Hoge{},
				reflect.TypeOf(Hoge{}),
				reflect.ValueOf(Hoge{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewObjectParser(tt.args.object); reflect.TypeOf(got.object) != reflect.TypeOf(tt.want.object) || got.objectType != tt.want.objectType || got.objectValue.Interface() != tt.want.objectValue.Interface() {
				t.Errorf("NewObjectParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectParser_NamedParam(t *testing.T) {
	now := time.Now()

	test1 := CreateArticleRequest{
		ArticleNumber:     1,
		Title:             "test-title",
		AuthorID:          "test-author",
		PublishedDateFrom: Date(now),
		PublishedDateTo:   Date(now),
	}

	test2 := CreateArticleRequest{}

	test3 := CreateArticleRequestNotOmitEmpty{
		ArticleNumber:     1,
		Title:             "test-title",
		AuthorID:          "test-author",
		PublishedDateFrom: Date(now),
		PublishedDateTo:   Date(now),
	}

	test4 := CreateArticleRequestNotOmitEmpty{}

	test5 := CreateArticleRequestFieldPtr{
		ArticleNumber:     allocateInt(1),
		Title:             allocateStr("test-title"),
		AuthorID:          allocateStr("test-author"),
		PublishedDateFrom: allocateDate(Date(now)),
		PublishedDateTo:   allocateDate(Date(now)),
	}

	test6 := CreateArticleRequestFieldPtr{}

	type fields struct {
		object      interface{}
		objectType  reflect.Type
		objectValue reflect.Value
	}
	type args struct {
		targetTag string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "",
			fields: fields{
				object:      test1,
				objectType:  reflect.TypeOf(test1),
				objectValue: reflect.ValueOf(test1),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{
				"article_number":      1,
				"title":               "test-title",
				"author_id":           "test-author",
				"published_date_from": now,
				"published_date_to":   now,
			},
		},
		{
			name: "",
			fields: fields{
				object:      test2,
				objectType:  reflect.TypeOf(test2),
				objectValue: reflect.ValueOf(test2),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{},
		},
		{
			name: "",
			fields: fields{
				object:      test3,
				objectType:  reflect.TypeOf(test3),
				objectValue: reflect.ValueOf(test3),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{
				"article_number":      1,
				"title":               "test-title",
				"author_id":           "test-author",
				"published_date_from": now,
				"published_date_to":   now,
			},
		},
		{
			name: "",
			fields: fields{
				object:      test4,
				objectType:  reflect.TypeOf(test4),
				objectValue: reflect.ValueOf(test4),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{
				"article_number":      0,
				"title":               "",
				"author_id":           "",
				"published_date_from": time.Time{},
				"published_date_to":   time.Time{},
			},
		},
		{
			name: "",
			fields: fields{
				object:      test5,
				objectType:  reflect.TypeOf(test5),
				objectValue: reflect.ValueOf(test5),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{
				"article_number":      1,
				"title":               "test-title",
				"author_id":           "test-author",
				"published_date_from": now,
				"published_date_to":   now,
			},
		},
		{
			name: "",
			fields: fields{
				object:      test6,
				objectType:  reflect.TypeOf(test6),
				objectValue: reflect.ValueOf(test6),
			},
			args: args{
				targetTag: "search",
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectParser := &ObjectParser{
				object:      tt.fields.object,
				objectType:  tt.fields.objectType,
				objectValue: tt.fields.objectValue,
			}
			if got := objectParser.TagValueMap(tt.args.targetTag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"\nObjectParser.TagValueMap() = %v "+
						"\n                     want = %v", got, tt.want,
				)
			}
		})
	}
}

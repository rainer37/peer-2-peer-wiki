package p2pwiki

import (
  "time"
)

path string = "/Path/to/file"

type Article struct {
  Title string,
  Content string,

  version int,
  timestamp int,  // unix timestamp
  size int  // bytes
}

func NewArticle(title string, content string) *Article {
  a := new(Article)
  a.Title = title
  a.Content = content
  a.version = 1
  a.timestamp = time.Now().Unix()
  a.size = len(title) + len(content)
  return a
}
// Read article from disk
func ReadArticle(title string, version int) *Article {
}
// Modify article
func (a *Article) Update(content string) *Article {
  b := NewArticle(a.Title, content)
  b.version++
  return b
}
// Returns all versions numbers of the article
func (a *Article) Versions() []int {

}
// Write article to disk
func (a *Article) Flush() (err Error) {

}
// Delete the article from disk
func (a *Article) Remove() (err Error) {

}

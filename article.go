package p2pwiki

import (
  "time"
  "itc"
)

path string = "/Path/to/file"

type Article struct {
  Title string,
  Content string,

  version *itc.Stamp,
  timestamp int,  // unix timestamp
  size int  // bytes
}

func NewArticle(title string, content string) *Article {
  a := new(Article)
  a.Title = title
  a.Content = content
  a.version = itc.NewStamp()
  a.timestamp = time.Now().Unix()
  a.size = len(title) + len(content)
  return a
}
// Read article from disk
func ReadArticle(title string) *Article {
}
// Write article to disk
func (a *Article) Flush() error {

}
// Delete the article from disk
func (a *Article) Remove() error {

}


// Modify article
func (a *Article) Modify(content string) {
  a.Content = content
  a.version.Event()
}

// Compare article Versions
func (a *Article) IsNewer(b *Article) (bool, error) {
  if a.Title != b.Title {
    return false, new Error("Cannot compare articles with different titles.")
  }
  return a.version.LEQ(b.version), nil
}

// Duplicate article and give it a new version id
func (a *Article) Replicate() *Article {
  b := a
  b.version = a.version.Fork()
  return b
}

// Combine an article with a replica
func (a *Article) Merge(b *Article) error {
  if a.Title != b.Title {
    return new Error("Can't merge articles with different titles.")
  }
  a.version = a.version.Join(b.version)

  return nil
}

// Convert the article to a byte[] for transmission
func (a *Article) Marshal() ([]byte, error) {

}

// Convert byte[] to article
func Unmarshal(data []byte) (*Article, error) {
  
}

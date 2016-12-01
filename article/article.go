package article

import (
  "fmt"
)

// TODO Deal with Vector Clock.
// TODO Need some way to push/pull articles.



type Article struct {
  Title string
  Hist Treedoc
  //Version VectorClock
  Buffer []OpLog
}


// Create a new empty article with specified title
func NewArticle(title string) {
  a := Article{}
  a.Title = title

  return &a
}

func OpenArticle(path string, title string) (*Article, error) {
  var a Article
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &a)

  return &a, err
}

// Write the article to the specified path as a JSON file. The file will be named
// with the title of the article.
func (a *Article) Save(path string) error {
  dat,err := json.Marshal(*a)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + a.Title + ".json", dat, 0644)
  return err
}


func (a *Article) Print() {
  atoms := a.Hist.Contents()
  fmt.Printf("\n%s\n", a.Title)
  fmt.Printf("%s\n", strings.Repeat("-", len(a.Title)))
  for _,paragraph := range atom {
    fmt.Printf("%s\n", paragraph)
  }
  fmt.Printf("\n\n")
}

func (a *Article) Insert(pos int, atom Atom, site Disambiguator) error {
  n,err := a.Hist.Insert(pos, atom, site)
  if err != nil {
    return err
  }

  // insert into log
  a.Buffer = append(a.Buffer, Operation{"insert", n})
}

func (a *Article) Delete(pos int, site Disambiguator) error {
  n,err := a.Hist.Delete(pos, site)
  if err != nil {
    return err
  }

  a.Buffer = append(a.Buffer, Operation{"delete", n})
}

// Replay all commands in the buffer
func (a *Article) FlushBuffer() error {
  var err error
  for _,op := range a.Buffer {
    switch op {
    case "insert":
      err = a.Hist.insertNode(op.Node)
    case "delete":
      err = a.Hist.deleteNode(op.Node)
    default:
      return fmt.Errorf("Article::FlushBuffer() - Invalid Command.")
    }
  }
  return err
}

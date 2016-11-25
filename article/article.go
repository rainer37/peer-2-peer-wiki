/*
  Package summary

 */
package article

import (
//  "log"
  "fmt"
//  "io"
  "github.com/nickbradley/p2pwiki/crdt"
  "encoding/json"
  "strings"
)

type Shared struct {
  Title string
  history crdt.Treedoc
}


func (a *Shared) Replay(log ReplayLog) error {
  for i,op := range log.Operations {
    switch op {
    case "insert":
      a.history.Insert(log.OpArgs[i].Text, log.OpArgs[i].Pos, log.Site)
    case "delete":
      a.history.Delete(log.OpArgs[i].Pos, log.Site)
    default:
      return fmt.Errorf("Unkown operation.")
    }
  }
  return nil
}




// Not persisted: user should push often
type Local struct {
  Title string
  Buffer []string
  Log ReplayLog
}

func NewLocal(title string, buffer []string) *Local {
  a := Local{title, buffer, ReplayLog{}}
  return &a
}
func OpenLocal(title string) (*Local, error) {
  return nil, nil
}

func (a *Local) Insert(pos int, text string) error {
  i := pos - 1
  a.Buffer = append(a.Buffer, "")
  copy(a.Buffer[i+1:], a.Buffer[i:])
  a.Buffer[i] = text

  a.Log.append("insert", OpArg{pos, text})
  return nil
}

func (a *Local) Delete(pos int) error {
  i := pos -1
  copy(a.Buffer[i:], a.Buffer[i+1:])
  a.Buffer[len(a.Buffer)-1] = ""
  a.Buffer = a.Buffer[:len(a.Buffer)-1]
  a.Log.append("delete", OpArg{pos,""})
  return nil
}

func (a *Local) Save() error {
  b,err := json.Marshal(*a)
  fmt.Println("JSON is", string(b[:]))
  return err
}

func (a *Local) String() string {
  return strings.Join(a.Buffer, "\n")
}






type ReplayLog struct {
  Site string
  Operations []string
  OpArgs []OpArg
}
func NewReplayLog(site string) ReplayLog {
  rlog := ReplayLog{site, []string{}, []OpArg{}}
  return rlog
}
func (r *ReplayLog) append(operation string, args OpArg) {
  r.Operations = append(r.Operations, operation)
  r.OpArgs = append(r.OpArgs, args)
}

type OpArg struct {
  Pos int
  Text string
}











/*
 * High-level summary of article here
 * Articles are composed of a series of paragraphs.
 * Paragraphs are the smallest unit tracked in the revision history.
 */


// Articles are paragraph-based

type Article struct {
  Title string  // treat this as an ID
  TitleHash string
  path string  // location on disk
  //Paragraphs []string  // for simplicity with the treedoc, store text by line
  Active bool  // indicates if article has been "deleted" (might need to be a vector clock)
  history crdt.Treedoc
  version crdt.DynamicVectorClock
}



// Create a new article and return a pointer to it
// If the article exists but is not active, re-initialize it
func NewArticle(title string, path string) *Article {
  a := Article{
    title,
    title,
    path,
    false,
    crdt.Treedoc{},
    crdt.DynamicVectorClock{},
  }

  return &a
}

// Open from disk
// func Open() (Article, error) {
//
// }

// Insert a paragraph into the article and update the revision history
func (a *Article) Insert(position int, text string) error {
  // TODO: Implment
  if position < 1 {
    return fmt.Errorf("Paragraph number must be greater than 0.")
  }
  err := a.history.Insert(text, position, "site1")
  if err != nil {
    return err
  }
  err = a.version.Increment("site1")
  if err != nil {
    fmt.Println("SDSD")
    a.version.Append("site1")
    a.version.Increment("site1")
  }
  // Insert in the correct position of the content array
  // Update the revision history
  // Increment vector clock

  return nil
}

func (a *Article) Print(showVersion bool) {
  paragraphs := a.history.Contents()
  fmt.Printf("\n%s\n", a.Title)
  fmt.Println("---------")
  for _,paragraph := range paragraphs {
    fmt.Printf("\n%s", paragraph)
    if showVersion {
      fmt.Printf(" [%d:%s]", a.version.Value("site2"), "site1")
    }
  }
  fmt.Printf("\n\n\n")
}

/*
// Remove a paragraph from the article and update the revision history
func (a *Article) RemoveParagraph(position int) error {
  // TODO: Implment
  if position < 1 {
    return fmt.Errorf("Paragraph number must be greater than 0.")
  }

  // Remove the specified element from the content array
  // Update the revison histroy
  // Increment vector clock

  return nil
}

// Set the active field to false
func (a *Article) Delete() {
  a.Active = false

  // Increment vector clock
}

// Delete the article from disk
func (a *Article) DeleteHard() error {
  // TODO: Implement

  // Remove from disk
  // Set for GC

  return nil
}

// Saves the article to disk
func (a *Article) Save() error {
  // TODO: Implement

  return nil
}

// Merge two articles by replacing the content array with the newer version and
// merging the revision history. The articles being merged must have the same title.
func (a *Article) Merge(b *Article) error {
  // TODO: Implement
  if a.Title != b.Title {
    return fmt.Errorf("Cannot merge articles with different titles.")
  }

  // Use the revision history to determine which article is newer and update accordingly
  // a > b: do nothing
  // a < b: overwrite a's content and merge b's history.

  return nil
}

// Compare article titles
func (a *Article) Equals(b *Article) bool {
  return a.Title == b.Title
}

// Compare article objects
func (a *Article) DeepEquals(b *Article) bool {
  return a == b  // TODO: implement deep equality
}

// Print the article
func (a *Article) ToString() string {
  str := fmt.Sprintln("*** %s ***", a.Title)
  for _,para := range a.Content {
    str += fmt.Sprintf("\t%s\n", para)
  }
  return str
}

// Print the article in raw format showing paragraph numbers
func (a *Article) PrintRaw() {
  fmt.Println(a.Title)
  for i,line := range a.Content {
    tmp := fmt.Sprintf("[%d]", i)
    fmt.Printf("%-5s\t%s\n", tmp, line)
  }
}

// Print the revision history of the article
func (a *Article) PrintHistory() {
  fmt.Println(a.history)
}








// // TODO: Move to the node package
// // Send the article over the specified connection
// func (a *Article) Send(conn) error {
//   // http://stackoverflow.com/questions/26030627/tcp-client-server-file-transfer-in-go
//   file,err := os.Open(a.name) // For read access.
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer file.Close() // make sure to close the file even if we panic.
//   n, err = io.Copy(connection, file)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(n, "bytes sent")
// }



// TODO Add garbage collection for articles with active = false. This is hard because
// there might be an update making it active again that hasn't reached the replica.
type Collection struct {
  name string
  articles map[[sha256.Size]byte]string
}
// name: full disk path where the store will be persisted
func NewCollection(name string) *Collection {
  c := new(Collection)
  c.name = name
  return c
}

// Adds an article to the collection
// Return an error if an article with the same title already exists
func (c *Collection) Add(a *Article) error {

}

// Determines whether an article is in the collection
func (c *Collection) Contains(key []byte) bool {

}

// Retrieve an article from the collection
func (c *Collection) Get(key []byte) (*Article, error) {

}

// Removes the specified article from the collection and deletes the article
func (c *Collection) Remove(key []byte) {
  // be sure to remove article from disk! (article.Delete())
}

// Removes all articles from the collection
func (c *Collection) Clear() {

}

// Persist store to disk
func (c *Collection) Save() error {
  // return error if can't write to disk
}
// Load store from disk into memory
func (c *Collection) Restore() error {
  // return error if can't read from disk
}
*/

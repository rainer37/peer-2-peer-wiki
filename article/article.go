/*
  Package summary

 */
package article

import (
  "log"
  "fmt"
  "io"
)






/*
 * High-level summary of article here
 * Articles are composed of a series of paragraphs.
 * Paragraphs are the smallest unit tracked in the revision history.
 */


// Articles are line-based

type Article struct {
  Title string  // treat this as an ID
  path string  // location on disk
  Paragraphs []string  // for simplicity with the treedoc, store text by line
  Active bool  // indicates if article has been "deleted" (might need to be a vector clock)
  history treedoc
}
// Create a new article and return a pointer to it
// If the article exists but is not active, re-initialize it
func NewArticle(...) *Article {
  // How to re-initialize so that it can be merged?
}

// Open from disk
func Open() (Article, error) {

}

// Insert a paragraph into the article and update the revision history
func (a *Article) InsertParagraph(position int, text string) error {
  // TODO: Implment
  if position < 1 {
    return fmt.Errorf("Paragraph number must be greater than 0.")
  }

  // Insert in the correct position of the content array
  // Update the revision history
  // Increment vector clock

  return nil
}

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

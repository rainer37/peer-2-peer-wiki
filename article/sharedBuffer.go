package article

import "fmt"

type SharedBuffer struct {
  Title string
  Hist Treedoc
}

func NewSharedBuffer(title string) *SharedBuffer {
  b := SharedBuffer{}
  b.Title = title

  return &b
}
func OpenSharedBuffer(title string, path string) *SharedBuffer {
  // TODO implement
  return nil
}


func (b *SharedBuffer) Replay(log OpLog) error {

  for i,op := range log.Operations {
    switch op {
    case "insert":
      b.Hist.Insert(log.OpArgs[i].Text, log.OpArgs[i].Pos, log.Site)
    case "delete":
      b.Hist.Delete(log.OpArgs[i].Pos, log.Site)
    default:
      return fmt.Errorf("Unkown operation in replay log.")
    }
  }

  return nil
}

func (b *SharedBuffer) Contents() []string {
  return b.Hist.Contents()
}

func (b *SharedBuffer) Save(path string) {
  // marshal file and write to disk
}

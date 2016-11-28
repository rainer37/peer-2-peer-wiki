package article

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  )

type SharedBuffer struct {
  Title string
  Hist Treedoc
}

func NewSharedBuffer(title string) *SharedBuffer {
  sb := SharedBuffer{}
  sb.Title = title

  return &sb
}
func OpenSharedBuffer(title string, path string) (*SharedBuffer, error) {
  var sb SharedBuffer
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &sb)

  return &sb, err
}


func (sb *SharedBuffer) Replay(log OpLog) error {

  for i,op := range log.Operations {
    switch op {
    case "insert":
      sb.Hist.Insert(log.OpArgs[i].Text, log.OpArgs[i].Pos, log.Site)
    case "delete":
      sb.Hist.Delete(log.OpArgs[i].Pos, log.Site)
    default:
      return fmt.Errorf("Unkown operation in replay log.")
    }
  }

  return nil
}

func (sb *SharedBuffer) Contents() []string {
  return sb.Hist.Contents()
}

func (sb *SharedBuffer) Save(path string) error {
  dat,err := json.Marshal(*sb)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + sb.Title + ".json", dat, 0644)
  return err
}

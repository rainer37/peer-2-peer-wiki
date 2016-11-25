package article

import (
  "fmt"
  "encoding/json"
  "strings"
  "io/ioutil"
)

type LocalBuffer struct {
  Title string
  Paras []string  // paragraphs
  Log OpLog  // log of operations
}

func NewLocalBuffer(title string, paragraphs []string, site string) *LocalBuffer {
  opLog := OpLog{}
  opLog.Site = site
  a := LocalBuffer{title, paragraphs, opLog}
  return &a
}

func OpenLocalBuffer(path string, title string) (*LocalBuffer, error) {
  var lb LocalBuffer
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &lb)

  return &lb, err
}

// Insert a paragraph at the location specified
func (a *LocalBuffer) Insert(pos int, text string) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(a.Paras):
    a.Paras = append(a.Paras, text)
  default:
    i := pos-1
    a.Paras = append(a.Paras, "")
    copy(a.Paras[i+1:], a.Paras[i:])
    a.Paras[i] = text
  }
  a.Log.append("insert", OpArg{pos, text})

  return nil
}

// Remove the paragraph at the specified position
func (a *LocalBuffer) Delete(pos int) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(a.Paras):
    _,a.Paras = a.Paras[len(a.Paras)-1], a.Paras[:len(a.Paras)-1]
  default:
    i := pos-1
    copy(a.Paras[i:], a.Paras[i+1:])
    a.Paras[len(a.Paras)-1] = ""
    a.Paras = a.Paras[:len(a.Paras)-1]
  }
  a.Log.append("delete", OpArg{pos,""})

  return nil
}

func (a *LocalBuffer) Save(path string) error {
  b,err := json.Marshal(*a)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + a.Title + ".json", b, 0644)
  return err
}


func (a *LocalBuffer) String() string {
  return strings.Join(a.Paras, "\n")
}



type OpLog struct {
  Site string
  Operations []string
  OpArgs []OpArg
}
func (r *OpLog) append(operation string, args OpArg) {
  r.Operations = append(r.Operations, operation)
  r.OpArgs = append(r.OpArgs, args)
}

type OpArg struct {
  Pos int
  Text string
}

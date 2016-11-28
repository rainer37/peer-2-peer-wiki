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
  opLog.Title = title
  opLog.Site = site
  lb := LocalBuffer{title, paragraphs, opLog}
  return &lb
}

func OpenLocalBuffer(path string, title string) (*LocalBuffer, error) {
  var lb LocalBuffer
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &lb)

  return &lb, err
}

// Insert a paragraph at the location specified
func (lb *LocalBuffer) Insert(pos int, text string) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(lb.Paras):
    lb.Paras = append(lb.Paras, text)
  default:
    i := pos-1
    lb.Paras = append(lb.Paras, "")
    copy(lb.Paras[i+1:], lb.Paras[i:])
    lb.Paras[i] = text
  }
  lb.Log.Append("insert", OpArg{pos, text})

  return nil
}

// Remove the paragraph at the specified position
func (lb *LocalBuffer) Delete(pos int) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(lb.Paras):
    _,lb.Paras = lb.Paras[len(lb.Paras)-1], lb.Paras[:len(lb.Paras)-1]
  default:
    i := pos-1
    copy(lb.Paras[i:], lb.Paras[i+1:])
    lb.Paras[len(lb.Paras)-1] = ""
    lb.Paras = lb.Paras[:len(lb.Paras)-1]
  }
  lb.Log.Append("delete", OpArg{pos,""})

  return nil
}

func (lb *LocalBuffer) Print() {
  fmt.Printf("\n%s\n", lb.Title)
  fmt.Printf("%s\n", strings.Repeat("-", len(lb.Title)))
  for _,paragraph := range lb.Paras {
    fmt.Printf("%s\n", paragraph)
  }
  fmt.Printf("\n\n")
}

func (lb *LocalBuffer) Save(path string) error {
  dat,err := json.Marshal(*lb)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + lb.Title + ".json", dat, 0644)
  return err
}


func (lb *LocalBuffer) String() string {
  return strings.Join(lb.Paras, "\n")
}


type OpLog struct {
  Title string
  Site string
  Operations []string
  OpArgs []OpArg
}
func (r *OpLog) Append(operation string, args OpArg) {
  r.Operations = append(r.Operations, operation)
  r.OpArgs = append(r.OpArgs, args)
}
func (r *OpLog) Remove(upto int) error {
  i := upto - 1

  switch {
  case upto < 1:
    fmt.Errorf("OpLog::Remove() - Invalid position.")
  case upto > len(r.Operations):
    i = len(r.Operations) - 1
  }

  // delete from Operations array
  copy(r.Operations[i:], r.Operations[i+1:])
  r.Operations[len(r.Operations)-1] = ""
  r.Operations = r.Operations[:len(r.Operations)-1]

  // delete from OpArgs array
  copy(r.OpArgs[i:], r.OpArgs[i+1:])
  r.OpArgs[len(r.OpArgs)-1] = OpArg{}
  r.OpArgs = r.OpArgs[:len(r.OpArgs)-1]

  return nil
}

type OpArg struct {
  Pos int
  Text string
}

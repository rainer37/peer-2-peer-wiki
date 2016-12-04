package article

import (
  "fmt"
)

func ConstructTreeManually() *Treedoc {
  var t Treedoc
  nodeA := Node{"A", "dA", false, nil, nil}
  nodeC := Node{"C", "dB", false, nil, nil}
  nodeE := Node{"E", "dA", false, nil, nil}
  nodeG := Node{"G", "dA", false, nil, nil}
  nodeH := Node{"H", "dB", false, nil, nil}
  nodeJ := Node{"J", "dB", false, nil, nil}
  nodeD := Node{"D", "dB", false, nil, nil}

  t1 := Treedoc{}
  t1.MiniNodes = append(t1.MiniNodes, &nodeA)

  t2 := Treedoc{}
  t2.MiniNodes = append(t2.MiniNodes, &nodeC)
  nodeB := Node{"B", "dA", false, nil, &t2}

  t3 := Treedoc{}
  t3.MiniNodes = append(t3.MiniNodes, &nodeE)

  t4 := Treedoc{}
  t4.MiniNodes = append(t4.MiniNodes, &nodeG)
  nodeF := Node{"F", "dA", false, nil, &t4}

  t5 := Treedoc{}
  t5.MiniNodes = append(t5.MiniNodes, &nodeH)
  nodeI := Node{"I", "dB", false, &t5, nil}

  t6 := Treedoc{}
  t6.MiniNodes = append(t6.MiniNodes, &nodeJ)


  t7 := Treedoc{}
  //nodeZ := Node{"Z", "dZ", false, nil, nil}
  t7.MiniNodes = append(t7.MiniNodes, &nodeB, &nodeD)
  t7.Left = &t1
  t7.Right = &t3

    // root
  t.MiniNodes = append(t.MiniNodes, &nodeF, &nodeI)
  t.Left = &t7
  t.Right = &t6

  return &t
}


func ExampleBuildTestTree() {
  t := ConstructTreeManually()
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I J]
}

func ExampleInsert1() {
  t := ConstructTreeManually()
  t.Insert(1, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [X A B C D E F G H I J]
}

func ExampleInsert2() {
  t := ConstructTreeManually()
  t.Insert(2, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A X B C D E F G H I J]
}

func ExampleInsert3() {
  t := ConstructTreeManually()
  t.Insert(3, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B X C D E F G H I J]
}

func ExampleInsert4() {
  t := ConstructTreeManually()
  t.Insert(4, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C X D E F G H I J]
}

func ExampleInsert5() {
  t := ConstructTreeManually()
  t.Insert(5, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D X E F G H I J]
}

func ExampleInsert6() {
  t := ConstructTreeManually()
  t.Insert(6, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E X F G H I J]
}

func ExampleInsert7() {
  t := ConstructTreeManually()
  t.Insert(7, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F X G H I J]
}

func ExampleInsert8() {
  t := ConstructTreeManually()
  t.Insert(8, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G X H I J]
}

func ExampleInsert9() {
  t := ConstructTreeManually()
  t.Insert(9, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H X I J]
}

func ExampleInsert10() {
  t := ConstructTreeManually()
  t.Insert(10, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I X J]
}

func ExampleInsert11() {
  t := ConstructTreeManually()
  t.Insert(11, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I J X]
}


func ExampleInsertNode() {
  t := ConstructTreeManually()
  t.insertNode(Path{PosId{Left,"dX"}}, &Node{"X", "dX", false, nil, nil})
  fmt.Println(t.Contents())

  // Output:
  // [A B C D X E F G H I J]
}


func ExampleInsertNode2() {
  t := ConstructTreeManually()
  t.insertNode(Path{PosId{Empty,"dX"}}, &Node{"X", "dX", false, nil, nil})
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I X J]
}


func ExampleDeleteNode1() {
  t := ConstructTreeManually()
  t.Delete(1, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [B C D E F G H I J]
}
func ExampleDeleteNode2() {
  t := ConstructTreeManually()
  t.Delete(2, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A C D E F G H I J]
}
func ExampleDeleteNode3() {
  t := ConstructTreeManually()
  t.Delete(3, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B D E F G H I J]
}
func ExampleDeleteNode4() {
  t := ConstructTreeManually()
  t.Delete(4, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C E F G H I J]
}
func ExampleDeleteNode5() {
  t := ConstructTreeManually()
  t.Delete(5, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D F G H I J]
}
func ExampleDeleteNode6() {
  t := ConstructTreeManually()
  t.Delete(6, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E G H I J]
}
func ExampleDeleteNode7() {
  t := ConstructTreeManually()
  t.Delete(7, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F H I J]
}
func ExampleDeleteNode8() {
  t := ConstructTreeManually()
  t.Delete(8, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G I J]
}
func ExampleDeleteNode9() {
  t := ConstructTreeManually()
  t.Delete(9, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H J]
}
func ExampleDeleteNode10() {
  t := ConstructTreeManually()
  t.Delete(10, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I]
}
func ExampleDeleteNode11() {
  t := ConstructTreeManually()
  t.Delete(11, "dA")
  fmt.Println(t.Contents())

  // Output:
  // [A B C D E F G H I]
}

func ExampleDeleteInsert1() {
  t := ConstructTreeManually()
  t.Delete(2, "dA")
  t.Insert(2, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A X C D E F G H I J]
}

func ExampleDeleteInsert2() {
  t := ConstructTreeManually()
  t.Delete(2, "dA")
  t.Delete(2, "dA")
  t.Insert(2, "X", "dX")
  fmt.Println(t.Contents())

  // Output:
  // [A X D E F G H I J]
}

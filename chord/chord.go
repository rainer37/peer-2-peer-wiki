package chord

import(
  "crypto/sha1"
  "math"
  "fmt"
  "net/rpc"
  "net"
	"time"
)

// TODO @rain - Pls implement these. Note that they will be called using RPC
// type chord interface {
//   // Lookup article on behalf of the client. Never explicitly return an error
//   // instead, set found to false.
//   Lookup(hTitle uint64, found *bool) error
//
//   // Lookup the article and return the contents as an array of paragraphs.
//   // Error should be explicity set if the article is not found since we can't
//   // differentiate between an empty article and an non-existent article using only
//   // the reply.
//   Pull(hTitle uint64, contents *[]string) error
//
//   // Send the operations log to the chord node resposible for hosting the article.
//   // The node should update the article's history and return the number of successful
//   // operations applied. Never return a error: if an error is encountered during
//   // replay, no further operations should be applied and the number of successful
//   // opertations should be returned.
//   Push(log OpLog, replayCount *int) error
// }



const m int = 6
const r int = 3 // degree of Fault Tolerance

var master_IP string = "127.0.0.1:1338"

type Article struct {
	Title string
	Content string
}

type Node struct {
	IP string
	ID uint64
	Successor_IP string
	Predecessor_IP string
	Fingers [m]Finger
	Data map[string]*Article
	Suc_list [r]Finger
}

type Finger struct {
	FID uint64
	FIP string
}



//万能struct
type Args struct {
	A string
	B string
}



func Hash(key string) uint64 {
  result := sha1.Sum([]byte(key))
  sum := uint64(0)
  for i := 0; i < len(result); i++ {
    sum += uint64(math.Pow(256.0, float64(i))) * uint64(result[i])
  }
  sum %= 64
  //fmt.Println("SHASUM of",key,": ",sum)
  return sum
}

func RPCall(ip string, args interface{}, reply interface{}, f string) error {
  client, err := rpc.Dial("tcp", ip)
  if err != nil {
    fmt.Println(err)
    return err
  }

  err = client.Call(f, args, reply)

  if err != nil {
    println(err)
    return err
  }

  return nil
}

func CheckErr(err error) {
  if err != nil {
    println(err)
  }
}

func RPCListen(ip string, n *Node) {
  addr, err := net.ResolveTCPAddr("tcp", ip)
  CheckErr(err)

  inbound, err := net.ListenTCP("tcp", addr)
  CheckErr(err)
  var art *Article = &Article{"Sample Title1", "12"}

  rpc.Register(n)
  rpc.Register(art)
  rpc.Accept(inbound)
}




func NewNode(ip string) *Node {
	return &Node{
		IP:		ip,
		ID:    	Hash(ip),
		Data:	make(map[string]*Article),
	}
}

func (a *Article) Get_Content(args *Args, reply *int) error {
	println(a.Content)
	return nil
}

func (n *Node) Report(args *Args, reply *int) error {
	println("Node:",n.ID,"@",n.IP,"Data Lenth:",len(n.Data),"Suc:",Hash(n.Successor_IP),"Pred:",Hash(n.Predecessor_IP))
	return nil
}

func (n *Node) Notify(args *Args, reply *string) error {
	//println(Hash(args.A), Hash(n.IP))
	if (n.Predecessor_IP == n.Successor_IP || Hash(n.Predecessor_IP) != Hash(args.A)) {
		n.Predecessor_IP = args.A
		println("Now my predecessor is:", args.A)
		for i,v := range n.Data {
			if (Hash(i) > Hash(n.IP) || Hash(i) <= Hash(args.A)) && Hash(n.IP) > Hash(args.A) ||
			(Hash(n.IP) < Hash(args.A) && (Hash(i) > Hash(n.IP) && Hash(i) <= Hash(args.A))){
				println("You should have",i,Hash(i))
				line := Args{i,v.Content}
				var rp int = 0
				err := RPCall(args.A, &line, &rp, "Node.Put")
				CheckErr(err)
				delete(n.Data, i)
			}
		}
	} else {
		println("You can't be my pred, seraching for someone else...")
	}
	return nil
}

func (n *Node) linear_search(target_id uint64, reply *string) {
	ip_hash := Hash(n.IP)
	suc_hash := Hash(n.Successor_IP)
	pred_hash := Hash(n.Predecessor_IP)

	println("Dude, looking for", target_id)

	if n.IP == n.Successor_IP {
		println("1Found! You should join",n.IP)
		*reply = n.IP
		return
	}

	if (target_id < ip_hash && target_id > pred_hash)||(target_id < ip_hash && ip_hash < pred_hash) {
		println("2Found! You should join",n.IP)
		*reply = n.IP
	} else if (target_id > ip_hash && suc_hash < ip_hash) || (target_id > ip_hash && suc_hash > target_id){
 		println("3Found! You should join",n.Successor_IP)
		*reply = n.Successor_IP
	}else {
		err := RPCall(n.Successor_IP, &target_id, reply, "Node.Linear_search_with_id")
		CheckErr(err)
	}
}

// search for successor ip with ip
func (n *Node) Linear_search_with_id(id *uint64, reply *string) error {
	n.linear_search(*id, reply)
	return nil
}

// search for successor ip with ip
func (n *Node) Linear_search_with_ip(args *Args, reply *string) error {
	println(args.A)
	target_id := Hash(args.A)
	n.linear_search(target_id, reply)
	return nil
}

func (n *Node) stablize() {
	var reply string = "No"
	line := Args{"",""}

	err := RPCall(n.Successor_IP, &line, &reply, "Node.GetPredecessor")

	if err != nil {
		println("Successor is gone...")
	}

	//println("My successor's pred is:",reply)

	if reply != n.IP && reply != "No" {
		println("Time to stablize...")

		n.Successor_IP = reply

		line.A = n.IP
		err := RPCall(n.Successor_IP, &line, &reply, "Node.Notify")
		CheckErr(err)
		println("Stablized...")
	} else {
		//println("Now i'm alone...")
	}
}

func (n *Node) finger_gen(){
	for i := 0; i < m; i++ {
		n.Fingers[i].FID = n.ID + uint64(math.Pow(2,float64(i)))
		println(n.Fingers[i].FID)
	}
}

func (n *Node) GetSuccessor(args *Args,reply *string) error {
	*reply = n.Successor_IP
	return nil
}

func (n *Node) GetPredecessor(args *Args,reply *string) error {
	*reply = n.Predecessor_IP
	return nil
}

func (n *Node) Get(args *Args, reply *string) error {

	if n.Data[args.A] != nil {
 		*reply = n.Data[args.A].Content
	} else {
		*reply = "No Such An Article Found"
	}

	return nil
}

func (n *Node) Put(args *Args, reply *int) error {
	println("putting ",args.A)

	if n.Data[args.A] != nil {
		n.Data[args.A].Content = args.B
	} else {
		n.Data[args.A] = &Article{args.A, args.B}
	}
	*reply = 1
	return nil
}

func (n *Node) Join(to_join string) {
	println(n.IP,"Joining...")
	n.finger_gen()

	var reply string
	line := Args{n.IP,""}

	err := RPCall(to_join, &line, &reply, "Node.Linear_search_with_ip") // findSuccessor()

	println("Will join",reply,Hash(reply))

	n.Successor_IP = reply
	n.Suc_list[0] = Finger{Hash(reply), reply}
	n.Predecessor_IP = "No"

	periodic(n)

	line.A = n.IP

	go RPCListen(n.IP, n)

    err = RPCall(reply, &line, &reply, "Node.Notify")
    CheckErr(err)

	for{
		// loop forever
	}
}

func (n *Node) println_suc_list() {
	for i := 0; i < r ; i++ {
		print(n.Suc_list[i].FID," ")
	}
}

// Actually this is not recoil.
// generate successor list.
func (n *Node) recoil() {
	n.Suc_list[0] = Finger{Hash(n.Successor_IP),n.Successor_IP}

	if n.IP == n.Successor_IP {
		println("Alone...")
		for i := 1; i < r ; i++ { n.Suc_list[i] = Finger{0,""} }
		return
	}


	var reply string = n.Successor_IP
	line := Args{"",""}

	for i := 0; i < r-1 ; i++ {
		err := RPCall(reply, &line, &reply, "Node.GetSuccessor")
		if err != nil {
			println(err)
			n.Suc_list[i+1] = Finger{0,""}

		} else {
			n.Suc_list[i+1] = Finger{Hash(reply),reply}
		}
	}
}

func (n *Node) ping() {
	var reply string
	line := Args{n.IP,""}
	err := RPCall(n.Successor_IP, &line, &reply, "Node.GetSuccessor") // findSuccessor()

	if err != nil {
		println("Successor is down...")
		if n.Suc_list[1].FIP != "" {
			err = RPCall(n.Suc_list[1].FIP, &line, &reply, "Node.Notify") // findSuccessor()
			if err != nil {
				println("Successor's successor is also down...")
			} else {
				n.Successor_IP = n.Suc_list[1].FIP
			}
		}
	}
}

func periodic(n *Node) {
	ticker := time.NewTicker(time.Millisecond * 3000)
    go func() {
        for range ticker.C {
            n.ping() // check if successor is still alive
            n.stablize()
            n.recoil() // generate successor list
            quick_report(n)
        }
    }()
}

func quick_report(n *Node) {
	//println("Succ:", n.Successor_IP, "Pred:", n.Predecessor_IP)
	for i,_ := range n.Data {
		print(Hash(i)," ")
	}
	print("Pred: ", Hash(n.Predecessor_IP), " Succ: ")
	n.println_suc_list()
	println()
}

func CreateRing() {

	first_node := NewNode(master_IP)
	first_node.Successor_IP = master_IP
	first_node.Predecessor_IP = "No"
	first_node.Suc_list[0] = Finger{Hash(master_IP), master_IP}

	first_node.Data["Sample"] = &Article{"Sample Title", "The family provided by fmt, however, are built to be in production code. They report predictably to  stdout, unless otherwise specified. They are more versatile (fmt.Fprint* can report to any io.Writer, such as os.Stdout, os.Stderr, or even a net.Conn type.) and are not implementation specific."}
	first_node.Data["Sample1"] = &Article{"Sample Title1", "1"}
	first_node.Data["Sample2"] = &Article{"Sample Title2", "2"}
	first_node.Data["Sample3"] = &Article{"Sample Title3", "3"}
	first_node.Data["Sample4"] = &Article{"Sample Title4", "4"}
	first_node.Data["Sample5"] = &Article{"Sample Title5", "5"}
	first_node.Data["Sample6"] = &Article{"Sample Title6", "6"}
	first_node.Data["Sample7"] = &Article{"Sample Title7", "7"}
	first_node.Data["Sample8"] = &Article{"Sample Title8", "8"}
	first_node.Data["Sample9"] = &Article{"Sample Title9", "9"}
	first_node.Data["Sample10"] = &Article{"Sample Title10", "10"}
	first_node.Data["Sample12"] = &Article{"Sample Title11", "11"}
	first_node.Data["Sample13"] = &Article{"Sample Title12", "12"}
	first_node.Data["Sample14"] = &Article{"Sample Title13", "13"}
	first_node.Data["Sample15"] = &Article{"Sample Title14", "14"}

	first_node.finger_gen()
	periodic(first_node)

	RPCListen(master_IP, first_node)
}

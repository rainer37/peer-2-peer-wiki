package chord

import(
	"crypto/sha1"
	"math"
	"fmt"
	"net/rpc"
	"net"
	"time"
	"os"
	"errors"
	"strings"
  	"github.com/nickbradley/p2pwiki/article"
  	"bitbucket.org/bestchai/dinv/dinvRT"
)

const m int = 6 // degree of max number of nodes supporting
const r int = 2 // degree of Fault Tolerance

var master_IP string = "127.0.0.1:1338"

type Node struct {
	IP string
	ID uint64
	Successor_IP string
	Predecessor_IP string
	Fingers [m]Finger
	Data map[string]*article.Article
	Suc_list [r]Finger
}

type Finger struct {
	FID uint64
	FIP string
}

type ArtPair struct {
	Title string
	Log article.OpLog
	Vec []byte
}

//万能struct
type Args struct {
	A string
	B string
}

type StrLog struct {
	Str string
	Log []byte
}

type ArtLog struct {
	Art *article.Article
	Log []byte
}

func unpack(log []byte) {
	var dummy []byte
	dinvRT.Unpack(log, &dummy)
}

// Hash the key into ID .
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

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

// general rpc
func RPCall(ip string, args interface{}, reply interface{}, f string) error {
  client, err := rpc.Dial("tcp", ip)
  if err != nil {
    //fmt.Println(err)
    return err
  }

  err = client.Call(f, args, reply)

  if err != nil {
    //println(err)
    return err
  }

  return nil
}

func CheckErr(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

func RPCListen(ip string, n *Node) {
  addr, err := net.ResolveTCPAddr("tcp", ip)
  CheckErr(err)

  inbound, err := net.ListenTCP("tcp", addr)
  CheckErr(err)

  rpc.Register(n)
  rpc.Accept(inbound)
}

func NewNode(ip string) *Node {
	return &Node{
		IP:		ip,
		ID:    	Hash(ip),
		Data:	make(map[string]*article.Article),
	}
}

func (n *Node) PushArticle(pair *ArtPair, log *[]byte) error {
	println("Client update articles....")

	unpack(pair.Vec)

	art,err := article.OpenArticle("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/", pair.Title)

	if err != nil {
		art = article.NewArticle(pair.Title)
	}

	art.Replay(pair.Log)
	art.Save("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/")
	//art.Print()
	n.Data[pair.Title] = art

	*log = dinvRT.Pack(nil)
	return nil
}

func (n *Node) Pull(title *StrLog, art *ArtLog) error {
	println("Pulling Article:", title.Str)

	unpack(title.Log)

	if n.Data[title.Str] == nil {
		return errors.New("No Such Article Exception")
	}

	a,err := article.OpenArticle("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/", n.Data[title.Str].Title)

	if err != nil {
		fmt.Println("Here: ",err)
		return err
	}

	art.Art = a
	art.Log = dinvRT.Pack(nil)

	println("@",art.Art.Title," pulled...")
	return nil
}

func (n *Node) Push(art *ArtLog, log *[]byte) error {

	//unpack(art.Log)
	//println("Putting", art.Title)

	if n.Data[art.Art.Title] == nil {
		n.Data[art.Art.Title] = art.Art
	}else{
		n.Data[art.Art.Title] = art.Art
	}

  	if exist,_ := exists("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/"); !exist {
 		os.MkdirAll("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/", 0777)
  	}

	art.Art.Save("./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/")

	//*log = dinvRT.Pack(nil)

	return nil
}

func (n *Node) Report(args *Args, reply *int) error {
	println("Node:",n.ID,"@",n.IP,"Data Lenth:",len(n.Data),"Suc:",Hash(n.Successor_IP),"Pred:",Hash(n.Predecessor_IP))
	return nil
}

// check if the key is in between
func is_key_in_between(kip uint64, pip uint64, sip uint64) bool {
	if sip > pip {
		if kip > pip && kip <= sip {
			return true
		}
	} else {
		if kip > pip || (kip+63) <= (sip+63){
			return true
		}
	}
	return false
}

// Notify my successor that i'm now your predecessor
func (n *Node) Notify(new_ip *StrLog, ret *[]byte) error {
	//println("Notified...")

	// no need to check if the new_ip is in correct range, Find already ensured that.
	println("Now my pred is", new_ip.Str, Hash(new_ip.Str))
	n.Predecessor_IP = new_ip.Str

	unpack(new_ip.Log)

	// Now start to divide the keys and make replications

	for i,v := range n.Data {
		if is_key_in_between(Hash(i), n.ID, Hash(new_ip.Str)) {
			//println("Key to give:", Hash(i))
			var a []byte

			// TODO: check error
			//RPCall(new_ip.Str, &ArtLog{v,dinvRT.Pack(nil)}, &a, "Node.Push")
			RPCall(new_ip.Str, &ArtLog{v,nil}, &a, "Node.Push")

			//unpack(a)

			delete(n.Data, i)
		}
	}

	*ret = dinvRT.Pack(nil)

	return nil
}


// TODO: check correctness
// when the node joining, reply the ip it should join to
func (n *Node) Find(new_ip *StrLog, ip_to_join *StrLog) error {
	println(new_ip.Str, "is joining...")

	unpack(new_ip.Log)
	// if i'm the only guy in the chord
	if n.IP == n.Successor_IP {
		println("1.It should join on me....")
		ip_to_join.Str = n.IP
		ip_to_join.Log = dinvRT.Pack(nil)
		return nil
	}

	// if i'm not alone, start the searching
	// Linear search for now
	// TODO: finger search

	// TWO CASES ON FIND
	// CASE ONE: n.ip < n.pred_ip
	// CASE TWO: n.ip >= n.pred_ip

	pip := Hash(n.Predecessor_IP)
	nip := Hash(new_ip.Str)

	if is_key_in_between(nip, pip, n.ID) {

			println("2.It should join on me....")
			ip_to_join.Str = n.IP
			ip_to_join.Log = dinvRT.Pack(nil)
			return nil

	} else {

			// First try the Finger search
			/*
			for i := 1; i<m; i++ {
				if n.Fingers[i].FID > nip {
					println("USING FINGER", i-1)
				}
			}
			*/
			RPCall(n.Successor_IP, &StrLog{new_ip.Str,dinvRT.Pack(nil)} , ip_to_join ,"Node.Find")
	}

	return nil
}

func (n *Node) Find_successor(nid *uint64, result *string) error {
	//println(*nid, "is what i look for...")

	if n.IP == n.Successor_IP {
		*result = n.IP
		return nil
	}

	pip := Hash(n.Predecessor_IP)
	nip := *nid

	if is_key_in_between(nip, pip, n.ID) {

			*result = n.IP
			return nil

	} else {

			RPCall(n.Successor_IP, nid , result ,"Node.Find_successor")
	}

	return nil
}

func (n *Node) Find_article(title *StrLog, ip *StrLog) error {
	println("Looking for article:", title.Str, Hash(title.Str))

	unpack(title.Log)

	// if i'm the only guy in the chord
	if n.IP == n.Successor_IP {
		ip.Str = n.IP
		ip.Log = dinvRT.Pack(nil)
		return nil
	}

	//RETURN THE CACHE ONE if THERE IS ONE IN THE LOOKUP PATH

	cacheDir := "./articles/cache/"+n.IP[strings.Index(n.IP, ":")+1:]+"/"
	println(cacheDir)

	_,err := article.OpenArticle(cacheDir, title.Str)

	if err == nil {
		ip.Str = n.IP
		ip.Log = dinvRT.Pack(nil)
		return nil
	}

	pip := Hash(n.Predecessor_IP)
	nip := Hash(title.Str)

	if is_key_in_between(nip, pip, n.ID) {

			ip.Str = n.IP
			ip.Log = dinvRT.Pack(nil)
			return nil

	} else {

			// finger search
			j := 0;
			for i:=0; i<m; i++ {
				if n.Fingers[i].FID >= nip {
					j = i - 1
					break
				}
			}
			//println("Cool i", j)

			if j == -1 {
				RPCall(n.Successor_IP, &StrLog{title.Str,dinvRT.Pack(nil)} , ip ,"Node.Find_article")
			} else {
				RPCall(n.Fingers[j].FIP, &StrLog{title.Str,dinvRT.Pack(nil)}, ip ,"Node.Find_article")
			}

			unpack(ip.Log)
			ip.Log = dinvRT.Pack(nil)
			//RPCall(n.Successor_IP, title , ip ,"Node.Find_article")
	}
	return nil
}

// check if my successor's pred is still me
// if not, notify the new_ip
func (n *Node) stablize() {

	// if no one in chord, no check then
	if n.Predecessor_IP == "" {
		println("I'm alone...")
		return
	}

	var succ_pred_ip string

	// TODO: check err
	RPCall(n.Successor_IP, &succ_pred_ip, &succ_pred_ip, "Node.GetPred")

	//println("My successor's pred is:", succ_pred_ip, "@")

	if succ_pred_ip != n.IP {

		println("Time to stablize....")
		var ret []byte
		RPCall(succ_pred_ip, &StrLog{n.IP,dinvRT.Pack(nil)}, &ret, "Node.Notify")
		unpack(ret)
		n.Successor_IP = succ_pred_ip

	}
}

func (n *Node) GetSuccessor(args *Args,reply *string) error {
	*reply = n.Successor_IP
	return nil
}

func (n *Node) GetPred(args *string,reply *string) error {
	//println("Hey pred is", n.Predecessor_IP)
	*reply = n.Predecessor_IP
	return nil
}

// the replicate will replicate my own pieces (excluding the pieces from my predessor) to my successor
func (n *Node) replicate() {

	if n.Successor_IP == n.IP { return }

	for i,v := range n.Data {
		if is_key_in_between(Hash(i), Hash(n.Predecessor_IP), n.ID){
			var a []byte

			println("Replicating", Hash(i))
			// TODO: check error
			RPCall(n.Successor_IP, &ArtLog{v,nil}, &a, "Node.Push")
		}
	}
}

// check if my successor is still alive by pinging it.
func (n *Node) ping() {
	var next string = n.IP
	var arg Args
	err := RPCall(n.Successor_IP,&arg,&next,"Node.GetSuccessor")

	if err != nil {
		fmt.Println("PING ERR: ",err)
		println("SUC gone...notify the next available agent...")
		n.reconcil()
	}
}

// find the next available successor in to system.
// if not, set the successro as myself and predecessor as nil
func (n *Node) reconcil() {
	var ret int = 0
	for i := 1; i < r; i++ {

		if n.Suc_list[i].FIP == n.IP {
			n.Predecessor_IP = ""
			n.Successor_IP = n.IP
			n.clear_sucs()
			return
		}
		err := RPCall(n.Suc_list[i].FIP, &n.IP, &ret, "Node.Notify")

		if err == nil {
			println("FOUND!!! my new successor is: ", n.Suc_list[i].FIP)
			n.Successor_IP = n.Suc_list[i].FIP
			return
		}
	}
}

// clear the successor list for reseting
func (n *Node) clear_sucs() {
	for i := 0; i < r; i++ {
		n.Suc_list[i].FIP = n.IP
	}
}

// generate the bounds in finger table.
func (n *Node) finger_gen(){
	for i := 0; i < m; i++ {
		n.Fingers[i].FID = (n.ID + uint64(math.Pow(2,float64(i)))) % uint64(math.Pow(2,float64(m)))
		//println("Finger",i,n.Fingers[i].FID)
		n.Fingers[i] = Finger{n.ID, n.IP}
	}
}

func get_finger_bound(base uint64, index int) uint64 {
	return (base + uint64(math.Pow(2,float64(index)))) % uint64(math.Pow(2,float64(m)))
}

// generate ips in finger table
func (n *Node) fix_fingers() {

	var result string = ""

	for i := 0; i < m; i++ {

		bound := get_finger_bound(n.ID, i)
		err := RPCall(n.Successor_IP, &bound, &result, "Node.Find_successor")

		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println("FINGER",i, result)
		n.Fingers[i] = Finger{Hash(result), result}

	}

	//println()
}

// generate the successor list by recursively linear checking.
func (n *Node) fix_successors() {

	var next string = n.IP
	var arg Args
	for i := 0; i < r; i++ {
		err := RPCall(next,&arg,&next,"Node.GetSuccessor")

		// if the successor is gone
		if err != nil {
			fmt.Println("FS ERR: ",err)
		} else {
			n.Suc_list[i] = Finger{Hash(next),next}
			//println("NEXT: ", next)
		}
	}
}

// periodic operations that every node does
// 1. check if the successor's predessor is myself to detect new node join
// 2. replicate the keys to successors.
// 3. generate a brief report
func periodic(n *Node) {
	ticker := time.NewTicker(time.Millisecond * 3000)
    go func() {
        for range ticker.C {

        	//@dump
	        //dinvRT.Dump(n.IP+":Successor&predessor", "S,P", n.Successor_IP,n.Predecessor_IP)

        	n.ping()
            n.stablize()
            n.replicate()
            n.fix_successors()
            n.fix_fingers()
            quick_report(n)

            for i := 0; i < m; i++ {
				//println("Finger",i,n.Fingers[i].FID)
			}
        }
    }()
}

func quick_report(n *Node) {

	print("IM: ||", n.ID, "|| ")
	for i,_ := range n.Data {
		print(Hash(i)," ")
	}
	if n.Predecessor_IP == "" {
		print("Pred: nil", " Succ: ",Hash(n.Successor_IP))
	}else{
		print("Pred: ", Hash(n.Predecessor_IP), " Succ: ",Hash(n.Successor_IP))
	}
	println()
}

func (n *Node) Join(to_join string) {
	println(n.IP," || ",Hash(n.IP)," ||", "is Joining...")
	//n.finger_gen()

	dinvRT.Initalize(n.IP)

	var ip_to_join StrLog

	err := RPCall(to_join, &StrLog{n.IP,dinvRT.Pack(nil)}, &ip_to_join, "Node.Find") // findSuccessor()

	unpack(ip_to_join.Log)

	if err != nil {
		println("Failing to join any nodes...")
		return
	}

	fmt.Println("I Will join",ip_to_join.Str,Hash(ip_to_join.Str))

	n.Successor_IP = ip_to_join.Str
	n.Predecessor_IP = ""
	n.finger_gen()

	periodic(n)

	go RPCListen(n.IP, n)
	var ret []byte

    err = RPCall(n.Successor_IP, &StrLog{n.IP, dinvRT.Pack(nil)}, &ret, "Node.Notify")
    CheckErr(err)

    unpack(ret)

	for{
		// loop forever
	}
}


func CreateRing(srvAddr string) {

	master_IP = srvAddr
	first_node := NewNode(master_IP)
	first_node.Successor_IP = master_IP
	first_node.Predecessor_IP = ""

	dinvRT.Initalize("MASTER")

	//println("H:", Hash("A1"))
	first_node.finger_gen()

	a1 := article.NewArticle("A1")
	a1.Insert(1, "B", "author")
	first_node.Data["A1"] = a1

	a2 := article.NewArticle("A2")
	a2.Insert(1, "B", master_IP)
	first_node.Data["A2"] = a2

	a3 := article.NewArticle("A3")
	a3.Insert(1, "First Sentence in A3", master_IP)
	first_node.Data["A3"] = a3

	a4 := article.NewArticle("A4")
	a4.Insert(1, "First Sentence in A4", master_IP)
	first_node.Data["A4"] = a4

	a7 := article.NewArticle("A7")
	a7.Insert(1, "First Sentence in A7", master_IP)
	first_node.Data["A7"] = a7

	a6 := article.NewArticle("A6")
	a6.Insert(1, "First Sentence in A6", master_IP)
	first_node.Data["A6"] = a6

  	if exist,_ := exists("./articles/cache/1338/"); !exist {
 		os.MkdirAll("./articles/cache/1338/", 0777)
  	}

  	a1.Log = article.OpLog{}

	a1.Save("./articles/cache/1338/");
	a2.Save("./articles/cache/1338/");
	a3.Save("./articles/cache/1338/");
	a4.Save("./articles/cache/1338/");
	a7.Save("./articles/cache/1338/");
	a6.Save("./articles/cache/1338/");

	periodic(first_node)

	RPCListen(master_IP, first_node)
}

package p2pwiki


type ArticleLookup struct {
  articles []*Article
}



type Peer struct {

}

func (this *Peer) locate(n int) []string {

}

// Replicate an article to a peer.
// Calls the remote peer's Article.Update method using RPC.
func replicate(article *Article, address string) (err Error) {
  artRep := article.replicate()

  msg := artReg.Marshal()

  // send msg to address
}

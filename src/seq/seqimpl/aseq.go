package seqimpl

// import "seq"

// type ASeq struct {
// 	AMeta
// 	hash int32
// }

// func (a *ASeq) String() string {
// 	return "TODO"
// }

// func (a *ASeq) Seq() seq.Seq {
// 	return a
// }

// func (a *ASeq) Equals(o interface{}) bool {
// 	if a == o {
// 		return true
// 	}
	
// 	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
// 	os, ok := o.(seq.Seqable)

// 	if !ok {
// 		return false
// 	}


// 	ms := os.Seq()
// 	for s := a.Seq(); s != nil; s = s.Next() {
// 		if ms == nil || ! sequtils.Equals(s.First(),ms.First()) {
// 			return false;
// 		}
// 		ms = ms.Next() 
// 	}

// 	return ms == nil;
// }



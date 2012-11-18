// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
)

func moreFromSeq(s iseq.Seq) iseq.Seq {
	sn := s.Next()
	if sn == nil {
		return CachedEmptyList
	}
	return sn
}

func smartCons(x, coll interface{}) iseq.Seq {

	if coll == nil {
		return NewPList1(x)
	}

	if s, ok := coll.(iseq.Seq); ok {
		return NewCons(x, s)
	}

	return NewCons(x, sequtil.ConvertToSeq(coll))
}

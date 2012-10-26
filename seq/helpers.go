// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
)

func moreFromSeq(s iseq.Seq) iseq.Seq {
	sn := s.Next()
	if sn == nil {
		return CachedEmptyList
	}
	return sn
}
package now

import (
	"time"
)

var JST = time.FixedZone("JST", 9*60*60)

var Now = func() time.Time {
	n := time.Now()
	nJST := n.In(JST)
	return nJST
}

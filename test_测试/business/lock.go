package business

import (
	"sync"
)

type Ldata struct {
	sync.RWMutex
	s string
}

// get data atomically
func (l *Ldata) Data() string {
	l.Lock()
	defer l.Unlock()

	return l.s
}

// set data atomically
func (l *Ldata) SetData(d string) {
	l.Lock()
	defer l.Unlock()
	l.s = d
}

//func main() {
//	var wg sync.WaitGroup
//	wg.Add(200)
//	var ss Ldata
//	for range [100]struct{}{} {
//		go func() {
//			time.Sleep(time.Second * time.Duration(rand.Intn(1000)) / 1000)
//
//			//log.Println(ss.Data())
//			wg.Done()
//		}()
//	}
//
//	for i := range [100]struct{}{} {
//		go func(i int) {
//			time.Sleep(time.Second * time.Duration(rand.Intn(1000)) / 1000)
//			s := fmt.Sprint("#", i)
//			//log.Println("====", s)
//
//			ss.SetData(s)
//			wg.Done()
//		}(i)
//	}
//
//	wg.Wait()
//
//}

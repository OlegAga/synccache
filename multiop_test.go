package synccache

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMultiGetsSetsRemovesAndUpdates(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	c := New(0, 0, "")
	keys := []string{"a", "b", "c"}
	wg := sync.WaitGroup{}
	for i := 0; i < 40; i++ {
		wg.Add(1)
		v := i
		op := rand.Intn(4)
		key := keys[rand.Intn(len(keys))]
		if op == 0 {
			go func(x int) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Millisecond)
				v, _ := c.Get(key)
				fmt.Printf("Get key %s value %v in %s \n", key, v, time.Now())
			}(i)
		} else if op == 1 {
			go func(x int) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Millisecond)
				c.Set(key, v, 5*time.Second)
				fmt.Printf("Set key %s value %v in %s \n", key, v, time.Now())
			}(i)
		} else if op == 2 {
			go func(x int) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Millisecond)
				c.Remove(key)
				fmt.Printf("Remove key %s in %s \n", key, time.Now())
			}(i)
		} else {
			go func(x int) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Millisecond)
				c.Update(key, v)
				fmt.Printf("Update key %s value %v in %s \n", key, v, time.Now())
			}(i)
		}
	}
	wg.Wait()
	items := strings.Split(c.Keys(), ",")
	for _, k := range items {
		v, _ := c.Get(k)
		fmt.Printf("Key %s = %v\n", k, v)
	}

}

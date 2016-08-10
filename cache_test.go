package synccache

import (
	"os"
	"testing"
	"time"
)

func TestSimpleGetsAndSets(t *testing.T) {
	c := New(0, 0, "")

	a, err := c.Get("a")
	if err == nil || a != nil {
		t.Error("Getting A err value that shouldn't eresultist:", a)
	}

	c.Set("a", "a", 0)
	c.Set("b", 1, 0)
	c.Set("c", 1.1, 0)

	result, err := c.Get("a")
	if err != nil {
		t.Error("a was not err while getting a")
	}
	if result == nil {
		t.Error("result for a is nil")
	} else if result != "a" {
		t.Error("a does not equal \"a\"; value:", result)
	}

	result, err = c.Get("b")
	if err != nil {
		t.Error("b was not err while getting b")
	}
	if result == nil {
		t.Error("result for b is nil")
	} else if result != 1 {
		t.Error("b does not equal 1; value:", result)
	}

	result, err = c.Get("c")
	if err != nil {
		t.Error("c was not err while getting c")
	}
	if result == nil {
		t.Error("result for c is nil")
	} else if result != 1.1 {
		t.Error("c does not equal 1.1; value:", result)
	}
}

func TestSaveDbToFile(t *testing.T) {
	c := New(0, 1*time.Second, "test.db")

	c.Set("a", "a", 0)

	time.Sleep(3 * time.Second)

	_, err := os.Open("test.db")
	if err != nil {
		t.Error("File not create after save")
	}
}

func TestUpDbFromFile(t *testing.T) {
	c := New(0, 0, "")

	c.Load("test.db")

	result, err := c.Get("a")
	if err != nil {
		t.Error("a was not err while getting a after load")
	}
	if result == nil {
		t.Error("result for a is nil after load")
	} else if result != "a" {
		t.Error("a does not equal \"a\"; value:", result, " after load")
	}
}

func TestExpiriedKeyAndCleanupDb(t *testing.T) {
	c := New(200*time.Millisecond, 0, "")

	c.Set("a", "a", 100*time.Millisecond)

	result, err := c.Get("a")
	if err != nil {
		t.Error("a was not err while getting a after 0 sec")
	}

	time.Sleep(101 * time.Millisecond)
	result, err = c.Get("a")
	if result != nil || err == nil {
		t.Error("a not expected as expired", err)
	}

	time.Sleep(201 * time.Millisecond)
	_, err = c.Get("a")
	if result != nil || err == nil {
		t.Error("a not expected as removed GC", err)
	}
}

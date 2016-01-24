package minpq

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMinPQ(t *testing.T) {
	Convey("Test the empty queue.", t, func() {
		pq := New(2, 0)
		So(pq, ShouldNotBeNil)
		So(pq.Contains("foo"), ShouldBeFalse)
		So(pq.Empty(), ShouldBeTrue)
		_, ok := pq.CurrentPriority("foo")
		So(ok, ShouldBeFalse)
		So(pq.Size(), ShouldBeZeroValue)
		So(pq.UpdatePriority("foo", 1), ShouldNotBeNil)
		val, _ := pq.Pop()
		So(val, ShouldBeNil)
	})

	Convey("Test insertion.", t, func() {
		pq := New(2, 10)
		pq.Insert("foo", 1)
		So(pq.Empty(), ShouldBeFalse)
		So(pq.Size(), ShouldEqual, 1)
		So(pq.Contains("foo"), ShouldBeTrue)
		pri, ok := pq.CurrentPriority("foo")
		So(ok, ShouldBeTrue)
		So(pri, ShouldEqual, 1)
	})

	Convey("Test pop.", t, func() {
		pq := New(2, 10)
		pq.Insert("foo", 1)
		pq.Insert("bar", 0)
		val, pri := pq.Pop()
		So(val.(string), ShouldEqual, "bar")
		So(pri, ShouldEqual, 0)
		val, pri = pq.Pop()
		So(val.(string), ShouldEqual, "foo")
		So(pri, ShouldEqual, 1)
		So(pq.Contains("foo"), ShouldBeFalse)
		So(pq.Contains("bar"), ShouldBeFalse)
	})

	Convey("Test update priority.", t, func() {
		pq := New(2, 10)
		pq.Insert("foo", 1)
		pq.Insert("bar", 0)
		So(pq.UpdatePriority("bar", 2), ShouldBeNil)
		pri, ok := pq.CurrentPriority("bar")
		So(ok, ShouldBeTrue)
		So(pri, ShouldEqual, 2)
		val, pri := pq.Pop()
		So(val, ShouldNotBeNil)
		So(val.(string), ShouldEqual, "foo")
		So(pri, ShouldEqual, 1)
	})

	Convey("Test update priority with a larger input.", t, func() {
		pq := New(2, 1000)
		for i := 0; i < 1000; i++ {
			pq.Insert(999-i, float32(999-i))
		}
		for i := 0; i < 1000; i++ {
			err := pq.UpdatePriority(i, float32(i+1))
			So(err, ShouldBeNil)
		}
		for i := 0; i < 1000; i++ {
			val, pri := pq.Pop()
			So(pri, ShouldEqual, i+1)
			So(val.(int), ShouldEqual, i)
		}

		pq = New(2, 1000)
		for i := 0; i < 1000; i++ {
			pq.Insert(999-i, float32(999-i))
		}
		for i := 0; i < 1000; i++ {
			err := pq.UpdatePriority(i, float32(i-1))
			So(err, ShouldBeNil)
		}
		for i := 0; i < 1000; i++ {
			val, pri := pq.Pop()
			So(pri, ShouldEqual, i-1)
			So(val.(int), ShouldEqual, i)
		}
	})
}

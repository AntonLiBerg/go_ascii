# Vars
name := value
var name = value

# functions
func name(params) (return,return..) {

}

# slices (list)
name := []type{}

# map
name := map[typekey]typeval{}
aval,okbool := name[key]

# struct
type NameStruct struct {
    name type
    ..
}
name := NameStruct{prop:val..}

# methods (funcs on type)
func (u StructName) NameFunc(){

}

# pointer
name := val
namePointer := *name

# errors
errors.New("msg")

# interface
type Name Interface{
    NameFunc() returnType
}
# channel
ch := make(chan returntype, storage)
blocks when full

ch <- write
read <-ch
close(ch)

# Select

select waits on multiple channel operations.

select {
case msg := <-ch:
    fmt.Println("received:", msg)
case err := <-errCh:
    fmt.Println("error:", err)
case <-time.After(1 * time.Second):
            fmt.Println("timed out")
    
}

It is like switch, but for channels.


# go routines (async)


go func() {
    ch <- "msg"
}


# testing

import "testing"



That shape matters. Go’s test runner finds functions named TestXxx that take *testing.T.
  t.Fatalf fails the test immediately. t.Errorf reports a failure but keeps running the test.
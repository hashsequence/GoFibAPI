package main

import(
  "strings"
  "sync"
  "net/http"
)
/********************************************************************
*Server Structs
*
********************************************************************/

//struct design for our Api Server
type FibServer struct {
    Server *http.Server
		mux *http.ServeMux
}

//config format
type Configuration struct {
	Address      string  //address of the url set to http://0.0.0.0:8080
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

/********************************************************************
*Handler Structs
*
********************************************************************/
//json struct that will be sent to the client with the current fibonacci number
type FibNum struct {
  CurrNum string
}

//empty structs for the handlers for current, next, previous, reason for this implementation is so that I can wrap each
//handler in the recovery wrapper to handle panics
type CurrentHandler struct {}
type NextHandler struct {}
type PreviousHandler struct {}

//instantiates handlers
func NewCurrentHandler() *CurrentHandler {
  return &CurrentHandler{}
}

func NewNextHandler() *NextHandler {
  return &NextHandler{}
}

func NewPreviousHandler() *PreviousHandler {
  return &PreviousHandler{}
}


/********************************************************************
*LargeInt
*Summary: stores integers in an array to store numbers greater than int64
*********************************************************************/
type LargeInt struct {
  	Val string
}

//instantiates struct that will store the number in a string
func NewLargeInt(vals ...string) LargeInt {
	str := ""
	if len(vals) == 0 {
		str = "0"
	} else {
		str = vals[0]
	}

	return LargeInt{str}
}

//adding function : adds two LargeInt and returns the value in a string
func (ls LargeInt) Add(o *LargeInt) string {
	num1, num2 := ls.Val, o.Val
	res := ""
	carry := byte(0)
	i := len(num1) - 1
	j := len(num2) - 1

	for i >= 0 || j >= 0 {
		sum := byte(0)
		if i >= 0 {
			sum += num1[i] - '0'
		}
		if j >= 0 {
			sum += num2[j] - '0'
		}
		sum += carry
		carry = sum / 10
		n := sum % 10
		res = string(n+'0') + res
		i--
		j--
	}
	if carry > 0 {
		res = string(carry+'0') + res
	}
	return res
}

//helper functions to implement subtract
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func IsSmaller(str1, str2 string) bool {
	n1, n2 := len(str1), len(str2)
	if n1 < n2 {
		return true
	}
	if n1 > n2 {
		return false
	}

	for i, _ := range str1 {
		if str1[i] < str2[i] {
			return true
		} else if str1[i] > str2[i] {
			return false
		}
	}
	return false
}

//Subtract function, subtract to LargeInts, handles negative numbers too
func (ls LargeInt) Subtract(o *LargeInt) string {
	var str1 string
	var str2 string
	isSmallerFlag := IsSmaller(ls.Val, o.Val)
	if isSmallerFlag {
		str1, str2 = o.Val, ls.Val

	} else {
		str1, str2 = ls.Val, o.Val
	}

	str := ""
	n1, n2 := len(str1), len(str2)
	diff := n1 - n2
	carry := 0

	for i := n2 - 1; i >= 0; i-- {
		sub := int(str1[i+diff]-'0') - int(str2[i]-'0') - carry
		if sub < 0 {
			sub = sub + 10
			carry = 1
		} else {
			carry = 0
		}
		str = str + string(byte(sub)+'0')
	}

	for i := n1 - n2 - 1; i >= 0; i-- {
		if str1[i] == '0' && carry > 0 {
			str = str + "9"
			continue
		}

		sub := int(str1[i]-'0') - carry
		if i > 0 || sub > 0 {
			str = str + string(byte(sub)+'0')
		}
		carry = 0
	}
	str =  strings.TrimLeft(Reverse(str),"0")
	if str == "" {
		str = "0"
	}
	if isSmallerFlag {
		str = "-" + str
	}
	return str

}

//returns number as string from our LargeInt struct
func (ls LargeInt) GetVal() string {
  return ls.Val
}

/********************************************************************
*ConcurrArrOfLargeInt
*Summary: synchronous array to store LargeInt used for our data store
*********************************************************************/

type ConcurrArrOfLargeInt struct{
  sync.RWMutex
  arr []LargeInt
}

func NewConcurrArr() *ConcurrArrOfLargeInt {
  c := &ConcurrArrOfLargeInt{arr : []LargeInt{}}
  return c
}


//sets the number at index of our array, is synchronous
func (c *ConcurrArrOfLargeInt) Set(key int, value LargeInt) {
  c.Lock()
  defer c.Unlock()
  if len(c.arr) == 0 {
    c.arr = append(c.arr, value)
    return
  }
  if key < len(c.arr) && key >= 0 {
    c.arr[key] = value
    return
  }
  if key == len(c.arr) {
    c.arr = append(c.arr, value)
    return
  }
}

//same implementation as Set except it does not uses lock, reason is that it is use in the
//ShiftForward and ShiftBackward Function to change the values in our array, and the
//locks are used there
func (c *ConcurrArrOfLargeInt) set(key int, value LargeInt) {
  if len(c.arr) == 0 {
    c.arr = append(c.arr, value)
    return
  }
  if key < len(c.arr) && key >= 0 {
    c.arr[key] = value
    return
  }
  if key == len(c.arr) {
    c.arr = append(c.arr, value)
    return
  }
}

//Getter function and returns LargeInt at index
func (c* ConcurrArrOfLargeInt) Get(key int) LargeInt {
  c.Lock()
  defer c.Unlock()
  if key < len(c.arr) && key >= 0 {
    return c.arr[key]
  }
  return NewLargeInt("nil")
}

//same implementation of get but does not uses lock
func (c* ConcurrArrOfLargeInt) get(key int) LargeInt {

  if key < len(c.arr) && key >= 0 {
    return c.arr[key]
  }
  return NewLargeInt("nil")
}

//shifts the sequence forward. If the sequence is 0, 1, 1 then it will become 1,1,2 with 2 being the current
//fibonacci number
func (c* ConcurrArrOfLargeInt) ShiftForward() string{
  c.Lock()
  defer c.Unlock()
  if c.get(1).GetVal() == "nil" && c.get(0).GetVal() == "nil" {
    c.set(1,NewLargeInt(c.get(2).GetVal()))
    c.set(2,NewLargeInt("1"))
  } else {
    c.set(0,NewLargeInt(c.get(1).GetVal()))
    c.set(1,NewLargeInt(c.get(2).GetVal()))
    val1 := NewLargeInt(c.get(1).GetVal())
    val2 := NewLargeInt(c.get(0).GetVal())
    c.set(2,NewLargeInt(val1.Add(&val2)))
  }

  if len(c.get(2).GetVal()) >=  MAX_LEN {
    c.set(0,NewLargeInt("nil"))
    c.set(1,NewLargeInt("nil"))
    c.set(2,NewLargeInt("0"))
  }
  return c.get(2).GetVal()

}
//shifts the sequence backward so 1,2,3 --> 1,1,2
func (c* ConcurrArrOfLargeInt) ShiftBackward() string {
  c.Lock()
  defer c.Unlock()
  if c.get(2).GetVal() == "0" {
      c.set(1,NewLargeInt("nil"))
  } else if c.get(2).GetVal() == "1" && c.get(1).GetVal() == "0" {
    c.set(2,NewLargeInt("0"))
    c.set(1,NewLargeInt("nil"))
  } else if c.get(2).GetVal() == "1" && c.get(1).GetVal() == "1" {
    c.set(1,NewLargeInt("0"))
    c.set(0,NewLargeInt("nil"))
  } else {
    c.set(2,NewLargeInt(c.get(1).GetVal()))
    c.set(1,NewLargeInt(c.get(0).GetVal()))
    val1 := NewLargeInt(c.get(2).GetVal())
    val2 := NewLargeInt(c.get(1).GetVal())
    c.set(0,NewLargeInt(val1.Subtract(&val2)))
  }
  if c.get(2).GetVal()[0:1]  == "-" {
    c.set(0,NewLargeInt("nil"))
  	c.set(1,NewLargeInt("nil"))
  	c.set(2,NewLargeInt("0"))
  }
  return c.get(2).GetVal()

}

//Resets the fibonacci Sequence
func (c* ConcurrArrOfLargeInt) Reset() string{
  c.Lock()
  defer c.Unlock()
  c.set(0,NewLargeInt("nil"))
  c.set(1,NewLargeInt("nil"))
  c.set(2,NewLargeInt("0"))
  return c.get(2).GetVal()

}

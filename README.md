Go Fib Api
===============================

Author: Avery Wong
DateCreated: 10/20/2019

I designed the Web Api using barebones Golang. I chose Go because it
is lightweight, fast, and makes it easy to program for concurrency.

The Api handles the 3 endpoints:

**current**
  Returns Json with the current Fibonacci number.

* **URL**
  /current

* **Method:**
`GET`

* **Data Params**
  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content:** `{ CurrNum : "5" }`

* **Error Response:**

    * **Code:** 500 StatusInternalServerError <br />

**next**
  Returns Json with the next Fibonacci number.
  If the number reaches MAX_LEN(set to 100 digits), then the number resets back to zero.

* **URL**
  /next

* **Method:**
`GET`

* **Data Params**
  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content:** `{ CurrNum : "8" }`

* **Error Response:**

    * **Code:** 500 StatusInternalServerError <br />

**previous**
  Returns Json with the previous Fibonacci Number.
  numbers stays at 0 if current number is 0, since there is no previous number before 0 in
  the sequence.

* **URL**
  /previous

* **Method:**
`GET`

* **Data Params**
  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content:** `{ CurrNum : "3" }`

* **Error Response:**

    * **Code:** 500 StatusInternalServerError <br />

Purpose
==============================
* The purpose of this project was to showcase how to:
  * make a simple Api and be able to deploy it using dep onto heroku
  * I wanted to able to dockerize it and push it on to dockerhub
  * add nginx for ssl and reverse proxy endpoint (using self signed ssl certificates)


Design
===============================

I made an assumption that the Fibonacci
number is shared among clients, so if 3 clients calls next and the current was
0, then the callers will get 1,1,2 with respect to order

I had to implement my own data store to store the shared variable and I had to
make it concurrent:

```Go
type LargeInt struct {
  	Val string
}

type ConcurrArrOfLargeInt struct{
  sync.RWMutex
  arr []LargeInt
}
```

Each value is stored into a LargeInt struct that stores the number in a string,
since Int64 only contains a number up to 9223372036854775807
The ConcurrArrOfLargeInt is an array that is thread safe, achieved through
the use of Go's builtin mutex, so access to the current Fibonacci number
is consistent. I set the MAX_LEN of the value to be 100 digits, and will be reset to 0 if it goes over.

Also to deal with panics(errors in Golang) I wrapped each handler for the
endpoints in a recovery wrapper function which will reset the Fibonacci Sequence
to 0 and logs the error into a log and continues.

The Algorithm for Fibonaccis Sequence is :

F[n] = F[n-1] + F[n-2]
where F[0] = 0

so I store F[n-2] F[n-1] F[n] in an Array using my ConcurrArrOfLargeInt and
LargeInt to hold the value

F[n-2] is stored in FibArr[0]

F[n-1] is stored in FibArr[1]

F[n-2] is stored in FibArr[2]

```Go
var  FibArr *ConcurrArrOfLargeInt
```

I made two versions in this repository

**goApiRunLocally**
This folder holds the source code to run locally.
Simply execute the command in the bash terminal (I run Ubuntu 16.04):

```console
aw@ubuntu:~/goApiRunLocally$ go build -o fibApi
```

Then run the fibApi executable to start the api locally.

```console
aw@ubuntu:~/goApiRunLocally$ ./fibApi
fibApi started at : 0.0.0.0:8080
```

The config.json has the configuration for the address, which can be edited if you want to choose a different port.

sample calls on web browser:

[http://localhost:8080/current](http://localhost:8080/current)

[http://localhost:8080/previous](http://localhost:8080/previous)

[http://localhost:8080/next](http://localhost:8080/next)

**goApiDeployedOnHeroku**

The api is deployed to Heroku here:

[https://fast-fjord-47876.herokuapp.com/](https://fast-fjord-47876.herokuapp.com/)

It is slow due to being a free account.

sample calls on web browser:

[https://fast-fjord-47876.herokuapp.com/current](https://fast-fjord-47876.herokuapp.com/current)

[https://fast-fjord-47876.herokuapp.com/next](https://fast-fjord-47876.herokuapp.com/next)

[https://fast-fjord-47876.herokuapp.com/previous](https://fast-fjord-47876.herokuapp.com/previous)

**apiCall**

This folder holds a Go script to make request to my api, ran 1000 requests and
averages around 600 ms for 1000 requests on my computer. I used this for testing purposes

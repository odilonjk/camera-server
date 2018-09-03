A very basic server to rec and download videos from webcam.

* Go: **1.11**

Dependencies:

* mux
* gocv

----

How to run?
    
    go run main.go

Start recording:

    POST: localhost:3000/api/v1/record/start
    
Stop recording:

    POST: localhost:3000/api/v1/record/stop
    
Downlod last recorded video:

    GET: localhost:3000/api/v1/record
    
List videos:

    Not implemented yet.
    
Download video by name:

    Not implemented yet.
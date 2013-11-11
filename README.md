GoWeb
=====

This is a simple http server in go. To a large part it mimics the http.FileServer that comes with GO! 
It's a static file server that adds your handler function whenever it doesn't find the file specified
from the URL. 
On line 74, replace "nil, nil" with "test, start" or whatever function you have in mind.

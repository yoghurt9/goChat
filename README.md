# Chat Example

This application shows how to use the
[websocket](https://github.com/gorilla/websocket) package to implement a simple [private]()
web chat application.
This project refers to [https://github.com/gorilla/websocket/tree/master/examples/chat].

## Running the example

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

Once you have Go up and running, you can download, build and run the example
using the following commands.


    cd ./chat
    go get github.com/gorilla/websocket
    go run *.go

To use the chat example, open http://localhost:8080/ in your browser, and 
enter an ID. Note that the ID of each client cannot be the same.After 
entering the main page, you can enter the ID of another client to chat.

## Server

The server application defines two types, `Client` and `Router`. The server
creates an instance of the `Client` type for each websocket connection. A
`Client` acts as an intermediary between the websocket connection and a single
instance of the `Router` type. The `Router` maintains a set of registered clients 
and send messages to the specified client.

The application runs one goroutine for the `Router` and two goroutines for each
`Client`. The goroutines communicate with each other using channels. The `Router`
has channels for registering clients, unregistering clients and routing messages.
A `Client` has a buffered channel of outbound messages. 

### Router

The application's `main` function starts the router's `run` method as a goroutine.
Clients send requests to the router using the `register`, `unregister` and
`work` channels.

The router registers clients by adding id(you type in your browser) as a key in 
the `clients` map. The map value is the client pointer.

The unregister code is a little more complicated. In addition to deleting the
client pointer from the `clients` map, the router closes the clients's `send`
channel to signal the client that no more messages will be sent to the client.

The router obtains the client through the message's toId and sending the
message to the client's `send` channel. If the client's `send` buffer is full,
then the router assumes that the client is dead or stuck. In this case, the router
unregisters the client and closes the websocket.

### Client

The `serveWs` function is registered by the application's `main` function as
an HTTP handler. The handler upgrades the HTTP connection to the WebSocket
protocol, creates a client, registers the client with the router and schedules the
client to be unregistered using a defer statement.

Next, the HTTP handler starts the client's `writePump` method as a goroutine.
This method transfers messages from the client's send channel to the websocket
connection. The writer method exits when the channel is closed by the router or
there's an error writing to the websocket connection.

Finally, the HTTP handler calls the client's `readPump` method. This method
transfers inbound messages from the websocket to the router.

WebSocket connections. The application ensures that these concurrency requirements
are met by executing all reads from the `readPump` goroutine and all writes from 
the `writePump` goroutine.

To improve efficiency under high load, the `writePump` function coalesces
pending chat messages in the `send` channel to a single WebSocket message. This
reduces the number of system calls and the amount of data sent over the
network.

## Frontend

Please enter an ID on document load, but not the same as other clients.
If websocket functionality is available, then the script opens a connection to
the server and registers a callback to handle messages from the server. The
callback appends the message to the chat log using the appendLog function.

To allow the user to manually scroll through the chat log without interruption
from new messages, the `appendLog` function checks the scroll position before
adding new content. If the chat log is scrolled to the bottom, then the
function scrolls new content into view after adding the content. Otherwise, the
scroll position is not changed.

The form handler writes the user input to the websocket and clears the input
field.

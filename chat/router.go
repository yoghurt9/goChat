package main

type Router struct {
	// Registered clients, key is userId.
	clients map[string]*Client

	// Inbound messages from the clients.
	work chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newRouter() *Router {
	return &Router{
		work:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (r *Router) run()  {
	for {
		select {
		case client := <-r.register:
			r.regFunc(client)
		case client := <-r.unregister:
			r.unregFunc(client)
		case message := <-r.work:
			r.route(message)
		}
	}
}

// Register the client and add the client to the map
func (r *Router) regFunc(client *Client) {
	r.clients[client.userId] = client
}

// Unregisters the client, remove the client from the map, and closes the websocket connection.
func (r *Router) unregFunc(client *Client) {
	if _, ok := r.clients[client.userId]; ok {
		delete(r.clients, client.userId)
		close(client.send)
	}
}

// Send a message to the specified client
func (r *Router) route(msg *Message) {
	client, ok := r.clients[msg.ToId]
	if ok {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(r.clients, client.userId)
		}
	}
}
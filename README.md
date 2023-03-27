<h1 style="text-align: center">SSE Dispatcher</h1>

![alt text](./sse_dispatcher.jpeg)

<h2 style="text-align: left">What is this project</h2>
<p>This project solves the problem when scaling with docker-compose a backend with server-sent events.</p> 
<p>It is composed of a backend and a dispatcher.</p>
<p>The backend is a simple HTTP server with two endpoints and a gRPC server:</p>
<ul>
<li>GET &nbsp;&nbsp;/sse?client=[identification of the client]: creates a channel for SSE messages</li>
<li>POST /push-message: push message to clients</li>
<li> gRPC server: listen for calls by the dispatcher</li>
</ul>
<p>The dispatcher is a gRPC server that receives the messages from the backend and sends back to the backends that were scaled by docker compose.</p>

<h2 style="text-align: left">How to use this project</h2>
<p>1) Create a docker image of the backend</p>

```docker image build ./backend --tag backend:v1.0.0```

<p>2) Create a docker image of the dispatcher</p>

```docker image build ./dispatcher --tag dispatcher:v1.0.0```

<p>3) Run docker compose</p>

```docker-compose up -d --remove-orphans --scale backend=2```
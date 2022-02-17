
# client-side Chat App

#### Video Demo: <>

#### Description

You and your friends can chat with each other over this web app. No data (such as name, chats) is saved on a server.

  

Serve this app and give the URL to your friends, and you are good to go.

  
  

run :

` go run cmd/web/*.go `

  

The default port is 8080, but you can change it simply by this command before running the app.

  
  

` export PORT=<your desired port>`

  

##### Technology

The Back-end stack is Golang, and The Front-end stack is Bootstrap and JS.
We chose simple libraries for this project.

- [Pat Muxer](github.com/bmizerany/pat) is used for muxer.

- [Jet Template Engine](github.com/CloudyKit/jet/v6) for template engine. We used a template engine to make frontend development easier.

- [Gorilla Websocket](github.com/gorilla/websocket) for handling WebSocket. Websocket has real-time communication capability that makes it suitable for a chat app.
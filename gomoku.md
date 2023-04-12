# GoMoKu Document

```
.
├── client
│   ├── background.gif
│   ├── blackStone.gif
│   ├── whiteStone.gif
│   ├── chessboard.js // chessboard script
│   └── index.html    // display and senders & handlers
│
├── manager.go // manage clients as well as chessboard in each
├── client.go  // client inner state: chessboard, WebSocket connection, two go routines
├── room.go    // chessboard
│
├── event.go   // senders & handlers
└── main.go    // program entry
```

## System Structure

```
Server <--{WebSocket}--> Clients
```

## CS Communication

> Use WebSocket to inform client/server of changes

* role message -> C
* name message -> S & C
* move message -> S & C
* result message -> C

1. *client* enters game, waiting until match with another *client*
2. if two clients match, *server* send *role message* to the two *client*s
3. *client*s change name and send *name message* to *server* &
   *server* informs the *client* pair of the *name change*
4. *client*s move and send *move message* to *server* &
   *server* informs the *client* pair of the *move change*<br>
   meanwhile, *server* maintains inner chessboard state, if the move causes wining,
   *sever* sends *result message* to the *client* pair

## Design

* Support communication concurrently: for each client, create two go routines
  for read from and write to WebSocket
* Event mechanism: if an event that has been registered in manager occurs, call
  corresponding handler; easy to scale
* Decomposition of front end and back end: easy to deploy and reuse


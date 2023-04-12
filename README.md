# gomoku

Work process:

* [x] how to use go websocket -- [mastering websockets](https://programmingpercy.tech/blog/mastering-websockets-with-go/)
* [x] replace the frontend code that can't be applied to gomoku with or without Vue CLi for gomoku
* [x] connect the wire between frontend and Go WebSocket
* [x] test & put into production

---

In general, gomoku needs:
* GUI for players which represents:
    * [x] player name
    * [x] play result
    * [x] gomoku chessboard
* APIs to communicate messages
    * [x] role message -> C
    * [x] name message -> S & C
    * [x] move message -> S & C
    * [x] result message -> C


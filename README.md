# gomoku

Work process:

* [x] how to use go websocket -- [mastering websockets](https://programmingpercy.tech/blog/mastering-websockets-with-go/)
* [ ] replace the frontend code that can't be applied to gomoku with or without Vue CLi for gomoku
* [ ] connect the wire between Vue CLi and Go WebSocket
* [ ] test & put into production

---

In general, gomoku needs:
* GUI for players which represents:
    * [x] player name
    * [x] play result
    * [x] gomoku chessboard
* APIs to communicate messages
    * [ ] name message -> S
    * [ ] move message -> S
    * [ ] result message -> C
    * [ ] room message -> S (optional in fact)


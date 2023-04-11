/**
 * Generate chessboard
 * */
function initChessboard(n) {
    let board = document.getElementById("board");
    for (let i = 0; i < n; i++) {
        let row = board.insertRow(i);
        for (let j = 0; j < n; j++) {
            let cell = row.insertCell(j);
            turnIdle(cell);

            cell.addEventListener("click", function () {
                // alert("cell ["+ getX(this) + "," + getY(this) + "] clicked");
                if (!matched) {
                    alert("please wait until matched")
                    return false
                }
                if (!inTurn) {
                    alert("please wait until your turn")
                    return false
                }
                if (this.className != null) {
                    alert("cannot put on an occupied pile")
                    return false
                }
                if (role === 1) {
                    turnWhite(this)
                } else {
                    turnBlack(this)
                }
                sendMove(getX(cell), getY(cell))
                return false
            });
        }
    }
}

function turnWhite(cell) {
    cell.className = "whiteStone";
}

function turnBlack(cell) {
    cell.className = "blackStone";
}

function turnIdle(cell) {
    cell.className = "background";
}

function getX(cell) {
    return cell.parentNode.rowIndex;
}

function getY(cell) {
    return cell.cellIndex;
}

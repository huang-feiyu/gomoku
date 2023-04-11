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
                if (cell.className === "blackStone") {
                    turnWhite(cell);
                } else if (cell.className === "whiteStone") {
                    turnBlack(cell);
                } else {
                    turnWhite(cell);
                }
            });

            // test start
            if (j % 3 === 0) {
                turnBlack(cell);
            } else if (j % 3 === 2) {
                turnWhite(cell);
            }
            // test end
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

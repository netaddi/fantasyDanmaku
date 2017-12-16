var blackboard = document.getElementById('blackboard');
var board = document.getElementById('board');
var blackBoardCSSHeight = '84%';

function openBlackboard() {
    blackboard.style.height = blackBoardCSSHeight;
    board.style.display = 'block';
}

function closeBlackboard() {
    blackboard.style.height = '0';
    board.style.display = 'none';
}

function getJSON(theUrl, callback)
{
    var request = new XMLHttpRequest();
    request.open('GET', theUrl, true);

    request.onload = function() {
        if (this.status >= 200 && this.status < 400) {
            var data = JSON.parse(this.response);
            callback(data)
        }
    };
    request.send();
}
var userList;
var prizeDrawing = false;
var prizeInterval;

function initializePrizeDraw(){
    getJSON('/getUserList', function (userData) {
        userList = userData;
        openBlackboard();
        document.getElementById('lottery').style.display = 'block';
    })
}

function closePrizeDraw() {
    document.getElementById('lottery').style.display = 'none';
    closeBlackboard();
}

function drawPrize() {
    if (prizeDrawing){
        prizeDrawing = false;
        clearInterval(prizeInterval);
    } else {
        prizeDrawing = true;
        setPrizeInterval();
    }
}

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function setPrizeInterval(){
    var userListLength = userList.length;
    var waitMS = 50;
    prizeInterval = setInterval(function (){
        var id = getRandomInt(0, userListLength - 1);
        document.getElementById('prizeRegCode').innerHTML = userList[id].regCode.toString();
        document.getElementById('prizeNickname').innerHTML = userList[id].nickname.toString();
    }, waitMS);
}
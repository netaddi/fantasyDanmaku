// danmaku object type:
// text:
const danmakuDivId = 'danmaku';

class danmakuRails{
    constructor(){
        this.rails = [[]]
    }
    addDanmaku(danmaku){

    }
}

function moveDanmaku(danmakuItem){
    var totalLength = window.innerWidth + danmakuItem.offsetWidth;
    var movedLength = 0;
    var timer = setInterval(function () {
        movedLength += speed;
        danmakuItem.style.transform = 'translateX(-' + movedLength + 'px)';
        if(movedLength > totalLength){
            clearInterval(timer);
            danmakuItem.remove();
        }
    }, 10);

}

function generateDanmaku(jsonMessage) {
    var danmakuDiv = document.getElementById(danmakuDivId);
    var newHTMLNode = document.createElement("div");
    newHTMLNode.classList.add("comment");
    newHTMLNode.innerHTML = jsonMessage['Text'];
    newHTMLNode.style.color = jsonMessage['Color'];
    newHTMLNode.style.left = innerWidth + 'px';
    danmakuDiv.appendChild(newHTMLNode);
    moveDanmaku(newHTMLNode);
}

// danmaku object type:
// text:
const danmakuDivId = 'danmaku';
function generateDanmaku(jsonMessage) {
    var danmakuDiv = document.getElementById(danmakuDivId);
    var newHTMLNode = document.createElement("div");
    newHTMLNode.classList.add("comment");
    newHTMLNode.innerHTML = jsonMessage['Text'];
    newHTMLNode.style.fontSize = jsonMessage['Size'];
    newHTMLNode.style.color = jsonMessage['Color'];
    newHTMLNode.style.fontFamily = '"Helvetica Neue",Helvetica,Arial,"PingFang SC","Hiragino Sans GB","WenQuanYi Micro Hei","Microsoft Yahei",sans-serif';
    newHTMLNode.style.left = 330;

    danmakuDiv.appendChild(newHTMLNode);
}

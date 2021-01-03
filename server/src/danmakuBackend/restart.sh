ps aux | grep danmaku | grep -v grep | awk '{print $2}' | xargs kill && \
go build && \
nohup ./danmakuBackend >> nohup.log &
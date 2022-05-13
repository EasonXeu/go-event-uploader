## This is a go sdk which can upload logs with a high performance

1. This clien imported the tool [Memory Analyser](https://github.com/go-echarts/statsview), we can use it to analyse performance convinently.
3. The call to the produce.SendLog() method is returned immediately, so that the call site won't be blocked.
4. The logs were not sent to backend successfully would be sent to retry queue. And all logs in the retry queue would be sent again in the future.
5. The client would be closed safely, which can make sure there is no data missing, all logs which couldn't be sent successfully before closing would be saved in the file system, so that they can be sent in the next runtime.

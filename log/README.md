# Introduction:
This is simple go server log pkg.
Support:
- log to console
- log to file
- log level control
- log rotate daily

# usage:
1. Output To console
    ```golang
    import "github.com/sgs/log"
    fun main() {
        defer log.Close()
        log.Debug("This is Debug msg")
        log.Trace("This is Trace msg")
        // default log level INFO
        log.Info("This is Info msg")
        log.Warn("This is Warn msg")
        log.Fatal("This is Fatal msg")
    }
    ```
2. Output To file
    ```golang
    import "github.com/sgs/log"
    fun main() {
        defer log.Close()
        log.SetOutput("/tmp/xxx.log")
        log.SetLogLevel(INFO)
        log.SetRotateDaily()
        log.Info("This is Debug msg")
    }
    ```
# Feature planning
-  log rotate by file size
-  set Time zone UTC/CST,default CST
-  set LOG LEVEL
-  log to both file and console?(Anybody would use it this way? k8s log follow console?)

# Bugs??
- some log before call "rotate_chan<-true", may output to new log file.
-  some log output console before call SetOutput(),may output to log file.
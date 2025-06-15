
func main() {
    done := make(chan int)

    go func(){
        fmt.Println("你好, 世界")
        <-done
    }()

    done <- 1
}

根据 Go 语言内存模型规范，对于从无缓冲 Channel 进行的接收，发生在对该 Channel 进行的发送完成之前。因此，后台线程 <-done 接收操作完成之后，main 线程的 done <- 1 发送操作才可能完成（从而退出 main、退出程序），而此时打印工作已经完成了。

上面的代码虽然可以正确同步，但是对管道的缓存大小太敏感：如果管道有缓存的话，就无法保证 main 退出之前后台线程能正常打印了。更好的做法是将管道的发送和接收方向调换一下，这样可以避免同步事件受管道缓存大小的影响：

func main() {
    done := make(chan int, 1) // 带缓存的管道

    go func(){
        fmt.Println("你好, 世界")
        done <- 1
    }()

    <-done
}

对于带缓冲的 Channel，对于 Channel 的第 K 个接收完成操作发生在第 K+C 个发送操作完成之前，其中 C 是 Channel 的缓存大小。虽然管道是带缓存的，main 线程接收完成是在后台线程发送开始但还未完成的时刻，此时打印工作也是已经完成的。
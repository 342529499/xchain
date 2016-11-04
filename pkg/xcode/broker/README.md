xcode为智能合约的执行引擎
 
其服务端代码在 /Users/michael/IdeaProjects/src/github.com/1851616111/xchain/pkg/xcode/server 中

其客户端代码在 /Users/michael/IdeaProjects/src/github.com/1851616111/xchain/pkg/xcode/broker 中


broker并没有采取主动向server注册的模式，因为需要先知道xcode的server的address.
为了省略向broker告知 xcode server address。采用server主动连接xcode的模式。
这也充分体现了，broker作为server的一个独立的执行进程，server应该展现充足的掌控力。

而broker的代码看起来也会更整洁。 
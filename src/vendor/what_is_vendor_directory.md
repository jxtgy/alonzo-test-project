
## vendor目录的作用
$GOPATH
|	src/
|	|	github.com/constabulary/example-gsftp/
|	|	|	cmd/
|	|	|	|	gsftp/
|	|	|	|	|	main.go
|	|	|	vendor/
|	|	|	|	github.com/pkg/sftp/
|	|	|	|	golang.org/x/crypto/ssh/
|	|	|	|	|	agent/

假设你的工程目录是上图所示，假设你正在github.com/constabulary/example-gsftp/cmd/gsftp/main.go这个文件里编写代码，你
打算引入这个包：github.com/constabulary/example-gsftp/vendor/golang.org/x/crypto/ssh
在代码里并不需要写成下面这样:
import "github.com/constabulary/example-gsftp/vendor/golang.org/x/crypto/ssh"
实际上你只需要写如下代码：
import "golang.org/x/crypto/ssh"

原理是：编译器会首先在vendor目录下寻找golang.org/x/crypto/ssh包，再去GOPATH目录寻找这个包。
vendor目录的作用就是为了避免出现特别长的import代码。

## reference
https://github.com/golang/go/issues/14566

The design of the "vendor" mechanism is what you should expect: https://golang.org/s/go15vendor
https://go.googlesource.com/proposal/+/master/design/25719-go15vendor.md


Adjusting that example to use the new vendor directory, the source tree would look like:

$GOPATH
|	src/
|	|	github.com/constabulary/example-gsftp/
|	|	|	cmd/
|	|	|	|	gsftp/
|	|	|	|	|	main.go
|	|	|	vendor/
|	|	|	|	github.com/pkg/sftp/
|	|	|	|	golang.org/x/crypto/ssh/
|	|	|	|	|	agent/
The file github.com/constabulary/example-gsftp/cmd/gsftp/main.go says:

import (
	...
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/pkg/sftp"
)
Because github.com/constabulary/example-gsftp/vendor/golang.org/x/crypto/ssh exists and the file being compiled is within the subtree rooted at github.com/constabulary/example-gsftp (the parent of the vendor directory), the source line:

import "golang.org/x/crypto/ssh"
is compiled as if it were:

import "github.com/constabulary/example-gsftp/vendor/golang.org/x/crypto/ssh"
(but this longer form is never written).

So the source code in github.com/constabulary/example-gsftp depends on the vendored copy of golang.org/x/crypto/ssh, not one elsewhere in $GOPATH.

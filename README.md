# GO_FLUTTER_ALIOSS

go-flutter 的阿里云oss插件，依照我的另一个flutter版插件写的 [aliossflutter](https://github.com/jlcool/aliossflutter)，接口可以在这里看

暂时只能使用secretInit初始化，实现了upload，download，进度，其他的慢慢加

## Usage

#### 1. 安装 https://github.com/jlcool/aliossflutter

#### 2. 在cmd/options.go中添加

Import as:

```go
import "github.com/jlcool/go_flutter_alioss"
```

```go
flutter.AddPlugin(&aliossFlutter.AliossFlutterPlugin{}),
```

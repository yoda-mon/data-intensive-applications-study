# Encoding

### Structure
```sh
.
├── README.md        
├── avro             avro作業用
├── cmd              go本体
├── config           schemaファイル格納
├── dist             binary格納
├── go.mod
├── go.sum
├── pom.xml
├── protocolbuffer   pb作業用
├── sample.json      大元のファイル
├── src              java本体
├── target           javaビルド先
└── thrift           thrift作業用
```

## Main
### JSON -> Java Object -> Binary
```sh
mvn clean package
java -jar target/encoding-1.0-SNAPSHOT-jar-with-dependencies.jar

```

`dist/`直下にjson, messagepack, thrift(binary protocol), protocol buffer, avroでエンコードされたファイルが生成される

### Binary -> Golang Object
上を行ったのち

```sh
go run cmd/main.go
```

でそれぞれデコードしてgoのオブジェクトに収め、Printする
## Memo
### MessagePack
#### Java
https://github.com/msgpack/msgpack-java
ラク

#### Go
http://ugorji.net/blog/go-codec-primer#decoding
この例を参考、ただ上手くStructへのマッピングができずお茶は濁す

### Thrift
#### Java
https://thrift.apache.org/
https://thrift-tutorial.readthedocs.io/en/latest/index.html
https://thrift.apache.org/docs/types

公式のMapperないようなので、Jacksonで直接マッピングさせる

```sh
brew install thrift # mac
thrift -r --gen java sample.thrift
```

`gen-java`が生成されるので、`/thrift`以下を`src/main/java`下に配置
#### Go
```sh
thrift -r --gen go sample.thrift
```

`gen-go`が生成されるので、`/thrift`以下を`cmd/`下に配置

### Protocol Buffers
#### Java
https://developers.google.com/protocol-buffers
https://github.com/FasterXML/jackson-dataformats-binary/tree/master/protobuf

新しいフォーマットだとrequiredとかoptionalではないみたいのでproto2でやる  

コード生成せずともJacksonがよろしくやってくれるのでラク

#### Go
デコード用のコード生成
```
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc --proto_path=config --go_out=cmd/protobuf/ --go_opt=paths=source_relative config/sample.proto
```

生成したコードを`cmd/`下に配置
(Goをインストールし直したら上手くいったがgo modulesと相性悪いみたい)
一回生成できたらあとは超ラク

### Avro

#### Java
テキストに書かれたスキーマはAvro IDLでスキーマなどを簡便に記述するためのDSLみたいなもんらしい（拡張子`.avdl`）  
スキーマ(`.avsc`)として利用するには変換する必要あり。

```sh
wget https://downloads.apache.org/avro/stable/java/avro-tools-1.9.2.jar
java -jar avro-tools-1.9.2.jar idl2schemata sample.avdl
```

これで`sample.avsc`ができる、あとはPBと同じ流れ

#### Go
https://github.com/linkedin/goavro
こちらもコード生成しなくても読めたのでラク


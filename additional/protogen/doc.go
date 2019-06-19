//go:generate /usr/local/bin/protoc --go_out $GOPATH/src $GOPATH/src/kuitest/additional/protogen/pkga/a.proto -I $GOPATH/src/gitlab.myteksi.net/gophers/go/commons/util/grab-kit/extproto -I $GOPATH/src/kuitest
//go:generate /usr/local/bin/protoc --go_out $GOPATH/src $GOPATH/src/kuitest/additional/protogen/pkgb/b.proto -I $GOPATH/src/gitlab.myteksi.net/gophers/go/commons/util/grab-kit/extproto -I $GOPATH/src/kuitest -I /Users/kui.chen/gopath/src/kuitest/additional/protogen/pkga
package main

func main() {

}

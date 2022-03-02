#export GOPATH=$HOME/go
#PATH=$PATH:$GOPATH/bin

protoc greet/greetpb/greet.proto  --go_out=plugins=grpc:.
protoc caluculator/caluculatorpb/caluculator.proto  --go_out=plugins=grpc:.
protoc blog/blogpb/blog.proto  --go_out=plugins=grpc:.

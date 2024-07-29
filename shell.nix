{ mkShell
, go
, gopls
, ffmpeg
, nixfmt-classic
, protoc-gen-go
, protobuf
, protoc-gen-go-grpc
, act
, go-migrate
, sqlc
, ...
}: mkShell {
  hardeningDisable = [ "all" ];
  nativeBuildInputs = [
    go
    protobuf
    go-migrate
    sqlc
    act
  ];
  buildInputs = [
    gopls
    protoc-gen-go-grpc
    protoc-gen-go
    ffmpeg
    nixfmt-classic
  ];

  GOROOT = "${go}/share/go";

}

{ mkShell, go, gopls, ffmpeg, nixfmt-classic, protoc-gen-go, protobuf, protoc-gen-go-grpc, ... }: mkShell {
  hardeningDisable = [ "all" ];
  nativeBuildInputs = [
    go
    protobuf
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

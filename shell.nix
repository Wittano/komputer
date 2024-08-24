{ mkShell
, go
, gopls
, ffmpeg
, nixfmt-classic
, act
, ...
}: mkShell {
  NIX_HARDENING_ENABLE = "";
  hardeningDisable = [ "all" ];
  nativeBuildInputs = [
    go
    act
  ];

  buildInputs = [
    gopls
    ffmpeg
    nixfmt-classic
  ];

  GOROOT = "${go}/share/go";
}

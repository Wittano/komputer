{ mkShell, go, gopls, ffmpeg, nixfmt-classic, ... }: mkShell {
  hardeningDisable = [ "all" ];
  buildInputs = [
    # Go
    go
    gopls

    # Runtime dependecies
    ffmpeg

    # Nixpkgs
    nixfmt-classic
  ];

  GOROOT = "${go}/share/go";

}

{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    gopls
    ffmpeg
    rnix-lsp
    nixfmt
    nixpkgs-fmt
  ];
}
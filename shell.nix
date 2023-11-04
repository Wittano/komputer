{ pkgs ? import <nixpkgs> {} }:
 pkgs.mkShell {
    buildInputs = with pkgs; [
        go
        ffmpeg
        rnix-lsp
        nixfmt
        nixpkgs-fmt
    ];
}
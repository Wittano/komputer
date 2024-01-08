{
  description = "komputer - discord bot";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
      };

      lib = nixpkgs.lib.extend (self: super: {
        my = import ./lib { inherit pkgs; lib = self; };
      });

      komputer = pkgs.callPackage ./default.nix { };
    in
    rec {
      inherit lib;

      packages.x86_64-linux.komputer = komputer;
      defaultPackage.x86_64-linux = komputer;

      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          ffmpeg
          rnix-lsp
          nixfmt
          nixpkgs-fmt
        ];
      };
    };
}


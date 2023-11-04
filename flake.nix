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
    in rec {
      inherit lib;

      packages.x86_64-linux.komputer = pkgs.callPackage ./default.nix {};

      defaultPackage.x86_64-linux = self.packages.x86_64-linux.komputer;

      devShell = import ./shell.nix { inherit pkgs; };
    };
}
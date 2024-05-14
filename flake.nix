{
  description = "komputer - discord bot";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          packages.default = pkgs.callPackage ./nixos/default.nix { };
          devShells.default = pkgs.callPackage ./nixos/shell.nix { };
        }
      ) // { nixosModules.default = ./nixos/module.nix; };
}


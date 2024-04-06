{
  description = "komputer - discord bot";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };

      lib = nixpkgs.lib;
    in {
      defaultPackage.${system} = pkgs.buildGoModule {
        pname = "komputer";
        version = "v1.1.0";

        src = ./.;

        vendorHash = "sha256-TAbmS8xYqleXIU0cCRiJVyC93jiXiLuLepMD4WcS7IQ=";
        CGO_ENABLED = 1;
        proxyVendor = true;

        nativeBuildInputs = with pkgs; [ gcc pkg-config libopus ];
        propagatedBuildInputs = with pkgs; [ ffmpeg opusfile ];

        preBuild = "go get layeh.com/gopus";

        meta = with lib; {
          homepage = "https://github.com/Wittano/komputer";
          description =
            "Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
          license = licenses.gpl3;
          maintainers = with maintainers; [ Wittano ];
          platforms = platforms.linux;
        };
      };
      devShells.${system}.default = pkgs.mkShell {
        hardeningDisable = [ "all" ];
        buildInputs = with pkgs; [
          # Go
          go
          gopls

          # Runtime dependecies
          ffmpeg

          # Nixpkgs
          nixfmt-classic
        ];
      };
    };
}


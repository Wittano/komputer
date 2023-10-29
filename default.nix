{ lib, buildGoModule, gcc, pkg-config, libopus, ffmpeg, opusfile  }: buildGoModule {
  pname = "komputer";
  version = "v1.0.0";

  src = ./.;

  vendorHash = "sha256-TAbmS8xYqleXIU0cCRiJVyC93jiXiLuLepMD4WcS7IQ=";

  nativeBuildInputs = [ gcc pkg-config libopus ];

  CGO_ENABLED = 1;

  proxyVendor = true;

  preBuild = ''
    go get layeh.com/gopus
  '';

  propagatedBuildInputs = [ ffmpeg opusfile ];

  meta = with lib; {
    homepage = "https://github.com/Wittano/komputer";
    description = "Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
    license = licenses.gpl3;
    maintainers = with maintainers; [ Wittano ];
    platforms = platforms.linux;
  };
}
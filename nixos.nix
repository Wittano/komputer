{ pkgs, lib, config, ... }:
with lib;
let
  cfg = config.komputer;
  komputer = pkgs.callPackage ./default.nix { };
in
{
  options = {
    komputer = {
      enable = mkEnableOption "Enable komputer discord bot";
      package = mkOption {
        type = types.package;
        default = komputer;
      };
      applicationIDSecretPath = mkOption {
        type = types.path;
        description = "Application ID for you local version of komputer bot";
      };
      tokenSecretPath = mkOption {
        type = types.path;
        description = "Path to file, that contain discord token for bot";
      };
      mongodbURISecretPath = mkOption {
        type = types.path;
        description = "Connection URI to your instance of mongodb";
      };
    };
  };

  config = mkIf (cfg.enable) {
    assertions = [
      {
        assertion = cfg.tokenSecretPath != "";
        message = "Option komputer.token is empty";
      }
      {
        assertion = cfg.applicationIDSecretPath != "";
        message = "Option komputer.applicationID is empty";
      }
      {
        assertion = cfg.mongodbURISecretPath != "";
        message = "Option komputer.mongodbURI is empty";
      }
    ];

    systemd.services.komputer = {
      description = "Komputer - Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
      wantedBy = [ "multi-user.target" ];
      path = cfg.package.propagatedBuildInputs or [ ];
      environment = {
        DISCORD_BOT_TOKEN_PATH = cfg.tokenSecretPath;
        APPLICATION_ID_PATH = cfg.applicationIDSecretPath;
        MONGODB_URI_PATH = cfg.mongodbURISecretPath;
      };
      script = "${cfg.package}/bin/komputer";
    };
  };
}


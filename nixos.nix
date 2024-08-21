{ pkgs, lib, config, ... }:
with lib;
let
  cfg = config.komputer;

  komputer = pkgs.callPackage ./../default.nix { };
in
{
  options = {
    komputer = {
      enable = mkEnableOption "Enable komputer discord bot";
      package = mkOption {
        type = types.package;
        default = komputer;
        description = "komputer package";
      };
      guildID = mkOption {
        type = types nullOr types.str;
        default = null;
        description = "Discord server id, that you deploy bot";
      };
      applicationID = mkOption {
        type = types.str;
        description = "Application ID for you local version of komputer bot";
      };
      token = mkOption {
        type = types.str;
        description = ''
          Discord token for bot. 
          <REMEMBER!>
          Your token never shouldn't be publish on any public git repository e.g. Github or Gitlab
        '';
      };
      mongodbURI = mkOption {
        type = types.str;
        description = "Connection URI to your instance of mongodb";
      };
    };
  };

  config = mkIf (cfg.enable) {
    assertions = [
      {
        assertion = cfg.token != "";
        message = "Option komputer.token is empty";
      }
      {
        assertion = cfg.applicationID != "";
        message = "Option komputer.applicationID is empty";
      }
      {
        assertion = cfg.mongodbURI != "";
        message = "Option komputer.mongodbURI is empty";
      }
    ];

    systemd.services.komputer = {
      description = "Komputer - Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
      wantedBy = [ "multi-user.target" ];
      path = cfg.package.propagatedBuildInputs or [];
      environment = {
        DISCORD_BOT_TOKEN = cfg.token;
        APPLICATION_ID = cfg.applicationID;
        MONGODB_URI = cfg.mongodbURI;
      } // attrsets.optionAtts (cfg.guildID != null) { SERVER_GUID = cfg.guildID; };
      script = "${cfg.package}/bin/komputer";
    };
  };
}


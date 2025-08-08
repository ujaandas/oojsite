{
  description = "Example development environment flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages = {
          oojsite = pkgs.buildGoModule {
            pname = "oojsite";
            version = "0.1.0";
            src = ./.;
            vendorHash = null;
            meta = with pkgs.lib; {
              description = "Go personal blogsite";
              homepage = "https://ujaandas.me";
              license = licenses.mit;
              maintainers = [ maintainers."ujaandas" ];
            };
          };
        };
        defaultPackage = packages.oojsite;
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
          ];
          shellHook = ''
            echo "Welcome to the development shell!"
          '';
        };
      }
    );
}

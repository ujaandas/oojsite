{
  description = "Flake for oojsite with standalone Tailwind build";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        oojsite = pkgs.buildGoModule {
          pname = "oojsite";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-gM37SLXNi4uY3uetmagNarbUvaFapQciajrguWVSd34=";
          buildInputs = with pkgs; [
            makeWrapper
            tailwindcss
          ];

          postInstall = ''
            wrapProgram $out/bin/oojsite \
              --prefix PATH : "${pkgs.tailwindcss}/bin"
          '';

        };
      in
      {
        defaultPackage = oojsite;
        packages.oojsite = oojsite;

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
            watchexec
          ];
        };
      }
    );
}

{
  description = "Flake for oojsite with standalone Tailwind build";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
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
          vendorHash = null;

          nativeBuildInputs = with pkgs; [ tailwindcss ];

          buildPhase = ''
            tailwindcss \
              -i ${./assets/static/css/styles.css} \
              -o $out/public/styles.css

            mkdir -p $out/bin
            go build -o $out/bin/oojsite .
          '';

          installPhase = ''
            $out/bin/oojsite --out $out/
          '';
        };

      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
            tailwindcss
            watchexec
          ];
        };

        defaultPackage = oojsite;
      }
    );
}

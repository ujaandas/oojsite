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
          vendorHash = "sha256-gM37SLXNi4uY3uetmagNarbUvaFapQciajrguWVSd34=";

          nativeBuildInputs = with pkgs; [ tailwindcss ];

          buildPhase = ''
            tailwindcss \
              --input ${./public/styles.css} \
              --output $out/public/styles.css \
              --content **/*.html \
              --minify

            mkdir -p $out/bin
            go build -o $out/bin/oojsite .
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

        defaultApp = {
          type = "app";
          program = "${pkgs.writeShellScriptBin "oojsite-run" ''
            tmpdir=$(mktemp -d)
            echo "Generating site into $tmpdir"

            mkdir -p "$tmpdir/public"
            cp ${oojsite}/public/styles.css "$tmpdir/public/styles.css"

            ${oojsite}/bin/oojsite --out "$tmpdir"

            rm -f ./out
            ln -s "$tmpdir" ./out
          ''}/bin/oojsite-run";
        };
      }
    );
}

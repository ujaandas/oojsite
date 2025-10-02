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
      rec {
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
        defaultApp = apps.oojsite-run;

        apps = {
          oojsite-run = {
            type = "app";
            program = "${pkgs.writeShellScriptBin "oojsite-run" ''
              outdir="./out"
              rm -rf $outdir
              echo "Generating site into $outdir"

              mkdir -p "$outdir/public"
              cp -r ${oojsite}/public/styles.css "$outdir/public/styles.css"

              ${oojsite}/bin/oojsite --out "$outdir"

              echo "View site at file://$outdir/index.html"
            ''}/bin/oojsite-run";
          };

          watch = {
            type = "app";
            program = "${pkgs.writeShellScriptBin "oojsite-watch" ''
              echo "Watching for changes and regenerating site..."

              watchexec \
                --restart \
                --clear \
                --watch ./public \
                --watch ./templates \
                --watch ./site \
                --exts css,html,go,md \
                -- nix run
            ''}/bin/oojsite-watch";
          };
        };
      }
    );
}

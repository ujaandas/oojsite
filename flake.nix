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

          nativeBuildInputs = with pkgs; [ tailwindcss ];

          buildPhase = ''
            tailwindcss \
              --input ./public/styles.css \
              --output $out/public/styles.css \
              --minify \
              --config ./tailwind.config.js


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

              ${oojsite}/bin/oojsite --out "$outdir"

              mkdir -p "$outdir/public"
              cp -r ${oojsite}/public/. "$outdir/public/"

              echo "View site at file://$outdir/index.html"
            ''}/bin/oojsite-run";
          };

          watch = {
            type = "app";
            program = "${pkgs.writeShellScriptBin "oojsite-watch" ''
              echo "Watching site templates and regenerating..."

              watchexec \
                --restart \
                --clear \
                --watch ./templates \
                --watch ./site \
                --watch ./public \
                --exts html,go,md,css \
                -- 'tmpdir=$(mktemp -d);
                  go run . --out $tmpdir || exit 1;

                  tailwindcss \
                    --input ./public/styles.css \
                    --output "$tmpdir/public/styles.css" \
                    --minify \
                    --config ./tailwind.config.js || exit 1;

                  cp -r $tmpdir/public/. ./out/public/
                '
            ''}/bin/oojsite-watch";
          };
        };
      }
    );
}

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

        compiledCss =
          pkgs.runCommand "tailwind-css"
            {
              buildInputs = [ pkgs.tailwindcss ];
            }
            ''
              mkdir -p $out
              tailwindcss \
                --input  ${./assets/static/css/tailwind.css} \
                --output $out/tailwind.css \
                --content ${./assets}/templates/*.html ${./assets}/content/**/*.md \
                --minify
            '';

        oojsite = pkgs.buildGoModule {
          pname = "oojsite";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;

          buildPhase = ''
            go build -o oojsite ./cmd/oojsite
          '';

          installPhase = ''
            mkdir -p $out/{bin,static/css}
            install -m755 oojsite $out/bin/oojsite
            install -m644 ${compiledCss}/tailwind.css \
              $out/static/css/tailwind.css
          '';
        };

        runWrapper = pkgs.writeShellScriptBin "oojsite" ''
          set -euo pipefail

          outdir="/tmp/oojsite-dev"

          echo "↪ Generating site into $outdir"
          rm -rf "$outdir"
          mkdir -p "$outdir"

          ${oojsite}/bin/oojsite --out="$outdir" "$@"

          mkdir -p "$outdir/static/css"
          cp ${oojsite}/static/css/tailwind.css "$outdir/static/css/"

          rm -rf ./public
          ln -s "$outdir" ./public
        '';
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

        packages.oojsite = oojsite;

        apps.oojsite = {
          type = "app";
          program = "${runWrapper}/bin/oojsite";
        };

        apps.dev = {
          type = "app";
          program = "${pkgs.writeShellScriptBin "oojsite-dev" ''
            set -euo pipefail
            echo "↪ Starting oojsite live dev mode…"
            watchexec -r -e md,html,css -- \
              nix run .#oojsite
          ''}/bin/oojsite-dev";
        };
      }
    );
}

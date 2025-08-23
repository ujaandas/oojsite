{
  description = "Flake with separate build vs run";

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
          pkgs.runCommand "tailwindcss-output"
            {
              buildInputs = [ pkgs.tailwindcss ];
            }
            ''
              mkdir -p $out/tailwind
              tailwindcss \
                --input    ${./assets/static/css/tailwind.css} \
                --output   $out/tailwind/tailwind.css \
                --content  ${./assets}/templates/*.html ${./assets}/content/**/*.md \
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
            mkdir -p $out/bin $out/static/css
            cp oojsite $out/bin/

            cp ${compiledCss}/tailwind/tailwind.css \
              $out/static/css/tailwind.css
          '';
        };

        oojsiteWrapped = pkgs.writeShellScriptBin "oojsite" ''
          set -ex

          tmpdir=$(mktemp -d -t oojsite-XXXXXX)

          ${oojsite}/bin/oojsite --out="$tmpdir" "$@"

          mkdir -p "$tmpdir/static/css"
          cp ${oojsite}/static/css/tailwind.css "$tmpdir/static/css/"

          rm -rf ./public
          ln -s "$tmpdir" ./public
        '';
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
            tailwindcss
          ];
          shellHook = ''
            echo "Welcome to the oojsite dev shell"
          '';
        };

        packages.oojsite = oojsite;

        apps.oojsite = {
          type = "app";
          program = "${oojsiteWrapped}/bin/oojsite";
        };
      }
    );
}

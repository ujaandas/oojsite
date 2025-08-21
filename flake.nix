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
      {
        defaultPackage = pkgs.buildGoModule {
          pname = "bloggor";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;

          nativeBuildInputs = [ pkgs.tailwindcss ];

          buildPhase = ''
            tailwindcss \
              --input internal/generate/static/css/styles.css \
              --output public/css/styles.css \
              --content content/**/*.md internal/render/templates/*.html \
              --minify

            go build -o bloggor ./cmd/bloggor
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp bloggor $out/bin/
          '';

          subPackages = [ "cmd/bloggor" ];
          meta = with pkgs.lib; {
            description = "Go personal blogsite";
            homepage = "https://ujaandas.me";
            license = licenses.mit;
            maintainers = [ maintainers."ujaandas" ];
          };
        };

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

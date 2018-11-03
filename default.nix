with import <nixpkgs> {};
let
 unstable = import <nixos-unstable> {};
 go = unstable.go_1_11;
in
stdenv.mkDerivation rec {
  name = "env";
  env = buildEnv { name = name; paths = buildInputs; };
  buildInputs = [
    go
    unstable.gotools
    unstable.godef
    unstable.golint
  ];
}

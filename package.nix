{buildGoModule}:
buildGoModule {
  pname = "lightspeed-dhl-adapter";
  version = "0.1.1";
  src = ./src;

  vendorHash = "sha256-B/NtKPAOHRcVq1VBK/L/kCQ04Fvyitfpt5c3M273I8M=";
  meta.mainProgram = "lightspeed-dhl";

  GO_TEST = "none";
}

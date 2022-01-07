require "formula"

class Mop < Formula
  homepage "https://github.com/mop-tracker/mop"
  head     "https://github.com/mop-tracker/mop.git"
  url      "https://github.com/mop-tracker/mop/archive/refs/tags/v1.0.0.tar.gz"
  sha1     "bc666ec165d08b43134f7ec0bf29083ad5466243" # Needs updating.

  depends_on "go" => :build

  def install
    system "go", "get", "github.com/nsf/termbox-go"
    system "go build cmd/mop.go"
    bin.install "mop"
  end
end

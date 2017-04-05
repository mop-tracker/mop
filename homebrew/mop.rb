require "formula"

class Mop < Formula
  homepage "https://github.com/mop-tracker/mop"
  head     "https://github.com/mop-tracker/mop.git"
  url      "https://github.com/mop-tracker/mop/archive/v0.1.0.tar.gz"
  sha1     "b3bd5b529430da22bfa1beeca435a25e72513e27"

  depends_on "go" => :build

  def install
    system "go", "get", "github.com/michaeldv/termbox-go"
    system "go build cmd/mop.go"
    bin.install "mop"
  end
end

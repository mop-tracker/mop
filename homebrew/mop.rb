require "formula"

class Mop < Formula
  homepage "https://github.com/brandleesee/mop"
  head     "https://github.com/brandleesee/mop.git"
  sha1     "b3bd5b529430da22bfa1beeca435a25e72513e27"

  depends_on "go" => :build

  def install
    system "go", "get", "github.com/brandleesee/termbox-go"
    system "go build cmd/mop.go"
    bin.install "mop"
  end
end

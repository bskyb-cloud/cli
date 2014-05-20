require 'formula'

class CfSrc < Formula
  homepage 'https://github.com/nimbus-cloud/cli'
  url 'https://github.com/nimbus-cloud/cli.git', :tag => 'v6.0.1'

  head 'https://github.com/nimbus-cloud/cli.git', :branch => 'master'

  depends_on 'go' => :build

  def install
    inreplace 'cf/app_constants.go', 'SHA', 'homebrew'
    inreplace 'cf/app_constants.go', 'BUILT_FROM_SOURCE', 'homebrew'
    system 'bin/build'
    bin.install 'out/cf'
    doc.install 'LICENSE'
  end

  test do
    system "#{bin}/cf"
  end
end

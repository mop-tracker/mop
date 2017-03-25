## Terminal Stocks

Terminal Stocks is a command-line utility that displays continuous up-to-date information about the U.S. markets and individual stocks in a terminal emulator. 

![](http://i.imgur.com/SkyRCpW.png)

### Updates

This repository assimilates the following fixes:

* [fix for market data parsing](https://github.com/mop-tracker/mop/pull/24/commits/641b65aa60f6a3090186fb3a87fc6c9b04e786fe)
* [correct display of timezone name](https://github.com/cuihantao/mop/commit/7037a19e20d87d2d02037fcfc4e22b4b686d75de)
* [red colour for negative numbers & better N/A handling](https://github.com/dalbert/mop/commit/ee4f79f97f2261a0fa4eb22f103c2c9ec645d38c)

### Installing Go Language in Ubuntu and derivatives

```bash
sudo apt -y install golang
mkdir -p ~/go
sudo mkdir -p /etc/profile.d/
sudo nautilus /etc/profile.d/goenv.sh
    export GOROOT=/usr/lib/go
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
source /etc/profile.d/goenv.sh
```

### Installing Mop in Ubuntu and derivatives

Terminal Stocks is implemented in Go and compiles down to a single executable file.

**Make sure your $GOPATH is set.** *see instructions above*

Build mop and install it in `` $GOPATH/bin `` :

```bash
go get github.com/brandleesee/mop
cd $GOPATH/src/github.com/brandleesee/mop
make install
```

### Usage

For demonstartion purposes Terminal Stocks comes preconfigured with a number of stock tickers. You can easily change the default list by using the following keyboard commands:

    +       Adds stocks to the list.
    -       Removes stocks from the list.
    o       Changes column sort order.
    g       Groups stocks by advancing/declining issues.
    q       Quits mop.
    ?       Displays help screen.
    esc     Quits mop.

When prompted please enter comma-delimited list of stock tickers. The list and other settings are stored in `` .moprc `` file in your `` $HOME `` directory.

### License

Copyright (c) 2013-2016 Michael Dvorkin. All Rights Reserved.
"mike" + "@dvorkin" + ".net" || "twitter.com/mid"

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

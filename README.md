# Terminal Stocks

**Tested on Windows 7 and Ubuntu 16.04 & 18.04.**

Terminal Stocks is a command-line utility that displays continuous up-to-date information about select markets and individual stocks in a terminal emulator. 

![](https://user-images.githubusercontent.com/698668/44194756-cf458a80-a0eb-11e8-93b4-3f8a3cdc5c7a.png)

## Installing Go Language in Ubuntu and derivatives

### Download Go in system-preferred folder structure

```bash
sudo apt -y install golang
mkdir -p ~/go
sudo mkdir -p /etc/profile.d/
```

### Create Go Environment

```bash
sudo nautilus /etc/profile.d/goenv.sh
```

### Put the following in the Go Environment file `` goenv.sh ``

```bash
    export GOROOT=/usr/lib/go
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

### Make Go Environment recognised by system
    
```bash
source /etc/profile.d/goenv.sh
```

## Installing Terminal Stocks in Ubuntu and derivatives

Terminal Stocks is implemented in Go and compiles down to a single executable file.

### Make sure your `` $GOPATH `` is set

```
see instructions above
```

### Build Terminal Stocks and install it in `` $GOPATH/bin ``:

```bash
go get github.com/brandleesee/TerminalStocks
cd $GOPATH/src/github.com/brandleesee/TerminalStocks
make install
```

![](https://i.imgur.com/qkT8SL7.png)

## Building TerminalStocks on Windows 7 x64

### Download Go

https://dl.google.com/go/go1.11.windows-amd64.msi

### Create Subfolders in $Users/USERNAME

```
C:\Users\CHANGE-TO-USERNAME\go\src\github.com\brandleesee
C:\Users\CHANGE-TO-USERNAME\go\src\github.com\mattn
C:\Users\CHANGE-TO-USERNAME\go\src\github.com\nsf 
```

### Extract Repositories to Respective Folders

![](https://i.imgur.com/WZmfQtq.png)
https://github.com/brandleesee/TerminalStocks/archive/master.zip  
C:\Users\blc\go\src\github.com\brandleesee\TerminalStocks  
*rename TerminalStocks-master to TerminalStocks*  

![](https://i.imgur.com/SjAhiWC.png)
https://github.com/mattn/go-runewidth/archive/master.zip  
C:\Users\blc\go\src\github.com\mattn\go-runewidth  
*rename go-runewidth-master to go-runewidth*  

![](https://i.imgur.com/cbUpBId.png)
https://github.com/nsf/termbox-go/archive/master.zip  
C:\Users\blc\go\src\github.com\nsf\termbox-go  
*rename termbox-go-master to termbox-go*  

**make sure that there are no duplicate folders**
**C:\Users\blc\go\src\github.com\brandleesee\TerminalStocks\TerminalStocks**

### Creating the Executable

Run `cmd` then:

```
cd C:\Users\**USERNAME**\go\src\github.com\brandleesee\TerminalStocks\cmd\TerminalStocks
go build
```

Starting the newly created executable will start TerminalStocks in cmd.

### Maximize Command Prompt Window

*this is done only once*

* Run `cmd`
* Right-click on Title Bar and select `Defaults`
* In **Options** tab, tick `QuickEdit mode`
* In **Layout** tab resize as required.

<table>
<tbody>
<tr>
<td align="center"><img src="https://i.imgur.com/QVzbqIT.png" /></td>
<td align="center"><img src="https://i.imgur.com/8UFs3Mg.png" /></td>
</tr>
</tbody>
</table>

## Usage

For demonstartion purposes Terminal Stocks comes preconfigured with a number of stock tickers. You can easily change the default list by using the following keyboard commands:

    +       Adds stocks to the list.
    -       Removes stocks from the list.
    o       Changes column sort order.
    g       Groups stocks by advancing/declining issues.
    q       Quits Terminal Stocks.
    ?       Displays help screen.
    esc     Quits Terminal Stocks.

When prompted please enter comma-delimited list of stock tickers. The list and other settings are stored in `` .TSrc `` file in your `` $HOME `` directory.

## License

Copyright (c) 2017-2018 Brandon Lee Camilleri. All Rights Reserved. github.com/brandleesee  
Copyright (c) 2013-2016 Michael Dvorkin. All Rights Reserved. "mike" + "@dvorkin" + ".net" || "twitter.com/mid"  

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

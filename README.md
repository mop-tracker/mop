### Mop: track stocks the hacker way ###
Mop is a command-line utility that displays continuous up-to-date
information about the U.S. markets and individual stocks. One
screenshot is worth a thousand words:

![Screenshot](https://raw.githubusercontent.com/mop-tracker/mop/master/doc/screenshot.png "Mop Screenshot")

### Installing Mop ###
Mop is implemented in Go and compiles down to a single executable file.

    # Make sure your $GOPATH is set.
    $ go get github.com/mop-tracker/mop
    $ cd $GOPATH/src/github.com/mop-tracker/mop
    $ make            # <-- Compile and run mop.
    $ make build      # <-- Build mop in current directory.
    $ make install    # <-- Build mop and install it in $GOPATH/bin.


### Using Mop ###
For demonstartion purposes Mop comes preconfigured with a number of
stock tickers. You can easily change the default list by using the
following keyboard commands:

    +       Add stocks to the list.
    -       Remove stocks from the list.
    o       Change column sort order.
    g       Group stocks by advancing/declining issues.
    ?       Display help screen.
    esc     Quit mop.

When prompted please enter comma-delimited list of stock tickers. The
list and other settings are stored in ``.moprc`` file in your ``$HOME``
directory.


### Contributing ###
Mop is my personal project I came up with to learn Go programming. Your
comments, suggestions, and contributions are welcome.

* Fork the project on Github.
* Make your feature addition or bug fix.
* Commit, do not change program version, or commit history.
* Email me commit URL (*do not* send pull requests).


### License ###
Copyright (c) 2013-2016 Michael Dvorkin. All Rights Reserved.
"mike" + "@dvorkin" + ".net" || "twitter.com/mid"

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

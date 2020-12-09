### Mop: track stocks the hacker way ###

Mop is a command-line utility that displays continuous up-to-date
information about the U.S. markets and individual stocks. One
screenshot is worth a thousand words:

![Screenshot](https://raw.githubusercontent.com/mop-tracker/mop/master/doc/screenshot.png "Mop Screenshot")

### Installing Mop ###
Mop is implemented in Go and compiles down to a single executable file.

    $ git clone https://github.com/mop-tracker/mop
    $ cd mop
    $ go build ./mop/cmd
    $ ./mop

### Using Mop ###
For demonstration purposes Mop comes preconfigured with a number of
stock tickers. You can easily change the default list by using the
following keyboard commands:

    +       Add stocks to the list.
    -       Remove stocks from the list.
    o       Change column sort order.
    g       Group stocks by advancing/declining issues.
    f       Set a filtering expression.
    F       Unset a filtering expression.
    ?       Display help screen.
    esc     Quit mop.

When prompted please enter comma-delimited list of stock tickers. The
list and other settings are stored in the profile file (default: ``.moprc`` in your ``$HOME`` directory)

### Expression-based Filtering
Mop has an in realtime expression-based filtering engine that is very easy to use.

At the main screen, press `f` and a prompt will appear. Write an expression that uses the stock properties.

Example:

```last <= 5```

This expression will make Mop show only the stocks whose `last` values are less than $5.

The available properties are: `last`, `change`, `changePercent`, `open`, `low`, `high`, `low52`, `high52`, `volume`, `avgVolume`, `pe`, `peX`, `dividend`, `yield`, `mktCap`, `mktCapX` and `advancing`.

The expression **must** return a boolean value, otherwise it will fail.

For detailed information about the syntax, please refer to [Knetic/govaluate#what-operators-and-types-does-this-support](https://github.com/Knetic/govaluate#what-operators-and-types-does-this-support).

To clear the filter, press `Shift+F`.

You can specify the profile you want to use by passing ``-profile <filename>`` to the command-line.

### Contributing ###

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/michaeldv/mop?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Mop is my personal project I came up with to learn Go programming. Your
comments, suggestions, and contributions are welcome.

* Fork the project on Github.
* Make your feature addition or bug fix.
* Pull requests accepted.
* Commit, do not change program version, or commit history.


### License ###
Copyright (c) 2013-2019 by Michael Dvorkin and contributors. All Rights Reserved.
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

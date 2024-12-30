### mop: track stocks the hacker way
A command-line utility that displays current information on select markets and individual stocks.

![image](https://user-images.githubusercontent.com/12674437/144474220-5f35f893-6860-4ba5-9c3a-58b80df16255.png)

### Installing mop from source

Ensure GO language is installed. Download from: https://go.dev/dl/ and the $GOPATH is set then:

```
git clone https://github.com/mop-tracker/mop
cd mop
go build ./cmd/mop
./mop
```

### Using mop
By default, mop refreshes quotes and market data on a 5 minute (600s) interval.  These values can be changed in the .moprc file.

For demonstration purposes mop comes preconfigured with a number of stock tickers. You can easily change the default list by using the following keyboard commands:

```
   +                  Add stocks to list
   -                  Remove stocks from list
   ? h H              Display this help screen
   f                  Set filtering expression
   F                  Unset filtering expression
   g G                Group stocks by advancing/declining issues
   o                  Change column sort order
   p P                Pause market data and stock updates
   t                  Toggle timestamp on/off
   Mouse Scroll       Scroll up/down
   PgUp/PgDn          Scroll up/down
   Up/Down arrows     Scroll up
   j J                Scroll up
   k K                Scroll down
   q esc              Quit mop
```

When prompted please enter comma-delimited list of stock tickers.

The list and other settings are stored in the profile file (default: ``.moprc`` in your ``$HOME`` directory).

### No Timestamp
![image](https://github.com/mop-tracker/mop/assets/12674437/f753ce40-c9b2-4ed1-bf51-2a79c83f3f1c)

### With Timestamp
![image](https://github.com/mop-tracker/mop/assets/12674437/8d732111-d25a-425f-bdbf-f0ada6e04b75)

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

### Options and settings

In `~/.moprc`:

```
    "RowShading": true,
```

will cause row shading on alternate lines.

### Contributing
* Pull requests accepted.

### License
Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
"mike" + "@dvorkin" + ".net" || "twitter.com/mid"

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

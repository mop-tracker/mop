// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main
import  (
        "fmt"
        "bytes"
        "regexp"
        "strings"
	"text/template"
)
type Market struct {
        Dow       map[string]string
        Nasdaq    map[string]string
        Sp500     map[string]string
        Advances  map[string]string
        Declines  map[string]string
        Unchanged map[string]string
        Highs     map[string]string
        Lows      map[string]string
}
//-----------------------------------------------------------------------------
func main() {
        html := `
...
<table id="yfimktsumm" border="0" cellspacing="0" cellpadding="0" class="rts" summary="Market Summary">
<thead>
<th class="first">Symbol</th>
<th>Last</th>
<th>Change</th>
</thead>
<tbody>
<tr class="e">
<td><a href="/q?s=%5EDJI">Dow</a></td><td class="idx">
<span id="yfs_l10_^dji">15,300.34</span></td>
<td class="cu"><span id="yfs_c10_^dji"><img width="10" height="14" style="margin-right:-2px;" border="0" src="http://l.yimg.com/os/mit/media/m/base/images/transparent-1093278.png" class="pos_arrow" alt="Up"> <b style="color:#008800;">75.65</b></span><span id="yfs_p20_^dji"><b style="color:#008800;"> (0.50%)</b></span></td>
</tr>
<tr class="">
<td><a href="/q?s=%5EIXIC">Nasdaq</a></td>
<td class="idx"><span id="yfs_l10_^ixic">3,504.26</span></td>
<td class="cu"><span id="yfs_c10_^ixic"><img width="10" height="14" style="margin-right:-2px;" border="0" src="http://l.yimg.com/os/mit/media/m/base/images/transparent-1093278.png" class="pos_arrow" alt="Up"> <b style="color:#008800;">19.43</b></span><span id="yfs_p20_^ixic"><b style="color:#008800;"> (0.56%)</b></span></td>
</tr>
<tr class="e">
<td><a href="/q?s=%5EGSPC">S&amp;P 500</a></td>
<td class="idx"><span id="yfs_l10_^gspc">1,652.32</span></td>
<td class="cu"><span id="yfs_c10_^gspc"><img width="10" height="14" style="margin-right:-2px;" border="0" src="http://l.yimg.com/os/mit/media/m/base/images/transparent-1093278.png" class="pos_arrow" alt="Up"> <b style="color:#008800;">11.86</b></span><span id="yfs_p20_^gspc"><b style="color:#008800;"> (0.72%)</b></span></td>
</tr>
<tr class="">
<td><a href="/q?s=%5ETNX">10-Yr Bond</a></td>
<td class="idx"><span id="yfs_l10_^tnx">2.63</span>%
                  </td>
<td class="idx"><span id="yfs_c10_^tnx"><img width="10" height="14" style="margin-right:-2px;" border="0" src="http://l.yimg.com/os/mit/media/m/base/images/transparent-1093278.png" class="neg_arrow" alt="Down"> <b style="color:#cc0000;">0.02</b></span></td>
</tr>
<tr class="e">
<td colspan="2"><a href="/q?s=%5ETV.N">NYSE Volume</a></td>
<td class="idx"><span id="yfs_lt0_^tv.n">3,490,723,250.00</span></td>
</tr>
<tr class="">
<td colspan="2"><a href="/q?s=%5ETV.O">Nasdaq Volume...</a></td>
<td class="idx"><span id="yfs_lt0_^tv.o">1,594,900,625.00</span></td>
</tr>
</tbody>
</table>
<div id="ma">
<strong>Indices:</strong> <a href="http://us.rd.yahoo.com/finance/finhome/usindices/*http://finance.yahoo.com/indices?e=dow_jones">US</a> - <a href="http://us.rd.yahoo.com/finance/finhome/worldindices/*http://finance.yahoo.com/intlindices?e=americas">World</a> | <a href="http://us.rd.yahoo.com/finance/finhome/mostactives/*http://finance.yahoo.com/actives?e=o">Most Actives</a>
</div>
</div>
</div>
<div class="tba"><h3>Advances &amp; Declines</h3></div>
<div class="ob">
<table id="yfimktsumm" border="0" cellspacing="0" cellpadding="0" class="rts">
<thead>
<th class="first">&nbsp;</th>
<th>NYSE</th>
<th>NASDAQ</th>
</thead>
<tbody>
<tr class="e">
<td class="first">Advances</td>
<td align="right">2,992
                                (72%)
                            </td>
<td align="right">1,445
                                (57%)
                            </td>
</tr>
<tr>
<td class="first">Declines</td>
<td align="right">1,040
                                (25%)
                            </td>
<td align="right">950
                                (38%)
                            </td>
</tr>
<tr class="e">
<td class="first">Unchanged</td>
<td align="right">113
                                (3%)
                            </td>
<td align="right">128
                                (5%)
                            </td>
</tr>
<tr>
<td class="first">Up Vol*</td>
<td align="right">2,582
        (74%)
      </td>
<td align="right">950
        (60%)
      </td>
</tr>
<tr class="e">
<td class="first">Down Vol*</td>
<td align="right">863
        (25%)
      </td>
<td align="right">625
        (39%)
      </td>
</tr>
<tr>
<td class="first">Unch. Vol*</td>
<td align="right">46
        (1%)
      </td>
<td align="right">20
        (1%)
      </td>
</tr>
<tr class="e">
<td class="first">New Hi's</td>
<td align="right">350</td>
<td align="right">314</td>
</tr>
<tr>
<td class="first">New Lo's</td>
<td align="right">117</td>
<td align="right">19</td>
</tr>
</tbody>
</table>
...
<table id="yfimktsumm" border="0" cellspacing="0" cellpadding="0" class="rts">
<thead>
<th class="first">NYSE</th>
<th>LAST</th>
<th>CHANGE</th>
</thead>
`
        start := strings.Index(html, `<table id="yfimktsumm"`)
        finish := strings.LastIndex(html, `<table id="yfimktsumm"`)
        html = strings.Replace(html[start:finish], "\n", "", -1)
        html = strings.Replace(html, "&amp;", "&", -1)

        const any     = `\s*<.+?>`
        const some    = `<.+?`
        const space   = `\s*`
        const arrow   = `"(Up|Down)">\s*`
        const price   = `([\d\.,]+)`
        const percent = `\(([\d\.,%]+)\)`

        regex := []string{
                "(Dow)",       any, price, some, arrow, any, price, some, percent, any,
                "(Nasdaq)",    any, price, some, arrow, any, price, some, percent, any,
                "(S&P 500)",   any, price, some, arrow, any, price, some, percent, any,
                "(Advances)",  any, price, space, percent, any, price, space, percent, any,
                "(Declines)",  any, price, space, percent, any, price, space, percent, any,
                "(Unchanged)", any, price, space, percent, any, price, space, percent, any,
                "(New Hi's)",  any, price, any, price, any,
                "(New Lo's)",  any, price, any, price, any,
        }

        re := regexp.MustCompile(strings.Join(regex, ""))
        matches := re.FindAllStringSubmatch(html, -1)

        if len(matches) > 0 {
                fmt.Printf("%d matches\n", len(matches[0]))
                for i, str := range matches[0][1:] {
                        fmt.Printf("%d) [%s]\n", i, str)
                }
        } else {
                println("No matches")
        }

        m := Market{
                Dow:       make(map[string]string),
                Nasdaq:    make(map[string]string),
                Sp500:     make(map[string]string),
                Advances:  make(map[string]string),
                Declines:  make(map[string]string),
                Unchanged: make(map[string]string),
                Highs:     make(map[string]string),
                Lows:      make(map[string]string),
        }
        m.Dow[`name`]          = matches[0][1]
        m.Dow[`latest`]        = matches[0][2]
        m.Dow[`change`]        = matches[0][4]
        if matches[0][3] == "Up" {
                m.Dow[`change`] = "+" + matches[0][4]
                m.Dow[`percent`] = "+" + matches[0][5]
        } else {
                m.Dow[`change`] = "-" + matches[0][4]
                m.Dow[`percent`] = "-" + matches[0][5]
        }

        m.Nasdaq[`name`]       = matches[0][6]
        m.Nasdaq[`latest`]     = matches[0][7]
        if matches[0][8] == "Up" {
                m.Nasdaq[`change`] = "+" + matches[0][9]
                m.Nasdaq[`percent`] = "+" + matches[0][10]
        } else {
                m.Nasdaq[`change`] = "-" + matches[0][9]
                m.Nasdaq[`percent`] = "-" + matches[0][10]
        }

        m.Sp500[`name`]        = matches[0][11]
        m.Sp500[`latest`]      = matches[0][12]
        if matches[0][13] == "Up" {
                m.Sp500[`change`] = "+" + matches[0][14]
                m.Sp500[`percent`] = "+" + matches[0][15]
        } else {
                m.Sp500[`change`] = "-" + matches[0][14]
                m.Sp500[`percent`] = "-" + matches[0][15]
        }

        m.Advances[`name`]     = matches[0][16]
        m.Advances[`nyse`]     = matches[0][17]
        m.Advances[`nysep`]    = matches[0][18]
        m.Advances[`nasdaq`]   = matches[0][19]
        m.Advances[`nasdaqp`]  = matches[0][20]

        m.Declines[`name`]     = matches[0][21]
        m.Declines[`nyse`]     = matches[0][22]
        m.Declines[`nysep`]    = matches[0][23]
        m.Declines[`nasdaq`]   = matches[0][24]
        m.Declines[`nasdaqp`]  = matches[0][25]

        m.Unchanged[`name`]    = matches[0][26]
        m.Unchanged[`nyse`]    = matches[0][27]
        m.Unchanged[`nysep`]   = matches[0][28]
        m.Unchanged[`nasdaq`]  = matches[0][29]
        m.Unchanged[`nasdaqp`] = matches[0][30]

        m.Highs[`name`]        = matches[0][31]
        m.Highs[`nyse`]        = matches[0][32]
        m.Highs[`nasdaq`]      = matches[0][33]
        m.Lows[`name`]         = matches[0][34]
        m.Lows[`nyse`]         = matches[0][35]
        m.Lows[`nasdaq`]       = matches[0][36]
        fmt.Printf("%q\n", m)
        println(Format(m))
}

//-----------------------------------------------------------------------------
func Format(m Market) string {
	markup := `{{.Dow.name}}: {{.Dow.change}} ({{.Dow.percent}}) at {{.Dow.latest}}, `
        markup += `{{.Sp500.name}}: {{.Sp500.change}} ({{.Sp500.percent}}) at {{.Sp500.latest}}, `
        markup += `{{.Nasdaq.name}}: {{.Nasdaq.change}} ({{.Nasdaq.percent}}) at {{.Nasdaq.latest}}`
        markup += "\n"
        markup += `{{.Advances.name}}: {{.Advances.nyse}} ({{.Advances.nysep}}) on NYSE and {{.Advances.nasdaq}} ({{.Advances.nasdaqp}}) on Nasdaq. `
        markup += `{{.Declines.name}}: {{.Declines.nyse}} ({{.Declines.nysep}}) on NYSE and {{.Declines.nasdaq}} ({{.Declines.nasdaqp}}) on Nasdaq`
        markup += "\n"
        markup += `New highs: {{.Highs.nyse}} on NYSE and {{.Highs.nasdaq}} on Nasdaq. `
        markup += `New lows: {{.Lows.nyse}} on NYSE and {{.Lows.nasdaq}} on Nasdaq.`
	template, err := template.New("screen").Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, m)
	if err != nil {
		panic(err)
	}

        return buffer.String()
}

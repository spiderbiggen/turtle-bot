package nyaa

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_parseFeed(t *testing.T) {
	inp := `
<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:nyaa="https://nyaa.si/xmlns/nyaa" version="2.0">
	<channel>
		<title>Nyaa - &#34;[SubsPlease] (1080p)&#34; - Torrent File RSS</title>
		<description>RSS Feed for &#34;[SubsPlease] (1080p)&#34;</description>
		<link>https://nyaa.si/</link>
		<atom:link href="https://nyaa.si/?page=rss" rel="self" type="application/rss+xml" />
		<item>
			<title>[SubsPlease] The Prince of Tennis II - U-17 World Cup - 03 (1080p) [4D6884BC].mkv</title>
				<link>https://nyaa.si/download/1554818.torrent</link>
				<guid isPermaLink="true">https://nyaa.si/view/1554818</guid>
				<pubDate>Wed, 20 Jul 2022 16:01:38 -0000</pubDate>

				<nyaa:seeders>188</nyaa:seeders>
				<nyaa:leechers>28</nyaa:leechers>
				<nyaa:downloads>379</nyaa:downloads>
				<nyaa:infoHash>8b2308090dab1cd8a8bd8059933405946d423273</nyaa:infoHash>
			<nyaa:categoryId>1_2</nyaa:categoryId>
			<nyaa:category>Anime - English-translated</nyaa:category>
			<nyaa:size>1.3 GiB</nyaa:size>
			<nyaa:comments>0</nyaa:comments>
			<nyaa:trusted>Yes</nyaa:trusted>
			<nyaa:remake>No</nyaa:remake>
			<description><![CDATA[<a href="https://nyaa.si/view/1554818">#1554818 | [SubsPlease] The Prince of Tennis II - U-17 World Cup - 03 (1080p) [4D6884BC].mkv</a> | 1.3 GiB | Anime - English-translated | 8b2308090dab1cd8a8bd8059933405946d423273]]></description>
		</item>
	</channel>
</rss>
`
	reader := strings.NewReader(inp)

	expectedDate := time.Date(2022, time.July, 20, 16, 1, 38, 0, time.UTC)
	expected := &rssResult{
		Episode: Episode{
			AnimeTitle: "The Prince of Tennis II - U-17 World Cup",
			Number:     0x3,
			Decimal:    0x0,
			Version:    0x0,
		},
		Download: Download{
			Comments:      "https://nyaa.si/view/1554818",
			Resolution:    "1080p",
			Torrent:       "https://nyaa.si/download/1554818.torrent",
			FileName:      "[SubsPlease] The Prince of Tennis II - U-17 World Cup - 03 (1080p) [4D6884BC].mkv",
			PublishedDate: &expectedDate,
		},
	}

	gotList, err := parseFeed(reader)
	if err != nil {
		t.Errorf("parseFeed() error = %v", err)
		return
	}
	if len(gotList) != 1 {
		t.Errorf("unexpected length got = %d, want = 1", len(gotList))
	}
	got := gotList[0]
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("parseFeed()\ngot  = %#v\nwant = %#v", got, expected)
	}
}

package main

// many of these are pulled directly from Technorati's Top 100 http://technorati.com/blogs/top100/
var builtinFeeds = []Feed{
	Feed{
		URL:     "http://www.reddit.com/.rss",
		Title:   "Reddit",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.theonion.com/theonion/daily",
		Title:   "The Onion",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.gawker.com/gawker/full",
		Title:   "Gawker",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/CalculatedRisk",
		Title:   "Calculated Risk",
		Default: false,
	},
	Feed{
		URL:     "http://www.guardian.co.uk/news/datablog/rss",
		Title:   "Guardian Datablog",
		Default: false,
	},
	Feed{
		URL:     "http://www.economist.com/blogs/babbage/index.xml",
		Title:   "Babbage",
		Default: false,
	},
	Feed{
		URL:     "http://www.wired.com/threatlevel/feed/",
		Title:   "Wired Threat Level",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.huffingtonpost.com/FeaturedPosts",
		Title:   "Huffington Post Featured",
		Default: false,
	},
	Feed{
		URL:     "http://www.tmz.com/rss.xml",
		Title:   "TMZ",
		Default: false,
	},
	Feed{
		URL:     "http://www.guardian.co.uk/world/us-news-blog/rss",
		Title:   "The Guardian US News",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/hotair/main",
		Title:   "Hot Air",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.mashable.com/Mashable",
		Title:   "Mashable",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/TechCrunch",
		Title:   "TechCrunch",
		Default: false,
	},
	Feed{
		URL:     "http://www.theblaze.com/stories/feed/",
		Title:   "The Blaze",
		Default: false,
	},
	Feed{
		URL:     "http://opinionator.blogs.nytimes.com/feed/",
		Title:   "The Opinionater",
		Default: false,
	},
	Feed{
		URL:     "http://politicalticker.blogs.cnn.com/feed/",
		Title:   "Political Ticker",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.nbcnews.com/feeds/topstories",
		Title:   "NBC News",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/foxnews/latest",
		Title:   "FOX News",
		Default: false,
	},
	Feed{
		URL:     "http://www.destructoid.com/?mode=atom",
		Title:   "Destructoid",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/cnet/NnTv",
		Title:   "CNet",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.boingboing.net/boingboing/iBag",
		Title:   "Boing Boing",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.gawker.com/lifehacker/full",
		Title:   "Lifehacker",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/43folders",
		Title:   "43 Folders",
		Default: false,
	},
	Feed{
		URL:     "http://www.buzzfeed.com/index.xml",
		Title:   "Buzz Feed",
		Default: false,
	},
	Feed{
		URL:     "http://www.engadget.com/rss.xml",
		Title:   "Engadget",
		Default: false,
	},
	Feed{
		URL:     "http://bits.blogs.nytimes.com/feed/",
		Title:   "Bits",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/Best-Of-The-Atlantic",
		Title:   "The Atlantic",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/Asymco",
		Title:   "Asymco",
		Default: false,
	},
	Feed{
		URL:     "http://rss.csmonitor.com/csmonitor/connectingthedots",
		Title:   "Connecting the Dots",
		Default: false,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/alistapart/main",
		Title:   "A List Apart",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/ommalik",
		Title:   "GigaOM",
		Default: true,
	},
	Feed{
		URL:     "http://www.marco.org/rss",
		Title:   "Marco Arment",
		Default: true,
	},
	Feed{
		URL:     "http://www.newyorker.com/online/blogs/newsdesk/rss.xml",
		Title:   "The New Yorker News Desk",
		Default: true,
	},
	Feed{
		URL:     "http://www.theverge.com/rss/index.xml",
		Title:   "The Verge",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.theonion.com/avclub/newswire/",
		Title:   "AV Club Newswire",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/marginalrevolution",
		Title:   "Marginal Revolution",
		Default: true,
	},
	Feed{
		URL:     "http://daringfireball.net/index.xml",
		Title:   "Daring Fireball",
		Default: true,
	},
	Feed{
		URL:     "http://scripting.com/rss.xml",
		Title:   "Scripting News",
		Default: true,
	},
	Feed{
		URL:     "http://buzzmachine.com/feed/",
		Title:   "Buzz Machine",
		Default: true,
	},
	Feed{
		URL:     "http://blogs.wsj.com/numbersguy/feed/",
		Title:   "The Number's Guy",
		Default: true,
	},
	Feed{
		URL:     "http://flowingdata.com/feed/",
		Title:   "Flowing Data",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.kottke.org/main",
		Title:   "Jason Kottke",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.washingtonpost.com/rss/rss_ezra-klein",
		Title:   "Wonk Blog",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/538dotcom",
		Title:   "538",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.arstechnica.com/arstechnica/features",
		Title:   "Ars Features",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.dashes.com/AnilDash",
		Title:   "Anil Dash",
		Default: true,
	},
	Feed{
		URL:     "http://feed.torrentfreak.com/Torrentfreak/",
		Title:   "Torrent Freak",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/blogspot/MKuf",
		Title:   "Official Google Blog",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/GoogleOperatingSystem",
		Title:   "Google Operating System",
		Default: true,
	},
	Feed{
		URL:     "http://feeds.feedburner.com/thebrowser/xrdJ",
		Title:   "The Browser",
		Default: true,
	},
}

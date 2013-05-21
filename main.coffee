if not KeyEvent?
    KeyEvent =
        DOM_VK_CANCEL: 3
        DOM_VK_HELP: 6
        DOM_VK_BACK_SPACE: 8
        DOM_VK_TAB: 9
        DOM_VK_CLEAR: 12
        DOM_VK_RETURN: 13
        DOM_VK_ENTER: 14
        DOM_VK_SHIFT: 16
        DOM_VK_CONTROL: 17
        DOM_VK_ALT: 18
        DOM_VK_PAUSE: 19
        DOM_VK_CAPS_LOCK: 20
        DOM_VK_ESCAPE: 27
        DOM_VK_SPACE: 32
        DOM_VK_PAGE_UP: 33
        DOM_VK_PAGE_DOWN: 34
        DOM_VK_END: 35
        DOM_VK_HOME: 36
        DOM_VK_LEFT: 37
        DOM_VK_UP: 38
        DOM_VK_RIGHT: 39
        DOM_VK_DOWN: 40
        DOM_VK_PRINTSCREEN: 44
        DOM_VK_INSERT: 45
        DOM_VK_DELETE: 46
        DOM_VK_0: 48
        DOM_VK_1: 49
        DOM_VK_2: 50
        DOM_VK_3: 51
        DOM_VK_4: 52
        DOM_VK_5: 53
        DOM_VK_6: 54
        DOM_VK_7: 55
        DOM_VK_8: 56
        DOM_VK_9: 57
        DOM_VK_SEMICOLON: 59
        DOM_VK_EQUALS: 61
        DOM_VK_A: 65
        DOM_VK_B: 66
        DOM_VK_C: 67
        DOM_VK_D: 68
        DOM_VK_E: 69
        DOM_VK_F: 70
        DOM_VK_G: 71
        DOM_VK_H: 72
        DOM_VK_I: 73
        DOM_VK_J: 74
        DOM_VK_K: 75
        DOM_VK_L: 76
        DOM_VK_M: 77
        DOM_VK_N: 78
        DOM_VK_O: 79
        DOM_VK_P: 80
        DOM_VK_Q: 81
        DOM_VK_R: 82
        DOM_VK_S: 83
        DOM_VK_T: 84
        DOM_VK_U: 85
        DOM_VK_V: 86
        DOM_VK_W: 87
        DOM_VK_X: 88
        DOM_VK_Y: 89
        DOM_VK_Z: 90
        DOM_VK_CONTEXT_MENU: 93
        DOM_VK_NUMPAD0: 96
        DOM_VK_NUMPAD1: 97
        DOM_VK_NUMPAD2: 98
        DOM_VK_NUMPAD3: 99
        DOM_VK_NUMPAD4: 100
        DOM_VK_NUMPAD5: 101
        DOM_VK_NUMPAD6: 102
        DOM_VK_NUMPAD7: 103
        DOM_VK_NUMPAD8: 104
        DOM_VK_NUMPAD9: 105
        DOM_VK_MULTIPLY: 106
        DOM_VK_ADD: 107
        DOM_VK_SEPARATOR: 108
        DOM_VK_SUBTRACT: 109
        DOM_VK_DECIMAL: 110
        DOM_VK_DIVIDE: 111
        DOM_VK_F1: 112
        DOM_VK_F2: 113
        DOM_VK_F3: 114
        DOM_VK_F4: 115
        DOM_VK_F5: 116
        DOM_VK_F6: 117
        DOM_VK_F7: 118
        DOM_VK_F8: 119
        DOM_VK_F9: 120
        DOM_VK_F10: 121
        DOM_VK_F11: 122
        DOM_VK_F12: 123
        DOM_VK_F13: 124
        DOM_VK_F14: 125
        DOM_VK_F15: 126
        DOM_VK_F16: 127
        DOM_VK_F17: 128
        DOM_VK_F18: 129
        DOM_VK_F19: 130
        DOM_VK_F20: 131
        DOM_VK_F21: 132
        DOM_VK_F22: 133
        DOM_VK_F23: 134
        DOM_VK_F24: 135
        DOM_VK_NUM_LOCK: 144
        DOM_VK_SCROLL_LOCK: 145
        DOM_VK_COMMA: 188
        DOM_VK_PERIOD: 190
        DOM_VK_SLASH: 191
        DOM_VK_BACK_QUOTE: 192
        DOM_VK_OPEN_BRACKET: 219
        DOM_VK_BACK_SLASH: 220
        DOM_VK_CLOSE_BRACKET: 221
        DOM_VK_QUOTE: 222
        DOM_VK_META: 224

OK = 200
LIST = 1
COUNT = 10
TIMEOUT = 64

$.ajaxSetup
	async: false

$.fn.exists = -> @length > 0

$.fn.scrollTo = (target, options, callback) ->
	if typeof options is 'function' and arguments_.length is 2
		callback = options
		options = target
	settings = $.extend(
		scrollTarget: target
		offsetTop: 50
		duration: 500
		easing: 'swing'
		, options)
	@each ->
		scrollPane = $(this)
		scrollTarget = (if (typeof settings.scrollTarget is 'number') then settings.scrollTarget else $(settings.scrollTarget))
		scrollY = (if (typeof scrollTarget is 'number') then scrollTarget else scrollTarget.offset().top + scrollPane.scrollTop() - parseInt(settings.offsetTop))
		scrollPane.animate
			scrollTop: scrollY
			, parseInt(settings.duration), settings.easing, ->
			callback.call this  if typeof callback is 'function'

$.extend
	postJSON: (url, data, callback) ->
		$.ajax
			type: 'POST'
			url: url
			data: JSON.stringify data
			success: callback
			dataType: 'json'
			contentType: 'application/json'
			processData: false
			async: false

show = (element) ->
	$.getJSON '/article?id=' + encodeURIComponent(element.attr('id'))
	element.children('.article-content').slideToggle()

hide = (element) ->
	element.children('.article-content').slideToggle()

addArticle = (data) ->
	$('<article/>').
		addClass('article').
		addClass('unread').
		attr('id', data['ID']).
		hide().
		append(
			$('<a/>').
				addClass('go').
				addClass('action').
				attr('target', '_blank').
				attr('href', '/article?url=' + data['URL'] + '&id=' + data['ID']).
				html('▶')
		).
		append(
			$('<header/>').
				addClass('article-header').
				append(
					$('<a/>').
						attr('target', '_blank').
						attr('href', data['URL']).
						html(data['Title']).
						addClass('article-link').
						click (event) ->
							event.preventDefault()
							#if event.ctrlKey
							window.open $('.current').children('.go').attr('href'), '_blank'
							#else
							#show($(this).parent().parent())
							false
				)
		).
		append(
			$('<div/>').
				addClass('article-content').
				html($.parseHTML(data['Summary'])).
				click (event) ->
					event.preventDefault()
					false
		)

addArticles = (object) ->
	list = []
	return list if not object?
	for article in object['Articles']
		continue if not article['ID']?
		element = addArticle article
		list.push element if element?
	list

makeCurrent = (articles, current) ->
	makeArticle articles
	$(document.getElementById(article.attr('id'))).addClass('current').show() for article in articles.slice(0, LIST)
	removeCurrent current
	if $('.current').index() is 0
		$('#prev').hide()
	else
		$('#prev').show()
	$('body').scrollTo($('.current').offset().top) if $('.current').exists()

makeArticle = (articles) ->
	for article in articles
		article.hide().appendTo '#articles'

nextArticle = (count, timeout, fun, current) ->
	$.getJSON('/article?output=json&count=' + count, (data) ->
		if data['URL']?
			timeout *= 2
			setTimeout nextArticle, timeout, count, timeout, fun, current
		else
			articles = addArticles data
			newarticles = []
			for article in articles
				if not $(document.getElementById(article.attr('id'))).exists()
					newarticles.push article
			if newarticles.length is 0
				timeout *= 2
				setTimeout nextArticle, timeout, count, timeout, fun, current
			else
				localArticles = localStorage.getObj('articles')
				localArticles = localArticles.concat(data.Articles)
				localStorage.setObj('articles', localArticles)
				fun(newarticles, current)
#			$('#next').hide()
#			$('.current').
#				css({position: 'fixed', top: $(document).height(), left: $('#articles').offset().left}).
#				show().
#				animate {top: $('#articles').offset().top}, ->
#								$(this).css({position: 'static'})
#								$('#next').show()
	).
		fail -> setTimeout nextArticle, timeout, count, timeout, fun, current

removeCurrent = (current) ->
	markAsRead current
	current.removeClass 'current'

next = ->
	if $('#next').is ':visible'
		$('#next').hide()
		current = $('.current')
		index = current.index()
		if index is -1
			index = $('.unread').first().index()
			if index is -1
				index = $('.read').last().index()
				if index is -1
					index = 0
		if index + LIST < $('#articles').children().length
			$('#articles').children().slice(index + 1, index + LIST + 1).addClass('current').show()
			removeCurrent current
			index = $('.current').index()
			if index is 0
				$('#prev').hide()
			else
				$('#prev').show()
			$('body').scrollTo($('.current').offset().top) if $('.current').exists()
		else
			nextArticle COUNT, TIMEOUT, makeCurrent, current, index
		setTimeout nextArticle, TIMEOUT, COUNT, TIMEOUT, makeArticle if index >= $('#articles').children().last().index() / 2
		$('#next').show()

prev = ->
	if $('#prev').is ':visible'
		$('#prev').hide()
		index = $('.current').index()
		if index >= LIST
			markAsRead $('.current')
			$('.current').removeClass 'current'
			$('#articles').children().slice(index - LIST, index).addClass('current').show()
#		$('.current').
#			show().
#			css({position: 'fixed'}).
#			animate {top: $(document).height(), left: $('#articles').offset().left}, ->
#				$(this).css({position: 'static'})
#				$('#next').show()
#				$('#prev').hide()
#					css({position: 'fixed', top: 0, left: $('#articles').offset().left}).
#					animate {top: $('#articles').offset().top}, ->
#						$(this).css({position: 'static'}).
#							addClass('current').show()
#						if $('.current').first().index() < LIST
#							$('#prev').hide()
#						else
#							$('#prev').show()
#				$('#next').show()
		if index - LIST is 0
			$('#prev').hide()
		else
			$('#prev').show()
#			$('#next').hide()

Storage.prototype.setObj = (key, obj) -> @setItem(key, JSON.stringify(obj))

Storage.prototype.getObj = (key) -> JSON.parse(@getItem(key))

markAsRead = (elements) ->
	localArticles = localStorage.getObj('articles')
	elements.each ->
		$(this).
			hide().
			removeClass('unread').
			addClass('read')
		for article, i in localArticles
			if article['ID'] is $(this).attr('id')
				localArticles.remove(i)
				break
	localStorage.setObj('articles', localArticles)

addFeed = (url) ->
	data = Feeds: [
		Subscribed: true
		URL: url
	]
	$.postJSON '/feed', data
	$(document.getElementById(url)).remove()
	location.reload()

removeFeed = (url) ->
	data = Feeds: [
		Subscribed: false
		URL: url
	]
	$.postJSON '/feed', data
	$(document.getElementById(url)).remove()
	location.reload()

Array.prototype.remove = (from, to) ->
	rest = this.slice((to || from) + 1 || this.length)
	this.length = from < 0 ? this.length + from : from
	return this.push.apply(this, rest)

$ ->
	$('<section/>').
		attr('id', 'articles').
		insertBefore('#next') if not $('#articles').exists()

	localArticles = localStorage.getObj('articles')
	if localArticles? and localArticles.length > 0
		articles = []
		for article in localArticles
			articles.push(addArticle(article))
		makeCurrent(articles, $('.current'))
	else
		localStorage.setObj('articles', [])

	$('#next').show()
	next() if not $('.unread').exists()

	$(window).keyup (event) ->
		switch event.which
			when KeyEvent.DOM_VK_J, KeyEvent.DOM_VK_N, KeyEvent.DOM_VK_NUMPAD2, KeyEvent.DOM_VK_NUMPAD3 #KeyEvent.DOM_VK_SPACE, KeyEvent.DOM_VK_PAGEDOWN, KeyEvent.DOM_VK_DOWN, 
				event.preventDefault()
				if event.shiftKey
					prev()
				else
					next()
			when KeyEvent.DOM_VK_PAGEUP, KeyEvent.DOM_VK_K, KeyEvent.DOM_VK_P, KeyEvent.DOM_VK_NUMPAD8, KeyEvent.DOM_VK_NUMPAD9 # KeyEvent.DOM_VK_UP,
				event.preventDefault()
				if event.shiftKey
					next()
				else
					prev()
			when KeyEvent.DOM_VK_RIGHT, KeyEvent.DOM_VK_NUMPAD5, KeyEvent.DOM_VK_ENTER, KeyEvent.DOM_VK_RETURN
				event.preventDefault()
				if event.shiftKey
					window.open $('.current').children('.go').attr('href'), '_blank'
				else
					if $('.current').children('.article-content').is(':visible')
						window.open $('.current').children('.go').attr('href'), '_blank'
					else
						show($('.current'))
			when KeyEvent.DOM_VK_LEFT, KeyEvent.DOM_VK_NUMPAD6
				event.preventDefault()
				hide($('.current')) if $('.current').children('.article-content').is(':visible')
			else
				return true
		return false

	$('#prev').click (event) ->
		event.preventDefault()
		prev()
		false

	$('#next').click (event) ->
		event.preventDefault()
		next()
		false

	$('.subscribe').click (event) ->
		event.preventDefault()
		addFeed event.currentTarget.parentNode.getAttribute 'id'
		false

	$('.unsubscribe').click (event) ->
		event.preventDefault()
		removeFeed event.currentTarget.parentNode.getAttribute 'id'
		false

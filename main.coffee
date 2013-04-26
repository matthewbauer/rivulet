OK = 200
NUMBER = 1
TIMEOUT = 1000

jQuery.fn.exists = -> @length > 0

jQuery.extend
	postJSON: (url, data, callback) ->
		jQuery.ajax
			type: 'POST'
			url: url
			data: JSON.stringify(data)
			success: callback
			dataType: 'json'
			contentType: 'application/json'
			processData: false
			async: false

show = (element) ->
	$.getJSON '/article?id=' + encodeURIComponent()
	element.children('.article-content').slideToggle()

addArticle = (data) ->
	if not $('#articles').exists()
		$('<section/>').
			attr('id', 'articles').
			insertBefore('#next')
	if $(document.getElementById(data['ID'])).exists()
		return false
	$('<article/>').
		addClass('article').
		addClass('current').
		addClass('unread').
		attr('id', data['ID']).
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
							show($(this).parent().parent())
							false
				)
		).
		append(
			$('<div/>').
				addClass('article-content').
				html(data['Summary'])
		).
		appendTo('#articles')
	$('#prev').show() if $('.read').exists()
	return true

addArticles = (object) ->
	$('#next').show()
	for article in object['Articles']
		added = addArticle article
		if not added
			return false
	return true

refresh = ->
	nextArticle NUMBER

nextArticle = (number) ->
	$.getJSON '/article?output=json&number=' + number, (data) ->
		if data['URL']? or not data['Articles']? or data['Articles'].length == 0
			$('#next').hide()
			refresh()
		else
			added = addArticles data
			if not added
				$('#next').hide()
				markAsRead $('.current')
				setTimeout(refresh, TIMEOUT)

next = ->
	$('#prev').show()
	if $('.current').last().is(':last-child')
		markAsRead $('.current')
		$('.current').removeClass('current')
		refresh()
	else
		n = $('.current').next()
		$('.current').removeClass('current')
		n.addClass('current').show()

prev = ->
	if not $('.current').first().is(':first-child')
		markAsRead $('.current.unread')
		p = $('.current').prev()
		$('.current').removeClass('current')
		p.addClass('current').show()
		if $('.current').first().is(':first-child')
			$('#prev').hide()


markAsRead = (elements) ->
	data = Articles: []
	elements.each ->
		data.Articles.push
			ID: $(this).attr('id')
			Read: true
	$.postJSON '/article?read=1', data
	for element in data.Articles
		$(document.getElementById(element['ID'])).
			removeClass('unread').
			addClass('read')

nextArticles = ->
	nextArticle NUMBER

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
		Subscribed: true
		URL: url
	]
	$.postJSON '/feed', data
	$(document.getElementById(url)).remove()
	location.reload()

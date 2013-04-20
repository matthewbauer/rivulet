OK = 200
ENTER = 13
SPACE = 32
PAGEDOWN = 34
NUMBER = 1

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

show = (element) -> element.children('.article-content').slideToggle()

addArticle = (data) ->
	if not $('#articles').exists()
		$('<section/>').
			attr('id', 'articles').
			insertBefore('#next')
	return if data['ID']?
	return if $(document.getElementById(data['ID'])).exists()
	$('<article/>').
		addClass('article').
		addClass('unread').
		attr('id', data['ID']).
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

addArticles = (object) ->
	if object.Articles? and object.Articles.length != 0
		addArticle article for article in object.Articles

nextArticle = (number) ->
	$.getJSON '/article?output=json&number=' + number, (data) ->
		if data['URL']?
			window.location.href = data['URL']
		else
			addArticles data

markAsRead = (data) -> $.postJSON '/article?read=1', data

next = ->
	markAllAsRead()
	window.setTimeout nextArticles, 500

markAllAsRead = ->
	data = Articles: []
	$('.unread').each ->
		data.Articles.push
			ID: $(this).attr('id')
			Read: true
	markAsRead data
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

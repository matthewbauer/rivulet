OK = 200
LIST = 1
NUMBER = 8
TIMEOUT = 1000

jQuery.fn.exists = -> @length > 0

jQuery.extend
	postJSON: (url, data, callback) ->
		jQuery.ajax
			type: 'POST'
			url: url
			data: JSON.stringify data
			success: callback
			dataType: 'json'
			contentType: 'application/json'
			processData: false
			async: false

show = (element) ->
	$.getJSON '/article?id=' + encodeURIComponent()
	element.children('.article-content').slideToggle()

addArticle = (data) ->
	$('<article/>').
		addClass('article').
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
		)

addArticles = (object) ->
	if not object?
		return []
	list = []
	for article in object['Articles']
		element = addArticle article
		if element?
			list.push element.hide()
	list

makeCurrent = (articles) ->
	for article in articles.slice(0, LIST)
		article.show().addClass('current')
	for article in articles
		if not $(document.getElementById(data['ID'])).exists()
			article.appendTo '#articles'

makeArticle = (articles) ->
	for article in articles
		article.appendTo '#articles'

nextArticle = (number, fun) ->
	$.getJSON '/article?output=json&number=' + number, (data) ->
		if data['URL']?
			setTimeout nextArticle, TIMEOUT, NUMBER, makeCurrent
#			$('<article/>').
#				html('the end of the river').
#				addClass('current').
#				appendTo('#articles')
		else
			articles = addArticles data
			fun(articles)
#			$('#next').hide()
#			articles.appendTo('body').show().
#				css({position: 'fixed', top: $(document).height(), left: $('#articles').offset().left}).
#				animate {top: $('#articles').offset().top}, ->
#								$(this).children().each -> $(this).appendTo('#articles')
#								$(this).remove()
#								$('#next').show()

next = ->
	if $('#next').is ':visible'
		$('#next').hide()
		if $('.current').last().index() + LIST + 1 <= $('#articles').children().length
			index = $('.current').last().index()
			markAsRead $('.current')
			$('.current').removeClass 'current'
			$('#articles').children().slice(index + 1, index + LIST + 1).addClass('current').show()
			if index > (NUMBER / 2)
				nextArticle 1, makeArticle
		else
			markAsRead $('.current')
			$('.current').removeClass 'current'
			nextArticle NUMBER, makeCurrent
		if $('.current').first().index() < LIST
			$('#prev').hide()
		else
			$('#prev').show()
		$('#next').show()

prev = ->
	if $('.current').first().index() >= LIST
		index = $('.current').first().index()
		markAsRead $('.current')
		$('.current').removeClass 'current'
		$('#articles').children().slice(index - LIST, index).addClass('current').show()
		if $('.current').first().index() < LIST
			$('#prev').hide()

markAsRead = (elements) ->
	elements.each ->
		$(this).
			hide().
			removeClass('unread').
			addClass('read')

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

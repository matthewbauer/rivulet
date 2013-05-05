OK = 200
LIST = 1
NUMBER = 1
TIMEOUT = 1024

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
	$.getJSON '/article?id=' + encodeURIComponent(element.attr('id'))
	element.children('.article-content').slideToggle()

hide = (element) ->
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
	list = []
	return list if not object?
	for article in object['Articles']
		element = addArticle article
		list.push element.hide() if element?
	list

makeCurrent = (articles) ->
	article.show().addClass('current') for article in articles.slice(0, LIST)
	for article in articles
		if not $(document.getElementById(article.attr('id'))).exists()
			article.appendTo '#articles'

makeArticle = (articles) ->
	for article in articles
		if not $(document.getElementById(article.attr('id'))).exists()
			article.appendTo '#articles'

nextArticle = (number, timeout, fun) ->
	$.getJSON '/article?output=json&number=' + number, (data) ->
		if data['URL']?
			setTimeout nextArticle, timeout * 2, NUMBER, timeout * 2, makeCurrent
		else
			articles = addArticles data
			#div = $('<div/>')
			#for article in articles.slice(0, LIST)
				#div.addClass('current').append(article)
			#$('#next').hide()
			#div.appendTo('body').
				#show().
				#css({position: 'fixed', top: $(document).height(), left: $('#articles').offset().left}).
				#animate {top: $('#articles').offset().top}, ->
								#$(this).children().appendTo($('#articles'))
								#$(this).remove()
								#$('#next').show()
			#for article in articles.slice(LIST + 1)
				#article.appendTo('#articles')
			fun(articles)

next = ->
	if $('#next').is ':visible'
		$('#next').hide()
		if $('.current').last().index() + LIST + 1 <= $('#articles').children().length
			index = $('.current').last().index()
			markAsRead $('.current')
			$('.current').removeClass 'current'
			$('#articles').children().slice(index + 1, index + LIST + 1).addClass('current').show()
			nextArticle 1, TIMEOUT, makeArticle if index > (NUMBER / 2)
		else
			markAsRead $('.current')
			$('.current').removeClass 'current'
			nextArticle NUMBER, TIMEOUT, makeCurrent
		$('#prev').show()
		$('#next').show()

prev = ->
	if $('.current').first().index() >= LIST
		index = $('.current').first().index()
		markAsRead $('.current')
		$('.current').removeClass 'current'
		$('#articles').children().slice(index - LIST, index).addClass('current').show()
		$('#prev').hide() if $('.current').first().index() < LIST

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

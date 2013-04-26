$('.article-link').click (event) ->
	event.preventDefault()
	show($(this).parent().parent())
	false

$(window).keydown (event) ->
	switch event.which
		when KeyEvent.DOM_VK_SPACE, KeyEvent.DOM_VK_PAGEDOWN, KeyEvent.DOM_VK_DOWN, KeyEvent.DOM_VK_J, KeyEvent.DOM_VK_N
			if event.shiftKey
				prev()
			else
				next()
		when KeyEvent.DOM_VK_PAGEUP, KeyEvent.DOM_VK_UP, KeyEvent.DOM_VK_K, KeyEvent.DOM_VK_P
			if event.shiftKey
				next()
			else
				prev()
		when KeyEvent.DOM_VK_ENTER, KeyEvent.DOM_VK_LEFT, KeyEvent.DOM_VK_RIGHT, KeyEvent.DOM_VK_RETURN
			if event.shiftKey
				window.open($('.unread').children('.go').attr('href'), '_blank')
			else
				show($('.unread'))
		else
			return true
	event.preventDefault()
	return false

$('#prev').hide()

$('#prev').click (event) ->
	event.preventDefault()
	prev()
	false

$('#next').click (event) ->
	event.preventDefault()
	next()
	false

nextArticles() if not $('.unread').exists()

$('.article-link').click (event) ->
	event.preventDefault()
	show $(this).parent().parent()
	false

$(window).keyup (event) ->
	switch event.which
		when KeyEvent.DOM_VK_SPACE, KeyEvent.DOM_VK_PAGEDOWN, KeyEvent.DOM_VK_DOWN, KeyEvent.DOM_VK_J, KeyEvent.DOM_VK_N, KeyEvent.DOM_VK_NUMPAD2, KeyEvent.DOM_VK_NUMPAD3
			event.preventDefault()
			if event.shiftKey
				prev()
			else
				next()
		when KeyEvent.DOM_VK_PAGEUP, KeyEvent.DOM_VK_UP, KeyEvent.DOM_VK_K, KeyEvent.DOM_VK_P, KeyEvent.DOM_VK_NUMPAD8, KeyEvent.DOM_VK_NUMPAD9
			event.preventDefault()
			if event.shiftKey
				next()
			else
				prev()
		when KeyEvent.DOM_VK_ENTER, KeyEvent.DOM_VK_RETURN
			event.preventDefault()
			if event.shiftKey
				window.open $('.current').children('.go').attr('href'), '_blank'
			else
				show($('.current'))
		when KeyEvent.DOM_VK_RIGHT, KeyEvent.DOM_VK_NUMPAD5
			event.preventDefault()
			if event.shiftKey
				window.open $('.current').children('.go').attr('href'), '_blank'
			else
				show($('.current')) if not $('.current').children('.article-content').is(':visible')
		when KeyEvent.DOM_VK_LEFT, KeyEvent.DOM_VK_NUMPAD6
			event.preventDefault()
			hide($('.current')) if $('.current').children('.article-content').is(':visible')
		else
			return true
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


$('<section/>').
	attr('id', 'articles').
	insertBefore('#next') if not $('#articles').exists()

next() if not $('.unread').exists()

$('#prev').hide()

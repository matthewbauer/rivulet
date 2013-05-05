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
		when KeyEvent.DOM_VK_ENTER, KeyEvent.DOM_VK_LEFT, KeyEvent.DOM_VK_RIGHT, KeyEvent.DOM_VK_RETURN, KeyEvent.DOM_VK_NUMPAD5, KeyEvent.DOM_VK_NUMPAD6
			event.preventDefault()
			if event.shiftKey
				window.open $('.unread').children('.go').attr('href'), '_blank'
			else
				show($('.unread'))
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

if not $('#articles').exists()
	$('<section/>').
		attr('id', 'articles').
		insertBefore('#next')

next() if not $('.unread').exists()

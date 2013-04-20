$('.article-link').click (event) ->
	event.preventDefault()
	show($(this).parent().parent())
	false

$(window).keypress (event) ->
	if event.which is SPACE or event.which is PAGEDOWN
		event.preventDefault()
		next()
		return false
	if event.which is ENTER
		event.preventDefault()
		show($('.unread'))
		return false
	true

$('#next').click (event) ->
	event.preventDefault()
	next()
	false

nextArticles() if not $('.unread').exists()

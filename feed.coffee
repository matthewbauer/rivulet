$('.subscribe').click (event) ->
	event.preventDefault()
	addFeed event.currentTarget.parentNode.getAttribute 'id'
	false

$('.unsubscribe').click (event) ->
	event.preventDefault()
	removeFeed event.currentTarget.parentNode.getAttribute 'id'
	false

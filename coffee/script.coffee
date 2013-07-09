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
COUNT = 5
TIMEOUT = 64
MAXERROR = 3

online = true

$.ajaxSetup
  async: false
  cache: false

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
      callback.call this if typeof callback is 'function'

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

Array.prototype.remove = (from, to) ->
  rest = this.slice((to || from) + 1 || this.length)
  this.length = from < 0 ? this.length + from : from
  this.push.apply this, rest

Storage.prototype.setObj = (key, obj) -> @setItem(key, JSON.stringify(obj))

Storage.prototype.getObj = (key) -> JSON.parse(@getItem(key))

parseURL = (url) ->
  a = document.createElement('a')
  a.href = url
  host: a.host
  hostname: a.hostname
  pathname: a.pathname
  port: a.port
  protocol: a.protocol
  search: a.search
  hash: a.hash

show = (element) ->
#  $.getJSON '/article?id=' + encodeURIComponent(element.attr('id'))
  element.children('.article-content').slideToggle()

hide = (element) ->

unsubscribe = (url) ->
  _gaq.push(['_trackEvent', 'Feeds', 'Unsubscribe', url])
  $.ajax(
    url: '/feed?url=' + url
    type: 'DELETE'
  )

subscribe = (url) ->
  _gaq.push(['_trackEvent', 'Feeds', 'Subscribe', url])
  $.postJSON '/feed?url=' + url

addArticle = (data) ->
  $('<article/>').
    addClass('article').
    addClass('unread').
    attr('id', data['ID']).
    hide().
    append(
        $('<span/>').
          addClass('label').
          addClass('feedtag').
          append(
            $('<a/>').
              addClass('feedname').
              attr('href', data['FeedURL']).
              attr('target', '_blank').
              html(data['FeedName']).
              click (event) ->
                event.preventDefault()
                window.open('/feed?url=' + $(this).attr('href'), '_self')
                false
          ).
          append(
            $('<a/>').
              addClass('remove').
              attr('title', 'unsubscribe').
              append(
                $('<i/>').
                  addClass('icon-minus-sign')
              ).click (event) ->
                event.preventDefault()
                feedurl = $(this).parent().find('.feedname').attr('href')
                feedname = $(this).parent().find('.feedname').html()
                $('#unsubscribe-alert').
                  find('.feed-name').html(feedname)
                $('#unsubscribe-alert').
                  find('.feed-url').html(feedurl).attr('href', feedurl)
                $('#unsubscribe-alert').modal('toggle')
                false
          )
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
              if $('.current').children('.article-content').is(':visible')
                window.open $('.current').children('.article-header').children('.article-link').attr('href'), '_blank'
              else
                $('.current').children('.article-content').show()
              false
        )
      ).
    append(
      $('<div/>').
        addClass('article-content').
        html($.parseHTML(data['Content']))
    )

offlineSetup = ->
  localArticles = localStorage.getObj 'articles'
  if localArticles? and localArticles.length > 0
    articles = []
    for article in localArticles
      articles.push addArticle(article)
    makeArticle articles
  else
    localStorage.setObj 'articles', []

addArticles = (object) ->
  list = []
  return list if not object?
  for article in object['Articles']
    continue if not article['ID']?
    element = addArticle article
    list.push element if element?
  list

makeArticle = (articles) ->
  for article in articles
    article.hide().appendTo '#articles'

makeCurrent = (articles, current) ->
  makeArticle articles
  for article in articles.slice(0, LIST)
    articleElement = document.getElementById(article.attr('id'))
    if articleElement?
      articleElement.classList.add('current')
  removeCurrent current
  if $('.current').index() is 0
    $('#prev').hide()
  else
    $('#prev').css 'display', 'block'
  _gaq.push(['_trackEvent', 'Articles', 'Next',  $('.current').attr('id')])
  $('body').scrollTo $('.current').offset().top if $('.current').exists()

nextArticle = (count, timeout, errornum, fun, current) ->
  if errornum > MAXERROR
    return
  if not current?
    current = $('.current')
  $.getJSON('/article?output=json&count=' + count, (data) ->
    _gaq.push(['_trackEvent', 'Articles', 'Get',  count])
    if data['URL'] is '/feed'
      if errornum < MAXERROR
        nextArticle count, timeout, errornum + 1, fun, current
        return
      else
        fun [addArticle {'Title': 'No more articles', 'URL': '/feed'}], current
      if fun is makeCurrent
        $('#next').hide()
    else if data['URL']?
      window.location = data['URL']
      timeout *= 2
    else
      articles = addArticles data
      ids = []
      newarticles = []
      for article in articles
        if not $(document.getElementById(article.attr('id'))).exists() and $.inArray(article.attr('id'), ids) == -1
          ids.push article.attr('id')
          newarticles.push article
      if newarticles.length is 0
        timeout *= 2
        nextArticle count, timeout, errornum + 1, fun, current
      else
        localArticles = localStorage.getObj('articles')
        localArticles = localArticles.concat(data.Articles)
        localStorage.setObj 'articles', localArticles
        fun newarticles, current
  ).fail((data) ->
    if online
      window.location.href = '/login'
  )

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
      else if index is 0
        index = -1
    if index + LIST < $('#articles').children().length
      $('#articles').children().slice(index + 1, index + LIST + 1).addClass('current').show()
      removeCurrent current
      _gaq.push(['_trackEvent', 'Articles', 'Next',  $('.current').attr('id')])
      index = $('.current').index()
      if index is 0
        $('#prev').hide()
      else
        $('#prev').css('display', 'block')
      $('#next').show()
      $('body').scrollTo $('.current').offset().top if $('.current').exists()
    else
      if online
        console.log 'getting more articles'
        nextArticle COUNT, TIMEOUT, 0, makeCurrent, current, index
        $('#next').show()
      else
        $('#next').hide()
    setTimeout nextArticle, TIMEOUT, COUNT, TIMEOUT, 0, makeArticle if index >= $('#articles').children().last().index() - 2

prev = ->
  if $('#prev').is ':visible'
    $('#prev').hide()
    index = $('.current').index()
    if index >= LIST
      markAsRead $('.current')
      $('.current').removeClass 'current'
      $('#articles').children().slice(index - LIST, index).addClass('current').show()
    _gaq.push(['_trackEvent', 'Articles', 'Previous',  $('.current').attr('id')])
    if index - LIST is 0
      $('#prev').hide()
    else
      $('#prev').css 'display', 'block'
    $('#next').show()
    $('body').scrollTo $('.current').offset().top if $('.current').exists()

markAsRead = (elements) ->
  localArticles = localStorage.getObj 'articles'
  elements.each ->
    $(this).
      hide().
      removeClass('unread').
      addClass('read')
    for article, i in localArticles
      if article['ID'] is $(this).attr 'id'
        localArticles.remove i
        break
  localStorage.setObj 'articles', localArticles

articles = () ->
  $('<section/>').
    attr('id', 'articles').
    insertBefore('#next') if not $('#articles').exists()

  offlineSetup()

  $('#next').show()
  if not $('.current').exists()
    next()
  else
    $('body').scrollTo $('.current').offset().top if $('.current').exists()

addfeed = ->
  $('#feedmodal').modal('toggle')

$ ->
  window.addEventListener 'offline', -> online = false
  window.addEventListener 'online', -> online = true
  online = navigator.onLine
  window.applicationCache.addEventListener 'error', -> online = false
  applicationCache.addEventListener 'updateready', -> window.location.reload()

  articles() if $('#next').exists()

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
          window.open $('.current').children('.article-header').children('.article-link').attr('href'), '_blank'
        else
          if $('.current').children('.article-content').is(':visible')
            window.open $('.current').children('.article-header').children('.article-link').attr('href'), '_blank'
          else
            show $('.current')
      when KeyEvent.DOM_VK_LEFT, KeyEvent.DOM_VK_NUMPAD6
        event.preventDefault()
        hide $('.current') if $('.current').children('.article-content').is ':visible'
      else
        return true
    return false

  $('#showbar').click (event) ->
    event.preventDefault()
    showbar()
    false

  $('#add').click (event) ->
    event.preventDefault()
    addfeed()
    false

  $('#prev').click (event) ->
    event.preventDefault()
    prev()
    false

  $('#next').click (event) ->
    event.preventDefault()
    next()
    false

  $('#addfeed').submit (event) ->
    event.preventDefault()
    url = $('#addfeed').find('input:first').val()
    subscribe(url)
    $('#feedmodal').hide()
    false

  $('#unsubscribe-ok').click (event) ->
    feedurl = $('#unsubscribe-alert').find('.feed-url').html()
    $('.current').removeClass('current unread').addClass('read').fadeOut('slow', ->
      next()
    )
    $('.article').each(->
      if feedurl is $(this).find('.feedname').attr('href')
        $(this).remove()
    )
    localArticles = localStorage.getObj('articles')
    newArticles = []
    for article in localArticles
      if article['FeedURL'] is not feedurl
        newArticles.push(article)
    localStorage.setObj 'articles', newArticles
    unsubscribe(feedurl)


srvCall = (data) ->
  try
    json = JSON.parse(data)
  catch err
    location.href = 'http://myrivulet.appspot.com/'
  location.href = json['URL']
jsonp = (src) ->
  s = document.createElement('script')
  old = document.getElementById('srvCall')
  old && document.body.removeChild(old)
  s.charset = 'UTF-8'
  s.id = 'srvCall'
  document.body.insertBefore(s, document.body.firstChild)
  s.src = src + '?output=json&count=1&callback=srvCall&' + new Date().getTime()

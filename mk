#!/bin/sh

echo 'making'

coffee --compile --output static coffee/*
#closure static/script.js > static/script.js

jade -o static jade
cp static/{articles,feeds,user}.html templates

lessc less/style.less static/style.css

yuicompressor static/script.js -o static/script.js
yuicompressor static/style.css -o static/style.css

cp icon/* static

#./bump


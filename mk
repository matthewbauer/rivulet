#!/bin/sh

echo 'making'

coffee --compile --output static coffee/*
#closure static/script.js > static/script.js

jade -o static jade
cp static/{articles,feeds,user}.html templates

lessc less/style.less static/style.css

yuicompressor static/script.js -o static/script.js
yuicompressor static/style.css -o static/style.css

#sizes='16 24 32 57 64 72 96 114 128 144 195 256'
#for size in $sizes
#do
#  inkscape -C -e static/icon$size.png -w $size -h $size static/icon.svg
#done

cp icon/* static

#./bump


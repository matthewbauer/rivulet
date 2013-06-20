#!/bin/sh

coffee --compile --lint --output static script.coffee
#closure static/script.js > static/script.js

lessc style.css static/style.css

yuicompressor static/script.js -o static/script.js
yuicompressor static/style.css -o static/style.css

cp logo* icon* static

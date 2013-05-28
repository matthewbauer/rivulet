#!/bin/sh

coffee --compile --lint --output static script.coffee
#closure static/script.js > static/script.js

lessc style.css static/style.css


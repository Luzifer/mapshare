#!/bin/bash
set -euxo pipefail

css_deps=(
	# Must-have order
	npm/bootstrap@4/dist/css/bootstrap.min.css
	npm/bootswatch@4/dist/darkly/bootstrap.min.css
	npm/bootstrap-vue@2/dist/bootstrap-vue.min.css

	# Other packages
	npm/leaflet@1.6.0/dist/leaflet.min.css
)

js_deps=(
	# Must-have order
	npm/vue@2/dist/vue.min.js
	npm/bootstrap-vue@2/dist/bootstrap-vue.min.js

	# Other packages
	npm/axios@0.19.0/dist/axios.min.js
	npm/leaflet@1.6.0/dist/leaflet.min.js
	npm/moment@2.24.0/min/moment.min.js
)

leaflet_images=(
	marker-icon.png
	marker-icon-2x.png
	marker-shadow.png
)

rm -rf frontend/leaflet
mkdir -p frontend/leaflet
for img in "${leaflet_images[@]}"; do
	curl -sSfLo "frontend/leaflet/${img}" "https://cdn.jsdelivr.net/npm/leaflet@1.5.1/dist/images/${img}"
done

IFS=','

curl -sSfLo frontend/combine.js "https://cdn.jsdelivr.net/combine/${js_deps[*]}"
curl -sSfLo frontend/combine.css "https://cdn.jsdelivr.net/combine/${css_deps[*]}"

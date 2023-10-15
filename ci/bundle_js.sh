#!/bin/bash
set -euxo pipefail

leaflet_ver="1.7.1"

css_deps=(
  # Must-have order
  npm/bootstrap@4/dist/css/bootstrap.min.css
  npm/bootswatch@4/dist/darkly/bootstrap.min.css
  npm/bootstrap-vue@2/dist/bootstrap-vue.min.css

  # Other packages
  npm/leaflet@${leaflet_ver}/dist/leaflet.min.css
)

js_deps=(
  # Must-have order
  npm/vue@2/dist/vue.min.js
  npm/bootstrap-vue@2/dist/bootstrap-vue.min.js
  npm/vue-i18n@8.21.0/dist/vue-i18n.min.js

  # Other packages
  npm/leaflet@${leaflet_ver}/dist/leaflet.min.js
  npm/moment@2.27.0/min/moment.min.js
)

leaflet_images=(
  layers.png
  layers-2x.png
  marker-icon.png
  marker-icon-2x.png
  marker-shadow.png
)

rm -rf frontend/images
mkdir -p frontend/images
for img in "${leaflet_images[@]}"; do
  curl -sSfLo "frontend/images/${img}" "https://cdn.jsdelivr.net/npm/leaflet@${leaflet_ver}/dist/images/${img}"
done

IFS=','

curl -sSfLo frontend/combine.js "https://cdn.jsdelivr.net/combine/${js_deps[*]}"
curl -sSfLo frontend/combine.css "https://cdn.jsdelivr.net/combine/${css_deps[*]}"

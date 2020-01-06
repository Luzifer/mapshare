[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/mapshare)](https://goreportcard.com/report/github.com/Luzifer/mapshare)
![](https://badges.fyi/github/license/Luzifer/mapshare)
![](https://badges.fyi/github/downloads/Luzifer/mapshare)
![](https://badges.fyi/github/latest-release/Luzifer/mapshare)
![](https://knut.in/project-status/mapshare)

# Luzifer / mapshare

This project is a very simple and data protecting alternative to sharing a location through Glympse or similar services.

You can setup your own instance in minutes, it does not require any database (even retained location data is dropped on restart of the service!) and you can share your location from a mobile browser. To view the location nothing more than a browser is required.

When sharing a location you have the choice to select whether the server should retain the location data (until restart) or just pipe it through. Retaining the data has the advantage new viewers (or viewers whose websocket has reconnected) instantly see your location. When not retaining data the data is received, sent to all connected sockets and afterwards instantly forgotten.

# htmltoebook
Converts html webpages to a readable ebook. Fetches the list of html web pages, strips out readable content of the pages and creates an ebook that can be used in phones or ebook-readers.

* Fetches and cleans up the webpages as readable paragraphs using [readability](github.com/go-shiori/go-readability) package.
* Currently supported output format is mobi
* Two user interfaces supported, cli and webui. Launches webui by default.
* By default ebook is generated in folder $HOME/Downloads/htmltoebook

![Screenshot](screenshot.png)

## TODO
* Support epub format, zip file
* Add system tray icon
* Add webview to link with platform webkit, without need to open link in external browser
* UI to generate range of URLs that could expand in a given numeric range
* Generate release binaries using Actions.
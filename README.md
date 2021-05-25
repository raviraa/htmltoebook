# htmltoebook
Converts html webpages to a readable ebook. Fetches the list of html web pages, strips out readable content of the pages and creates an ebook that can be used in phones or ebook-readers.

* Fetches and cleans up the webpages as readable paragraphs using [readability](github.com/go-shiori/go-readability) package.
* Currently supported output format is epub
* Two user interfaces supported, cli and webui. Launches webui by default.
* By default ebook is generated in folder $HOME/Downloads/htmltoebook
* Settings are stored in configuration file at $HOME/.htmltoebook.toml
* Range of links can be specified. Ex. `http://example/page{{1-3}}/content`
* Built releases [available](https://github.com/raviraa/htmltoebook/releases) for linux, windows and mac darwin. Only linux version is tested.

![Screenshot](screenshot.png)

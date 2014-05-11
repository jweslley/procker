# Procker

Procker is a tool for managing Procfile-based applications.

Also, Procker can be used as a library to build other applications which require process's management.


## Installation

    go get github.com/jweslley/procker
    make


## Getting started

1. Write a `Procfile`
2. *Optional:* write a `.env` file
3. Run the app using `procker`:

    procker start

For more information, use `procker help`.


## Resources

* [Procfile](https://devcenter.heroku.com/articles/procfile)
* [Store config in the environment](http://www.12factor.net/config)


## Alternatives

* [foreman](https://github.com/ddollar/foreman)
* [forego](https://github.com/ddollar/forego)
* [shoreman](https://github.com/hecticjeff/shoreman)
* [honcho](https://github.com/nickstenning/honcho)
* [norman](https://github.com/josh/norman)


## Bugs and Feedback

If you discover any bugs or have some idea, feel free to create an issue on GitHub:

http://github.com/jweslley/procker/issues


## License

MIT License. Copyright (c) 2014 [Jonhnny Weslley](<http://www.jonhnnyweslley.net>)

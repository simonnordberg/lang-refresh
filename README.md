# Language refresh

This repository contains the same simple application implemented in different languages. The goal is simple: refresh
and in some cases learn interesting, useful and sometimes pointless programming languages.

# The application

The application is simple: a command line utility to import/create [Google Photos](https://photos.google.com/) albums
exported using [Google Takeout](https://takeout.google.com/) to [PhotoPrism](https://www.photoprism.app/). The album
import is done separately once the images have already been
[imported and indexed](https://docs.photoprism.app/developer-guide/media/import/).

The application was chosen because 1) I needed to get the job done, and 2) it is reasonably trivial while still
requiring a reasonably [diverse set of features](#features-and-considerations).

# Application flow

The high-level flow of the application are as follows:

* Manage application configuration, e.g. using the command line, environment variables or configuration file,
  e.g. API endpoint, username, password, local import/photo directory
* API authentication and authorization
* List album directories, image and metadata files from local drive
* Parse metadata.json containing album identifiers, slugs etc.
* Use the API to create the album if it doesn't already exist
* Enumerate over all images and generate a SHA1 hash used to query the API
* Add the already imported image to the album

# Features and considerations

In addition to the application flow per-se, I want to make sure to check the following aspects for each language:

* As close to a language idiomatic implementation as possible
* Managing and parsing command line options and environment variable
* Calling a JSON API
* Logging
* Building and parsing JSON files and content
* File system navigation and management
* Build scripts and containerization (e.g. Docker)

# Languages

This list will be modified as get inspiration for what I want to pursue next.

| Language |   Status    |
|:--------:|:-----------:|
|  Python  |     v1      |
|   Rust   |     v0.5    |
|   Java   | Not started |
|   Ruby   |     v1      |
|    Go    | Not started |
|  Elixir  | Not started |
| Haskell  | Not started |
|  Kotlin  | Not started |

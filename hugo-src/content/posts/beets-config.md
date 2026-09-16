+++
date = '2024-03-27T10:00:00-06:00'
draft = true
title = 'Managing music locally with beets'
+++

> "Out, out, damned Spotify!"
>
> — <cite>[Richard Stallman](https://www.stallman.org/spotify.html), *What's bad about: Spotify*</cite>

Although the convenience of streaming cannot be understated, there's something special (if not a bit hipster-ish) about curating a personal music library of files. I started seriously listening to music using iTunes, and it (along with a trusty iPod Touch) shaped how I interacted and related to the act of listening to music as the medium of local files is the message. Hence, the on-demand infinite possibilities of Spotify not only put a bajillion of songs at my fingertips,[^missing-music] but my switch to Linux obsoleted my previous preferred software. Enter an open-source and appropriately decoupled duo: cmus, the C music player and **[beets](https://beets.io/)**, a music metadata manager and our subject for today.

[^missing-music]: Notable omissions in Spotify's catalog include Jaco Pastorious's solo work and a high-school obsession of mine titled "Synthesocietal."

Beets has a lot of features and typical of CLI clients, many of them are hidden behind `--help` flags, man pages, and domain-specific knowledge. We'll get to those, don't fear, but to make the journey as smooth as possible, I first have to point out the water that we're swimming in. iTunes hid behind its glossy GUI a tight coupling between metadata[^data] and the music[^metadata] whereby the pertinent information like the Artist, title, album were written to and stored on each music file such that the presence of metadata implied the presence of the file and conversely. However, beets stores the metadata for the music independently of the files in a decoupled fashion.

[^data]: data about data

[^metadata]: that is data about data, the latter of which comprises my sweet collection of mp3's, flac files, etc.

Here is a breakdown of my current `config.yaml` and the philosophy behind it.

## The global foundation
My music lives in `~/music`, and I keep the database in the standard `.config` location. Multi-threading is enabled (`threaded: true`) to speed up imports (mainly loudness normalization), and I’ve opted for a colorful UI because my terminal supports it.

```yaml
directory: ~/music
library: ~/.config/beets/library.db
threaded: true
ui:
  color: true
```

## Import
When I import new music, I prefer to start from scratch. My configuration removes all existing metadata (`from_scratch: true`) and re-identifies everything using **MusicBrainz**. This ensures total consistency across my collection.

I also prioritize English transliterations for metadata if they exist, which keeps the library navigable and helps as much textual content renders as possible.

```yaml
import:
  autotag: true
  move: false # I prefer to manually manage file moves if needed
  write: true
  from_scratch: true
  languages: en
```

## Path Logic and Organization
Beets allows for incredibly powerful path formatting. I organize my library by `$albumartist/$album/`, with logic to handle multi-disc releases automatically.

```yaml
paths:
  default: $albumartist/$album/%if{$multidisc,Disc $disc/}$track $title
```

## Plugin Highlights
The real power of `beets` lies in its plugins. Here are the ones I find indispensable:

### 1. The "Zero" Plugin
I use `zero` to strip out metadata fields I don't want, such as embedded lyrics, comments, or generic genres. This keeps the file tags lean and focused only on the information I care about.

```yaml
zero:
  auto: true
  fields: images lyrics comments genre
```

### 2. Acoustic Fingerprinting with "Chroma"
The `chroma` plugin uses the **AcoustID** project to identify songs based on their actual audio waveform, which is almost always required when importing pirated music from the Russians.

```yaml
chroma:
  auto: true
```

### 3. ReplayGain for Consistent Volume
I use the `replaygain` plugin with an `ffmpeg` backend to normalize the perceived loudness of my tracks. I’ve disabled `parallel_on_import` to ensure that metadata is written correctly to the files themselves, rather than just the database.

```yaml
replaygain:
  # cat /proc/cpuinfo for your machine's max
  threads: 8
  backend: ffmpeg
  parallel_on_import: false
```

### 4. Integration with ListenBrainz
Finally, I sync my listening history to **[ListenBrainz](https://listenbrainz.org/)**, an open-source alternative to Last.fm that integrates perfectly with the MusicBrainz ecosystem.

<!-- TODO(author): this is the only plugin without a config snippet. Expand: what gets synced (plays, loves, "now playing"?) and add the matching listenbrainz block (username/token) like the other plugins. -->

## Full Configuration File

The snippets above are only the highlights. You can download the complete configuration here:

<!-- TODO(author): decide whether to call out the other plugins the full config enables (scrub, info, unimported, mbsync, autobpm, edit) or leave the download as the catch-all. -->

- [config.yaml](/downloads/config.yaml)

## Conclusion

Setting up `beets` is an investment in your music library. It takes time to dial in the configuration, but once you do, the result is a perfectly tagged, organized, and searchable collection.

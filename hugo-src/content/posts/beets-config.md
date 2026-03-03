+++
date = '2024-03-27T10:00:00-06:00'
draft = true
title = 'Managing local music libraries with beets'
+++

# Managing local music libraries with beets

In an era of streaming, maintaining a high-quality local music library can feel like a lost art. For those who still prefer the control of local files, **[beets](https://beets.io/)** is the definitive command-line tool for orchestrating a perfect library. It’s more than just a tagger—it’s a media management ecosystem.

Here is a breakdown of my current `config.yaml` and the philosophy behind it.

## The Global Foundation
My library lives in `~/music`, and I keep the database in the standard `.config` location. Multi-threading is enabled (`threaded: true`) to speed up the interface, and I’ve opted for a colorful UI to keep the CLI experience readable.

```yaml
directory: ~/music
library: ~/.config/beets/library.db
threaded: true
ui:
  color: true
```

## Import Philosophy: A Clean Slate
When I import new music, I prefer to start from scratch. My configuration removes all existing metadata (`from_scratch: true`) and re-identifies everything using **MusicBrainz**. This ensures total consistency across my collection.

I also prioritize English transliterations for metadata if they exist, which keeps the library navigable regardless of the artist's origin.

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
Sometimes metadata is so broken that name-based matching fails. The `chroma` plugin uses the **AcoustID** project to identify songs based on their actual audio waveform.

```yaml
chroma:
  auto: true
```

### 3. ReplayGain for Consistent Volume
I use the `replaygain` plugin with an `ffmpeg` backend to normalize the perceived loudness of my tracks. I’ve disabled `parallel_on_import` to ensure that metadata is written correctly to the files themselves, rather than just the database.

```yaml
replaygain:
  threads: 4
  backend: ffmpeg
  parallel_on_import: false
```

### 4. Integration with ListenBrainz
Finally, I sync my listening history to **[ListenBrainz](https://listenbrainz.org/)**, an open-source alternative to Last.fm that integrates perfectly with the MusicBrainz ecosystem.

## Conclusion
## Full Configuration File
Setting up `beets` is an investment in your music library. It takes time to dial in the configuration, but once you do, the result is a perfectly tagged, organized, and searchable collection that puts streaming services to shame.

You can download my full configuration file here:
- [config.yaml](/downloads/config.yaml)

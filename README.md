# Subsyncer

Subsyncer is a tool for automatic synchronization of subtitle files.  

It works by comparing an input subtitle file to a "reference" subtitle file,
assumed to be synchronized properly. The subtitle files are analyzed, using machine
translation and indexing techniques, and the input subtitle file is then re-synchronized
so that it matches the reference subtitle.
 
## How to use

```sh
subsyncer --input-file=$HOME/MyMovie/MyMovie.srt \
          --input-lang=heb \
          --ref-file=$HOME/MyMovie/MyMovie.eng.srt \
          --ref-lang=eng
```

Note: it doesn't work as of yet, this is still work-in-progress :)
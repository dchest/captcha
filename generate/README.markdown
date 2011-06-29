How to create sounds for other languages
========================================

* Record sounds for 0-9 and put them into a directory with language names.
  Speak fast enough to make sound files small.  Make sure the level of sound is
  the same as in the provided samples for English (this is important for making
  captchas harder to break). Same files in 8 KHz 8-bit PCM WAV format.  (To do
  this in Audacity, set "Project Rate (Hz)" at the lower left corner to 8000,
  then File > Export, select "Other uncompressed files", click Options...,
  select "WAV (Microsoft)" for Header, and "Unsigned 8 bit PCM" for Encoding.)

  If you're not sure if your sounds are okay or how to save them properly, just
  save one of them into any format (MP3 is okay), and send it to me
  <dmitry@codingrobots.com>. I'll check it, and if it's okay, I'll ask you for
  other sounds, and process them myself (in this case, you can stop reading.)

* Put 0.wav - 9.wav into the subdirectory with language name (e.g. "ua").

* Open main.go and edit "var langs" on line 21 to include the new directory
  name.

* make && ./generate

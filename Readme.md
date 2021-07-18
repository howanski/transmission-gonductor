# **Transmission Gonductor**


# Introduction

This small server-side tool was created across two weeks worth of afternoons to fullfill simple needs:
- I'd like it to be as light as possible to run on passively cooled machine, either x86 or ARM
- It needs to conduct Transmission's torrent downloads so the files will be downloaded in alphabetical order, i.e. "Ubuntu introduction S01E01", then "Ubuntu introduction S01E02" and so on.

I've decided to use Go language as I had no previous experience in using it and I thought it would be fun to try something new.

Honestly strong typing of variables in Go might be a bit disadvantege when communicating via JSON as marshalling and unmarshalling JSON string the proper way is boring and time-consuming, so I've done some lazy hacks to get rid of declaring bunch of new Types just to throw string across TCP/IP.

# Future
Possibly I'll redeclare Transmission JSON interfaces the proper way if I'll have time to merge functions from my [earlier project](https://github.com/howanski/network-limit-watcher-extension), as I'd like it to run on server also, preferably on container.

# Running
> Linux

By default just run `server_run.sh` script, make sure you have go installed in your system.

> Windows

I think i'll make some compiler script later, it shouldn't be hard.

> Docker

Docker instructions will be added soon, I need some time to sleep too :)

# Usage
User Interface is hosted on [Port 8080](http://localhost:8080/) by default.

User Interface is simple bootstrap form split on 3 parts on accordion - simply fill values, tick functions you want to on or off, and click "Save" - you will be prompted if something's not right.

Interface will update every 10-15 seconds to inform you about server status.
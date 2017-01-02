# rmdups -- remove duplicate copies of data files and images

`rmdups` creates a index of files on your disk, and is then able to show
duplicated data files and images. Once the duplicates are found, the
corresponding messages are printed out, so that you can investigate.

# How to use rmdups

# Background

For the last 10 years or so I've used around 3 computers, several
hard-drives, USB drives and other storage media. During this time I've
copied files from one disk to another, many times, often unnecessarily.
When it came to making a backup, I wasn't sure which files I should backup,
so I often ended up copying the same files to differnet locations many
times.

# Notes

**XXX**: some more work is needed at dealing with duplicates. Since removing
is quite likely dangerous, it'd be good to be able to preview files before
removing them. Right now you specify `-remove` which in fact doesn't remove
anything, but prints `rm <file>` message instead. You can capture the output
via `script` and then `grep rm` on resulting file, and pipe it to `/bin/sh`.
Target solution should probably be to preview the files, or make a shell
script with `rm` commands.

# TODO

Change to `xxhash`

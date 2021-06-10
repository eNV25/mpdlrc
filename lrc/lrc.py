"""
This module parses LRC files. Lines are in the format `[mm:ss.xx]text`, where  mm  is minutes,
ss  is seconds and  xx  is  one-hundredths of a second.  text  is the rest of the line, and
holds the actual lyric.

A `Line` class with attributes `time` and `text` is used to enapsulate an LRC line.
"""

import re as _re


PATTERN = r"\[(?P<mm>\d\d)[:](?P<ss>\d\d)[.](?P<xx>\d\d)\](?:\[\d\d:\d\d.\d\d\])*(?P<text>.*)"
REGEX = _re.compile(PATTERN)


class Line:
    def __init__(self, mm, ss, xx, text):
        self.time = round(int(mm) * 60 + int(ss) + int(xx) / 100, 2)
        self.text = text

    def __str__(self):
        return f"Line(time={self.time:.2f}, text={self.text!r})"


def parser(string):
    for match in REGEX.finditer(string):
        if match:
            mm, ss, xx, text = match.groups()
            yield Line(mm, ss, xx, text)

# test
if __name__ == "__main__":
    with open("./path/to/file.lrc") as file:
        lyrics = file.read()
    for line in parser(lyrics):
        print(line)

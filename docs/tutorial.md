# ⊧ dupi -- tutorial

## Install

```
go install github.com/go-air/dupi@latest
```

(If you would like binary distributions, please let us know on the 
[issue tracker](https://github.com/go-air/dupi/issues))

## What to expect and what not to expect.

Dupi is designed to identify all "duplicates" in a set of documents.  What a
"duplicate" means in this context can vary a bit, so let's first rule out some
possible mis-use cases:

Dupi is not designed to find all similar words in a set of words, nor is it
designed to show "similar" documents to a given document, particularly in a way
which takes into account semantic content, such as considering themes or 
synonyms.

A "duplicate" for dupi means, more or less, the same text being present in more
than one document.  Dupi can be used for helping to identify plagiarism or
copyright violations accross large sets of documents, for example.

Text processing is however inherently noisy: the same text may be formatted
differently, use different line-wraping techniques, or different uses of
capitalisation, or OCR errors, etc.  For this reason, dupi is built in a way
which allows extending the idea of what is a "duplicate" in various ways and to
various kinds of documents.  

Out of the box, dupi is configured simply to find common subsequences of text
accross sets of documents.  This tutorial  addresses this use case.  More
sophisticated usage and development is needed for other use cases.

Let's get started.

## Overview

In the following, we will

1. Look at the command line usage.
1. Create a dupi index from a set of documents.
1. Extract duplicate blots from the set of documents.
1. Examine the blots.
1. Append more documents to the index.
1. Query the set of documents for things similar to a given document.

## CLI

Dupi provides a command line interface which is common these days and
takes the form `dupi <verb> [args]` where `verb` tells dupi what to do
and `[args]` provides information about the object of the verb or modifiers.

```
dupi -h
usage: dupi [global opts] <verb> <args>
verbs are:
        index                               paths
        extract       extract from the index root
        blot                         blot [files]
        unblot                      unblot <blot>
        inspect           inspect the root index.

global options:
        -r                default=""               index root

To get help on a verb, try dupi <verb> -h.
```

## Creating an index

To create an index just run 'dupi index' and provide it with a list of 
files or directories.  Dupi will traverse all subdirectories and add
each file.  The files should be text files.

Example:
```
dupi index .
```

## Extracting Duplicates

```
dupi extract
```

## Appending to the index

```
dupi index -a /path/to/new/docs
```

## Blotting

Sometimes it might be interesting to see if a file has a blot.  Dupi
provides the ability to blot files using the same mechanism as is
used in the index.

```
dupi blot file
```


## Querying the index

Dupi provides primitives for unblotting, which takes a blot and
reconstructs the corresponding text and instances.  This is still
rudimentary, but here are some examples.

```
dupi extract | awk '{print $1}' | xargs dupi unblot
```

Or 

```
dupi blot file | xargs dupi unblot
```

## Conclusion

We have shown some basic usage of dupi.  As dupi is in early stages 
of development, some things may be added or changed, we will try to
keep this document up to date.   Feel free to help improve our
documentation using [issues](https://github.com/go-air/dupi/issues) or 
[pull requests](https://github.com/go-air/dupi/pulls).





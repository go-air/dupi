# ⊧ dupi -- tutorial

## Install

```
go install github.com/go-air/dupi@latest
```

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

1. Create a dupi index from a set of documents.
1. Extract duplicates from the set of documents.
1. Append more documents to the index.
1. Query the set of documents for things similar to a given document.

## Creating an index

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

## Querying the index

```
dupi like /path/to/unknown/docs
```

## Conclusion





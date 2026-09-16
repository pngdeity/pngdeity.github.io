+++
date = '2025-01-29T10:00:00-06:00'
draft = true
title = 'Stealing GitHub Commits'
+++

> "You Wouldn't Steal a Car"
>
> <cite>MPAA Anti-Piracy PSA, 2004</cite>

Lmao the mad lad did it https://github.com/anthropics/claude-code/issues/65710

Here's a summary of one version of the article:

So, GitHub is one of if not the most popular code forge website. Features on GitHub (like their contribution heat map) have become widely recognized and well-known in tech and tech-adjacent circles. However, there are so non-obvious caveats about how the contribution graph works. One, those contributions on the graph count for code that is slotted into main/master branches only. Unmerged changes don't count. Two, up to this point I've only mentioned that a user's contributions get reflected in the contribution graph. For contributions that are independent of git (e.g. Pull requests), GitHub has a 100% accuracy in recording a user's work. GitHub does not necessarily have a 100% accuracy in recording contribution types from git, namely commits. This is important because of how GitHub matches commits with users. \<more technical explanation\> As a result, GitHub-native technologies (like the contribution graph can work as a signal for developer productivity, but the selective counts of contributions makes it kind of noisy for new hires.

Why does GitHub display the AuthorDate vs the CommitDate (or vice versa) in the web UI on a repository's branch's commit log page [e.g.](https://github.com/pngdeity/eoh-2026-judging/commits/main/)? I think these are git object headers, which also have something to do with september 17, 2001

An alternative framing \-\>

You used this during your GLUG lightning talk.

There's a maximum length for the technical explanation section, so decide the essential background information to provide enough context for the middle and end.

Roughly, here's how git internals are reflected on GitHub's UI, and here are the problems with those UI choices, and here's some theorizing about GitHub's rationale/motivations.

1. Git and GitHub operation
   1. GitHub browsing
      1. GUI perspective
   2. Here's how GitHub works with the demonstration of the .patch url

2. Problems from past choices
   1. The abstract case
      1. Look for the "committer" field of a commit object.
   2. Concrete issues
      1. Stealing commits for fun and profit(?)
      2.

Sources:

https://git-scm.com/docs/git-config\#Documentation/git-config.txt-username

The problem of non functioning emails, whether now-dead (big example is .edu emails) or never functioning (which would be even easier to steal since no email will be received). This is made even easier by GitHub's integrated search.

Contribution graph counts

Theorizing about GitHub's Rationale

For example, like that they don't verify that your account is old enough to be the source of some commits. Or, that there doesn't seem to be any specific way to get your email back from another malicious user.

Unhoused information (relevant and should probably be integrated):

What percentage of developers use GPG signing? When did GitHub start showing verified badges and offer "vigilant mode"?

How could GitHub prevent this problem?

The GitHub URL scheme is [https://github.com/pngdeity/pngdeity.github.io/commit/$SHA256](https://github.com/pngdeity/pngdeity.github.io/commit/$SHA256) where $SHA256 is a 41 character SHA 256 hash

Asides (not relevant but could be integrated):

The "Mon Sep 17 00:00:00 2001" string in commit objects.

The "H" in GitHub is always capitalized.

Personal project idea: GitHub with only git stuff.

You can also find a users PGP key by appending ".gpg" to their page.

+++
date = '2026-09-09T10:00:00-05:00'
draft = true
title = 'Licensing Software in the Open'
+++

> "Possession is nine-tenths of the law."
>
> — a pearl of legal wisdom

<!-- TODO(author): 1-2 sentence intro. What prompted the talk, and what should the reader take away? -->

## What "open source" means

[The Open Source Definition](https://opensource.org/osd), maintained by the Open Source Initiative, is the generally accepted (but not legally binding) definition of "open source." In plain language, an open source license requires that:

- The cost of the license be zero.
- The source be available.
- Users be treated equally, without discrimination — including based on the intended use of the source.
- End users have the right to modify the source.

<!-- TODO(author): your slides footnote these four criteria with caveats (Distribution*, Source*, modifications*). If those distinctions matter, spell them out here. -->

## Comparing popular licenses

The Open Source Initiative [lists the licenses](https://opensource.org/licenses?categories=popular-strong-community) that are popular or have a strong community behind them. Here is a comparison of a common selection:

<!-- TODO(author): this table is reconstructed from your slides — verify the wording and add anything you trimmed. -->

| License | Copyleft (reciprocity) | Network distribution | Patent grants | Attribution & endorsement |
| :-- | :-- | :-- | :-- | :-- |
| **MIT** | Permissive: code can be closed and monetized with zero restrictions. | Standard: SaaS usage does not trigger source sharing. | Implicit: relies purely on copyright law. | Standard: keep original copyright notices intact. |
| **Apache v2** | Permissive: can be absorbed into closed-source products. | Standard: SaaS usage does not trigger source sharing. | Explicit & retaliatory: instantly revokes the license if the user sues for patent infringement. | Standard: must also retain external "NOTICE" files, if provided. |
| **Mozilla Public (MPL)** | Weak (file-level): modified MPL files must remain open, but can be mixed with proprietary files. | Standard: SaaS usage does not trigger source sharing. | Explicit: grants patent rights to users for the contributed code. | Standard: keep original copyright notices intact. |
| **LGPLv3** | Weak (library-level): the library itself must remain open, but proprietary code can dynamically link against it. | Standard: SaaS usage does not trigger source sharing. | Explicit & retaliatory: inherits robust patent protections from the GPL framework. | Standard: keep original copyright notices intact. |
| **GPLv3** | Strong: the entire derivative software program must inherit the GPLv3 license. | Standard: SaaS usage does not trigger source sharing. | Explicit & retaliatory: protects users against patent trolls. | Standard: keep original copyright notices intact. |
| **AGPLv3** | Strong: the entire derivative software program must inherit the AGPL license. | Network copyleft: SaaS/cloud interaction triggers source sharing to remote users. | Explicit & retaliatory: protects users against patent trolls, using the GPLv3 framework. | Standard: keep original copyright notices intact. |

Of trivial note, the NCSA/University of Illinois license is not considered current.

## No license means "all rights reserved"

Per US copyright law, a work without a license — or with a malformed, self-contradictory licensing statement — remains the sole right of the creator or copyright holder. In other words: all rights reserved.

GitHub's Terms of Service are the narrow exception. [Section D.5](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service#5-license-grant-to-other-users) says that making a repository public grants other GitHub users limited rights to the code: viewing it on GitHub and forking it through the GitHub UI. GitHub's ToS then rightly encourages maintainers to [adopt a license](https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/adding-a-license-to-a-repository#including-an-open-source-license-in-your-repository).

## Case study: Tandemn Labs

<!-- TODO(author): decision point. This section names Tandemn Labs, its founding engineer Quan Hao Ng, and outside contributors, and argues that specific repos are not actually open source. Keep as-is, soften it, or anonymize? Note the published slides contain the same material. -->

[Tandemn Labs](https://github.com/tandemn-labs#open-source) is a startup with a GitHub organization, employees, and public repositories. Per the company's README, Tandemn's software is "fully open source."

Its projects tell a different story:

- **Tandemn Docs** — [repository](https://github.com/Tandemn-Labs/Tandemn-docs) is unlicensed. The documentation then leads to the following repository.
- **Tandemn System** — [repository](https://github.com/Tandemn-Labs/tandemn-system/) is licensed, but contains a git submodule pointing at `LLM_Placement_Solver`, which is [unlicensed](https://github.com/Tandemn-Labs/LLM_placement_solver). Since that is ambiguous, I would err on the side of calling it a malformed licensing statement, i.e. all rights reserved by Tandemn Labs, and therefore *not* open source.
- **tandemn-vllm** — [repository](https://github.com/Tandemn-Labs/tandemn-vllm) is unlicensed and has contributions from both Tandemn employees (for example, the GitHub user `orangeng`) and outside contributors. With no explicit contributor guidelines, the employee contributions remain under the company's license while each outside contribution remains under its author's exclusive license. The public code therefore contains sections licensed to different, mutually exclusive parties. In an extreme hypothetical, if Tandemn Labs deployed an outside contribution for monetary gain, the outside contributor could seek compensation, because they still hold the rights to their code.

<!-- TODO(author): correct any of the above and confirm the repos are still in this state (licenses may have been added since the talk). -->

## Final thought

License your code — hopefully with an open-source license.

<!-- TODO(author): expand this into the takeaway. -->

## Original Presentation

This post is based on a talk given to <!-- TODO(author): name the group and link it (gnulug/meetings on GitHub?) --> on September 9, 2026. You can view the original slides below:

<embed src="/downloads/open-source-licensing-slides.pdf" type="application/pdf" width="100%" height="600px" />

[Download the original PDF](/downloads/open-source-licensing-slides.pdf)

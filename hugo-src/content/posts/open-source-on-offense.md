+++
date = '2026-09-16T10:00:00-05:00'
draft = true
title = 'Open Source on Offense'
+++

> "Linux is a cancer."
>
> <cite>Steve Ballmer, Microsoft CEO, 2001</cite>

Why on Earth would Microsoft fund the creation and maintenance of an entirely free and open source operating system?

The answer is not that anyone at Microsoft read the GNU Manifesto and decided they liked it. It's that open source can be a weapon, and Microsoft has gotten very good at wielding it. The company that once called Linux a cancer[^ballmer] now treats it as raw material for the oldest trick in technology economics: [commoditize your complement](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/). Give away the layer you don't need, and keep the one you do.

[^ballmer]: The fuller quote, from a 2001 interview: "Linux is a cancer that attaches itself in an intellectual property sense to everything it touches." Notice what he objected to: the GPL's copyleft, not the quality of the kernel.

<!-- TODO(author): optional personal on-ramp. If you want to say how you first ran into this story, or why it gets under your skin, this is the place. Otherwise the argument stands on its own. -->

Once you see the pattern, a lot of corporate open source stops looking selfless, and a lot of otherwise baffling decisions start to make sense. Enter SONiC: a network operating system Microsoft built, gave away, and now leans on to keep the switching market a buyer's market.

## What SONiC is

**SONiC** — "Software for Open Networking in the Cloud" — is a [free and open source network operating system](https://en.wikipedia.org/wiki/SONiC_(operating_system)) built on Debian Linux, and it runs on the switches inside a data center. Microsoft developed it and [open-sourced it in 2016](https://azure.microsoft.com/en-us/blog/microsoft-showcases-software-for-open-networking-in-the-cloud-sonic/), building on an earlier release called the [Azure Cloud Switch](https://azure.microsoft.com/en-us/blog/microsoft-showcases-the-azure-cloud-switch-acs/), and on the Open Compute Project's [Switch Abstraction Interface](https://github.com/opencomputeproject/SAI) (SAI).

Two decisions do most of the work:

- **The software is decoupled from the hardware.** SONiC talks to the switch through SAI, so the same network OS runs on [more than 100 different switches](https://sonicfoundation.dev/) from a pile of different vendors and ASIC families.
- **The system is containerized.** Network functions run as microservices instead of one monolithic image, which lets operators swap and upgrade pieces without touching the rest.

In April 2022, Microsoft [handed governance to the Linux Foundation](https://www.linuxfoundation.org/press/press-release/software-for-open-networking-in-the-cloud-sonic-moves-to-the-linux-foundation). "We created it as open source so the entire networking ecosystem would grow stronger," said Dave Maltz, a corporate vice president for Azure networking. He added that SONiC already runs on millions of ports in the networks of cloud scalers, enterprises, and fintechs.

Here's the part I can't get over: the project's premier members today include Microsoft, Google, Alibaba, Broadcom, Dell, and NVIDIA, alongside **Cisco, Arista, and Nokia** — the incumbents whose proprietary switch software SONiC was built to displace.

## Commoditizing the complement

Time for the economics, because they are the entire story. In 2002, Joel Spolsky gave this pattern a name: [commoditize your complement](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/).

A *complement* is a product people buy alongside yours. Demand for your product rises when the price of its complements falls, so your strategic interest is to drive those prices down toward marginal cost — and there is no cheaper price than free.

Spolsky's examples are the familiar ones: IBM documenting the PC so the add-in market would commoditize, and Microsoft licensing MS-DOS to every clone-maker so the PC itself became a commodity. In both cases the point was never the layer being commoditized; it was the layer that stayed scarce. For IBM that layer was the PC, briefly. For Microsoft it was the operating system, for a long time.

Microsoft is running the same play with switches. Data-center networking had been a vertically integrated, high-margin business: you bought the box and the software together, and the software had names like Cisco IOS, Juniper Junos, and Arista EOS. Microsoft, meanwhile, is one of the largest buyers of networking gear on the planet. Open-source the network OS, standardize the interface to the silicon, and the box becomes a commodity you buy the way you buy servers — from white-box and ODM vendors, on merchant silicon, at commodity prices. Whatever Microsoft could have earned selling a network OS is rounding error next to what it saves buying switches for Azure.

## Giving away the standard

Commoditizing the software was phase one. Phase two was to stop owning it.

A project controlled by Microsoft would always be suspect, especially to the vendors it was commoditizing. By moving SONiC to the Linux Foundation, Microsoft made it neutral enough that competitors could adopt it and contribute to it without having to trust Redmond. That is how a Microsoft project became an industry standard: not by winning the market, but by declining to own the thing.

This isn't new either. Eric S. Raymond described the same maneuver in 1999, in a section of *The Cathedral and the Bazaar* called "[Open Source as a Strategic Weapon](http://www.catb.org/~esr/writings/cathedral-bazaar/magic-cauldron/ar01s11.html)": DEC funding the X Window System to "reset the competition" against Sun, and vendors funding Apache so that none of them had to beat Microsoft alone.

## The license is a means, not a value

Here is the uncomfortable part for anyone who came to free software for the freedom. The license is a mechanism, not a motive. Copyleft and permissive licensing are levers with different properties: copyleft stops a layer from being enclosed by anyone, and permissive licensing maximizes adoption. SONiC ships under a [mix of GPL and Apache](https://en.wikipedia.org/wiki/SONiC_(operating_system)), chosen for what it accomplishes rather than for what it professes.

The test of whether openness is a value or a means is simple: does it survive when it stops paying? Usually it does not. MongoDB moved to the SSPL in 2018, Elastic followed in 2021 before adding the AGPL back in 2024, HashiCorp relicensed Terraform in 2023 and got forked as OpenTofu, Redis left the BSD license in 2024, and Wizards of the Coast tried to revoke the open D&D license in 2023 once it had done its job. The modern move is to build the escape hatch in advance: a permissive license plus a contributor license agreement, so the company keeps the option to relicense later.

None of which means the contributors are cynical, or that the code isn't genuinely free. Both things can be true at once: the software is free, and the freedom is on the company's terms.

<!-- TODO(author): your paragraph. Say plainly what you think the uncanny part is — that "open" here is a tactic whose value is instrumental. -->

## Further reading

- Joel Spolsky, [Strategy Letter V: The Economics of Open Source](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/) (2002)
- Eric S. Raymond, [Open Source as a Strategic Weapon](http://www.catb.org/~esr/writings/cathedral-bazaar/magic-cauldron/ar01s11.html), in *The Cathedral and the Bazaar* (1999)
- Gwern Branwen, [Laws of Tech: Commoditize Your Complement](https://gwern.net/complement) (2018–2022)
- Linux Foundation, [SONiC Moves to the Linux Foundation](https://www.linuxfoundation.org/press/press-release/software-for-open-networking-in-the-cloud-sonic-moves-to-the-linux-foundation) (2022)
- [SONiC Foundation](https://sonicfoundation.dev/) and [SONiC on Wikipedia](https://en.wikipedia.org/wiki/SONiC_(operating_system))

<!-- TODO(author): still needed — your personal stake in the topic, and a closing line. The middle is factual scaffolding you can rewrite in your voice. -->

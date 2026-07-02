# KLM KT34A Setup Guide

This guide documents setup and teardown for the club's 4-element KLM KT-34A HF yagi antenna. This antenna is the "little brother" of the 6-element KLM KT36A.

The KT34 packs most of the punch of its big brother, but with significantly easier setup and teardown. It has two driven elements, driven 180 degrees out of phase with each other, one director and one reflector. It has a 3KW rated 4:1 balun serving the 50W match. Phasing on the driven elements is achieved with straps shorting the left side of one driven element to the right side of the other, and vice-versa.

---
**Table of Contents**
* Table of Contents
{:toc}
---


## Overview

This page covers how to put together our club's KLM KT34. It is specific to our club, in that we have a set of markings on the antenna that help us repeat the process.

This guide does not cover the full process of building the antenna for the first time from the parts that come out of the box(es).


## Tools Required

- 9/16 socket, for mast and mounting plate bolts
- 7/16 socket, for the truss system
- 1/2 inch socket, for the mounting plate bolts
- 3/8 socket, for the phasing strap mounts
- 11/32 socket, for the balun mount
- 5/16 socket, for hose clamps (a flathead screwdriver would work but be much less pleasant)

NOTE: Maybe we should consider standardizing these a bit better. We could easily optimize out two of these required socket sizes (11/32 and 1/2)

## Markings and Modifications

There are several markings on the antenna which tell us how to disassemble/reassemble the antenna:

- There is a sticker with an arrow on the front of the boom, pointing towards the side of the antenna with the directors
- Each element has black and red electrical tape on the left and right sides of the element (respectively)
- There are numbers etched into each element, indicating how far back from the director they should be.
- There is also paint indicating where the truss pieces mount to the boom

## Unpacking

The antenna is usually stored on the roof of the PAARA box trailer. It's stored as:

- 4 elements
- Boom bundle:
	- Truss rope
	- Mast mounting plate

The "boom bundle" is taped together with electrical tape.

This antenna has only one boom section, and it has no mast adapter. Instead, on the bottom of the mounting plate, there are U-Bolts which go directly onto a mast stub of some kind.

## Assembly

The KT34 comes together relatively quickly. This is because it has fewer parts than its big brother.

1. Put the single piece boom on a few sawhorses
2. Line up the elements in the right order and put them in the right place. There are etched numbers in each element. 1 is the director and should be closest to the direction arrow, 2 is the front driven element, etc. There is spray paint indicating where they should go, but refer to the data sheet (TODO: Link) which has measurements in case of ambiguity.
3. Attach the elements with the hose clamps. Unfortunately, the drive point of the hose clamps are not all facing the same direction. Do not use the hose clam to determine the left/right alignment of the elements. Additionally, the linear loading sections do not all face the right way. Trust the red/black electrical tape! Use a 5/16" socket to tighten the hose clamps. 

IMPORTANT: it is very easy to damage hose clamps by over-tightening. Especially if using power tools, be extra careful when tightening.

4. Place sawhorses under a couple of elements (ideally two on each side) to keep the antenna from rolling
5. Attach the phasing lines. The phasing lines go across the two driven elements, from the left side on the front element to the right on the rear, and vice-versa. There is a 3D-printed spacer (de KK6ISP) which mostly keeps them from shorting to each other. The bolts are 3/8".

IMPORTANT: The phasing straps connect to the outside set of bolts on each driven element, not the inside (TODO: Photo).

IMPORTANT: Make sure the phasings straps are not shorting (touching) anything other than the bolts they are connected to. Tape them in place.

6. Attach the Balun. The Balun has a set of U-Bolts for mounting. It should point with the coax facing towards the mast. The balun wires have ring terminals which connect to the outside set of bolts on the matching element of the forward driven eleent. These bolts are 11/32

7. Loosen the bolts attaching the mounting plate to the boom so that it can rotate not quite freely. This will make it easier to catch the mount when it's time.

8. Level the elements. Choose one (usually the front) to act as the reference, and, keeping one level on the reference, use a second level to check that all the elements are level to each other.


## Mounting

For instructions on preparing a tower to accept an antenna, see one of the tower setup guides.

The KT34A slides onto a mast with a set of U-bolts, not to be confused with U-boats, which are German submarines from the 1940s.

First, fully loosen the U-bolts. There are eight nuts, two on each of four U-bolts. They need to be loosened almost all the way for a normal mast to fit through.

Make sure that when you slide the mast through, the saddles are on the opposite side of the mast from the curved part of the U-bolt. Push the antenna down a bit, then attach the truss mount above the antenna on the mast. Tighten down the truss mount first, and then push the antenna down the mast until the truss is taut. You don't want slack, but you also don't want the truss cable to become a piano string.

IMPORTANT: If you are going to rotate the antenna as part of hoisting, make sure the mast plate is loosened so it doesn't stay sideways.

Tie a string to one element on each side of the boom, to keep the antenna parallel to the ground while hoisting until the bolts are tightened, if you will rotate the antenna.

Before hoisting the antenna, make sure the coax is connected. Tape the coax to the boom to reduce stress on it and keep it from bouncing around and interfering with anything.


## Teardown and Packing

Teardown is essentially the reverse of assembly.

1. First, follow any tower-specific lowering instructions. This will be a bit different for different kinds of towers. In any case, remove the feedline before taking the antenna off of the tower

2. Place the antenna on a set of sawhorses. Put at least one sawhorse under an element on each side of the antenna to keep it from rolling.

3. Remove the truss, balun, and phasing straps before pulling elements or hardware apart.

4. Remove elements by loosening hose clamps. The hose clamps on the two elements on the end do not need to be fully opened. The two middle elements do.

5. Tape the phasing straps and mounting plate to the boom, as well as the truss rope.

We typically pack the elements on top of the box trailer. The elements get wired down to the ladder rack in pairs, and the boom gets wired down on its own.


## Common Gotchas and FAQs

- The Balun mounts near the driven element which has the extra cross-component. The Balun leads mount to the bolts on the cross-compoonent. This is the front driven element. Photo:
- The phasing lines mount to the outside set of bolts on the two driven elements.
- The truss wants enough tension to hold shape, but not so much that it bends the boom or distorts the driven elements.

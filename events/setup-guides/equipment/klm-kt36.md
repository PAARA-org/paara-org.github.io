# KLM KT36 Setup Guide

This guide documents setup and teardown for the club's 6-element KLM KT-36 HF yagi antenna, affectionately known as the "tri bander".

The KT-36 is an interesting antenna: a multi-band HF Yagi with 9-10 dB of gain across the three bands. It has served us well at field day.

---
**Table of Contents**
* Table of Contents
{:toc}
---

## Overview

This page covers how to put together our club's KLM KT36. It is specific to our club, in that we have a set of markings on the antenna that help us repeat the process.

This guide does not cover the full process of building the antenna for the first time from the parts that come out of the box(es).

The KT36 is a beast. It has two driven elements, driven 180 degrees out of phase with each other, three directors and one reflector. It has a 3KW rated 4:1 balun serving the 50W match. Phasing on the driven elements is achieved with straps shorting the left side of one driven element to the right side of the other, and vice-versa.


## Tools Required

- 9/16 socket, for mast and mounting plate bolts (TODO: are the mast u-bolts 9/16 or 1/2?)
- 7/16 socket, for bolts holding the pieces of the boom together and the truss system.
- 3/8 socket, for balun mounts and phasing strap mounts
- 5/16 socket, for hose clamps (a flathead screwdriver would work but be much less pleasant)

## Markings and Modifications

There are several markings on the antenna which tell us how to disassemble/reassemble the antenna:

- There are stripes of green electrical tape just to the outside of the bolts where we take the boom apart. We only split the boom in half.
- There is a sticker with an arrow on the front of the boom, pointing towards the side of the antenna with the directors
- Each element has black and red electrical tape on the left and right sides of the element (respectively)
- There is color coded spray paint indicating which elements go where. There are actually several layers of spray paint on the antenna. The newest set is orange, blue, yellow/blue, and purple/yellow.
- There is also paint indicating where the truss pieces mount to the boom
- Each element (will have) has numbers etched into them, indicating how far from the front director they should be installed.

## Unpacking

The antenna is usually stored on the roof of the PAARA box trailer. It's stored as:

- 6 elements (color coded)
- Boom bundle:
	- 2 halves of the boom
	- Mast
	- Truss cable
	- Mast mounting plate

The "boom bundle" is taped together with electrical tape.

## Assembly

The KT36 assembly process is relatively straightforward. It's similar to the KT34, but with a couple extra elements.

One interesting note is that the second element from the front is much shorter and stubbier than the others. It is only active on the 10m band, and is not linearly loaded like the others.

The biggest mechanical difference between assembling the KT34 and KT36 is the mounting system. The KT36 has a short mast that stays with it. This means that it needs an adapter to attach to a tower with a mast.

1. After unpacking, place the two pieces of the boom on sawhorses. They attach with two 7/16" bolts. There is a sticker with an arrow on the forward boom section. 
2. Line up the elements in the right order and put them in the right place. There are etched numbers in each element. 1 is the director and should be closest to the direction arrow, 2 is the short 10m director, 3 is the front driven element, etc. There is spray paint indicating where they should go, but refer to the data sheet (TODO: Link) which has measurements in case of ambiguity. (TODO: add etched numbers to the elements next time we disassemble)
3. Attach the elements with the hose clamps. Unfortunately, the drive point of the hose clamps are not all facing the same direction. Do not use the hose clam to determine the left/right alignment of the elements. Additionally, the linear loading sections do not all face the right way. Trust the red/black electrical tape! Use a 5/16" socket to tighten the hose clamps. 

IMPORTANT: it is very easy to damage hose clamps by over-tightening. Especially if using power tools, be extra careful when tightening.

4. Place sawhorses under a couple of elements (ideally two on each side) to keep the antenna from rolling

5. Attach the phasing lines. The phasing lines go across the two driven elements, from the left side on the front element to the right on the rear, and vice-versa. There are 3D-printed spacers (de KK6ISP) which mostly keeps them from shorting to each other. The bolts are 3/8".

IMPORTANT: The phasing straps connect to the outside set of bolts on each driven element, not the inside (TODO: Photo).

IMPORTANT: Make sure the phasings straps are not shorting (touching) anything other than the bolts they are connected to. Tape them in place.

6. Attach the Balun. The Balun has a set of U-Bolts for mounting. It should point with the coax facing towards the mast. The balun wires have ring terminals which connect to the outside set of bolts on the matching element of the forward driven eleent. These bolts are 3/8".

7. Attach, or reposition the truss clamps. These each have one 7/16" bolt, and there are paint marks for where they should go.

8. Attach the mast to the mast mounting plate. There are a set of U-bolts (either 3/8" or 7/16" TODO to check this) that attach the mast to the mounting plate. It should be mounted high enough (that is, enough above the antenna boom) to allow the truss to be taut, but not so high that it's too tight. Tighten the truss onto the mast before tigheting the mast onto the boom.

9. Loosen the bolts attaching the mounting plate to the boom so that it can rotate not quite freely. This will make it easier to catch the mount when it's time. (only applicable if mounting in a way that requires rotation such as the tower trailer)

10. Level the elements. Choose one (usually the front) to act as the reference, and, keeping one level on the reference, use a second level to check that all the elements are level to each other.


## Mounting

For instructions on preparing a tower to accept an antenna, see one of the tower setup guides.

If you want to mount the KT36A directly to a mast on a tower, you should skip step 8, and refer to the KT34 mounting instructions. 

If you are mounting the KT36A using its mast (such as to a tower with an adapter or with a through-hole mounting system), skip step 9 above. Place the bottom of the mast into the hole that can accept it and tighten down. This how we typically mount the KT36A to the AB621.

IMPORTANT: If you are going to rotate the antenna as part of hoisting, make sure the mast plate is loosened so it doesn't stay sideways.

Tie a string to one element on each side of the boom, to keep the antenna parallel to the ground while hoisting until the bolts are tightened, if you will rotate the antenna.

Before hoisting the antenna, make sure the coax is connected. Tape the coax to the boom to reduce stress on it and keep it from bouncing around and interfering with anything.

## Teardown

Teardown is essentially the reverse of assembly.

1. First, follow any tower-specific lowering instructions. This will be a bit different for different kinds of towers. In any case, remove the feedline before taking the antenna off of the tower

2. Place the antenna on a set of sawhorses. Put at least one sawhorse under an element on each side of the antenna to keep it from rolling.

3. Remove the truss, balun, and phasing straps before pulling elements or hardware apart.

4. Remove the mast by loosening the U-bolts and pulling it through.

5. Remove elements by loosening hose clamps. The hose clamps on the two elements on the end do not need to be fully opened. The four middle elements do.

6. Break the boom in half by removing the bolts. Put the bolts back in their holes so they don't get lost.


7. Tape the two boom pieces together. Tape the mast to the two boom pieces. Tape the phasing straps and mounting plate to the boom pieces, as well as the truss rope.


## Common Gotchas and FAQs

- The Balun mounts near the driven element which has the extra cross-component. The Balun leads mount to the bolts on the cross-compoonent. This is the rear driven element. Photo:
- The phasing lines mount to the outside set of bolts on the two driven elements.
- When you rotate the mast ahead of mounting, rotate it so that the truss stays above the elements, and the end of the mast which will go into the tower is slightly below the elements. This is to the left, facing towards the front of the antenna.

# APRS Information & Standards

* [aprs.fi](https://aprs.fi)
    * Check APRS traffic (beacons and messages), track one station on the map.
* [WinAPRS & MacAPRS](http://www.winaprs.org/)
    * Fairly capable software, half-way decent mapping support. $50 to register.
* [UI-View](http://www.ui-view.org/)
    * Easiest-to-use APRS software, but no zoomable maps without purchase of Precision Mapping.
    * APRS for Microsoft MapPoint. Great maps, but very limited features. $104 with MapPoint, $30 without.
* [AGWTracker](http://www.elcom.gr/sv2agw/tracker.htm)
    * Newest Windows APRS software. Great GeoTIFF map support, several cool features. $50 to register.

## APRS Map Servers

For those that don’t have APRS software installed or even a means of receiving packet radio transmissions, [APRS.fi](https://aprs.fi){:target="_blank"} and [findU](http://findu.com/){:target="_blank"} offer Web-based mapping of APRS stations.

## Maps and Data Files for WinAPRS

* [Latest FCC Callsign Database](http://wireless.fcc.gov/uls/data/complete/l_amat.zip)
    * This is updated once a week on weekends (so downloading it every monday will keep you current), but before you download it, be forewarned that it is huge—66.1 MB at last check! Download this file, extract its contents from the archive, place them in the directory of your choice. Next, launch WinAPRS and select CallBook DataBase from the Settings menu, select FCC 2000 & Later from the top drop-down listbox, enter the full path to the EN.dat and EN.idx files in the next two fields, respectively, and finally, press the Call Index button. Wait for the process to finish, and then press the OK button.
* Bay Area Elevation Maps (DEMs)
    * These data files contain fairly coarse, fairly old elevation data that WinAPRS can use to display a color shaded base map that will give you a rough idea of the terrain, while still seeing the roads and labels from the TIGER maps. Download these files, uncompress them with WinZip (or something similar), then place them in your c:\Program Files\WinAPRS\DEMS\ directory.

## ARDF-Capable APRS Software

* WinAPRS & MacAPRS
    * While WinAPRS has built-in direction finding features, the documentation and operation thereof is confusing and insufficient. But this will be our first testbed platform for group foxhunts. The software's problems include lack of a single non-modal (or otherwise) dialog that can be used to create and update DF reports, an inability to actually send out DF reports (due to the poor documentation), and a few others.
* Xastir
    * From the specs, this looks promising, in that it supports a wide variety of maps, is free, and can supposedly create sophisticated DF reports which even include the angle of uncertainty (think a pie-wedge instead of a line) to make the composite DF reports more useful.
* DosAPRS
    * The original is limited in that it only supports the rudimentary DosAPRS map format, but it offers the most complete support for the APRS DF features.
* UI-View
    * This software has no built-in DF features, though a plug-in, DF-Plot, is available for older versions that will plot three DF reports and create an APRS object at their intersection. Supposedly, DF-Plot will even respond to a non-standard DF report format in APRS beacon comment fields: %bbb%aaa, where bbb is the bearing in degrees, and aaa is the accuracy, also in degrees.
* GPSS
    * This has got to be the most convoluted piece of software I have ever tried. It might work great, but I was not able to follow the documentation well enough to figure out how to get it to import maps or do anything else useful.

## TNCs and Sound Card Interfaces

* Packet Engines
    * A packet engine is software that replicates the functionality of a real TNC. The one K6WEB has used is
* Homebrew Sound Card Interfaces
    * [Stephen Smith’s sound card interface](http://members.aol.com/wa8lmf/ham/tonekeyer.htm){:target="_blank"} designs offer a good explanation about the level conversions that are necessary.

## APRS and GPS-Related Standards

* APRS Specifications
    * For much of its life, APRS was loosely documented in the many text files that are part of the APRSdos documentation files. Recently, it has been formalized as a specification in cooperation with TAPR: [APRS Protocol Specification 1.01](ftp://ftp.tapr.org/aprssig/aprsspec/spec/aprs101/APRS101.pdf) (PDF file, adopted August 31, 2000).
* [Height Above Average Terrain (HAAT)](http://www.qsl.net/miarc/haat_inst.html)
    * An international standard used for determining signal coverage of radio signals. The antenna height above the average terrain elevations from 3.2 to 16 kilometers (2 to 10 miles) from the antenna for the eight directions spaced evenly for each 45° of azimuth starting with true north. Note: In general, a different antenna height above average terrain will be determined in each direction from the antenna. The average of these eight heights is the antenna height above average terrain. In some cases, such as seashore, fewer than eight directions may be used. For international coordination, it is officially measured in meters, even by the FCC in the USA.
* Main National Frequency: 144.39 MHz
    * This is the frequency that stations and digipeaters listen and beacon on.
* 2.5mm Stereo Phone Plug/Jack Serial Port
    * (Discouraged because it is a poor choice of connector for mobile use because the spring metal connectors in stereo phone jacks can cause intermitent connections when exposed to vibration such as you might find in a car, and because its three conductors are not enough for two-way data and power, which needs four pins.) This scheme is used for connecting serial GPS receivers to radios and other ham equipment, the [2.5mm stereo phone jack serial port standard](http://www.aprs.net/vm/DOS/MICE.HTM){:target="_blank"} was documented by Bob Bruninga, in his APRSdos documentation as the means of connecting GPS devices to Mic-E encoders. According to this documentation, on the 2.5mm jack (TNC end), Tip is used for TXD (data out), Ring is used for RXD (data in), and Sleve is used for GND (ground). This wiring pattern is used by the Pocket Tracker, the Kenwood TMD-700A, the Kenwood . Any references to the use of a 2-conductor mono phone plug/jack for GPS units is antiquated, as it supports only one-way data transmission, although the Baker2Vegas group found they had fewer problems with vibration using the mono connectors.
* TNC Packet Jack
    * A 6-pin mini-DIN connector that handles audio in, audio out, PTT, and squelch signals.
* Wireless Industry Mounting Bracket Standards
    * If you need to mount a GPS receiver, APRS display, radio, or other equipment in your car, there is an industry-standard mounting hole pattern. The AMPS hole pattern was created by the Association of Mobile Phone Systems?, and consists of four holes in a rectangular pattern measuring .750 inches by 1.188 inches. An alternative hole pattern was created by NEC (National Electronic Commission)? Mounting brackets using these hole patterns are available from [Ram Mounting Systems](http://www.ram-mount.com/){:target="_blank"}, [Pro.Fit International](http://www.pro-fit-intl.com/){:target="_blank"}.
* NMEA 0183 Interface Standard Version 2.0 (NMEA-0183 2.0
    * `The National Marine Electronics Association’s GPS interface standard` The standard plain text format that most GPS receivers output via a serial port. There are different types of GPS sentences, including:
        * `$GPGGA`
            * For position and height. Looks like: $GPGGA,102705,5157.9762,N,00029.3256,W,1,04,2.0,75.7,M,47.6,M,,*62
        * `$GPGLL`
            * For position only. Looks like: $GPGLL,2554.459,N,08020.187,W,154027.281,A
        * `$GPVTG`
            * For velocity and course Looks like: $GPVTG,318.7,T,,M,35.1,N,65.0,K*69
        * `$GPRMC`
            * Position, Course and speed. Looks like: $GPRMC,063909,A,3349.4302,N,11700.3721,W,43.022,89.3,291099,13.6,E*52
        * `$GPWPL`—For sending Waypoints to stand-alone trackers
        * `$PMGLB` —Magelan Altitude
        * `$PGRMZ` —Garmin Altitude
* Sending Waypoints to GPS Displays
    * Some APRS software can send stations and objects from your local map out as special APRS packets that will appear as destination waypoints. The intent is that this waypoint will be your destination, and APRSdos will draw a line from that station’s reported position to the waypoint.
* Overlaying Images From Internet
    * WinAPRS can retrieve images from the Internet and overlay them on existing maps, using coordinate references contained in a “.geo” file. One example of this feature is overlaying current weather satellite or weather radar images, so that you can see what the cloud cover or rain is like.
* Tiny Web Pages
    * If a local node is running the right software, it can respond to query packets with requested information, if available.
* Latitude/Longitude Calculators
    * [Lat/Lon Coordinates can be expressed in one of three ways](http://www.maptools.com/UsingLatLon/Formats.html){:target="_blank"}: decimal degrees, degrees minutes and seconds, or degrees and decimal minutes. Several sites offer conversion between the first two type, including the [FCC Lat/Lon Calculator](http://www.fcc.gov/mb/audio/bickel/DDDMMSS-decimal.html){:target="_blank"} and the [Efficacy.net Lat/Lon Converter](http://efficacy.net/experiments/convert_lat_long/){:target="_blank"}. Also of interest is the [AMSAT GridSquare Converter](http://www.amsat.org/amsat/toys/gridconv.html){:target="_blank"}. APRS coordinates are reported using Degrees and Decimal Minutes, in a concatenated format that looks like DDDMM.MM, so that 12215.79 represents 122.26330 degrees.

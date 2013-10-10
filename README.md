Serpent
=======

Serpent encryption cipher

Overview
--------

This is a Go implementation of the Serpent encryption cipher designed by
Ross Anderson, Eli Biham, Lars Knudsen and based on the python version
written by Frank Stajano, Cambridge University Computer Laboratory 
http://www.cl.cam.ac.uk/~fms27/

The main encryption and decryption functions are working for both normal
and bitslice mode. These functions use a 128-bit Bitstring and require a
256-bit bitstring for a key.


Active work
-----------

I am now working on functions to make this usable for
strings of text. Feel free to fork a copy and use as required. Pull
requests are welcome too.

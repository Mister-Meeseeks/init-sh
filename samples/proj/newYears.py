#!/usr/bin/env initPy

import sys
import myProj.newYears as nye

nCount = int(sys.argv[1]) \
    if len(sys.argv) > 1 else 10

nye.countdown(nCount)

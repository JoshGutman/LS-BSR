#!/usr/bin/env python 

from multiprocessing import Pool

def mp_shell(func, params, numProc):

	p = Pool(numProc)

	p.map(func, params)

	p.terminate()

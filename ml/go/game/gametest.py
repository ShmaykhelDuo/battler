
# python wrapper for package github.com/ShmaykhelDuo/battler/internal/game/gametest within overall package game
# This is what you import to use the package.
# File is generated by gopy. Do not edit.
# gopy pkg -exclude=model,ml1,ml2,moveml ../../internal/game ../../internal/bot context encoding/json

# the following is required to enable dlopen to open the _go.so file
import os,sys,inspect,collections
try:
	import collections.abc as _collections_abc
except ImportError:
	_collections_abc = collections

cwd = os.getcwd()
currentdir = os.path.dirname(os.path.abspath(inspect.getfile(inspect.currentframe())))
os.chdir(currentdir)
from . import _game
from . import go

os.chdir(cwd)

# to use this code in your end-user python file, import it as follows:
# from game import gametest
# and then refer to everything using gametest. prefix
# packages imported by this package listed below:

from . import game



# ---- Types ---


#---- Enums from Go (collections of consts with same type) ---


#---- Constants from Go: Python can only ask that you please don't change these! ---


# ---- Global Variables: can only use functions to access ---
def EffectDescExpirable():
	"""
	EffectDescExpirable Gets Go Variable: gametest.EffectDescExpirable
	
	"""
	return game.EffectDescription(handle=_game.gametest_EffectDescExpirable())

def Set_EffectDescExpirable(value):
	"""
	Set_EffectDescExpirable Sets Go Variable: gametest.EffectDescExpirable
	
	"""
	if isinstance(value, go.GoClass):
		_game.gametest_Set_EffectDescExpirable(value.handle)
	else:
		_game.gametest_Set_EffectDescExpirable(value)



# ---- Interfaces ---


# ---- Structs ---

# Python type for struct gametest.EffectExpirable
class EffectExpirable(go.GoClass):
	""""""
	def __init__(self, *args, **kwargs):
		"""
		handle=A Go-side object is always initialized with an explicit handle=arg
		otherwise parameters can be unnamed in order of field names or named fields
		in which case a new Go object is constructed first
		"""
		if len(kwargs) == 1 and 'handle' in kwargs:
			self.handle = kwargs['handle']
			_game.IncRef(self.handle)
		elif len(args) == 1 and isinstance(args[0], go.GoClass):
			self.handle = args[0].handle
			_game.IncRef(self.handle)
		else:
			self.handle = _game.gametest_EffectExpirable_CTor()
			_game.IncRef(self.handle)
	def __del__(self):
		_game.DecRef(self.handle)
	def __str__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'gametest.EffectExpirable{'
		first = True
		for v in pr:
			if callable(v[1]):
				continue
			if first:
				first = False
			else:
				sv += ', '
			sv += v[0] + '=' + str(v[1])
		return sv + '}'
	def __repr__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'gametest.EffectExpirable ( '
		for v in pr:
			if not callable(v[1]):
				sv += v[0] + '=' + str(v[1]) + ', '
		return sv + ')'
	def Desc(self):
		"""Desc() object
		
		Desc returns the effect's description.
		"""
		return game.EffectDescription(handle=_game.gametest_EffectExpirable_Desc(self.handle))
	def Clone(self):
		"""Clone() object
		
		Clone returns a clone of the effect.
		"""
		return game.Effect(handle=_game.gametest_EffectExpirable_Clone(self.handle))
	def Expire(self, goRun=False):
		"""Expire() """
		_game.gametest_EffectExpirable_Expire(self.handle, goRun)
	def HasExpired(self, turnState):
		"""HasExpired(object turnState) bool
		
		HasExpired reports whether the effect has expired.
		"""
		return _game.gametest_EffectExpirable_HasExpired(self.handle, turnState.handle)


# ---- Slices ---


# ---- Maps ---


# ---- Constructors ---
def NewEffectExpirable(expired):
	"""NewEffectExpirable(bool expired) object"""
	return EffectExpirable(handle=_game.gametest_NewEffectExpirable(expired))


# ---- Functions ---




# python wrapper for package github.com/ShmaykhelDuo/battler/internal/game/common within overall package game
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
# from game import common
# and then refer to everything using common. prefix
# packages imported by this package listed below:

from . import game



# ---- Types ---


#---- Enums from Go (collections of consts with same type) ---


#---- Constants from Go: Python can only ask that you please don't change these! ---


# ---- Global Variables: can only use functions to access ---


# ---- Interfaces ---


# ---- Structs ---

# Python type for struct common.Collectible
class Collectible(go.GoClass):
	"""Collectible is a mixin that allows storing amounts.\n"""
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
			self.handle = _game.common_Collectible_CTor()
			_game.IncRef(self.handle)
	def __del__(self):
		_game.DecRef(self.handle)
	def __str__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'common.Collectible{'
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
		sv = 'common.Collectible ( '
		for v in pr:
			if not callable(v[1]):
				sv += v[0] + '=' + str(v[1]) + ', '
		return sv + ')'
	def Amount(self):
		"""Amount() int
		
		Amount returns the collectible's amount.
		"""
		return _game.common_Collectible_Amount(self.handle)
	def Increase(self, amount, goRun=False):
		"""Increase(int amount) 
		
		Increase increases the collectible's amount.
		"""
		_game.common_Collectible_Increase(self.handle, amount, goRun)
	def Decrease(self, amount, goRun=False):
		"""Decrease(int amount) 
		
		Decrease decreases the collectible's amount.
		"""
		_game.common_Collectible_Decrease(self.handle, amount, goRun)
	def HasExpired(self, turnState):
		"""HasExpired(object turnState) bool
		
		HasExpired reports whether the effect has expired.
		"""
		return _game.common_Collectible_HasExpired(self.handle, turnState.handle)

# Python type for struct common.DurationExpirable
class DurationExpirable(go.GoClass):
	"""DurationExpirable is a mixin that allows expiring effects after specified turn.\n"""
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
			self.handle = _game.common_DurationExpirable_CTor()
			_game.IncRef(self.handle)
	def __del__(self):
		_game.DecRef(self.handle)
	def __str__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'common.DurationExpirable{'
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
		sv = 'common.DurationExpirable ( '
		for v in pr:
			if not callable(v[1]):
				sv += v[0] + '=' + str(v[1]) + ', '
		return sv + ')'
	def TurnsLeft(self, turnState):
		"""TurnsLeft(object turnState) int
		
		TurnsLeft returns the number of turns left until this effect expires.
		"""
		return _game.common_DurationExpirable_TurnsLeft(self.handle, turnState.handle)
	def HasExpired(self, turnState):
		"""HasExpired(object turnState) bool
		
		HasExpired reports whether the effect has expired.
		"""
		return _game.common_DurationExpirable_HasExpired(self.handle, turnState.handle)


# ---- Slices ---


# ---- Maps ---


# ---- Constructors ---
def NewCollectible(amount):
	"""NewCollectible(int amount) object
	
	NewCollectible returns a new collectible with specified amount.
	"""
	return Collectible(handle=_game.common_NewCollectible(amount))
def NewDurationExpirable(expCtx):
	"""NewDurationExpirable(object expCtx) object
	
	NewDurationExpirable returns a [DurationExpirable] with specified expiry turn.
	"""
	return DurationExpirable(handle=_game.common_NewDurationExpirable(expCtx.handle))


# ---- Functions ---



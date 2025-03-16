
# python wrapper for package github.com/ShmaykhelDuo/battler/internal/bot within overall package game
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
# from game import bot
# and then refer to everything using bot. prefix
# packages imported by this package listed below:

from . import context
from . import match



# ---- Types ---


#---- Enums from Go (collections of consts with same type) ---


#---- Constants from Go: Python can only ask that you please don't change these! ---


# ---- Global Variables: can only use functions to access ---


# ---- Interfaces ---


# ---- Structs ---

# Python type for struct bot.RevAdapter
class RevAdapter(go.GoClass):
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
			self.handle = _game.bot_RevAdapter_CTor()
			_game.IncRef(self.handle)
	def __del__(self):
		_game.DecRef(self.handle)
	def __str__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'bot.RevAdapter{'
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
		sv = 'bot.RevAdapter ( '
		for v in pr:
			if not callable(v[1]):
				sv += v[0] + '=' + str(v[1]) + ', '
		return sv + ')'
	def GetStateInit(self):
		"""GetStateInit() object"""
		return match.GameState(handle=_game.bot_RevAdapter_GetStateInit(self.handle))
	def GetState(self, skill):
		"""GetState(int skill) object res, str err"""
		return match.GameState(handle=_game.bot_RevAdapter_GetState(self.handle, skill))
	def SendState(self, ctx, state):
		"""SendState(object ctx, object state) str"""
		return _game.bot_RevAdapter_SendState(self.handle, ctx.handle, state.handle)
	def SendError(self, ctx, err):
		"""SendError(object ctx, str err) str"""
		return _game.bot_RevAdapter_SendError(self.handle, ctx.handle, err)
	def SendEnd(self, ctx):
		"""SendEnd(object ctx) str"""
		return _game.bot_RevAdapter_SendEnd(self.handle, ctx.handle)
	def RequestSkill(self, ctx):
		"""RequestSkill(object ctx) int, str"""
		return _game.bot_RevAdapter_RequestSkill(self.handle, ctx.handle)

# Python type for struct bot.RandomBot
class RandomBot(go.GoClass):
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
			self.handle = _game.bot_RandomBot_CTor()
			_game.IncRef(self.handle)
	def __del__(self):
		_game.DecRef(self.handle)
	def __str__(self):
		pr = [(p, getattr(self, p)) for p in dir(self) if not p.startswith('__')]
		sv = 'bot.RandomBot{'
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
		sv = 'bot.RandomBot ( '
		for v in pr:
			if not callable(v[1]):
				sv += v[0] + '=' + str(v[1]) + ', '
		return sv + ')'
	def SendState(self, ctx, state):
		"""SendState(object ctx, object state) str"""
		return _game.bot_RandomBot_SendState(self.handle, ctx.handle, state.handle)
	def SendError(self, ctx, err):
		"""SendError(object ctx, str err) str"""
		return _game.bot_RandomBot_SendError(self.handle, ctx.handle, err)
	def SendEnd(self, ctx):
		"""SendEnd(object ctx) str"""
		return _game.bot_RandomBot_SendEnd(self.handle, ctx.handle)
	def RequestSkill(self, ctx):
		"""RequestSkill(object ctx) int, str"""
		return _game.bot_RandomBot_RequestSkill(self.handle, ctx.handle)


# ---- Slices ---


# ---- Maps ---


# ---- Constructors ---
def NewAdapter():
	"""NewAdapter() object"""
	return RevAdapter(handle=_game.bot_NewAdapter())


# ---- Functions ---



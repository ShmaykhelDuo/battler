import math
from matplotlib.ticker import Formatter
import matplotlib as mpl
import numpy as np

class EngFormatter(Formatter):
    """
    Format axis values using engineering prefixes to represent powers
    of 1000, plus a specified unit, e.g., 10 MHz instead of 1e7.
    """

    # The SI engineering prefixes
    ENG_PREFIXES = {
        -30: "кв",
        -27: "рн",
        -24: "и",
        -21: "з",
        -18: "а",
        -15: "ф",
        -12: "п",
         -9: "н",
         -6: "мк",
         -3: "м",
          0: "",
          3: "к",
          6: "М",
          9: "Г",
         12: "Т",
         15: "П",
         18: "Э",
         21: "З",
         24: "И",
         27: "Рн",
         30: "Кв"
    }

    def __init__(self, unit="", places=None, sep=" ", *, usetex=None,
                 useMathText=None):
        r"""
        Parameters
        ----------
        unit : str, default: ""
            Unit symbol to use, suitable for use with single-letter
            representations of powers of 1000. For example, 'Hz' or 'm'.

        places : int, default: None
            Precision with which to display the number, specified in
            digits after the decimal point (there will be between one
            and three digits before the decimal point). If it is None,
            the formatting falls back to the floating point format '%g',
            which displays up to 6 *significant* digits, i.e. the equivalent
            value for *places* varies between 0 and 5 (inclusive).

        sep : str, default: " "
            Separator used between the value and the prefix/unit. For
            example, one get '3.14 mV' if ``sep`` is " " (default) and
            '3.14mV' if ``sep`` is "". Besides the default behavior, some
            other useful options may be:

            * ``sep=""`` to append directly the prefix/unit to the value;
            * ``sep="\N{THIN SPACE}"`` (``U+2009``);
            * ``sep="\N{NARROW NO-BREAK SPACE}"`` (``U+202F``);
            * ``sep="\N{NO-BREAK SPACE}"`` (``U+00A0``).

        usetex : bool, default: :rc:`text.usetex`
            To enable/disable the use of TeX's math mode for rendering the
            numbers in the formatter.

        useMathText : bool, default: :rc:`axes.formatter.use_mathtext`
            To enable/disable the use mathtext for rendering the numbers in
            the formatter.
        """
        self.unit = unit
        self.places = places
        self.sep = sep
        self.set_usetex(usetex)
        self.set_useMathText(useMathText)

    def get_usetex(self):
        return self._usetex

    def set_usetex(self, val):
        if val is None:
            self._usetex = mpl.rcParams['text.usetex']
        else:
            self._usetex = val

    usetex = property(fget=get_usetex, fset=set_usetex)

    def get_useMathText(self):
        return self._useMathText

    def set_useMathText(self, val):
        if val is None:
            self._useMathText = mpl.rcParams['axes.formatter.use_mathtext']
        else:
            self._useMathText = val

    useMathText = property(fget=get_useMathText, fset=set_useMathText)

    def __call__(self, x, pos=None):
        s = f"{self.format_eng(x)}{self.unit}"
        # Remove the trailing separator when there is neither prefix nor unit
        if self.sep and s.endswith(self.sep):
            s = s[:-len(self.sep)]
        return self.fix_minus(s)

    def format_eng(self, num):
        """
        Format a number in engineering notation, appending a letter
        representing the power of 1000 of the original number.
        Some examples:

        >>> format_eng(0)        # for self.places = 0
        '0'

        >>> format_eng(1000000)  # for self.places = 1
        '1.0 M'

        >>> format_eng(-1e-6)  # for self.places = 2
        '-1.00 \N{MICRO SIGN}'
        """
        sign = 1
        fmt = "g" if self.places is None else f".{self.places:d}f"

        if num < 0:
            sign = -1
            num = -num

        if num != 0:
            pow10 = int(math.floor(math.log10(num) / 3) * 3)
        else:
            pow10 = 0
            # Force num to zero, to avoid inconsistencies like
            # format_eng(-0) = "0" and format_eng(0.0) = "0"
            # but format_eng(-0.0) = "-0.0"
            num = 0.0

        pow10 = np.clip(pow10, min(self.ENG_PREFIXES), max(self.ENG_PREFIXES))

        mant = sign * num / (10.0 ** pow10)
        # Taking care of the cases like 999.9..., which may be rounded to 1000
        # instead of 1 k.  Beware of the corner case of values that are beyond
        # the range of SI prefixes (i.e. > 'Y').
        if (abs(float(format(mant, fmt))) >= 1000
                and pow10 < max(self.ENG_PREFIXES)):
            mant /= 1000
            pow10 += 3

        prefix = self.ENG_PREFIXES[int(pow10)]
        if self._usetex or self._useMathText:
            formatted = f"${mant:{fmt}}${self.sep}{prefix}"
        else:
            formatted = f"{mant:{fmt}}{self.sep}{prefix}"

        return formatted

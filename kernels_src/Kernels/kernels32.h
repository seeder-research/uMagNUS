
#if defined(__REAL_IS_DOUBLE__)
  #if defined(cl_khr_fp64) || defined(cl_amd_fp64)
    #define realint_t long
    #define realsum_t long
    #define real_t double
    #define real_t2 double2
    #define real_t3 double3
    #define real_t4 double4
    #if defined(cl_amd_fp64)
      #pragma OPENCL EXTENSION cl_amd_fp64 : enable
    #elif defined(cl_khr_fp64)
      #pragma OPENCL EXTENSION cl_khr_fp64 : enable
    #endif // cl_*_fp64
    #define AS_INT as_long
    #define AS_REAL as_double
  #endif
#else
  #define realint_t int
  #define realsum_t float
  #define real_t float
  #define real_t2 float2
  #define real_t3 float3
  #define real_t4 float4
  #define AS_INT as_int
  #define AS_REAL as_float
#endif // __REAL_IS_DOUBLE__
#ifndef _CONSTANTS_H_
#define _CONSTANTS_H_

#define PI     3.1415926535897932384626433
#define MU0    (4*PI*1e-7)        // Permeability of vacuum in Tm/A
#define QE     1.60217646E-19     // Electron charge in C
#define MUB    9.2740091523E-24   // Bohr magneton in J/T
#define GAMMA0 1.7595e11          // Gyromagnetic ratio of electron, in rad/Ts
#define HBAR   1.05457173E-34

#endif // _CONSTANTS_H_

#ifndef NULL
#define NULL ((void*)0)
#endif // NULL
/* Copyright (C) 1997-2014 Free Software Foundation, Inc.
   This file is part of the GNU C Library.

   The GNU C Library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   The GNU C Library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with the GNU C Library; if not, see
   <http://www.gnu.org/licenses/>.  */

/*
 *	ISO C99: 7.18 Integer types <stdint.h>
 */

#ifndef _STDINT_H
#define _STDINT_H	1

/* features.h */

#ifndef	_FEATURES_H
#define	_FEATURES_H	1

/* These are defined by the user (or the compiler)
   to specify the desired environment:

   __STRICT_ANSI__	ISO Standard C.
   _ISOC99_SOURCE	Extensions to ISO C89 from ISO C99.
   _ISOC11_SOURCE	Extensions to ISO C99 from ISO C11.
   _POSIX_SOURCE	IEEE Std 1003.1.
   _POSIX_C_SOURCE	If ==1, like _POSIX_SOURCE; if >=2 add IEEE Std 1003.2;
			if >=199309L, add IEEE Std 1003.1b-1993;
			if >=199506L, add IEEE Std 1003.1c-1995;
			if >=200112L, all of IEEE 1003.1-2004
			if >=200809L, all of IEEE 1003.1-2008
   _XOPEN_SOURCE	Includes POSIX and XPG things.  Set to 500 if
			Single Unix conformance is wanted, to 600 for the
			sixth revision, to 700 for the seventh revision.
   _XOPEN_SOURCE_EXTENDED XPG things and X/Open Unix extensions.
   _LARGEFILE_SOURCE	Some more functions for correct standard I/O.
   _LARGEFILE64_SOURCE	Additional functionality from LFS for large files.
   _FILE_OFFSET_BITS=N	Select default filesystem interface.
   _BSD_SOURCE		ISO C, POSIX, and 4.3BSD things.
   _SVID_SOURCE		ISO C, POSIX, and SVID things.
   _ATFILE_SOURCE	Additional *at interfaces.
   _GNU_SOURCE		All of the above, plus GNU extensions.
   _DEFAULT_SOURCE	The default set of features (taking precedence over
			__STRICT_ANSI__).
   _REENTRANT		Select additionally reentrant object.
   _THREAD_SAFE		Same as _REENTRANT, often used by other systems.
   _FORTIFY_SOURCE	If set to numeric value > 0 additional security
			measures are defined, according to level.

   The '-ansi' switch to the GNU C compiler, and standards conformance
   options such as '-std=c99', define __STRICT_ANSI__.  If none of
   these are defined, or if _DEFAULT_SOURCE is defined, the default is
   to have _SVID_SOURCE, _BSD_SOURCE, and _POSIX_SOURCE set to one and
   _POSIX_C_SOURCE set to 200809L.  If more than one of these are
   defined, they accumulate.  For example __STRICT_ANSI__,
   _POSIX_SOURCE and _POSIX_C_SOURCE together give you ISO C, 1003.1,
   and 1003.2, but nothing else.

   These are defined by this file and are used by the
   header files to decide what to declare or define:

   __USE_ISOC11		Define ISO C11 things.
   __USE_ISOC99		Define ISO C99 things.
   __USE_ISOC95		Define ISO C90 AMD1 (C95) things.
   __USE_POSIX		Define IEEE Std 1003.1 things.
   __USE_POSIX2		Define IEEE Std 1003.2 things.
   __USE_POSIX199309	Define IEEE Std 1003.1, and .1b things.
   __USE_POSIX199506	Define IEEE Std 1003.1, .1b, .1c and .1i things.
   __USE_XOPEN		Define XPG things.
   __USE_XOPEN_EXTENDED	Define X/Open Unix things.
   __USE_UNIX98		Define Single Unix V2 things.
   __USE_XOPEN2K        Define XPG6 things.
   __USE_XOPEN2KXSI     Define XPG6 XSI things.
   __USE_XOPEN2K8       Define XPG7 things.
   __USE_XOPEN2K8XSI    Define XPG7 XSI things.
   __USE_LARGEFILE	Define correct standard I/O things.
   __USE_LARGEFILE64	Define LFS things with separate names.
   __USE_FILE_OFFSET64	Define 64bit interface as default.
   __USE_BSD		Define 4.3BSD things.
   __USE_SVID		Define SVID things.
   __USE_MISC		Define things common to BSD and System V Unix.
   __USE_ATFILE		Define *at interfaces and AT_* constants for them.
   __USE_GNU		Define GNU extensions.
   __USE_REENTRANT	Define reentrant/thread-safe *_r functions.
   __USE_FORTIFY_LEVEL	Additional security measures used, according to level.

   The macros '__GNU_LIBRARY__', '__GLIBC__', and '__GLIBC_MINOR__' are
   defined by this file unconditionally.  '__GNU_LIBRARY__' is provided
   only for compatibility.  All new code should use the other symbols
   to test for features.

   All macros listed above as possibly being defined by this file are
   explicitly undefined if they are not explicitly defined.
   Feature-test macros that are not defined by the user or compiler
   but are implied by the other feature-test macros defined (or by the
   lack of any definitions) are defined by the file.  */


/* Undefine everything, so we get a clean slate.  */
#undef	__USE_ISOC11
#undef	__USE_ISOC99
#undef	__USE_ISOC95
#undef	__USE_ISOCXX11
#undef	__USE_POSIX
#undef	__USE_POSIX2
#undef	__USE_POSIX199309
#undef	__USE_POSIX199506
#undef	__USE_XOPEN
#undef	__USE_XOPEN_EXTENDED
#undef	__USE_UNIX98
#undef	__USE_XOPEN2K
#undef	__USE_XOPEN2KXSI
#undef	__USE_XOPEN2K8
#undef	__USE_XOPEN2K8XSI
#undef	__USE_LARGEFILE
#undef	__USE_LARGEFILE64
#undef	__USE_FILE_OFFSET64
#undef	__USE_BSD
#undef	__USE_SVID
#undef	__USE_MISC
#undef	__USE_ATFILE
#undef	__USE_GNU
#undef	__USE_REENTRANT
#undef	__USE_FORTIFY_LEVEL
#undef	__KERNEL_STRICT_NAMES

/* Suppress kernel-name space pollution unless user expressedly asks
   for it.  */
#ifndef _LOOSE_KERNEL_NAMES
# define __KERNEL_STRICT_NAMES
#endif

/* Convenience macros to test the versions of glibc and gcc.
   Use them like this:
   #if __GNUC_PREREQ (2,8)
   ... code requiring gcc 2.8 or later ...
   #endif
   Note - they won't work for gcc1 or glibc1, since the _MINOR macros
   were not defined then.  */
#if defined __GNUC__ && defined __GNUC_MINOR__
# define __GNUC_PREREQ(maj, min) \
	((__GNUC__ << 16) + __GNUC_MINOR__ >= ((maj) << 16) + (min))
#else
# define __GNUC_PREREQ(maj, min) 0
#endif


/* If _GNU_SOURCE was defined by the user, turn on all the other features.  */
#ifdef _GNU_SOURCE
# undef  _ISOC95_SOURCE
# define _ISOC95_SOURCE	1
# undef  _ISOC99_SOURCE
# define _ISOC99_SOURCE	1
# undef  _ISOC11_SOURCE
# define _ISOC11_SOURCE	1
# undef  _POSIX_SOURCE
# define _POSIX_SOURCE	1
# undef  _POSIX_C_SOURCE
# define _POSIX_C_SOURCE	200809L
# undef  _XOPEN_SOURCE
# define _XOPEN_SOURCE	700
# undef  _XOPEN_SOURCE_EXTENDED
# define _XOPEN_SOURCE_EXTENDED	1
# undef	 _LARGEFILE64_SOURCE
# define _LARGEFILE64_SOURCE	1
# undef  _DEFAULT_SOURCE
# define _DEFAULT_SOURCE	1
# undef  _BSD_SOURCE
# define _BSD_SOURCE	1
# undef  _SVID_SOURCE
# define _SVID_SOURCE	1
# undef  _ATFILE_SOURCE
# define _ATFILE_SOURCE	1
#endif

/* If nothing (other than _GNU_SOURCE and _DEFAULT_SOURCE) is defined,
   define _DEFAULT_SOURCE, _BSD_SOURCE and _SVID_SOURCE.  */
#if (defined _DEFAULT_SOURCE					\
     || (!defined __STRICT_ANSI__				\
	 && !defined _ISOC99_SOURCE				\
	 && !defined _POSIX_SOURCE && !defined _POSIX_C_SOURCE	\
	 && !defined _XOPEN_SOURCE				\
	 && !defined _BSD_SOURCE && !defined _SVID_SOURCE))
# undef  _DEFAULT_SOURCE
# define _DEFAULT_SOURCE	1
# undef  _BSD_SOURCE
# define _BSD_SOURCE	1
# undef  _SVID_SOURCE
# define _SVID_SOURCE	1
#endif

/* This is to enable the ISO C11 extension.  */
#if (defined _ISOC11_SOURCE \
     || (defined __STDC_VERSION__ && __STDC_VERSION__ >= 201112L))
# define __USE_ISOC11	1
#endif

/* This is to enable the ISO C99 extension.  */
#if (defined _ISOC99_SOURCE || defined _ISOC11_SOURCE \
     || (defined __STDC_VERSION__ && __STDC_VERSION__ >= 199901L))
# define __USE_ISOC99	1
#endif

/* This is to enable the ISO C90 Amendment 1:1995 extension.  */
#if (defined _ISOC99_SOURCE || defined _ISOC11_SOURCE \
     || (defined __STDC_VERSION__ && __STDC_VERSION__ >= 199409L))
# define __USE_ISOC95	1
#endif

/* This is to enable compatibility for ISO C++11.

   So far g++ does not provide a macro.  Check the temporary macro for
   now, too.  */
#if ((defined __cplusplus && __cplusplus >= 201103L)			      \
     || defined __GXX_EXPERIMENTAL_CXX0X__)
# define __USE_ISOCXX11	1
#endif

/* If none of the ANSI/POSIX macros are defined, or if _DEFAULT_SOURCE
   is defined, use POSIX.1-2008 (or another version depending on
   _XOPEN_SOURCE).  */
#ifdef _DEFAULT_SOURCE
# if !defined _POSIX_SOURCE && !defined _POSIX_C_SOURCE
#  define __USE_POSIX_IMPLICITLY	1
# endif
# undef  _POSIX_SOURCE
# define _POSIX_SOURCE	1
# undef  _POSIX_C_SOURCE
# define _POSIX_C_SOURCE	200809L
#endif
#if ((!defined __STRICT_ANSI__ || (_XOPEN_SOURCE - 0) >= 500) && \
     !defined _POSIX_SOURCE && !defined _POSIX_C_SOURCE)
# define _POSIX_SOURCE	1
# if defined _XOPEN_SOURCE && (_XOPEN_SOURCE - 0) < 500
#  define _POSIX_C_SOURCE	2
# elif defined _XOPEN_SOURCE && (_XOPEN_SOURCE - 0) < 600
#  define _POSIX_C_SOURCE	199506L
# elif defined _XOPEN_SOURCE && (_XOPEN_SOURCE - 0) < 700
#  define _POSIX_C_SOURCE	200112L
# else
#  define _POSIX_C_SOURCE	200809L
# endif
# define __USE_POSIX_IMPLICITLY	1
#endif

#if defined _POSIX_SOURCE || _POSIX_C_SOURCE >= 1 || defined _XOPEN_SOURCE
# define __USE_POSIX	1
#endif

#if defined _POSIX_C_SOURCE && _POSIX_C_SOURCE >= 2 || defined _XOPEN_SOURCE
# define __USE_POSIX2	1
#endif

#if (_POSIX_C_SOURCE - 0) >= 199309L
# define __USE_POSIX199309	1
#endif

#if (_POSIX_C_SOURCE - 0) >= 199506L
# define __USE_POSIX199506	1
#endif

#if (_POSIX_C_SOURCE - 0) >= 200112L
# define __USE_XOPEN2K		1
# undef __USE_ISOC95
# define __USE_ISOC95		1
# undef __USE_ISOC99
# define __USE_ISOC99		1
#endif

#if (_POSIX_C_SOURCE - 0) >= 200809L
# define __USE_XOPEN2K8		1
# undef  _ATFILE_SOURCE
# define _ATFILE_SOURCE	1
#endif

#ifdef	_XOPEN_SOURCE
# define __USE_XOPEN	1
# if (_XOPEN_SOURCE - 0) >= 500
#  define __USE_XOPEN_EXTENDED	1
#  define __USE_UNIX98	1
#  undef _LARGEFILE_SOURCE
#  define _LARGEFILE_SOURCE	1
#  if (_XOPEN_SOURCE - 0) >= 600
#   if (_XOPEN_SOURCE - 0) >= 700
#    define __USE_XOPEN2K8	1
#    define __USE_XOPEN2K8XSI	1
#   endif
#   define __USE_XOPEN2K	1
#   define __USE_XOPEN2KXSI	1
#   undef __USE_ISOC95
#   define __USE_ISOC95		1
#   undef __USE_ISOC99
#   define __USE_ISOC99		1
#  endif
# else
#  ifdef _XOPEN_SOURCE_EXTENDED
#   define __USE_XOPEN_EXTENDED	1
#  endif
# endif
#endif

#ifdef _LARGEFILE_SOURCE
# define __USE_LARGEFILE	1
#endif

#ifdef _LARGEFILE64_SOURCE
# define __USE_LARGEFILE64	1
#endif

#if defined _FILE_OFFSET_BITS && _FILE_OFFSET_BITS == 64
# define __USE_FILE_OFFSET64	1
#endif

#if defined _BSD_SOURCE || defined _SVID_SOURCE
# define __USE_MISC	1
#endif

#ifdef	_BSD_SOURCE
# define __USE_BSD	1
#endif

#ifdef	_SVID_SOURCE
# define __USE_SVID	1
#endif

#ifdef	_ATFILE_SOURCE
# define __USE_ATFILE	1
#endif

#ifdef	_GNU_SOURCE
# define __USE_GNU	1
#endif

#if defined _REENTRANT || defined _THREAD_SAFE
# define __USE_REENTRANT	1
#endif

#if defined _FORTIFY_SOURCE && _FORTIFY_SOURCE > 0 \
    && __GNUC_PREREQ (4, 1) && defined __OPTIMIZE__ && __OPTIMIZE__ > 0
# if _FORTIFY_SOURCE > 1
#  define __USE_FORTIFY_LEVEL 2
# else
#  define __USE_FORTIFY_LEVEL 1
# endif
#else
# define __USE_FORTIFY_LEVEL 0
#endif

/* Get definitions of __STDC_* predefined macros, if the compiler has
   not preincluded this header automatically.  */
/* stdc-predef.h */

#ifndef	_STDC_PREDEF_H
#define	_STDC_PREDEF_H	1

#ifdef __GCC_IEC_559
# if __GCC_IEC_559 > 0
#  define __STDC_IEC_559__		1
# endif
#else
# define __STDC_IEC_559__		1
#endif

#ifdef __GCC_IEC_559_COMPLEX
# if __GCC_IEC_559_COMPLEX > 0
#  define __STDC_IEC_559_COMPLEX__	1
# endif
#else
# define __STDC_IEC_559_COMPLEX__	1
#endif

/* wchar_t uses ISO/IEC 10646 (2nd ed., published 2011-03-15) /
   Unicode 6.0.  */
#define __STDC_ISO_10646__		201103L

/* We do not support C11 <threads.h>.  */
#define __STDC_NO_THREADS__		1

#endif

/* stdc-predef.h */

/* This macro indicates that the installed library is the GNU C Library.
   For historic reasons the value now is 6 and this will stay from now
   on.  The use of this variable is deprecated.  Use __GLIBC__ and
   __GLIBC_MINOR__ now (see below) when you want to test for a specific
   GNU C library version and use the values in <gnu/lib-names.h> to get
   the sonames of the shared libraries.  */
#undef  __GNU_LIBRARY__
#define __GNU_LIBRARY__ 6

/* Major and minor version number of the GNU C library package.  Use
   these macros to test for features in specific releases.  */
#define	__GLIBC__	2
#define	__GLIBC_MINOR__	19

#define __GLIBC_PREREQ(maj, min) \
	((__GLIBC__ << 16) + __GLIBC_MINOR__ >= ((maj) << 16) + (min))

/* This is here only because every header file already includes this one.  */
#ifndef __ASSEMBLER__
# ifndef _SYS_CDEFS_H
/* cdefs.h */
#ifndef	_SYS_CDEFS_H
#define	_SYS_CDEFS_H	1

/* The GNU libc does not support any K&R compilers or the traditional mode
   of ISO C compilers anymore.  Check for some of the combinations not
   anymore supported.  */
#if defined __GNUC__ && !defined __STDC__
# error "You need a ISO C conforming compiler to use the glibc headers"
#endif

/* Some user header file might have defined this before.  */
#undef	__P
#undef	__PMT

#ifdef __GNUC__

/* All functions, except those with callbacks or those that
   synchronize memory, are leaf functions.  */
# if __GNUC_PREREQ (4, 6) && !defined _LIBC
#  define __LEAF , __leaf__
#  define __LEAF_ATTR __attribute__ ((__leaf__))
# else
#  define __LEAF
#  define __LEAF_ATTR
# endif

/* GCC can always grok prototypes.  For C++ programs we add throw()
   to help it optimize the function calls.  But this works only with
   gcc 2.8.x and egcs.  For gcc 3.2 and up we even mark C functions
   as non-throwing using a function attribute since programs can use
   the -fexceptions options for C code as well.  */
# if !defined __cplusplus && __GNUC_PREREQ (3, 3)
#  define __THROW	__attribute__ ((__nothrow__ __LEAF))
#  define __THROWNL	__attribute__ ((__nothrow__))
#  define __NTH(fct)	__attribute__ ((__nothrow__ __LEAF)) fct
# else
#  if defined __cplusplus && __GNUC_PREREQ (2,8)
#   define __THROW	throw ()
#   define __THROWNL	throw ()
#   define __NTH(fct)	__LEAF_ATTR fct throw ()
#  else
#   define __THROW
#   define __THROWNL
#   define __NTH(fct)	fct
#  endif
# endif

#else	/* Not GCC.  */

# define __inline		/* No inline functions.  */

# define __THROW
# define __THROWNL
# define __NTH(fct)	fct

#endif	/* GCC.  */

/* These two macros are not used in glibc anymore.  They are kept here
   only because some other projects expect the macros to be defined.  */
#define __P(args)	args
#define __PMT(args)	args

/* For these things, GCC behaves the ANSI way normally,
   and the non-ANSI way under -traditional.  */

#define __CONCAT(x,y)	x ## y
#define __STRING(x)	#x

/* This is not a typedef so 'const __ptr_t' does the right thing.  */
#define __ptr_t void *
#define __long_double_t  long double


/* C++ needs to know that types and declarations are C, not C++.  */
#ifdef	__cplusplus
# define __BEGIN_DECLS	extern "C" {
# define __END_DECLS	}
#else
# define __BEGIN_DECLS
# define __END_DECLS
#endif


/* The standard library needs the functions from the ISO C90 standard
   in the std namespace.  At the same time we want to be safe for
   future changes and we include the ISO C99 code in the non-standard
   namespace __c99.  The C++ wrapper header take case of adding the
   definitions to the global namespace.  */
#if defined __cplusplus && defined _GLIBCPP_USE_NAMESPACES
# define __BEGIN_NAMESPACE_STD	namespace std {
# define __END_NAMESPACE_STD	}
# define __USING_NAMESPACE_STD(name) using std::name;
# define __BEGIN_NAMESPACE_C99	namespace __c99 {
# define __END_NAMESPACE_C99	}
# define __USING_NAMESPACE_C99(name) using __c99::name;
#else
/* For compatibility we do not add the declarations into any
   namespace.  They will end up in the global namespace which is what
   old code expects.  */
# define __BEGIN_NAMESPACE_STD
# define __END_NAMESPACE_STD
# define __USING_NAMESPACE_STD(name)
# define __BEGIN_NAMESPACE_C99
# define __END_NAMESPACE_C99
# define __USING_NAMESPACE_C99(name)
#endif


/* Fortify support.  */
#define __bos(ptr) __builtin_object_size (ptr, __USE_FORTIFY_LEVEL > 1)
#define __bos0(ptr) __builtin_object_size (ptr, 0)
#define __fortify_function __extern_always_inline __attribute_artificial__

#if __GNUC_PREREQ (4,3)
# define __warndecl(name, msg) \
  extern void name (void) __attribute__((__warning__ (msg)))
# define __warnattr(msg) __attribute__((__warning__ (msg)))
# define __errordecl(name, msg) \
  extern void name (void) __attribute__((__error__ (msg)))
#else
# define __warndecl(name, msg) extern void name (void)
# define __warnattr(msg)
# define __errordecl(name, msg) extern void name (void)
#endif

/* Support for flexible arrays.  */
#if __GNUC_PREREQ (2,97)
/* GCC 2.97 supports C99 flexible array members.  */
# define __flexarr	[]
#else
# ifdef __GNUC__
#  define __flexarr	[0]
# else
#  if defined __STDC_VERSION__ && __STDC_VERSION__ >= 199901L
#   define __flexarr	[]
#  else
/* Some other non-C99 compiler.  Approximate with [1].  */
#   define __flexarr	[1]
#  endif
# endif
#endif


/* __asm__ ("xyz") is used throughout the headers to rename functions
   at the assembly language level.  This is wrapped by the __REDIRECT
   macro, in order to support compilers that can do this some other
   way.  When compilers don't support asm-names at all, we have to do
   preprocessor tricks instead (which don't have exactly the right
   semantics, but it's the best we can do).

   Example:
   int __REDIRECT(setpgrp, (__pid_t pid, __pid_t pgrp), setpgid); */

#if defined __GNUC__ && __GNUC__ >= 2

# define __REDIRECT(name, proto, alias) name proto __asm__ (__ASMNAME (#alias))
# ifdef __cplusplus
#  define __REDIRECT_NTH(name, proto, alias) \
     name proto __THROW __asm__ (__ASMNAME (#alias))
#  define __REDIRECT_NTHNL(name, proto, alias) \
     name proto __THROWNL __asm__ (__ASMNAME (#alias))
# else
#  define __REDIRECT_NTH(name, proto, alias) \
     name proto __asm__ (__ASMNAME (#alias)) __THROW
#  define __REDIRECT_NTHNL(name, proto, alias) \
     name proto __asm__ (__ASMNAME (#alias)) __THROWNL
# endif
# define __ASMNAME(cname)  __ASMNAME2 (__USER_LABEL_PREFIX__, cname)
# define __ASMNAME2(prefix, cname) __STRING (prefix) cname

/*
#elif __SOME_OTHER_COMPILER__

# define __REDIRECT(name, proto, alias) name proto; \
	_Pragma("let " #name " = " #alias)
*/
#endif

/* GCC has various useful declarations that can be made with the
   '__attribute__' syntax.  All of the ways we use this do fine if
   they are omitted for compilers that don't understand it. */
#if !defined __GNUC__ || __GNUC__ < 2
# define __attribute__(xyz)	/* Ignore */
#endif

/* At some point during the gcc 2.96 development the 'malloc' attribute
   for functions was introduced.  We don't want to use it unconditionally
   (although this would be possible) since it generates warnings.  */
#if __GNUC_PREREQ (2,96)
# define __attribute_malloc__ __attribute__ ((__malloc__))
#else
# define __attribute_malloc__ /* Ignore */
#endif

/* Tell the compiler which arguments to an allocation function
   indicate the size of the allocation.  */
#if __GNUC_PREREQ (4, 3)
# define __attribute_alloc_size__(params) \
  __attribute__ ((__alloc_size__ params))
#else
# define __attribute_alloc_size__(params) /* Ignore.  */
#endif

/* At some point during the gcc 2.96 development the 'pure' attribute
   for functions was introduced.  We don't want to use it unconditionally
   (although this would be possible) since it generates warnings.  */
#if __GNUC_PREREQ (2,96)
# define __attribute_pure__ __attribute__ ((__pure__))
#else
# define __attribute_pure__ /* Ignore */
#endif

/* This declaration tells the compiler that the value is constant.  */
#if __GNUC_PREREQ (2,5)
# define __attribute_const__ __attribute__ ((__const__))
#else
# define __attribute_const__ /* Ignore */
#endif

/* At some point during the gcc 3.1 development the 'used' attribute
   for functions was introduced.  We don't want to use it unconditionally
   (although this would be possible) since it generates warnings.  */
#if __GNUC_PREREQ (3,1)
# define __attribute_used__ __attribute__ ((__used__))
# define __attribute_noinline__ __attribute__ ((__noinline__))
#else
# define __attribute_used__ __attribute__ ((__unused__))
# define __attribute_noinline__ /* Ignore */
#endif

/* gcc allows marking deprecated functions.  */
#if __GNUC_PREREQ (3,2)
# define __attribute_deprecated__ __attribute__ ((__deprecated__))
#else
# define __attribute_deprecated__ /* Ignore */
#endif

/* At some point during the gcc 2.8 development the 'format_arg' attribute
   for functions was introduced.  We don't want to use it unconditionally
   (although this would be possible) since it generates warnings.
   If several 'format_arg' attributes are given for the same function, in
   gcc-3.0 and older, all but the last one are ignored.  In newer gccs,
   all designated arguments are considered.  */
#if __GNUC_PREREQ (2,8)
# define __attribute_format_arg__(x) __attribute__ ((__format_arg__ (x)))
#else
# define __attribute_format_arg__(x) /* Ignore */
#endif

/* At some point during the gcc 2.97 development the 'strfmon' format
   attribute for functions was introduced.  We don't want to use it
   unconditionally (although this would be possible) since it
   generates warnings.  */
#if __GNUC_PREREQ (2,97)
# define __attribute_format_strfmon__(a,b) \
  __attribute__ ((__format__ (__strfmon__, a, b)))
#else
# define __attribute_format_strfmon__(a,b) /* Ignore */
#endif

/* The nonull function attribute allows to mark pointer parameters which
   must not be NULL.  */
#if __GNUC_PREREQ (3,3)
# define __nonnull(params) __attribute__ ((__nonnull__ params))
#else
# define __nonnull(params)
#endif

/* If fortification mode, we warn about unused results of certain
   function calls which can lead to problems.  */
#if __GNUC_PREREQ (3,4)
# define __attribute_warn_unused_result__ \
   __attribute__ ((__warn_unused_result__))
# if __USE_FORTIFY_LEVEL > 0
#  define __wur __attribute_warn_unused_result__
# endif
#else
# define __attribute_warn_unused_result__ /* empty */
#endif
#ifndef __wur
# define __wur /* Ignore */
#endif

/* Forces a function to be always inlined.  */
#if __GNUC_PREREQ (3,2)
# define __always_inline __inline __attribute__ ((__always_inline__))
#else
# define __always_inline __inline
#endif

/* Associate error messages with the source location of the call site rather
   than with the source location inside the function.  */
#if __GNUC_PREREQ (4,3)
# define __attribute_artificial__ __attribute__ ((__artificial__))
#else
# define __attribute_artificial__ /* Ignore */
#endif

#ifdef __GNUC__
/* One of these will be defined if the __gnu_inline__ attribute is
   available.  In C++, __GNUC_GNU_INLINE__ will be defined even though
   __inline does not use the GNU inlining rules.  If neither macro is
   defined, this version of GCC only supports GNU inline semantics. */
# if defined __GNUC_STDC_INLINE__ || defined __GNUC_GNU_INLINE__
#  define __extern_inline extern __inline __attribute__ ((__gnu_inline__))
#  define __extern_always_inline \
  extern __always_inline __attribute__ ((__gnu_inline__))
# else
#  define __extern_inline extern __inline
#  define __extern_always_inline extern __always_inline
# endif
#else /* Not GCC.  */
# define __extern_inline  /* Ignore */
# define __extern_always_inline /* Ignore */
#endif

/* GCC 4.3 and above allow passing all anonymous arguments of an
   __extern_always_inline function to some other vararg function.  */
#if __GNUC_PREREQ (4,3)
# define __va_arg_pack() __builtin_va_arg_pack ()
# define __va_arg_pack_len() __builtin_va_arg_pack_len ()
#endif

/* It is possible to compile containing GCC extensions even if GCC is
   run in pedantic mode if the uses are carefully marked using the
   '__extension__' keyword.  But this is not generally available before
   version 2.8.  */
#if !__GNUC_PREREQ (2,8)
# define __extension__		/* Ignore */
#endif

/* __restrict is known in EGCS 1.2 and above. */
#if !__GNUC_PREREQ (2,92)
# define __restrict	/* Ignore */
#endif

/* ISO C99 also allows to declare arrays as non-overlapping.  The syntax is
     array_name[restrict]
   GCC 3.1 supports this.  */
#if __GNUC_PREREQ (3,1) && !defined __GNUG__
# define __restrict_arr	__restrict
#else
# ifdef __GNUC__
#  define __restrict_arr	/* Not supported in old GCC.  */
# else
#  if defined __STDC_VERSION__ && __STDC_VERSION__ >= 199901L
#   define __restrict_arr	restrict
#  else
/* Some other non-C99 compiler.  */
#   define __restrict_arr	/* Not supported.  */
#  endif
# endif
#endif

#if __GNUC__ >= 3
# define __glibc_unlikely(cond)	__builtin_expect ((cond), 0)
# define __glibc_likely(cond)	__builtin_expect ((cond), 1)
#else
# define __glibc_unlikely(cond)	(cond)
# define __glibc_likely(cond)	(cond)
#endif

/* bits/wordsize.h */
/* Determine the wordsize from the preprocessor defines.  */

#if !defined __x86_64__
#	define __x86_64__
#endif
#if defined __ILP32__
#	undef __ILP32__
#endif
#if defined __x86_64__ && !defined __ILP32__
# define __WORDSIZE	64
#else
# define __WORDSIZE	32
#endif

#ifdef __x86_64__
# define __WORDSIZE_TIME64_COMPAT32	1
/* Both x86-64 and x32 use the 64-bit system call interface.  */
# define __SYSCALL_WORDSIZE		64
#endif /* bits/wordsize.h */



#if defined __LONG_DOUBLE_MATH_OPTIONAL && defined __NO_LONG_DOUBLE_MATH
# define __LDBL_COMPAT 1
# ifdef __REDIRECT
#  define __LDBL_REDIR1(name, proto, alias) __REDIRECT (name, proto, alias)
#  define __LDBL_REDIR(name, proto) \
  __LDBL_REDIR1 (name, proto, __nldbl_##name)
#  define __LDBL_REDIR1_NTH(name, proto, alias) __REDIRECT_NTH (name, proto, alias)
#  define __LDBL_REDIR_NTH(name, proto) \
  __LDBL_REDIR1_NTH (name, proto, __nldbl_##name)
#  define __LDBL_REDIR1_DECL(name, alias) \
  extern __typeof (name) name __asm (__ASMNAME (#alias));
#  define __LDBL_REDIR_DECL(name) \
  extern __typeof (name) name __asm (__ASMNAME ("__nldbl_" #name));
#  define __REDIRECT_LDBL(name, proto, alias) \
  __LDBL_REDIR1 (name, proto, __nldbl_##alias)
#  define __REDIRECT_NTH_LDBL(name, proto, alias) \
  __LDBL_REDIR1_NTH (name, proto, __nldbl_##alias)
# endif
#endif
#if !defined __LDBL_COMPAT || !defined __REDIRECT
# define __LDBL_REDIR1(name, proto, alias) name proto
# define __LDBL_REDIR(name, proto) name proto
# define __LDBL_REDIR1_NTH(name, proto, alias) name proto __THROW
# define __LDBL_REDIR_NTH(name, proto) name proto __THROW
# define __LDBL_REDIR_DECL(name)
# ifdef __REDIRECT
#  define __REDIRECT_LDBL(name, proto, alias) __REDIRECT (name, proto, alias)
#  define __REDIRECT_NTH_LDBL(name, proto, alias) \
  __REDIRECT_NTH (name, proto, alias)
# endif
#endif

#endif	 /* sys/cdefs.h */
# endif

/* If we don't have __REDIRECT, prototypes will be missing if
   __USE_FILE_OFFSET64 but not __USE_LARGEFILE[64]. */
# if defined __USE_FILE_OFFSET64 && !defined __REDIRECT
#  define __USE_LARGEFILE	1
#  define __USE_LARGEFILE64	1
# endif

#endif	/* !ASSEMBLER */

/* Decide whether we can define 'extern inline' functions in headers.  */
#if __GNUC_PREREQ (2, 7) && defined __OPTIMIZE__ \
    && !defined __OPTIMIZE_SIZE__ && !defined __NO_INLINE__ \
    && defined __extern_inline
# define __USE_EXTERN_INLINES	1
#endif


/* This is here only because every header file already includes this one.
   Get the definitions of all the appropriate '__stub_FUNCTION' symbols.
   <gnu/stubs.h> contains '#define __stub_FUNCTION' when FUNCTION is a stub
   that will always return failure (and set errno to ENOSYS).  */
/* stubs.h */

#if !defined __x86_64__
/* stubs-32.h */

#ifdef _LIBC
# error Applications may not define the macro _LIBC
#endif

#define __stub_chflags
#define __stub_fattach
#define __stub_fchflags
#define __stub_fdetach
#define __stub_gtty
#define __stub_lchmod
#define __stub_revoke
#define __stub_setlogin
#define __stub_sigreturn
#define __stub_sstk
#define __stub_stty

/* stubs-32.h */
#endif /* !defined __x86_64__ */
#if defined __x86_64__ && defined __LP64__
/* stubs-64.h */

#ifdef _LIBC
# error Applications may not define the macro _LIBC
#endif

#define __stub_bdflush
#define __stub_chflags
#define __stub_fattach
#define __stub_fchflags
#define __stub_fdetach
#define __stub_getmsg
#define __stub_gtty
#define __stub_lchmod
#define __stub_putmsg
#define __stub_revoke
#define __stub_setlogin
#define __stub_sigreturn
#define __stub_sstk
#define __stub_stty

/* stubs-64.h */
#endif /* defined __x86_64__ && defined __LP64__ */
#if defined __x86_64__ && defined __ILP32__
/* stubs-x32.h */

#ifdef _LIBC
# error Applications may not define the macro _LIBC
#endif

#define __stub_bdflush
#define __stub_chflags
#define __stub_create_module
#define __stub_fattach
#define __stub_fchflags
#define __stub_fdetach
#define __stub_get_kernel_syms
#define __stub_getmsg
#define __stub_gtty
#define __stub_lchmod
#define __stub_nfsservctl
#define __stub_putmsg
#define __stub_query_module
#define __stub_revoke
#define __stub_setlogin
#define __stub_sigreturn
#define __stub_sstk
#define __stub_stty
#define __stub_uselib

/* stubs-x32.h */
#endif /* defined __x86_64__ && defined __ILP32__ */

/* stubs.h */


#endif	/* features.h  */

/* bits/wchar.h */
#ifndef _BITS_WCHAR_H
#define _BITS_WCHAR_H	1

#ifdef __WCHAR_MAX__
# define __WCHAR_MAX	__WCHAR_MAX__
#elif L'\0' - 1 > 0
# define __WCHAR_MAX	(0xffffffffu + L'\0')
#else
# define __WCHAR_MAX	(0x7fffffff + L'\0')
#endif

#ifdef __WCHAR_MIN__
# define __WCHAR_MIN	__WCHAR_MIN__
#elif L'\0' - 1 > 0
# define __WCHAR_MIN	(L'\0' + 0)
#else
# define __WCHAR_MIN	(-__WCHAR_MAX - 1)
#endif

#endif	/* bits/wchar.h */

/* Exact integral types.  */

/* Signed.  */

/* There is some amount of overlap with <sys/types.h> as known by inet code */
#ifndef __int8_t_defined
# define __int8_t_defined
typedef signed char		int8_t;
typedef short int		int16_t;
typedef int			int32_t;
# if __WORDSIZE == 64
typedef long int		int64_t;
# else
__extension__
typedef long long int		int64_t;
# endif
#endif

/* Unsigned.  */
typedef unsigned char		uint8_t;
typedef unsigned short int	uint16_t;
#ifndef __uint32_t_defined
typedef unsigned int		uint32_t;
# define __uint32_t_defined
#endif
#if __WORDSIZE == 64
typedef unsigned long int	uint64_t;
#else
__extension__
typedef unsigned long long int	uint64_t;
#endif


/* Small types.  */

/* Signed.  */
typedef signed char		int_least8_t;
typedef short int		int_least16_t;
typedef int			int_least32_t;
#if __WORDSIZE == 64
typedef long int		int_least64_t;
#else
__extension__
typedef long long int		int_least64_t;
#endif

/* Unsigned.  */
typedef unsigned char		uint_least8_t;
typedef unsigned short int	uint_least16_t;
typedef unsigned int		uint_least32_t;
#if __WORDSIZE == 64
typedef unsigned long int	uint_least64_t;
#else
__extension__
typedef unsigned long long int	uint_least64_t;
#endif


/* Fast types.  */

/* Signed.  */
typedef signed char		int_fast8_t;
#if __WORDSIZE == 64
typedef long int		int_fast16_t;
typedef long int		int_fast32_t;
typedef long int		int_fast64_t;
#else
typedef int			int_fast16_t;
typedef int			int_fast32_t;
__extension__
typedef long long int		int_fast64_t;
#endif

/* Unsigned.  */
typedef unsigned char		uint_fast8_t;
#if __WORDSIZE == 64
typedef unsigned long int	uint_fast16_t;
typedef unsigned long int	uint_fast32_t;
typedef unsigned long int	uint_fast64_t;
#else
typedef unsigned int		uint_fast16_t;
typedef unsigned int		uint_fast32_t;
__extension__
typedef unsigned long long int	uint_fast64_t;
#endif


/* Types for 'void *' pointers.  */
/* #if __WORDSIZE == 64
# ifndef __intptr_t_defined
typedef long int		intptr_t;
#  define __intptr_t_defined
# endif
typedef unsigned long int	uintptr_t;
#else
# ifndef __intptr_t_defined
typedef int			intptr_t;
#  define __intptr_t_defined
# endif
typedef unsigned int		uintptr_t;
#endif
*/

/* Largest integral types.  */
#if __WORDSIZE == 64
typedef long int		intmax_t;
typedef unsigned long int	uintmax_t;
#else
__extension__
typedef long long int		intmax_t;
__extension__
typedef unsigned long long int	uintmax_t;
#endif


# if __WORDSIZE == 64
#  define __INT64_C(c)	c ## L
#  define __UINT64_C(c)	c ## UL
# else
#  define __INT64_C(c)	c ## LL
#  define __UINT64_C(c)	c ## ULL
# endif

/* Limits of integral types.  */

/* Minimum of signed integral types.  */
# define INT8_MIN		(-128)
# define INT16_MIN		(-32767-1)
# define INT32_MIN		(-2147483647-1)
# define INT64_MIN		(-__INT64_C(9223372036854775807)-1)
/* Maximum of signed integral types.  */
# define INT8_MAX		(127)
# define INT16_MAX		(32767)
# define INT32_MAX		(2147483647)
# define INT64_MAX		(__INT64_C(9223372036854775807))

/* Maximum of unsigned integral types.  */
# define UINT8_MAX		(255)
# define UINT16_MAX		(65535)
# define UINT32_MAX		(4294967295U)
# define UINT64_MAX		(__UINT64_C(18446744073709551615))


/* Minimum of signed integral types having a minimum size.  */
# define INT_LEAST8_MIN		(-128)
# define INT_LEAST16_MIN	(-32767-1)
# define INT_LEAST32_MIN	(-2147483647-1)
# define INT_LEAST64_MIN	(-__INT64_C(9223372036854775807)-1)
/* Maximum of signed integral types having a minimum size.  */
# define INT_LEAST8_MAX		(127)
# define INT_LEAST16_MAX	(32767)
# define INT_LEAST32_MAX	(2147483647)
# define INT_LEAST64_MAX	(__INT64_C(9223372036854775807))

/* Maximum of unsigned integral types having a minimum size.  */
# define UINT_LEAST8_MAX	(255)
# define UINT_LEAST16_MAX	(65535)
# define UINT_LEAST32_MAX	(4294967295U)
# define UINT_LEAST64_MAX	(__UINT64_C(18446744073709551615))


/* Minimum of fast signed integral types having a minimum size.  */
# define INT_FAST8_MIN		(-128)
# if __WORDSIZE == 64
#  define INT_FAST16_MIN	(-9223372036854775807L-1)
#  define INT_FAST32_MIN	(-9223372036854775807L-1)
# else
#  define INT_FAST16_MIN	(-2147483647-1)
#  define INT_FAST32_MIN	(-2147483647-1)
# endif
# define INT_FAST64_MIN		(-__INT64_C(9223372036854775807)-1)
/* Maximum of fast signed integral types having a minimum size.  */
# define INT_FAST8_MAX		(127)
# if __WORDSIZE == 64
#  define INT_FAST16_MAX	(9223372036854775807L)
#  define INT_FAST32_MAX	(9223372036854775807L)
# else
#  define INT_FAST16_MAX	(2147483647)
#  define INT_FAST32_MAX	(2147483647)
# endif
# define INT_FAST64_MAX		(__INT64_C(9223372036854775807))

/* Maximum of fast unsigned integral types having a minimum size.  */
# define UINT_FAST8_MAX		(255)
# if __WORDSIZE == 64
#  define UINT_FAST16_MAX	(18446744073709551615UL)
#  define UINT_FAST32_MAX	(18446744073709551615UL)
# else
#  define UINT_FAST16_MAX	(4294967295U)
#  define UINT_FAST32_MAX	(4294967295U)
# endif
# define UINT_FAST64_MAX	(__UINT64_C(18446744073709551615))


/* Values to test for integral types holding 'void *' pointer.  */
# if __WORDSIZE == 64
#  define INTPTR_MIN		(-9223372036854775807L-1)
#  define INTPTR_MAX		(9223372036854775807L)
#  define UINTPTR_MAX		(18446744073709551615UL)
# else
#  define INTPTR_MIN		(-2147483647-1)
#  define INTPTR_MAX		(2147483647)
#  define UINTPTR_MAX		(4294967295U)
# endif


/* Minimum for largest signed integral type.  */
# define INTMAX_MIN		(-__INT64_C(9223372036854775807)-1)
/* Maximum for largest signed integral type.  */
# define INTMAX_MAX		(__INT64_C(9223372036854775807))

/* Maximum for largest unsigned integral type.  */
# define UINTMAX_MAX		(__UINT64_C(18446744073709551615))


/* Limits of other integer types.  */

/* Limits of 'ptrdiff_t' type.  */
# if __WORDSIZE == 64
#  define PTRDIFF_MIN		(-9223372036854775807L-1)
#  define PTRDIFF_MAX		(9223372036854775807L)
# else
#  define PTRDIFF_MIN		(-2147483647-1)
#  define PTRDIFF_MAX		(2147483647)
# endif

/* Limits of 'sig_atomic_t'.  */
# define SIG_ATOMIC_MIN		(-2147483647-1)
# define SIG_ATOMIC_MAX		(2147483647)

/* Limit of 'size_t' type.  */
# if __WORDSIZE == 64
#  define SIZE_MAX		(18446744073709551615UL)
# else
#  define SIZE_MAX		(4294967295U)
# endif

/* Limits of 'wchar_t'.  */
# ifndef WCHAR_MIN
/* These constants might also be defined in <wchar.h>.  */
#  define WCHAR_MIN		__WCHAR_MIN
#  define WCHAR_MAX		__WCHAR_MAX
# endif

/* Limits of 'wint_t'.  */
# define WINT_MIN		(0u)
# define WINT_MAX		(4294967295u)

/* Signed.  */
# define INT8_C(c)	c
# define INT16_C(c)	c
# define INT32_C(c)	c
# if __WORDSIZE == 64
#  define INT64_C(c)	c ## L
# else
#  define INT64_C(c)	c ## LL
# endif

/* Unsigned.  */
# define UINT8_C(c)	c
# define UINT16_C(c)	c
# define UINT32_C(c)	c ## U
# if __WORDSIZE == 64
#  define UINT64_C(c)	c ## UL
# else
#  define UINT64_C(c)	c ## ULL
# endif

/* Maximal type.  */
# if __WORDSIZE == 64
#  define INTMAX_C(c)	c ## L
#  define UINTMAX_C(c)	c ## UL
# else
#  define INTMAX_C(c)	c ## LL
#  define UINTMAX_C(c)	c ## ULL
# endif

#endif /* stdint.h */
#ifndef _STENCIL_H_
#define _STENCIL_H_

// 3D array indexing
#define index(ix,iy,iz,Nx,Ny,Nz) ( ( (iz)*(Ny) + (iy) ) * (Nx) + (ix) )
#define idx(ix,iy,iz) ( index((ix),(iy),(iz),(Nx),(Ny),(Nz)) )


// modulo used for PBC wrap around
#define MOD(n, M) ( (( (n) % (M) ) + (M) ) % (M)  )

// have PBC in x, y or z?
#define PBCx (PBC & 1)
#define PBCy (PBC & 2)
#define PBCz (PBC & 4)

// clamp or wrap index at boundary, depending on PBC
// hclamp*: clamps on upper side (index+1)
// lclamp*: clamps on lower side (index-1)
// *clampx: clamps along x
// ...
#define hclampx(ix) (PBCx? MOD(ix, Nx) : min((ix), Nx-1))
#define lclampx(ix) (PBCx? MOD(ix, Nx) : max((ix), 0))

#define hclampy(iy) (PBCy? MOD(iy, Ny) : min((iy), Ny-1))
#define lclampy(iy) (PBCy? MOD(iy, Ny) : max((iy), 0))

#define hclampz(iz) (PBCz? MOD(iz, Nz) : min((iz), Nz-1))
#define lclampz(iz) (PBCz? MOD(iz, Nz) : max((iz), 0))


#endif // _STENCIL_H_
// This file implements common functions on float3 (vector).
// Author: Mykola Dvornik, Arne Vansteenkiste

#ifndef _FLOAT3_H_
#define _FLOAT3_H_

// converting set of 3 floats into a 3-component vector
static inline real_t3 make_float3(real_t a, real_t b, real_t c) {
	return (real_t3) {a, b, c};
}

// length of the 3-components vector
static inline real_t len(real_t3 a) {
	return length(a);
}

// returns a normalized copy of the 3-components vector
static inline real_t3 normalized(real_t3 a){
	real_t veclen = (len(a) != 0.0f) ? ( 1.0f / len(a) ) : 0.0f;
	return veclen * a;
}

// square
static inline real_t pow2(real_t x){
	return x * x;
}


// pow(x, 3)
static inline real_t pow3(real_t x){
	return x * x * x;
}


// pow(x, 4)
static inline real_t pow4(real_t x){
	float s = x*x;
	return s*s;
}

#define is0(m) ( dot(m, m) == 0.0f )

#endif // _FLOAT3_H_
#ifndef _EXCHANGE_H_
#define _EXCHANGE_H_

// indexing in symmetric matrix
#define symidx(i, j) ( (j<=i)? ( (((i)*((i)+1)) /2 )+(j) )  :  ( (((j)*((j)+1)) /2 )+(i) ) )

#endif // _EXCHANGE_H_
#ifndef _ATOMICF_H_
#define _ATOMICF_H_

// Atomic max of abs value.
static inline void atomicFmaxabs(volatile __global real_t* a, real_t b){
    b = fabs(b);
    atomic_max((__global int*)(a), *((int*)(&b)));
}

#endif // _ATOMICF_H_
#ifndef _REDUCE_H_
#define _REDUCE_H_

#if defined(__REAL_IS_DOUBLE__)
    #pragma OPENCL EXTENSION cl_khr_int64_base_atomics : enable
    #pragma OPENCL EXTENSION cl_khr_int64_extended_atomics : enable
#else
    #pragma OPENCL EXTENSION cl_khr_global_int32_base_atomics : enable
    #pragma OPENCL EXTENSION cl_khr_global_int32_extended_atomics : enable
#endif // __REAL_IS_DOUBLE__

#define __REDUCE_REG_COUNT__ 16

#if defined(__REAL_IS_DOUBLE__)

static inline void atomicAdd_r(volatile __global double* addr, double val) {
    union{
        unsigned long u64;
        double f64;
    } next, expected, current;
    current.f64 = *addr;
    do {
        next.f64 = (expected.f64 = current.f64) + val;
        current.u64 = atom_cmpxchg( (volatile __global unsigned long *) addr,
            expected.u64, next.u64);
    } while (current.u64 != expected.u64);
}

static inline void atomicMax_r(volatile __global double* addr, double val) {
    atom_max( (volatile __global unsigned long *) addr, as_ulong(val));
}

#else

static inline void atomicAdd_r(volatile __global float* addr, float val) {
    union{
        unsigned int u32;
        float f32;
    } next, expected, current;
    current.f32 = *addr;
    do {
        next.f32 = (expected.f32 = current.f32) + val;
        current.u32 = atomic_cmpxchg( (volatile __global unsigned int *) addr,
            expected.u32, next.u32);
    } while (current.u32 != expected.u32);
}

static inline void atomicMax_r(volatile __global float* addr, float val) {
    atomic_max( (volatile __global unsigned int *) addr, as_uint(val));
}

#endif // __REAL_IS_DOUBLE__

#endif // _REDUCE_H_
#ifndef _AMUL_H_
#define _AMUL_H_

// Returns mul * arr[i], or mul when arr == NULL;
static inline real_t amul(__global real_t *arr, float mul, int i) {
	return (arr == NULL)? (mul) : (mul * arr[i]);
}

// Returns m * a[i], or m when a == NULL;
static inline real_t3 vmul(__global real_t *ax, __global real_t *ay, __global real_t *az,
                             real_t  mx,          real_t  my,          real_t  mz, int i) {
    return make_float3(amul(ax, mx, i),
                       amul(ay, my, i),
                       amul(az, mz, i));
}

// Returns 1/Msat, or 0 when Msat == 0.
static inline real_t inv_Msat(__global real_t* Ms_, real_t Ms_mul, int i) {
    real_t ms = amul(Ms_, Ms_mul, i);
    if (ms == 0.0f) {
        return 0.0f;
    } else {
        return 1.0f / ms;
    }
}
#endif // _AMUL_H_
#ifndef __RNG_COMMON_H__
#define __RNG_COMMON_H__
// Taken from PhD thesis of Thomas Luu (Department of Mathematics
// at University College of London)
static inline float normcdfinv_(float u) {
	float	v;
	float	p;
	float	q;
	float	ushift;
	float   tmp;

    if ((u < 0.0f) || (u > 1.0f)) {
        return FLT_MIN;
    } else if ((u == 0.0f) || (u == 1.0f)) {
        return 0.0f;
    } else {
        tmp = u;
    }

    ushift = tmp - 0.5f;

    v = copysign(ushift, 0.0f);

    if (v < 0.499433f) {
        v = rsqrt((-tmp*tmp) + tmp);
        v *= 0.5f;
        p = 0.001732781974270904f;
        p = p * v + 0.1788417306083325f;
        p = p * v + 2.804338363421083f;
        p = p * v + 9.35716893191325f;
		p = p * v + 5.283080058166861f;
		p = p * v + 0.07885390444279965f;
		p *= ushift;
		q = 0.0001796248328874524f;
		q = q * v + 0.02398533988976253f;
		q = q * v + 0.4893072798067982f;
		q = q * v + 2.406460595830034f;
		q = q * v + 3.142947488363618f;
    } else {
        if (ushift > 0.0f) {
            tmp = 1.0f - tmp;
        }
        v = log2(tmp+tmp);
        v *= -0.6931471805599453f;
        if (v < 22.0f) {
            p = 0.000382438382914666f;
            p = p * v + 0.03679041341785685f;
            p = p * v + 0.5242351532484291f;
            p = p * v + 1.21642047402659f;
            q = 9.14019972725528e-6f;
            q = q * v + 0.003523083799369908f;
            q = q * v + 0.126802543865968f;
            q = q * v + 0.8502031783957995f;
        } else {
            p = 0.00001016962895771568f;
            p = p * v + 0.003330096951634844f;
            p = p * v + 0.1540146885433827f;
            p = p * v + 1.045480394868638f;
            q = 1.303450553973082e-7f;
            q = q * v + 0.0001728926914526662f;
            q = q * v + 0.02031866871146244f;
            q = q * v + 0.3977137974626933f;
        }
        p *= copysign(v, ushift);
    }
    q = q * v + 1.0f;
    v = 1.0f / q;
    return p * v;
}

// auxiliary function to convert a pair of uint32 to a single-
// precision float in (0, 1)
static inline float uint2float(uint a, uint b) {
    uint num1 = a;
    uint num2 = b;
    uint finalNum = 0;
    uint expo = 32;
    for (;expo > 0; expo--) {
        uint flag0 = num1 & 0x80000000;
        num1 <<= 1;
        if (flag0 != 0) {
            break;
        }
    }
    uint maskbits = 0x007fffff;
    finalNum ^= (num2 & maskbits);
    uint newExpo = 94 + expo;
    finalNum ^= (newExpo << 23);
    return as_float(finalNum); // return value
}

#if defined(__REAL_IS_DOUBLE__)
// auxiliary function to convert a pair of uint64 to a double-
// precision float in (0, 1)
static inline double ulong2double(ulong a, ulong b) {
    ulong num1 = a;
    ulong num2 = b;
    ulong finalNum = 0;
    ulong expo = 64;
    for (;expo > 0; expo--) {
        ulong flag0 = num1 & 0x8000000000000000;
        num1 <<= 1;
        if (flag0 != 0) {
            break;
        }
    }
    ulong maskbits = 0x000fffffffffffff;
    finalNum ^= (num2 & maskbits);
    ulong newExpo = 958 + expo;
    finalNum ^= (newExpo << 52);
    return as_double(finalNum); // return value
}
#endif // __REAL_IS_DOUBLE__

static inline void boxMuller(float* in, float* out, uint offset) {
    uint u1idx = 2*offset;
    uint u2idx = u1idx+1;
    out[u1idx] = sqrt( -2.0f * log(in[u1idx]) ) * cospi(2.0f * in[u2idx]);
    out[u2idx] = sqrt( -2.0f * log(in[u1idx]) ) * sinpi(2.0f * in[u2idx]);
}

#if defined(__REAL_IS_DOUBLE__)
static inline void boxMuller64(double* in, double* out, uint offset) {
    uint u1idx = 2*offset;
    uint u2idx = u1idx+1;
    out[u1idx] = sqrt( (real_t)(-2.0) * log(in[u1idx]) ) * cospi((real_t)(2.0) * in[u2idx]);
    out[u2idx] = sqrt( (real_t)(-2.0) * log(in[u1idx]) ) * sinpi((real_t)(2.0) * in[u2idx]);
}
#endif // __REAL_IS_DOUBLE__

#if defined(__REAL_IS_DOUBLE__)
#define XORWOW_NORM_double 2.328306549295727688e-10
#endif // __REAL_IS_DOUBLE__

#endif // __RNG_COMMON_H__
#ifndef __RNGTHREEFRY_H__
#define __RNGTHREEFRY_H__
#define THREEFRY_ELEMENTS_PER_BLOCK 256
#define SKEIN_KS_PARITY64 0x1BD11BDAA9FC1A22
#define SKEIN_KS_PARITY32 0x1BD11BDA
constant int THREEFRY2X32_ROTATION[] = {13, 15, 26,  6, 17, 29, 16, 24};
constant int THREEFRY2X64_ROTATION[] = {16, 42, 12, 31, 16, 32, 24, 21};

constant int THREEFRY4X32_ROTATION_0[] = {10, 26};
constant int THREEFRY4X32_ROTATION_1[] = {11, 21};
constant int THREEFRY4X32_ROTATION_2[] = {13, 27};
constant int THREEFRY4X32_ROTATION_3[] = {23,  5};
constant int THREEFRY4X32_ROTATION_4[] = { 6, 20};
constant int THREEFRY4X32_ROTATION_5[] = {17, 11};
constant int THREEFRY4X32_ROTATION_6[] = {25, 10};
constant int THREEFRY4X32_ROTATION_7[] = {18, 20};

constant int THREEFRY4X64_ROTATION_0[] = {14, 16};
constant int THREEFRY4X64_ROTATION_1[] = {52, 57};
constant int THREEFRY4X64_ROTATION_2[] = {23, 40};
constant int THREEFRY4X64_ROTATION_3[] = { 5, 37};
constant int THREEFRY4X64_ROTATION_4[] = {25, 33};
constant int THREEFRY4X64_ROTATION_5[] = {46, 12};
constant int THREEFRY4X64_ROTATION_6[] = {58, 22};
constant int THREEFRY4X64_ROTATION_7[] = {32, 32};


/**
State of threefry RNG.
*/
typedef struct{
    uint counter[4];
    uint result[4];
    uint key[4];
    uint tracker;
} threefry_state;

typedef struct{
    ulong counter[4];
    ulong result[4];
    ulong key[4];
    ulong tracker;
} threefry64_state;

static inline ulong RotL64(ulong x, uint N){
    return (x << (N & 63)) | (x >> ((64 - N) & 63));
}

static inline ulong RotL32(uint x, uint N){
    return (x << (N & 31)) | (x >> ((32 - N) & 31));
}

static inline void threefry_round(threefry_state* state){
    uint ks[5]; //
    ks[4] = SKEIN_KS_PARITY32;

    // Unrolled for loop
    ks[0] = state->key[0];
    state->result[0]  = state->counter[0];
    ks[4] ^= state->key[0];
    ks[1] = state->key[1];
    state->result[1]  = state->counter[1];
    ks[4] ^= state->key[1];
    ks[2] = state->key[2];
    state->result[2]  = state->counter[2];
    ks[4] ^= state->key[2];
    ks[3] = state->key[3];
    state->result[3]  = state->counter[3];
    ks[4] ^= state->key[3];

    /* Insert initial key before round 0 */
    state->result[0] += ks[0];
    state->result[1] += ks[1];
    state->result[2] += ks[2];
    state->result[3] += ks[3];

    /* First round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* Second round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* Third round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* Fourth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=1) */
    state->result[0] += ks[1];
    state->result[1] += ks[2];
    state->result[2] += ks[3];
    state->result[3] += ks[4];
    state->result[3] += 1; /* X[4-1] += r */

    /* Fifth round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_4[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_4[1]);
    state->result[3] ^= state->result[2];

    /* Sixth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_5[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_5[1]);
    state->result[1] ^= state->result[2];

    /* Seventh round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_6[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_6[1]);
    state->result[3] ^= state->result[2];

    /* Eighth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_7[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_7[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=2) */
    state->result[0] += ks[2];
    state->result[1] += ks[3];
    state->result[2] += ks[4];
    state->result[3] += ks[0];
    state->result[3] += 2; /* X[4-1] += r */

    /* 9-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* 10-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* 11-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* 12-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=3) */
    state->result[0] += ks[3];
    state->result[1] += ks[4];
    state->result[2] += ks[0];
    state->result[3] += ks[1];
    state->result[3] += 3; /* X[4-1] += r */

    /* 13-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_4[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_4[1]);
    state->result[3] ^= state->result[2];

    /* 14-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_5[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_5[1]);
    state->result[1] ^= state->result[2];

    /* 15-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_6[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_6[1]);
    state->result[3] ^= state->result[2];

    /* 16-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_7[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_7[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=4) */
    state->result[0] += ks[4];
    state->result[1] += ks[0];
    state->result[2] += ks[1];
    state->result[3] += ks[2];
    state->result[3] += 4; /* X[4-1] += r */

    /* 17-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* 18-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* 19-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* 20-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL32(state->result[3], THREEFRY4X32_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL32(state->result[1], THREEFRY4X32_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=5) */
    state->result[0] += ks[0];
    state->result[1] += ks[1];
    state->result[2] += ks[2];
    state->result[3] += ks[3];
    state->result[3] += 5; /* X[4-1] += r */

}

static inline void threefry64_round(threefry64_state* state){
    ulong ks[5]; //
    ks[4] = SKEIN_KS_PARITY64;

    // Unrolled for loop
    ks[0] = state->key[0];
    state->result[0]  = state->counter[0];
    ks[4] ^= state->key[0];
    ks[1] = state->key[1];
    state->result[1]  = state->counter[1];
    ks[4] ^= state->key[1];
    ks[2] = state->key[2];
    state->result[2]  = state->counter[2];
    ks[4] ^= state->key[2];
    ks[3] = state->key[3];
    state->result[3]  = state->counter[3];
    ks[4] ^= state->key[3];

    /* Insert initial key before round 0 */
    state->result[0] += ks[0];
    state->result[1] += ks[1];
    state->result[2] += ks[2];
    state->result[3] += ks[3];

    /* First round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* Second round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* Third round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* Fourth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=1) */
    state->result[0] += ks[1];
    state->result[1] += ks[2];
    state->result[2] += ks[3];
    state->result[3] += ks[4];
    state->result[3] += 1; /* X[4-1] += r */

    /* Fifth round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_4[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_4[1]);
    state->result[3] ^= state->result[2];

    /* Sixth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_5[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_5[1]);
    state->result[1] ^= state->result[2];

    /* Seventh round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_6[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_6[1]);
    state->result[3] ^= state->result[2];

    /* Eighth round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_7[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_7[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=2) */
    state->result[0] += ks[2];
    state->result[1] += ks[3];
    state->result[2] += ks[4];
    state->result[3] += ks[0];
    state->result[3] += 2; /* X[4-1] += r */

    /* 9-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* 10-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* 11-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* 12-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=3) */
    state->result[0] += ks[3];
    state->result[1] += ks[4];
    state->result[2] += ks[0];
    state->result[3] += ks[1];
    state->result[3] += 3; /* X[4-1] += r */

    /* 13-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_4[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_4[1]);
    state->result[3] ^= state->result[2];

    /* 14-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_5[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_5[1]);
    state->result[1] ^= state->result[2];

    /* 15-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_6[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_6[1]);
    state->result[3] ^= state->result[2];

    /* 16-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_7[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_7[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=4) */
    state->result[0] += ks[4];
    state->result[1] += ks[0];
    state->result[2] += ks[1];
    state->result[3] += ks[2];
    state->result[3] += 4; /* X[4-1] += r */

    /* 17-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_0[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_0[1]);
    state->result[3] ^= state->result[2];

    /* 18-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_1[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_1[1]);
    state->result[1] ^= state->result[2];

    /* 19-th round */
    state->result[0] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_2[0]);
    state->result[1] ^= state->result[0];
    state->result[2] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_2[1]);
    state->result[3] ^= state->result[2];

    /* 20-th round */
    state->result[0] += state->result[3];
    state->result[3] = RotL64(state->result[3], THREEFRY4X64_ROTATION_3[0]);
    state->result[3] ^= state->result[0];
    state->result[2] += state->result[1];
    state->result[1] = RotL64(state->result[1], THREEFRY4X64_ROTATION_3[1]);
    state->result[1] ^= state->result[2];

    /* InjectKey(r=5) */
    state->result[0] += ks[0];
    state->result[1] += ks[1];
    state->result[2] += ks[2];
    state->result[3] += ks[3];
    state->result[3] += 5; /* X[4-1] += r */

}
#endif // __RNGTHREEFRY_H__
#ifndef __RNGXORWOW_H__
#define __RNGXORWOW_H__
/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

#define XORWOW_FLOAT_MULTI 2.3283064e-10f
#define XORWOW_DOUBLE2_MULTI 2.328306549295727688e-10
#define XORWOW_DOUBLE_MULTI 5.4210108624275221700372640e-20

// defines from rocRAND for skipping XORWOW
#define XORWOW_N 5
#define XORWOW_M 32
#define XORWOW_SIZE (XORWOW_M * XORWOW_N * XORWOW_N)
#define XORWOW_JUMP_MATRICES 32
#define XORWOW_JUMP_LOG2 2

static inline void copy_vec(uint* dst, const uint* src) {
    for (int i = 0; i < XORWOW_N; i++) {
        dst[i] = src[i];
    }
}

static inline void mul_mat_vec_inplace(__global uint* m, uint* v) {
    uint r[XORWOW_N] = { 0 };
    for (int ij = 0; ij < XORWOW_N * XORWOW_M; ij++) {
        const int i = ij / XORWOW_M;
        const int j = ij % XORWOW_M;
        const uint b = (v[i] & (1 << j)) ? 0xffffffff : 0x0;
        for (int k = 0; k < XORWOW_N; k++) {
            r[k] ^= b & m[i * XORWOW_M * XORWOW_N + j * XORWOW_N + k];
        }
    }
    copy_vec(v, r);
}

static inline void xorwow_jump(ulong v, __global uint* jump_matrices, uint* xorwow_state) {
    ulong vi = v;
    uint mi = 0;
    while (vi > 0) {
        const uint is = (uint)(vi) & ((1 << XORWOW_JUMP_LOG2) - 1);
        for (uint i = 0; i < is; i++) {
            mul_mat_vec_inplace(&jump_matrices[mi*XORWOW_SIZE], xorwow_state);
        }
        mi++;
        vi >>= XORWOW_JUMP_LOG2;
    }
}

static inline void xorwow_discard(ulong offset, uint* xorwow_state, __global uint* h_xorwow_jump_matrices) {
    xorwow_jump(offset, h_xorwow_jump_matrices, xorwow_state);
}

static inline void xorwow_discard_subsequence(ulong subsequence, uint* xorwow_state, __global uint* h_xorwow_sequence_jump_matrices) {
    xorwow_jump(subsequence, h_xorwow_sequence_jump_matrices, xorwow_state);
}
#endif // __RNGXORWOW_H__
#ifndef _SUM_H_
#define _SUM_H_

static inline real_t sum(real_t a, real_t b){
	return a + b;
}

#endif // _SUM_H_
// Copy src (size S, smaller) into dst (size D, larger),
// and multiply by Bsat * vol
__kernel void
copypadmul2(__global real_t* __restrict dst,    int     Dx, int Dy, int Dz,
            __global real_t* __restrict src,    int     Sx, int Sy, int Sz,
            __global real_t* __restrict Ms_, real_t Ms_mul,
            __global real_t* __restrict vol) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix<Sx) && (iy<Sy) && (iz<Sz)) {
        int        sI = index(ix, iy, iz, Sx, Sy, Sz);  // source index
        real_t tmpFac = amul(Ms_, Ms_mul, sI);
        real_t   Bsat = MU0 * tmpFac;
        real_t      v = amul(vol, (real_t)1.0, sI);

        dst[index(ix, iy, iz, Dx, Dy, Dz)] = Bsat * v * src[sI];
    }
}
// Copy src (size S, larger) to dst (size D, smaller)
__kernel void
copyunpad(__global real_t* __restrict dst, int Dx, int Dy, int Dz,
          __global real_t* __restrict src, int Sx, int Sy, int Sz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix<Dx) && (iy<Dy) && (iz<Dz)) {
        dst[index(ix, iy, iz, Dx, Dy, Dz)] = src[index(ix, iy, iz, Sx, Sy, Sz)];
    }
}
// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
__kernel void
crop(__global real_t* __restrict  dst, int   Dx, int   Dy, int Dz,
     __global real_t* __restrict  src, int   Sx, int   Sy, int Sz,
                             int Offx, int Offy, int Offz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix<Dx) && (iy<Dy) && (iz<Dz)) {
        dst[index(ix, iy, iz, Dx, Dy, Dz)] = src[index(ix+Offx, iy+Offy, iz+Offz, Sx, Sy, Sz)];
    }
}
// add cubic anisotropy field to B.
// B:      effective field in T
// m:      reduced magnetization (unit length)
// Ms:     saturation magnetization in A/m.
// K1:     Kc1 in J/m3
// K2:     Kc2 in T/m3
// C1, C2: anisotropy axes
//
// based on http://www.southampton.ac.uk/~fangohr/software/oxs_cubic8.html
__kernel void
addcubicanisotropy2(__global real_t* __restrict   Bx, __global real_t* __restrict      By, __global real_t* __restrict Bz,
                    __global real_t* __restrict   mx, __global real_t* __restrict      my, __global real_t* __restrict mz,
                    __global real_t* __restrict  Ms_,                      real_t  Ms_mul,
                    __global real_t* __restrict  k1_,                      real_t  k1_mul,
                    __global real_t* __restrict  k2_,                      real_t  k2_mul,
                    __global real_t* __restrict  k3_,                      real_t  k3_mul,
                    __global real_t* __restrict c1x_,                      real_t c1x_mul,
                    __global real_t* __restrict c1y_,                      real_t c1y_mul,
                    __global real_t* __restrict c1z_,                      real_t c1z_mul,
                    __global real_t* __restrict c2x_,                      real_t c2x_mul,
                    __global real_t* __restrict c2y_,                      real_t c2y_mul,
                    __global real_t* __restrict c2z_,                      real_t c2z_mul,
                                            int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t    k1 = amul(k1_, k1_mul, i);
        real_t    k2 = amul(k2_, k2_mul, i);
        real_t    k3 = amul(k3_, k3_mul, i);

        k1 *= invMs;
        k2 *= invMs;
        k3 *= invMs;

        real_t  u1x = (c1x_ == NULL) ? c1x_mul : (c1x_mul * c1x_[i]);
        real_t3  u1 = normalized(vmul(c1x_, c1y_, c1z_, c1x_mul, c1y_mul, c1z_mul, i));
        real_t3  u2 = normalized(vmul(c2x_, c2y_, c2z_, c2x_mul, c2y_mul, c2z_mul, i));
        real_t3  u3 = cross(u1, u2); // 3rd axis perpendicular to u1,u2
        real_t3   m = make_float3(mx[i], my[i], mz[i]);

        real_t u1m = dot(u1, m);
        real_t u2m = dot(u2, m);
        real_t u3m = dot(u3, m);

        real_t3 B = (real_t)-2.0*k1*((pow2(u2m) + pow2(u3m)) * (    (u1m) * u1) +
                                     (pow2(u1m) + pow2(u3m)) * (    (u2m) * u2) +
                                     (pow2(u1m) + pow2(u2m)) * (    (u3m) * u3))-
                    (real_t)2.0f*k2*((pow2(u2m) * pow2(u3m)) * (    (u1m) * u1) +
                                     (pow2(u1m) * pow2(u3m)) * (    (u2m) * u2) +
                                     (pow2(u1m) * pow2(u2m)) * (    (u3m) * u3))-
                    (real_t)4.0f*k3*((pow4(u2m) + pow4(u3m)) * (pow3(u1m) * u1) +
                                     (pow4(u1m) + pow4(u3m)) * (pow3(u2m) * u2) +
                                     (pow4(u1m) + pow4(u2m)) * (pow3(u3m) * u3));

        Bx[i] += B.x;
        By[i] += B.y;
        Bz[i] += B.z;
    }
}
// dst[i] = a[i] / b[i]
__kernel void
pointwise_div(__global real_t* __restrict dst, __global real_t* __restrict a, __global real_t* __restrict b, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (b[i] == (real_t)0.0) ? (real_t)0.0 : (a[i] / b[i]);
    }
}
// dst[i] = a[i] / b[i]
__kernel void
divide(__global real_t* __restrict dst, __global real_t* __restrict a, __global real_t* __restrict b, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        dst[i] = ((a[i] == 0) || (b[i] == 0)) ? (real_t)0.0 : (a[i] / b[i]);
    }
}
// Exchange + Dzyaloshinskii-Moriya interaction according to
// Bagdanov and Röβler, PRL 87, 3, 2001. eq.8 (out-of-plane symmetry breaking).
// Taking into account proper boundary conditions.
// m: normalized magnetization
// H: effective field in Tesla
// D: dmi strength / Msat, in Tesla*m
// A: Aex/Msat
__kernel void
adddmi(__global real_t* __restrict     Hx, __global real_t* __restrict     Hy, __global  real_t* __restrict      Hz,
       __global real_t* __restrict     mx, __global real_t* __restrict     my, __global  real_t* __restrict      mz,
       __global real_t* __restrict    Ms_,                      real_t Ms_mul,
       __global real_t* __restrict aLUT2d, __global real_t* __restrict dLUT2d, __global uint8_t* __restrict regions,
                            real_t     cx,                      real_t     cy,                        real_t      cz,
                               int     Nx,                         int     Ny,                           int      Nz,
                           uint8_t    PBC,                     uint8_t OpenBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int      I = idx(ix, iy, iz);                  // central cell index
    real_t3  h = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);       // add to H
    real_t3 m0 = make_float3(mx[I], my[I], mz[I]); // central m
    uint8_t r0 = regions[I];
    int i_;                                        // neighbor index

    if(is0(m0)) {
        return;
    }

    // x derivatives (along length)
    {
        real_t3 m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // left neighbor
        i_ = idx(lclampx(ix-1), iy, iz);               // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-1 >= 0) || PBCx) {
            m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int    r1 = is0(m1)? r0 : regions[i_];              // don't use inter region params if m1=0
        real_t A1 = aLUT2d[symidx(r0, r1)];                 // inter-region Aex
        real_t D1 = dLUT2d[symidx(r0, r1)];                 // inter-region Dex
        if ((!is0(m1)) || (!OpenBC)) {                     // do nothing at an open boundary
            if (is0(m1)) {                                 // neighbor missing
                m1.x = m0.x - (-cx * (0.5f*D1/A1) * m0.z); // extrapolate missing m from Neumann BC's
                m1.y = m0.y;
                m1.z = m0.z + (-cx * (0.5f*D1/A1) * m0.x);
            }
            h   += (2.0f*A1/(cx*cx)) * (m1 - m0);          // exchange
            h.x += (D1/cx)*(- m1.z);
            h.z -= (D1/cx)*(- m1.x);
        }
    }

    {
        real_t3 m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);         // right neighbor
        i_ = idx(hclampx(ix+1), iy, iz);
        if ((ix+1 < Nx) || PBCx) {
            m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int    r2 = is0(m2)? r0 : regions[i_];
        real_t A2 = aLUT2d[symidx(r0, r2)];
        real_t D2 = dLUT2d[symidx(r0, r2)];
        if ((!is0(m2)) || (!OpenBC)) {
            if (is0(m2)) {
                m2.x = m0.x - (cx * (0.5f*D2/A2) * m0.z);
                m2.y = m0.y;
                m2.z = m0.z + (cx * (0.5f*D2/A2) * m0.x);
            }
            h   += (2.0f*A2/(cx*cx)) * (m2 - m0);
            h.x += (D2/cx)*(m2.z);
            h.z -= (D2/cx)*(m2.x);
        }
    }

    // y derivatives (along height)
    {
        real_t3 m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-1), iz);
        if ((iy-1 >= 0) || PBCy) {
            m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int    r1 = is0(m1)? r0 : regions[i_];
        real_t A1 = aLUT2d[symidx(r0, r1)];
        real_t D1 = dLUT2d[symidx(r0, r1)];
        if ((!is0(m1)) || (!OpenBC)) {
            if (is0(m1)) {
                m1.x = m0.x;
                m1.y = m0.y - (-cy * (0.5f*D1/A1) * m0.z);
                m1.z = m0.z + (-cy * (0.5f*D1/A1) * m0.y);
            }
            h   += (2.0f*A1/(cy*cy)) * (m1 - m0);
            h.y += (D1/cy)*(- m1.z);
            h.z -= (D1/cy)*(- m1.y);
        }
    }

    {
        real_t3 m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+1), iz);
        if  (iy+1 < Ny || PBCy) {
            m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int    r2 = is0(m2)? r0 : regions[i_];
        real_t A2 = aLUT2d[symidx(r0, r2)];
        real_t D2 = dLUT2d[symidx(r0, r2)];
        if ((!is0(m2)) || (!OpenBC)) {
            if (is0(m2)) {
                m2.x = m0.x;
                m2.y = m0.y - (cy * (0.5f*D2/A2) * m0.z);
                m2.z = m0.z + (cy * (0.5f*D2/A2) * m0.y);
            }
            h   += (2.0f*A2/(cy*cy)) * (m2 - m0);
            h.y += (D2/cy)*(m2.z);
            h.z -= (D2/cy)*(m2.y);
        }
    }

    // only take vertical derivative for 3D sim
    if (Nz != 1) {
        // bottom neighbor
        {
                    i_ = idx(ix, iy, lclampz(iz-1));
            real_t3 m1 = make_float3(mx[i_], my[i_], mz[i_]);
                    m1 = ( is0(m1)? m0: m1 );                   // Neumann BC
             real_t A1 = aLUT2d[symidx(r0, regions[i_])];
                    h += (2.0f*A1/(cz*cz)) * (m1 - m0);         // Exchange only
        }

        // top neighbor
        {
                    i_ = idx(ix, iy, hclampz(iz+1));
            real_t3 m2 = make_float3(mx[i_], my[i_], mz[i_]);
                    m2 = (is0(m2)) ? m0: m2;
            real_t  A2 = aLUT2d[symidx(r0, regions[i_])];
                    h += (2.0f*A2/(cz*cz)) * (m2 - m0);
        }
    }

    // write back, result is H + Hdmi + Hex
    real_t invMs = inv_Msat(Ms_, Ms_mul, I);

    Hx[I] += h.x*invMs;
    Hy[I] += h.y*invMs;
    Hz[I] += h.z*invMs;
}

// Note on boundary conditions.
//
// We need the derivative and laplacian of m in point A, but e.g. C lies out of the boundaries.
// We use the boundary condition in B (derivative of the magnetization) to extrapolate m to point C:
//     m_C = m_A + (dm/dx)|_B * cellsize
//
// When point C is inside the boundary, we just use its actual value.
//
// Then we can take the central derivative in A:
//     (dm/dx)|_A = (m_C - m_D) / (2*cellsize)
// And the laplacian:
//     lapl(m)|_A = (m_C + m_D - 2*m_A) / (cellsize^2)
//
// All these operations should be second order as they involve only central derivatives.
//
//    ------------------------------------------------------------------ *
//   |                                                   |             C |
//   |                                                   |          **   |
//   |                                                   |        ***    |
//   |                                                   |     ***       |
//   |                                                   |   ***         |
//   |                                                   | ***           |
//   |                                                   B               |
//   |                                               *** |               |
//   |                                            ***    |               |
//   |                                         ****      |               |
//   |                                     ****          |               |
//   |                                  ****             |               |
//   |                              ** A                 |               |
//   |                         *****                     |               |
//   |                   ******                          |               |
//   |          *********                                |               |
//   |D ********                                         |               |
//   |                                                   |               |
//   +----------------+----------------+-----------------+---------------+
//  -1              -0.5               0               0.5               1
//                                 x
// Exchange + Dzyaloshinskii-Moriya interaction for bulk material.
// Energy:
//
//     E  = D M . rot(M)
//
// Effective field:
//
//     Hx = 2A/Bs nabla²Mx + 2D/Bs dzMy - 2D/Bs dyMz
//     Hy = 2A/Bs nabla²My + 2D/Bs dxMz - 2D/Bs dzMx
//     Hz = 2A/Bs nabla²Mz + 2D/Bs dyMx - 2D/Bs dxMy
//
// Boundary conditions:
//
//             2A dxMx = 0
//      D Mz + 2A dxMy = 0
//     -D My + 2A dxMz = 0
//
//     -D Mz + 2A dyMx = 0
//             2A dyMy = 0
//      D Mx + 2A dyMz = 0
//
//      D My + 2A dzMx = 0
//     -D Mx + 2A dzMy = 0
//             2A dzMz = 0
//
__kernel void
adddmibulk(__global  real_t* __restrict      Hx, __global real_t* __restrict     Hy, __global real_t* __restrict Hz,
           __global  real_t* __restrict      mx, __global real_t* __restrict     my, __global real_t* __restrict mz,
           __global  real_t* __restrict     Ms_,                      real_t Ms_mul,
           __global  real_t* __restrict  aLUT2d, __global real_t* __restrict DLUT2d,
           __global uint8_t* __restrict regions,
                                 real_t      cx,                      real_t     cy,                      real_t cz,
                                    int      Nx,                         int     Ny,                         int Nz,
                                uint8_t     PBC,                     uint8_t OpenBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int      I = idx(ix, iy, iz);                   // central cell index
    real_t3  h = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);       // add to H
    real_t3 m0 = make_float3(mx[I], my[I], mz[I]); // central m
    uint8_t r0 = regions[I];
    int i_;                                         // neighbor index

    if(is0(m0)) {
        return;
    }

    // x derivatives (along length)
    {
        real_t3 m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // left neighbor
                i_ = idx(lclampx(ix-1), iy, iz);        // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-1 >= 0) || PBCx) {
            m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int      r1 = is0(m1)? r0 : regions[i_];
        real_t    A = aLUT2d[symidx(r0, r1)];
        real_t    D = DLUT2d[symidx(r0, r1)];
        real_t D_2A = D/(2.0f*A);
        if ((!is0(m1)) || (!OpenBC)) {                 // do nothing at an open boundary
            if (is0(m1)) {                             // neighbor missing
                m1.x = m0.x;
                m1.y = m0.y - (-cx * D_2A * m0.z);
                m1.z = m0.z + (-cx * D_2A * m0.y);
            }
            h   += (2.0f*A/(cx*cx)) * (m1 - m0);       // exchange
            h.y += (D/cx)*(-m1.z);
            h.z -= (D/cx)*(-m1.y);
        }
    }


    {
        real_t3 m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);   // right neighbor
        i_ = idx(hclampx(ix+1), iy, iz);
        if ((ix+1 < Nx) || PBCx) {
            m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int      r1 = is0(m2) ? r0 : regions[i_];
        real_t    A = aLUT2d[symidx(r0, r1)];
        real_t    D = DLUT2d[symidx(r0, r1)];
        real_t D_2A = D/(2.0f*A);
        if ((!is0(m2)) || (!OpenBC)) {
            if (is0(m2)) {
                m2.x = m0.x;
                m2.y = m0.y - (+cx * D_2A * m0.z);
                m2.z = m0.z + (+cx * D_2A * m0.y);
            }
            h   += (2.0f*A/(cx*cx)) * (m2 - m0);
            h.y += (D/cx)*(m2.z);
            h.z -= (D/cx)*(m2.y);
        }
    }

    // y derivatives (along height)
    {
        real_t3 m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-1), iz);
        if ((iy-1 >= 0) || PBCy) {
            m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int      r1 = is0(m1) ? r0 : regions[i_];
        real_t    A = aLUT2d[symidx(r0, r1)];
        real_t    D = DLUT2d[symidx(r0, r1)];
        real_t D_2A = D/(2.0f*A);
        if ((!is0(m1)) || (!OpenBC)) {
            if (is0(m1)) {
                m1.x = m0.x + (-cy * D_2A * m0.z);
                m1.y = m0.y;
                m1.z = m0.z - (-cy * D_2A * m0.x);
            }
            h   += (2.0f*A/(cy*cy)) * (m1 - m0);
            h.x -= (D/cy)*(-m1.z);
            h.z += (D/cy)*(-m1.x);
        }
    }

    {
        real_t3 m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+1), iz);
        if  ((iy+1 < Ny) || PBCy) {
            m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }
        int      r1 = is0(m2) ? r0 : regions[i_];
        real_t    A = aLUT2d[symidx(r0, r1)];
        real_t    D = DLUT2d[symidx(r0, r1)];
        real_t D_2A = D/(2.0f*A);
        if ((!is0(m2)) || (!OpenBC)) {
            if (is0(m2)) {
                m2.x = m0.x + (+cy * D_2A * m0.z);
                m2.y = m0.y;
                m2.z = m0.z - (+cy * D_2A * m0.x);
            }
            h   += (2.0f*A/(cy*cy)) * (m2 - m0);
            h.x -= (D/cy)*(m2.z);
            h.z += (D/cy)*(m2.x);
        }
    }

    // only take vertical derivative for 3D sim
    if (Nz != 1) {
        // bottom neighbor
        {
            real_t3 m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
            i_ = idx(ix, iy, lclampz(iz-1));
            if ((iz-1 >= 0) || PBCz) {
                m1 = make_float3(mx[i_], my[i_], mz[i_]);
            }
            int      r1 = is0(m1) ? r0 : regions[i_];
            real_t    A = aLUT2d[symidx(r0, r1)];
            real_t    D = DLUT2d[symidx(r0, r1)];
            real_t D_2A = D/(2.0f*A);
            if ((!is0(m1)) || (!OpenBC)) {
                if (is0(m1)) {
                    m1.x = m0.x - (-cz * D_2A * m0.y);
                    m1.y = m0.y + (-cz * D_2A * m0.x);
                    m1.z = m0.z;
                }
                h   += (2.0f*A/(cz*cz)) * (m1 - m0);
                h.x += (D/cz)*(- m1.y);
                h.y -= (D/cz)*(- m1.x);
            }
        }

        // top neighbor
        {
            real_t3 m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
            i_ = idx(ix, iy, hclampz(iz+1));
            if ((iz+1 < Nz) || PBCz) {
                m2 = make_float3(mx[i_], my[i_], mz[i_]);
            }
            int      r1 = is0(m2) ? r0 : regions[i_];
            real_t    A = aLUT2d[symidx(r0, r1)];
            real_t    D = DLUT2d[symidx(r0, r1)];
            real_t D_2A = D/(2.0f*A);
            if ((!is0(m2)) || (!OpenBC)) {
                if (is0(m2)) {
                    m2.x = m0.x - (+cz * D_2A * m0.y);
                    m2.y = m0.y + (+cz * D_2A * m0.x);
                    m2.z = m0.z;
                }
                h   += (2.0f*A/(cz*cz)) * (m2 - m0);
                h.x += (D/cz)*(m2.y );
                h.y -= (D/cz)*(m2.x );
            }
        }
    }

    // write back, result is H + Hdmi + Hex
    real_t invMs = inv_Msat(Ms_, Ms_mul, I);

    Hx[I] += h.x*invMs;
    Hy[I] += h.y*invMs;
    Hz[I] += h.z*invMs;
}

// Note on boundary conditions.
//
// We need the derivative and laplacian of m in point A, but e.g. C lies out of the boundaries.
// We use the boundary condition in B (derivative of the magnetization) to extrapolate m to point C:
//     m_C = m_A + (dm/dx)|_B * cellsize
//
// When point C is inside the boundary, we just use its actual value.
//
// Then we can take the central derivative in A:
//     (dm/dx)|_A = (m_C - m_D) / (2*cellsize)
// And the laplacian:
//     lapl(m)|_A = (m_C + m_D - 2*m_A) / (cellsize^2)
//
// All these operations should be second order as they involve only central derivatives.
//
//    ------------------------------------------------------------------ *
//   |                                                   |             C |
//   |                                                   |          **   |
//   |                                                   |        ***    |
//   |                                                   |     ***       |
//   |                                                   |   ***         |
//   |                                                   | ***           |
//   |                                                   B               |
//   |                                               *** |               |
//   |                                            ***    |               |
//   |                                         ****      |               |
//   |                                     ****          |               |
//   |                                  ****             |               |
//   |                              ** A                 |               |
//   |                         *****                     |               |
//   |                   ******                          |               |
//   |          *********                                |               |
//   |D ********                                         |               |
//   |                                                   |               |
//   +----------------+----------------+-----------------+---------------+
//  -1              -0.5               0               0.5               1
//                                 x
// dst += prefactor * dot(a,b)
__kernel void
dotproduct(__global real_t* __restrict dst,                      real_t prefactor,
           __global real_t* __restrict  ax, __global real_t* __restrict        ay, __global real_t* __restrict az,
           __global real_t* __restrict  bx, __global real_t* __restrict        by, __global real_t* __restrict bz,
                                   int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t3 A = {ax[i], ay[i], az[i]};
        real_t3 B = {bx[i], by[i], bz[i]};

        dst[i] += prefactor * dot(A, B);
    }
}
__kernel void
crossproduct(__global real_t* __restrict dstx, __global real_t* __restrict dsty, __global real_t* __restrict dstz,
             __global real_t* __restrict   ax, __global real_t* __restrict   ay, __global real_t* __restrict   az,
             __global real_t* __restrict   bx, __global real_t* __restrict   by, __global real_t* __restrict   bz,
                                     int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        real_t3   A = {ax[i], ay[i], az[i]};
        real_t3   B = {bx[i], by[i], bz[i]};
        real_t3 AxB = cross(A, B);

        dstx[i] = AxB.x;
        dsty[i] = AxB.y;
        dstz[i] = AxB.z;
    }
}
// Add exchange field to Beff.
//     m: normalized magnetization
//     B: effective field in Tesla
//     Aex_red: Aex / (Msat * 1e18 m2)
__kernel void
addexchange(__global real_t* __restrict     Bx, __global  real_t* __restrict      By, __global real_t* __restrict Bz,
            __global real_t* __restrict     mx, __global  real_t* __restrict      my, __global real_t* __restrict mz,
            __global real_t* __restrict    Ms_,                       real_t  Ms_mul,
            __global real_t* __restrict aLUT2d, __global uint8_t* __restrict regions,
                                 real_t     wx,                       real_t      wy,                      real_t wz,
                                    int     Nx,                          int      Ny,                         int Nz,
                                uint8_t    PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int      I = idx(ix, iy, iz);
    real_t3 m0 = make_float3(mx[I], my[I], mz[I]);

    if (is0(m0)) {
        return;
    }

    uint8_t r0 = regions[I];
    real_t3  B = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);

    int     i_; // neighbor index
    real_t3 m_; // neighbor mag
    real_t  a__; // inter-cell exchange stiffness

    // left neighbor
    i_  = idx(lclampx(ix-1), iy, iz);           // clamps or wraps index according to PBC
    m_  = make_float3(mx[i_], my[i_], mz[i_]);  // load m
    m_  = ( is0(m_)? m0: m_ );                  // replace missing non-boundary neighbor
    a__ = aLUT2d[symidx(r0, regions[i_])];
    B  += wx * a__ *(m_ - m0);

    // right neighbor
    i_  = idx(hclampx(ix+1), iy, iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    B  += wx * a__ *(m_ - m0);

    // back neighbor
    i_  = idx(ix, lclampy(iy-1), iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    B  += wy * a__ *(m_ - m0);

    // front neighbor
    i_  = idx(ix, hclampy(iy+1), iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    B  += wy * a__ *(m_ - m0);

    // only take vertical derivative for 3D sim
    if (Nz != 1) {
        // bottom neighbor
        i_  = idx(ix, iy, lclampz(iz-1));
        m_  = make_float3(mx[i_], my[i_], mz[i_]);
        m_  = ( is0(m_)? m0: m_ );
        a__ = aLUT2d[symidx(r0, regions[i_])];
        B  += wz * a__ *(m_ - m0);

        // top neighbor
        i_  = idx(ix, iy, hclampz(iz+1));
        m_  = make_float3(mx[i_], my[i_], mz[i_]);
        m_  = ( is0(m_)? m0: m_ );
        a__ = aLUT2d[symidx(r0, regions[i_])];
        B  += wz * a__ *(m_ - m0);
    }

    real_t invMs = inv_Msat(Ms_, Ms_mul, I);

    Bx[I] += B.x*invMs;
    By[I] += B.y*invMs;
    Bz[I] += B.z*invMs;
}
// Finds the average exchange strength around each cell, for debugging.
__kernel void
exchangedecode(__global real_t* __restrict dst, __global real_t* __restrict aLUT2d, __global uint8_t* __restrict regions,
                                    real_t  wx,                      real_t     wy,                       real_t      wz,
                                       int  Nx,                         int     Ny,                          int      Nz,
                                   uint8_t PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int      I = idx(ix, iy, iz);
    uint8_t r0 = regions[I];

    int i_;    // neighbor index
    real_t avg = (real_t)0.0;

    // left neighbor
    i_   = idx(lclampx(ix-1), iy, iz);           // clamps or wraps index according to PBC
    avg += aLUT2d[symidx(r0, regions[i_])];

    // right neighbor
    i_   = idx(hclampx(ix+1), iy, iz);
    avg += aLUT2d[symidx(r0, regions[i_])];

    // back neighbor
    i_   = idx(ix, lclampy(iy-1), iz);
    avg += aLUT2d[symidx(r0, regions[i_])];

    // front neighbor
    i_   = idx(ix, hclampy(iy+1), iz);
    avg += aLUT2d[symidx(r0, regions[i_])];

    // only take vertical derivative for 3D sim
    if (Nz != 1) {
        // bottom neighbor
        i_   = idx(ix, iy, lclampz(iz-1));
        avg += aLUT2d[symidx(r0, regions[i_])];

        // top neighbor
        i_   = idx(ix, iy, hclampz(iz+1));
        avg += aLUT2d[symidx(r0, regions[i_])];
    }

    dst[I] = avg;
}
__kernel void
kernmulC(__global real_t* __restrict fftM, __global real_t* __restrict fftK, int Nx, int Ny) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if ((ix>= Nx) || (iy>=Ny)) {
        return;
    }

    int I = iy*Nx + ix;
    int e = 2 * I;

    real_t reM = fftM[e  ];
    real_t imM = fftM[e+1];
    real_t reK = fftK[e  ];
    real_t imK = fftK[e+1];

    fftM[e  ] = reM * reK - imM * imK;
    fftM[e+1] = reM * imK + imM * reK;
}
// 2D XY (in-plane) micromagnetic kernel multiplication:
// |Mx| = |Kxx Kxy| * |Mx|
// |My|   |Kyx Kyy|   |My|
// Using the same symmetries as kernmulrsymm3d.cl
__kernel void
kernmulRSymm2Dxy(__global real_t* __restrict  fftMx, __global real_t* __restrict  fftMy,
                 __global real_t* __restrict fftKxx, __global real_t* __restrict fftKyy, __global real_t* __restrict fftKxy,
                                         int     Nx,                         int     Ny) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if ((ix>= Nx) || (iy>=Ny)) {
        return;
    }

    int I = iy*Nx + ix;
    int e = 2 * I;

    real_t reMx = fftMx[e  ];
    real_t imMx = fftMx[e+1];
    real_t reMy = fftMy[e  ];
    real_t imMy = fftMy[e+1];

    // symmetry factor
    real_t fxy = (real_t)1.0;
    if (iy > Ny/2) {
         iy = Ny-iy;
        fxy = -fxy;
    }
    I = iy*Nx + ix;

    real_t Kxx = fftKxx[I];
    real_t Kyy = fftKyy[I];
    real_t Kxy = fxy * fftKxy[I];

    fftMx[e  ] = reMx * Kxx + reMy * Kxy;
    fftMx[e+1] = imMx * Kxx + imMy * Kxy;
    fftMy[e  ] = reMx * Kxy + reMy * Kyy;
    fftMy[e+1] = imMx * Kxy + imMy * Kyy;
}
// 2D Z (out-of-plane only) micromagnetic kernel multiplication:
// Mz = Kzz * Mz
// Using the same symmetries as kernmulrsymm3d.cl
__kernel void
kernmulRSymm2Dz(__global real_t* __restrict fftMz, __global real_t* __restrict fftKzz, int Nx, int Ny) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if ((ix>= Nx) || (iy>=Ny)) {
        return;
    }

    int I = iy*Nx + ix;
    int e = 2 * I;

    real_t reMz = fftMz[e  ];
    real_t imMz = fftMz[e+1];

    if (iy > Ny/2) {
        iy = Ny-iy;
    }
    I = iy*Nx + ix;

    real_t Kzz = fftKzz[I];

    fftMz[e  ] = reMz * Kzz;
    fftMz[e+1] = imMz * Kzz;
}
// 3D micromagnetic kernel multiplication:
//
// |Mx|   |Kxx Kxy Kxz|   |Mx|
// |My| = |Kxy Kyy Kyz| * |My|
// |Mz|   |Kxz Kyz Kzz|   |Mz|
//
// ~kernel has mirror symmetry along Y and Z-axis,
// apart form first row,
// and is only stored (roughly) half:
//
// K11, K22, K02:
// xxxxx
// aaaaa
// bbbbb
// ....
// bbbbb
// aaaaa
//
// K12:
// xxxxx
// aaaaa
// bbbbb
// ...
// -bbbb
// -aaaa

__kernel void
kernmulRSymm3D(__global real_t* __restrict  fftMx, __global real_t* __restrict  fftMy, __global real_t* __restrict  fftMz,
               __global real_t* __restrict fftKxx, __global real_t* __restrict fftKyy, __global real_t* __restrict fftKzz,
               __global real_t* __restrict fftKyz, __global real_t* __restrict fftKxz, __global real_t* __restrict fftKxy,
                                       int     Nx,                         int     Ny,                         int     Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix>= Nx) || (iy>= Ny) || (iz>=Nz)) {
        return;
    }

    // fetch (complex) FFT'ed magnetization
    int I = (iz*Ny + iy)*Nx + ix;
    int e = 2 * I;

    real_t reMx = fftMx[e  ];
    real_t imMx = fftMx[e+1];
    real_t reMy = fftMy[e  ];
    real_t imMy = fftMy[e+1];
    real_t reMz = fftMz[e  ];
    real_t imMz = fftMz[e+1];

    // fetch kernel

    // minus signs are added to some elements if
    // reconstructed from symmetry.
    real_t signYZ = (real_t)1.0;
    real_t signXZ = (real_t)1.0;
    real_t signXY = (real_t)1.0;

    // use symmetry to fetch from redundant parts:
    // mirror index into first quadrant and set signs.
    if (iy > Ny/2) {
            iy = Ny-iy;
        signYZ = -signYZ;
        signXY = -signXY;
    }
    if (iz > Nz/2) {
            iz = Nz-iz;
        signYZ = -signYZ;
        signXZ = -signXZ;
    }

    // fetch kernel element from non-redundant part
    // and apply minus signs for mirrored parts.
    I = (iz*(Ny/2+1) + iy)*Nx + ix; // Ny/2+1: only half is stored

    real_t Kxx = fftKxx[I];
    real_t Kyy = fftKyy[I];
    real_t Kzz = fftKzz[I];
    real_t Kyz = fftKyz[I] * signYZ;
    real_t Kxz = fftKxz[I] * signXZ;
    real_t Kxy = fftKxy[I] * signXY;

    // m * K matrix multiplication, overwrite m with result.
    fftMx[e  ] = reMx * Kxx + reMy * Kxy + reMz * Kxz;
    fftMx[e+1] = imMx * Kxx + imMy * Kxy + imMz * Kxz;
    fftMy[e  ] = reMx * Kxy + reMy * Kyy + reMz * Kyz;
    fftMy[e+1] = imMx * Kxy + imMy * Kyy + imMz * Kyz;
    fftMz[e  ] = reMx * Kxz + reMy * Kyz + reMz * Kzz;
    fftMz[e+1] = imMx * Kxz + imMy * Kyz + imMz * Kzz;
}
// Landau-Lifshitz torque without precession
__kernel void
llnoprecess(__global real_t* __restrict tx, __global real_t* __restrict ty, __global real_t* __restrict tz,
            __global real_t* __restrict mx, __global real_t* __restrict my, __global real_t* __restrict mz,
            __global real_t* __restrict hx, __global real_t* __restrict hy, __global real_t* __restrict hz,
                                    int  N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3 m = {mx[i], my[i], mz[i]};
        real_t3 H = {hx[i], hy[i], hz[i]};

        real_t3    mxH = cross(m, H);
        real_t3 torque = -cross(m, mxH);

        tx[i] = torque.x;
        ty[i] = torque.y;
        tz[i] = torque.z;
    }
}
// Landau-Lifshitz torque.
__kernel void
lltorque2(__global real_t* __restrict     tx, __global real_t* __restrict        ty, __global real_t* __restrict tz,
          __global real_t* __restrict     mx, __global real_t* __restrict        my, __global real_t* __restrict mz,
          __global real_t* __restrict     hx, __global real_t* __restrict        hy, __global real_t* __restrict hz,
          __global real_t* __restrict alpha_,                      real_t alpha_mul,                         int  N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3    m = {mx[i], my[i], mz[i]};
        real_t3    H = {hx[i], hy[i], hz[i]};
        real_t alpha = amul(alpha_, alpha_mul, i);

        real_t3    mxH = cross(m, H);
        real_t    gilb = (real_t)-1.0 / ((real_t)1.0 + alpha * alpha);
        real_t3 torque = gilb * (mxH + alpha * cross(m, mxH));

        tx[i] = torque.x;
        ty[i] = torque.y;
        tz[i] = torque.z;
    }
}
// dst[i] = fac1*src1[i] + fac2*src2[i];
__kernel void
madd2(__global real_t* __restrict   dst,
      __global real_t* __restrict  src1, real_t fac1,
      __global real_t* __restrict  src2, real_t fac2, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = fac1*src1[i] + fac2*src2[i];
    }
}
// dst[i] = fac1 * src1[i] + fac2 * src2[i] + fac3 * src3[i]
__kernel void
madd3(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i] + fac3 * src3[i]);
        // parens for better accuracy heun solver.
    }
}
// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4
__kernel void
madd4(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3,
      __global real_t* __restrict src4, real_t fac4, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]);
    }
}
// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5
__kernel void
madd5(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3,
      __global real_t* __restrict src4, real_t fac4,
      __global real_t* __restrict src5, real_t fac5, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]);
    }
}
// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5 + src6[i] * fac6
__kernel void
madd6(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3,
      __global real_t* __restrict src4, real_t fac4,
      __global real_t* __restrict src5, real_t fac5,
      __global real_t* __restrict src6, real_t fac6, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]) + (fac6 * src6[i]);
    }
}
// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5 + src6[i] * fac6 + src7[i] * fac7
__kernel void
madd7(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3,
      __global real_t* __restrict src4, real_t fac4,
      __global real_t* __restrict src5, real_t fac5,
      __global real_t* __restrict src6, real_t fac6,
      __global real_t* __restrict src7, real_t fac7, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]) + (fac6 * src6[i]) + (fac7 * src7[i]);
    }
}
// See maxangle.go for more details.
__kernel void
setmaxangle(__global real_t* __restrict    dst,
            __global real_t* __restrict     mx, __global  real_t* __restrict      my, __global real_t* __restrict mz,
            __global real_t* __restrict aLUT2d, __global uint8_t* __restrict regions,
                                    int     Nx,                          int      Ny,                         int Nz,
                                uint8_t    PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int      I = idx(ix, iy, iz);
    real_t3 m0 = make_float3(mx[I], my[I], mz[I]);

    if (is0(m0)) {
        return;
    }

    uint8_t    r0 = regions[I];
    real_t  angle = 0.0f;

    int      i_; // neighbor index
    real_t3  m_; // neighbor mag
    real_t  a__; // inter-cell exchange stiffness

    // left neighbor
    i_  = idx(lclampx(ix-1), iy, iz);           // clamps or wraps index according to PBC
    m_  = make_float3(mx[i_], my[i_], mz[i_]);  // load m
    m_  = ( is0(m_)? m0: m_ );                  // replace missing non-boundary neighbor
    a__ = aLUT2d[symidx(r0, regions[i_])];
    if (a__ != 0) {
        angle = max(angle, acos(dot(m_,m0)));
    }

    // right neighbor
    i_  = idx(hclampx(ix+1), iy, iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    if (a__ != 0) {
        angle = max(angle, acos(dot(m_,m0)));
    }

    // back neighbor
    i_  = idx(ix, lclampy(iy-1), iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    if (a__ != 0) {
        angle = max(angle, acos(dot(m_,m0)));
    }

    // front neighbor
    i_  = idx(ix, hclampy(iy+1), iz);
    m_  = make_float3(mx[i_], my[i_], mz[i_]);
    m_  = ( is0(m_)? m0: m_ );
    a__ = aLUT2d[symidx(r0, regions[i_])];
    if (a__ != 0) {
        angle = max(angle, acos(dot(m_,m0)));
    }

    // only take vertical derivative for 3D sim
    if (Nz != 1) {
        // bottom neighbor
        i_  = idx(ix, iy, lclampz(iz-1));
        m_  = make_float3(mx[i_], my[i_], mz[i_]);
        m_  = ( is0(m_)? m0: m_ );
        a__ = aLUT2d[symidx(r0, regions[i_])];
        if (a__ != 0) {
            angle = max(angle, acos(dot(m_,m0)));
        }

        // top neighbor
        i_  = idx(ix, iy, hclampz(iz+1));
        m_  = make_float3(mx[i_], my[i_], mz[i_]);
        m_  = ( is0(m_)? m0: m_ );
        a__ = aLUT2d[symidx(r0, regions[i_])];
        if (a__ != 0) {
            angle = max(angle, acos(dot(m_,m0)));
        }
    }

    dst[I] = angle;
}
// Steepest descent energy minimizer
__kernel void
minimize(__global real_t* __restrict  mx, __global real_t* __restrict  my, __global real_t* __restrict  mz,
         __global real_t* __restrict m0x, __global real_t* __restrict m0y, __global real_t* __restrict m0z,
         __global real_t* __restrict  tx, __global real_t* __restrict  ty, __global real_t* __restrict  tz,
                              real_t  dt,                         int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3 m0 = {m0x[i], m0y[i], m0z[i]};
        real_t3  t = {tx[i], ty[i], tz[i]};

        real_t       t2 = dt*dt*dot(t, t);
        real_t3  result = ((real_t)4.0 - t2) * m0 + (real_t)4.0 * dt * t;
        real_t  divisor = (real_t)4.0 + t2;

        mx[i] = result.x / divisor;
        my[i] = result.y / divisor;
        mz[i] = result.z / divisor;
    }
}
// dst[i] = a[i] * b[i]
__kernel void
mul(__global real_t* __restrict  dst, __global real_t* __restrict  a, __global real_t* __restrict b, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = a[i] * b[i];
    }
}
// normalize vector {vx, vy, vz} to unit length, unless length or vol are zero.
__kernel void
normalize2(__global real_t* __restrict vx, __global real_t* __restrict vy, __global real_t* __restrict vz, __global real_t* __restrict vol, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t  v = (vol == NULL) ? (real_t)1.0 : vol[i];
        real_t3 V = {v*vx[i], v*vy[i], v*vz[i]};

        V = normalize(V);
        if (v == (real_t)0.0) {
            vx[i] = 0.0;
            vy[i] = 0.0;
            vz[i] = 0.0;
        } else {
            vx[i] = V.x;
            vy[i] = V.y;
            vz[i] = V.z;
        }
    }
}
// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013

__kernel void
addoommfslonczewskitorque(__global real_t* __restrict            tx, __global real_t* __restrict              ty, __global real_t* __restrict tz,
                          __global real_t* __restrict            mx, __global real_t* __restrict              my, __global real_t* __restrict mz,
                          __global real_t* __restrict           Ms_,                      real_t          Ms_mul,
                          __global real_t* __restrict           jz_,                      real_t          jz_mul,
                          __global real_t* __restrict           px_,                      real_t          px_mul,
                          __global real_t* __restrict           py_,                      real_t          py_mul,
                          __global real_t* __restrict           pz_,                      real_t           pz_mul,
                          __global real_t* __restrict        alpha_,                      real_t        alpha_mul,
                          __global real_t* __restrict         pfix_,                      real_t         pfix_mul,
                          __global real_t* __restrict        pfree_,                      real_t        pfree_mul,
                          __global real_t* __restrict    lambdafix_,                      real_t    lambdafix_mul,
                          __global real_t* __restrict   lambdafree_,                      real_t   lambdafree_mul,
                          __global real_t* __restrict epsilonPrime_,                      real_t epsilonPrime_mul,
                          __global real_t* __restrict          flt_,                      real_t          flt_mul,
                                                  int             N) {

    int     I = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int I = gid; I < N; I += gsize) {

        real_t3 m = make_float3(mx[I], my[I], mz[I]);
        real_t  J = amul(jz_, jz_mul, I);
        real_t3 p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, I));

        real_t  Ms           = amul(Ms_, Ms_mul, I);
        real_t  alpha        = amul(alpha_, alpha_mul, I);
        real_t  flt          = amul(flt_, flt_mul, I);
        real_t  pfix         = amul(pfix_, pfix_mul, I);
        real_t  pfree        = amul(pfree_, pfree_mul, I);
        real_t  lambdafix    = amul(lambdafix_, lambdafix_mul, I);
        real_t  lambdafree   = amul(lambdafree_, lambdafix_mul, I);
        real_t  epsilonPrime = amul(epsilonPrime_, epsilonPrime_mul, I);

        if ((J == (real_t)0.0) || (Ms == (real_t)0.0)) {
            return;
        }

        real_t            beta = (HBAR / QE) * (J / ((real_t)2.0 *flt*Ms) );
        real_t      lambdafix2 = lambdafix * lambdafix;
        real_t     lambdafree2 = lambdafree * lambdafree;
        real_t  lambdafreePlus = sqrt(lambdafree2 + (real_t)1.0);
        real_t   lambdafixPlus = sqrt( lambdafix2 + (real_t)1.0);
        real_t lambdafreeMinus = sqrt(lambdafree2 - (real_t)1.0);
        real_t  lambdafixMinus = sqrt( lambdafix2 - (real_t)1.0);
        real_t      plus_ratio = lambdafreePlus / lambdafixPlus;
        real_t     minus_ratio = (real_t)1.0;

        if (lambdafreeMinus > (real_t)0.0) {
            minus_ratio = lambdafixMinus / lambdafreeMinus;
        }

        // Compute q_plus and q_minus
        real_t  plus_factor = pfix * lambdafix2 * plus_ratio;
        real_t minus_factor = pfree * lambdafree2 * minus_ratio;
        real_t       q_plus = plus_factor + minus_factor;
        real_t      q_minus = plus_factor - minus_factor;
        real_t       lplus2 = lambdafreePlus * lambdafixPlus;
        real_t      lminus2 = lambdafreeMinus * lambdafixMinus;
        real_t        pdotm = dot(p, m);
        real_t       A_plus = lplus2 + (lminus2 * pdotm);
        real_t      A_minus = lplus2 - (lminus2 * pdotm);
        real_t      epsilon = (q_plus / A_plus) - (q_minus / A_minus);

        real_t A = beta * epsilon;
        real_t B = beta * epsilonPrime;

        real_t gilb     = (real_t)1.0 / ((real_t)1.0 + alpha * alpha);
        real_t mxpxmFac = gilb * (A + alpha * B);
        real_t pxmFac   = gilb * (B - alpha * A);

        real_t3 pxm      = cross(p, m);
        real_t3 mxpxm    = cross(m, pxm);

        tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
        ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
        tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
    }
}
// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013
__kernel void
addtworegionoommfslonczewskitorque( __global real_t* __restrict            tx, __global real_t* __restrict               ty, __global real_t* __restrict      tz,
                                    __global real_t* __restrict            mx, __global real_t* __restrict               my, __global real_t* __restrict      mz,
                                    __global real_t* __restrict           Ms_,                      real_t           Ms_mul,
                                   __global uint8_t* __restrict       regions,
                                                        uint8_t       regionA,                     uint8_t          regionB,
                                                            int       strideX,                         int          strideY,                         int strideZ,
                                                            int            Nx,                         int               Ny,                         int      Nz,
                                                         real_t            j_,
                                                         real_t        alpha_,
                                                         real_t         pfix_,                      real_t           pfree_,
                                                         real_t    lambdafix_,                      real_t      lambdafree_,
                                                         real_t epsilonPrime_,
                                                         real_t          flt_) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int I = idx(ix, iy, iz);
    if (regions[I] != regionA) {
        return;
    }

    real_t3  m0 = make_float3(mx[I], my[I], mz[I]);
    real_t  Ms0 = amul(Ms_, Ms_mul, I);
    if (is0(m0) || (Ms0 == 0.0f)) {
        return;
    }

    int cX = ix + strideX;
    int cY = iy + strideY;
    int cZ = iz + strideZ;
    if ((cX >= Nx) || (cY >= Ny) || (cZ >= Nz)) {
        return;
    }

    int i_ = idx(cX, cY, cZ); // "neighbor" index
    if (regions[i_] != regionB) {
        return;
    }

    real_t3  m1 = make_float3(mx[i_], my[i_], mz[i_]); // "neighbor" mag
    real_t  Ms1 = amul(Ms_, Ms_mul, i_);
    if (is0(m1) || (Ms1 == 0.0f)) {
        return;
    }

    if ((j_ == 0.0f) || (Ms0 == 0.0f) || (Ms1 == 0.0f)) {
        return;
    }

    // Calculate for cell belonging to regionA
    real_t           beta0 = (HBAR / QE) * (j_ / (2.0f *flt_) );
    real_t            beta = beta0 / Ms0;
    real_t      lambdafix2 = lambdafix_ * lambdafix_;
    real_t     lambdafree2 = lambdafree_ * lambdafree_;
    real_t  lambdafreePlus = sqrt(lambdafree2 + 1.0f);
    real_t   lambdafixPlus = sqrt(lambdafix2 + 1.0f);
    real_t lambdafreeMinus = sqrt(lambdafree2 - 1.0f);
    real_t  lambdafixMinus = sqrt(lambdafix2 - 1.0f);
    real_t      plus_ratio = lambdafreePlus / lambdafixPlus;
    real_t     minus_ratio = 1.0f;

    if (lambdafreeMinus > 0.0f) {
        minus_ratio = lambdafixMinus / lambdafreeMinus;
    }

    // Compute q_plus and q_minus
    real_t  plus_factor = pfix_ * lambdafix2 * plus_ratio;
    real_t minus_factor = pfree_ * lambdafree2 * minus_ratio;
    real_t       q_plus = plus_factor + minus_factor;
    real_t      q_minus = plus_factor - minus_factor;
    real_t       lplus2 = lambdafreePlus * lambdafixPlus;
    real_t      lminus2 = lambdafreeMinus * lambdafixMinus;
    real_t        pdotm = dot(m1, m0);
    real_t       A_plus = lplus2 + (lminus2 * pdotm);
    real_t      A_minus = lplus2 - (lminus2 * pdotm);
    real_t      epsilon = (q_plus / A_plus) - (q_minus / A_minus);

    real_t A = beta * epsilon;
    real_t B = beta * epsilonPrime_;

    real_t     gilb = 1.0f / (1.0f + alpha_ * alpha_);
    real_t mxpxmFac = gilb * (A + alpha_ * B);
    real_t   pxmFac = gilb * (B - alpha_ * A);

    real_t3   pxm = cross(m1, m0);
    real_t3 mxpxm = cross(m0, pxm);

    tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
    ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
    tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;

    // Now calculate for cell in regionB
    beta        = beta0 / Ms1;
    plus_ratio  = lambdafixPlus / lambdafreePlus;
    minus_ratio = 1.0f;

    if (lambdafixMinus > 0.0f) {
        minus_ratio = lambdafreeMinus / lambdafixMinus;
    }

    // Compute q_plus and q_minus
    plus_factor  = pfree_ * lambdafree2 * plus_ratio;
    minus_factor = pfix_ * lambdafix2 * minus_ratio;
    q_plus       = plus_factor + minus_factor;
    q_minus      = plus_factor - minus_factor;
    epsilon      = (q_plus / A_plus) - (q_minus / A_minus);

    A = beta * epsilon;
    B = beta * epsilonPrime_;

    mxpxmFac = gilb * (A + alpha_ * B);
    pxmFac   = gilb * (B - alpha_ * A);

    mxpxm = cross(m1, pxm);

    tx[i_] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
    ty[i_] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
    tz[i_] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
}
__kernel void
reducedot(         __global real_t* __restrict     src1,
                   __global real_t* __restrict     src2,
          volatile __global real_t* __restrict      dst,
                            real_t              initVal,
                               int                    n,
          volatile __local  real_t*            scratch){

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int            i = get_group_id(0)*grp_sz + local_idx;
    int       stride = get_global_size(0);
    // Initialize ring accumulator for intermediate results
    real_t accum[__REDUCE_REG_COUNT__];
    for (unsigned int s = 0; s < __REDUCE_REG_COUNT__; s++) {
        accum[s] = 0.0;
    }
    accum[0] = initVal;
    unsigned int itr = 0;

    // Read from global memory and accumulate in workitem ring accumulator
    while (i < n) {
        accum[itr] += src1[i] * src2[i]; // Load value from global buffer into ring accumulator

        // Update pointer to ring accumulator
        itr++;
        if (itr >= __REDUCE_REG_COUNT__) {
            itr = 0;
        }

        // Update pointer to next global value
        i += stride;
    }

    // All elements in global buffer have been picked up
    // Reduce intermediate results and add atomically to global buffer

    // Reduce value in ring buffer
    for (unsigned int s1 = (__REDUCE_REG_COUNT__ >> 1); s1 > 1; s1 >>= 1) {
        for (unsigned int s2 = 0; s2 < s1; s2++) {
            accum[s2] += accum[s2+s1];
        }
    }

    // Reduce in local buffer
    scratch[local_idx] = accum[0] + accum[1];

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll loop for remaining 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] += smem[local_idx + 32];
        smem[local_idx] += smem[local_idx + 16];
        smem[local_idx] += smem[local_idx +  8];
        smem[local_idx] += smem[local_idx +  4];
        smem[local_idx] += smem[local_idx +  2];
        smem[local_idx] += smem[local_idx +  1];
    }

    // Add atomically to global buffer
    if (local_idx == 0) {
        atomicAdd_r(dst, scratch[0]);
    }

}
__kernel void
reducemaxabs(         __global real_t* __restrict     src,
             volatile __global real_t* __restrict     dst,
                               real_t             initVal,
                                  int                   n,
             volatile  __local real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int            i = get_group_id(0)*grp_sz + local_idx;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (i < (unsigned int)(n)) {
        mine = fmax(mine, src[i]);
        i += stride;
    }

    // Load workitem value into local buffer and synchronize
    scratch[local_idx] = mine;
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce using lor loop
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
        }

        // Synchronize workitems before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll for loop that executes within one unit that works on 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
	atomicMax_r(dst, mine);
    }

}
__kernel void
reducemaxdiff(         __global real_t* __restrict    src1,
                       __global real_t* __restrict    src2,
              volatile __global real_t* __restrict     dst,
                                real_t             initVal,
                                   int                   n,
              volatile __local  real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int            i = get_group_id(0)*grp_sz + local_idx;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (i < (unsigned int)(n)) {
        mine = fmax(mine, fabs(src1[i] - src2[i]));
        i += stride;
    }

    // Load workitem value into local buffer and synchronize
    scratch[local_idx] = mine;
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce using lor loop
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
        }

        // Synchronize workitems before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll for loop that executes within one unit that works on 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
        atomicMax_r(dst, mine);
    }

}
__kernel void
reducemaxvecdiff2(         __global real_t* __restrict      x1,
                           __global real_t* __restrict      y1,
                           __global real_t* __restrict      z1,
                           __global real_t* __restrict      x2,
                           __global real_t* __restrict      y2,
                           __global real_t* __restrict      z2,
                  volatile __global real_t* __restrict     dst,
                                    real_t             initVal,
                                       int                   n,
                  volatile __local  real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int            i = get_group_id(0)*grp_sz + local_idx;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (i < n) {
        real_t3 v = distance((real_t3){x1[i], y1[i], z1[i]}, (real_t3){x2[i], y2[i], z2[i]});
        mine = fmax(mine, dot(v, v));
        i += stride;
    }

    // Load workitem value into local buffer and synchronize
    scratch[local_idx] = mine;
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce using lor loop
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
        }

        // Synchronize workitems before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll for loop that executes within one unit that works on 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
        atomicMax_r(dst, mine);
    }

}
__kernel void
reducemaxvecnorm2(         __global real_t* __restrict       x,
                           __global real_t* __restrict       y,
                           __global real_t* __restrict       z,
                  volatile __global real_t* __restrict     dst,
                                    real_t             initVal,
                                       int                   n,
                  volatile __local  real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int            i = get_group_id(0)*grp_sz + local_idx;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (i < n) {
        real_t3 v = (real_t3){x[i], y[i], z[i]};
        mine = fmax(mine, dot(v, v));
        i += stride;
    }

    // Load workitem value into local buffer and synchronize
    scratch[local_idx] = mine;
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce using lor loop
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
        }

        // Synchronize workitems before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll for loop that executes within one unit that works on 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
        atomicMax_r(dst, mine);
    }

}
__kernel void
reducesum(         __global real_t*    __restrict     src,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                      n,
          volatile __local  real_t*               scratch){

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int            i = get_group_id(0)*grp_sz + local_idx;
    int       stride = get_global_size(0);
    // Initialize ring accumulator for intermediate results
    real_t accum[__REDUCE_REG_COUNT__];
    for (unsigned int s = 0; s < __REDUCE_REG_COUNT__; s++) {
        accum[s] = 0.0;
    }
    accum[0] = initVal;
    unsigned int itr = 0;

    // Read from global memory and accumulate in workitem ring accumulator
    while (i < n) {
        accum[itr] += src[i]; // Load value from global buffer into ring accumulator

        // Update pointer to ring accumulator
        itr++;
        if (itr >= __REDUCE_REG_COUNT__) {
            itr = 0;
        }

        // Update pointer to next global value
        i += stride;
    }

    // All elements in global buffer have been picked up
    // Reduce intermediate results and add atomically to global buffer

    // Reduce value in ring buffer
    for (unsigned int s1 = (__REDUCE_REG_COUNT__ >> 1); s1 > 1; s1 >>= 1) {
        for (unsigned int s2 = 0; s2 < s1; s2++) {
            accum[s2] += accum[s2+s1];
        }
    }

    // Reduce in local buffer
    scratch[local_idx] = accum[0] + accum[1];

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll loop for remaining 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] += smem[local_idx + 32];
        smem[local_idx] += smem[local_idx + 16];
        smem[local_idx] += smem[local_idx +  8];
        smem[local_idx] += smem[local_idx +  4];
        smem[local_idx] += smem[local_idx +  2];
        smem[local_idx] += smem[local_idx +  1];
    }

    // Add atomically to global buffer
    if (local_idx == 0) {
        atomicAdd_r(dst, scratch[0]);
    }

}
// add region-based vector to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionaddv(__global  real_t* __restrict    dstx, __global real_t* __restrict dsty, __global real_t* __restrict dstz,
           __global  real_t* __restrict    LUTx, __global real_t* __restrict LUTy, __global real_t* __restrict LUTz,
           __global uint8_t* __restrict regions,
                                    int       N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        uint8_t r = regions[i];
        dstx[i] += LUTx[r];
        dsty[i] += LUTy[r];
        dstz[i] += LUTz[r];
    }
}
// add region-based scalar to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionadds(__global  real_t* __restrict     dst,
           __global  real_t* __restrict     LUT,
           __global uint8_t* __restrict regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        uint8_t r = regions[i];
        dst[i] += LUT[r];
    }
}
// decode the regions+LUT pair into an uncompressed array
__kernel void
regiondecode(__global real_t* __restrict dst, __global real_t* __restrict LUT, __global uint8_t* regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = LUT[regions[i]];
    }
}
__kernel void
regionselect(__global real_t* __restrict dst, __global real_t* __restrict src, __global uint8_t* regions, uint8_t region, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = ((regions[i] == region) ? src[i]: (real_t)0.0);
    }
}
// Select and resize one layer for interactive output
__kernel void
resize(__global real_t* __restrict   dst, int     Dx, int     Dy, int Dz,
       __global real_t* __restrict   src, int     Sx, int     Sy, int Sz,
                               int layer, int scalex, int scaley) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if (ix<Dx && iy<Dy) {

        real_t sum = (real_t)0.0;
        real_t   n = (real_t)0.0;

        for (int J=0; J<scaley; J++) {
            int j2 = iy*scaley+J;

            for (int K=0; K<scalex; K++) {
                int k2 = ix*scalex+K;

                if ((j2 < Sy) && (k2 < Sx)) {
                    sum += src[(layer*Sy + j2)*Sx + k2];
                    n += (real_t)1.0;
                }
            }
        }
        dst[iy*Dx + ix] = sum / n;
    }
}
// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shiftbytes(__global uint8_t* __restrict dst, __global uint8_t* __restrict src,
                                    int  Nx,                          int  Ny, int Nz, int shx, uint8_t clampV) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix < Nx) && (iy < Ny) && (iz < Nz)) {
        int ix2 = ix-shx;
        uint8_t newval;
        if ((ix2 < 0) || (ix2 >= Nx)) {
            newval = clampV;
        } else {
            newval = src[idx(ix2, iy, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}
// shift dst by shy cells (positive or negative) along Y-axis.
__kernel void
shiftbytesy(__global uint8_t* __restrict dst, __global uint8_t* __restrict src,
                                     int  Nx,                          int  Ny, int Nz, int shy, uint8_t clampV) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix < Nx) && (iy < Ny) && (iz < Nz)) {
        int iy2 = iy-shy;
        uint8_t newval;
        if ((iy2 < 0) || (iy2 >= Ny)) {
            newval = clampV;
        } else {
            newval = src[idx(ix, iy2, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}
// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shiftx(__global real_t* __restrict dst, __global real_t* __restrict src,
                               int  Nx,                         int  Ny, int Nz, int shx, real_t clampL, real_t clampR) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix < Nx) && (iy < Ny) && (iz < Nz)) {
        int ix2 = ix-shx;
        real_t newval;
        if (ix2 < 0) {
            newval = clampL;
        } else if (ix2 >= Nx) {
            newval = clampR;
        } else {
            newval = src[idx(ix2, iy, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}
// shift dst by shy cells (positive or negative) along Y-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shifty(__global real_t* __restrict dst, __global real_t* __restrict src,
                               int  Nx,                         int  Ny, int Nz, int shy, real_t clampL, real_t clampR) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix < Nx) && (iy < Ny) && (iz < Nz)) {
        int iy2 = iy-shy;
        real_t newval;
        if (iy2 < 0) {
            newval = clampL;
        } else if (iy2 >= Ny) {
            newval = clampR;
        } else {
            newval = src[idx(ix, iy2, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}
// shift dst by shz cells (positive or negative) along Z-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shiftz(__global real_t* __restrict dst, __global real_t* __restrict src,
                               int  Nx,                         int  Ny, int Nz, int shz, real_t clampL, real_t clampR) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix < Nx) && (iy < Ny) && (iz < Nz)) {
        int iz2 = iz-shz;
        real_t newval;
        if (iz2 < 0) {
            newval = clampL;
        } else if (iz2 >= Nz) {
            newval = clampR;
        } else {
            newval = src[idx(ix, iy, iz2)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}
// Add magneto-elastic coupling field to B.
// H = - δUmel / δM, 
// where Umel is magneto-elastic energy denstiy given by the eq. (12.18) of Gurevich&Melkov "Magnetization Oscillations and Waves", CRC Press, 1996
__kernel void
addmagnetoelasticfield(__global real_t* __restrict   Bx, __global real_t* __restrict      By, __global real_t* __restrict  Bz,
                       __global real_t* __restrict   mx, __global real_t* __restrict      my, __global real_t* __restrict  mz,
                       __global real_t* __restrict exx_,                      real_t exx_mul,
                       __global real_t* __restrict eyy_,                      real_t eyy_mul,
                       __global real_t* __restrict ezz_,                      real_t ezz_mul,
                       __global real_t* __restrict exy_,                      real_t exy_mul,
                       __global real_t* __restrict exz_,                      real_t exz_mul,
                       __global real_t* __restrict eyz_,                      real_t eyz_mul,
                       __global real_t* __restrict  B1_,                      real_t  B1_mul, 
                       __global real_t* __restrict  B2_,                      real_t  B2_mul,
                       __global real_t* __restrict  Ms_,                      real_t  Ms_mul,
                                               int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int I = gid; I < N; I += gsize) {

        real_t Exx = amul(exx_, exx_mul, I);
        real_t Eyy = amul(eyy_, eyy_mul, I);
        real_t Ezz = amul(ezz_, ezz_mul, I);

        real_t Exy = amul(exy_, exy_mul, I);
        real_t Eyx = Exy;

        real_t Exz = amul(exz_, exz_mul, I);
        real_t Ezx = Exz;

        real_t Eyz = amul(eyz_, eyz_mul, I);
        real_t Ezy = Eyz;

        real_t invMs = inv_Msat(Ms_, Ms_mul, I);

        real_t B1 = amul(B1_, B1_mul, I) * invMs;
        real_t B2 = amul(B2_, B2_mul, I) * invMs;

        real_t3 m = {mx[I], my[I], mz[I]};

        Bx[I] += -((real_t)2.0*B1*m.x*Exx + B2*(m.y*Exy + m.z*Exz));
        By[I] += -((real_t)2.0*B1*m.y*Eyy + B2*(m.x*Eyx + m.z*Eyz));
        Bz[I] += -((real_t)2.0*B1*m.z*Ezz + B2*(m.x*Ezx + m.y*Ezy));
    }
}
// Calculate magneto-elastic force density
// fmelp = Σ ∂σpq / ∂xq (q = x, y, z) , σpq = ∂Umel / ∂epq, 
// where epq is the strain tensor and 
// Umel is the magneto-elastic energy density given by the eq. (12.18) of Gurevich&Melkov "Magnetization Oscillations and Waves", CRC Press, 1996
__kernel void
getmagnetoelasticforce(__global real_t* __restrict   fx, __global real_t* __restrict     fy, __global real_t* __restrict   fz,
                       __global real_t* __restrict   mx, __global real_t* __restrict     my, __global real_t* __restrict   mz,
                       __global real_t* __restrict  B1_,                      real_t B1_mul, 
                       __global real_t* __restrict  B2_,                      real_t B2_mul,
                                            real_t rcsx,                      real_t   rcsy,                      real_t rcsz,
                                               int   Nx,                         int     Ny,                         int   Nz, 
                                           uint8_t  PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz))
    {
        return;
    }

    int        I = idx(ix, iy, iz);                  // central cell index
    real_t3   m0 = make_float3(mx[I], my[I], mz[I]); // +0
    real_t3 dmdx = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // ∂m/∂x
    real_t3 dmdy = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // ∂m/∂y
    real_t3 dmdz = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // ∂m/∂z
    int i_;                                          // neighbor index

    // ∂m/∂x
    {    
        real_t3 m_m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // -2
        i_ = idx(lclampx(ix-2), iy, iz);                 // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-2 >= 0) || PBCx)
        {
            m_m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // -1
        i_ = idx(lclampx(ix-1), iy, iz);                 // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-1 >= 0) || PBCx)
        {
            m_m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // +1
        i_ = idx(hclampx(ix+1), iy, iz);
        if ((ix+1 < Nx) || PBCx)
        {
            m_p1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // +2
        i_ = idx(hclampx(ix+2), iy, iz);
        if ((ix+2 < Nx) || PBCx)
        {
            m_p2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        if (is0(m_p1) && is0(m_m1))                                        //  +0
        {
            dmdx = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);                          // --1-- zero
        }
        else if ((is0(m_m2) | is0(m_p2)) && !is0(m_p1) && !is0(m_m1))
        {
            dmdx = 0.5f * (m_p1 - m_m1);                                   // -111-, 1111-, -1111 central difference,  ε ~ h^2
        }
        else if (is0(m_p1) && is0(m_m2))
        {
            dmdx =  m0 - m_m1;                                             // -11-- backward difference, ε ~ h^1
        }
        else if (is0(m_m1) && is0(m_p2))
        {
            dmdx = -m0 + m_p1;                                             // --11- forward difference,  ε ~ h^1
        }
        else if (!is0(m_m2) && is0(m_p1))
        {
            dmdx =  0.5f * m_m2 - 2.0f * m_m1 + 1.5f * m0;                 // 111-- backward difference, ε ~ h^2
        }
        else if (!is0(m_p2) && is0(m_m1))
        {
            dmdx = -0.5f * m_p2 + 2.0f * m_p1 - 1.5f * m0;                 // --111 forward difference,  ε ~ h^2
        }
        else
        {
            dmdx = (2.0f/3.0f)*(m_p1 - m_m1) + (1.0f/12.0f)*(m_m2 - m_p2); // 11111 central difference,  ε ~ h^4
        }
    }

    // ∂m/∂y
    {
        real_t3 m_m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-2), iz);
        if ((iy-2 >= 0) || PBCy)
        {
            m_m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-1), iz);
        if ((iy-1 >= 0) || PBCy)
        {
            m_m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+1), iz);
        if  ((iy+1 < Ny) || PBCy)
        {
            m_p1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+2), iz);
        if  (iy+2 < Ny || PBCy)
        {
            m_p2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        if (is0(m_p1) && is0(m_m1))                                        //  +0
        {
            dmdy = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);                          // --1-- zero
        }
        else if ((is0(m_m2) | is0(m_p2)) && !is0(m_p1) && !is0(m_m1))
        {
            dmdy = 0.5f * (m_p1 - m_m1);                                   // -111-, 1111-, -1111 central difference,  ε ~ h^2
        }
        else if (is0(m_p1) && is0(m_m2))
        {
            dmdy =  m0 - m_m1;                                             // -11-- backward difference, ε ~ h^1
        }
        else if (is0(m_m1) && is0(m_p2))
        {
            dmdy = -m0 + m_p1;                                             // --11- forward difference,  ε ~ h^1
        }
        else if (!is0(m_m2) && is0(m_p1))
        {
            dmdy =  0.5f * m_m2 - 2.0f * m_m1 + 1.5f * m0;                 // 111-- backward difference, ε ~ h^2
        }
        else if (!is0(m_p2) && is0(m_m1))
        {
            dmdy = -0.5f * m_p2 + 2.0f * m_p1 - 1.5f * m0;                 // --111 forward difference,  ε ~ h^2
        }
        else
        {
            dmdy = (2.0f/3.0f)*(m_p1 - m_m1) + (1.0f/12.0f)*(m_m2 - m_p2); // 11111 central difference,  ε ~ h^4
        }
    }


    // ∂u/∂z
    {
        real_t3 m_m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, iy, lclampz(iz-2));
        if ((iz-2 >= 0) || PBCz)
        {
            m_m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, iy, lclampz(iz-1));
        if ((iz-1 >= 0) || PBCz)
        {
            m_m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, iy, hclampz(iz+1));
        if  ((iz+1 < Nz) || PBCz)
        {
            m_p1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, iy, hclampz(iz+2));
        if  ((iz+2 < Nz) || PBCz)
        {
            m_p2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        if (is0(m_p1) && is0(m_m1))                                        //  +0
        {
            dmdz = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);                          // --1-- zero
        }
        else if ((is0(m_m2) | is0(m_p2)) && !is0(m_p1) && !is0(m_m1))
        {
            dmdz = 0.5f * (m_p1 - m_m1);                                   // -111-, 1111-, -1111 central difference,  ε ~ h^2
        }
        else if (is0(m_p1) && is0(m_m2))
        {
            dmdz =  m0 - m_m1;                                             // -11-- backward difference, ε ~ h^1
        }
        else if (is0(m_m1) && is0(m_p2))
        {
            dmdz = -m0 + m_p1;                                             // --11- forward difference,  ε ~ h^1
        }
        else if (!is0(m_m2) && is0(m_p1))
        {
            dmdz =  0.5f * m_m2 - 2.0f * m_m1 + 1.5f * m0;                 // 111-- backward difference, ε ~ h^2
        }
        else if (!is0(m_p2) && is0(m_m1))
        {
            dmdz = -0.5f * m_p2 + 2.0f * m_p1 - 1.5f * m0;                 // --111 forward difference,  ε ~ h^2
        }
        else
        {
            dmdz = (2.0f/3.0f)*(m_p1 - m_m1) + (1.0f/12.0f)*(m_m2 - m_p2); // 11111 central difference,  ε ~ h^4
        }
    }

    dmdx *= rcsx;
    dmdy *= rcsy;
    dmdz *= rcsz;

    real_t B1 = amul(B1_, B1_mul, I);
    real_t B2 = amul(B2_, B2_mul, I);

    fx[I] = 2.0f*B1*m0.x*dmdx.x + B2*(m0.x*(dmdy.y + dmdz.z) + m0.y*dmdy.x + m0.z*dmdz.x);
    fy[I] = 2.0f*B1*m0.y*dmdy.y + B2*(m0.x*dmdx.y + m0.y*(dmdx.x + dmdz.z) + m0.z*dmdz.y);
    fz[I] = 2.0f*B1*m0.z*dmdz.z + B2*(m0.x*dmdx.z + m0.y*dmdy.z + m0.z*(dmdx.x + dmdy.y));
}
// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013, 2016

__kernel void
addslonczewskitorque2(__global real_t* __restrict                tx, __global real_t* __restrict             ty, __global real_t* __restrict tz,
                      __global real_t* __restrict                mx, __global real_t* __restrict             my, __global real_t* __restrict mz,
                      __global real_t* __restrict               Ms_,                      real_t         Ms_mul,
                      __global real_t* __restrict               jz_,                      real_t         jz_mul,
                      __global real_t* __restrict               px_,                      real_t         px_mul,
                      __global real_t* __restrict               py_,                      real_t         py_mul,
                      __global real_t* __restrict               pz_,                      real_t         pz_mul,
                      __global real_t* __restrict            alpha_,                      real_t      alpha_mul,
                      __global real_t* __restrict              pol_,                      real_t        pol_mul,
                      __global real_t* __restrict           lambda_,                      real_t     lambda_mul,
                      __global real_t* __restrict         epsPrime_,                      real_t   epsPrime_mul,
                      __global real_t* __restrict        thickness_,                      real_t  thickness_mul,
                                           real_t     meshThickness,
                                           real_t freeLayerPosition,
                                              int                 N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3            m = make_float3(mx[i], my[i], mz[i]);
        real_t             J = amul(jz_, jz_mul, i);
        real_t3            p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, i));
        real_t            Ms = amul(Ms_, Ms_mul, i);
        real_t         alpha = amul(alpha_, alpha_mul, i);
        real_t           pol = amul(pol_, pol_mul, i);
        real_t        lambda = amul(lambda_, lambda_mul, i);
        real_t  epsilonPrime = amul(epsPrime_, epsPrime_mul, i);
        real_t     thickness = amul(thickness_, thickness_mul, i);

        if (thickness == (real_t)0.0) { // if thickness is not set, use the thickness of the mesh instead
            thickness = meshThickness;
        }
        thickness *= freeLayerPosition; // switch sign if fixedlayer is at the bottom

        if (J == (real_t)0.0 || Ms == (real_t)0.0) {
            return;
        }

        real_t    beta = (HBAR / QE) * (J / (thickness*Ms) );
        real_t lambda2 = lambda * lambda;
        real_t epsilon = pol * lambda2 / ((lambda2 + (real_t)1.0) + (lambda2 - (real_t)1.0) * dot(p, m));

        real_t A = beta * epsilon;
        real_t B = beta * epsilonPrime;

        real_t     gilb = (real_t)1.0 / ((real_t)1.0 + alpha * alpha);
        real_t mxpxmFac = gilb * (A + alpha * B);
        real_t   pxmFac = gilb * (B - alpha * A);

        real_t3   pxm = cross(p, m);
        real_t3 mxpxm = cross(m, pxm);

        tx[i] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
        ty[i] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
        tz[i] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
    }
}
// Add two region exchange field to Beff.
// The cells of the regions are separated
// by the displacement vector defined by
// real_t3{strideX*cellsize[X], strideY*cellsize[Y], strideZ*cellsize[Z]}
//        m: normalized magnetization
//        B: effective field in Tesla
//  sig_eff: bilinear exchange coefficient (with cell discretization) in J / m^3
// sig2_eff: biquadratic exchange coefficient (with cell discretization) in J / m^3
__kernel void
tworegionexchange_field( __global real_t* __restrict      Bx, __global real_t* __restrict       By, __global real_t* __restrict      Bz,
                         __global real_t* __restrict      mx, __global real_t* __restrict       my, __global real_t* __restrict      mz,
                         __global real_t* __restrict     Ms_,                      real_t   Ms_mul,
                        __global uint8_t* __restrict regions,
                                             uint8_t regionA,                     uint8_t  regionB,
                                                 int strideX,                         int  strideY,                         int strideZ,
                                              real_t sig_eff,                      real_t sig2_eff,
                                                 int      Nx,                         int       Ny,                         int      Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int I = idx(ix, iy, iz);
    if (regions[I] != regionA) {
        return;
    }

    real_t3  m0 = make_float3(mx[I], my[I], mz[I]);
    real_t  Ms0 = amul(Ms_, Ms_mul, I);

    if (is0(m0) || (Ms0 == 0.0f)) {
        return;
    }

    int cX = ix + strideX;
    int cY = iy + strideY;
    int cZ = iz + strideZ;

    if ((cX >= Nx) || (cY >= Ny) || (cZ >= Nz)) {
        return;
    }

    int i_ = idx(cX, cY, cZ); // "neighbor" index
    if (regions[i_] != regionB) {
        return;
    }

    real_t3  m1 = make_float3(mx[i_], my[i_], mz[i_]); // "neighbor" mag
    real_t  Ms1 = amul(Ms_, Ms_mul, i_);
    if (is0(m1) || (Ms1 == 0.0f)) {
        return;
    }

    real_t3     B = m0 - m1;
    real_t   dot1 = dot(m0, m1);
    real_t    fac = 2.0f * (sig_eff + 2.0f * sig2_eff * dot1);
    real_t  invMs = inv_Msat(Ms_, Ms_mul, I);

    if (Bx != NULL) {
        Bx[I]  -= B.x * (fac*invMs);
        Bx[i_] += B.x * (fac*invMs);
    }
    if (By != NULL) {
        By[I]  -= B.y * (fac*invMs);
        By[i_] += B.y * (fac*invMs);
    }
    if (Bz != NULL) {
        Bz[I]  -= B.z * (fac*invMs);
        Bz[i_] += B.z * (fac*invMs);
    }
}
// Add two region exchange energy to Edens.
// The cells of the regions are separated
// by the displacement vector
// real_t3{strideX*cellsize[X], strideY*cellsize[Y], strideZ*cellsize[Z]}
//        m: normalized magnetization
//    Edens: energy density in J / m^3
//  sig_eff: bilinear exchange coefficient (with cell discretization) in J / m^3
// sig2_eff: biquadratic exchange coefficient (with cell discretization) in J / m^3
__kernel void
tworegionexchange_edens( __global real_t* __restrict   Edens,
                         __global real_t* __restrict      mx, __global real_t* __restrict       my, __global real_t* __restrict      mz,
                         __global real_t* __restrict     Ms_,                      real_t   Ms_mul,
                        __global uint8_t* __restrict regions,
                                             uint8_t regionA,                     uint8_t  regionB,
                                                 int strideX,                         int  strideY,                         int strideZ,
                                              real_t sig_eff,                      real_t sig2_eff,
                                                 int      Nx,                         int       Ny,                         int      Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int I = idx(ix, iy, iz);
    if (regions[I] != regionA) {
        return;
    }

    real_t3  m0 = make_float3(mx[I], my[I], mz[I]);
    real_t  Ms0 = amul(Ms_, Ms_mul, I);

    if (is0(m0) || (Ms0 == 0.0f)) {
        return;
    }

    int cX = ix + strideX;
    int cY = iy + strideY;
    int cZ = iz + strideZ;

    if ((cX >= Nx) || (cY >= Ny) || (cZ >= Nz)) {
        return;
    }

    int i_ = idx(cX, cY, cZ); // "neighbor" index
    if (regions[i_] != regionB) {
        return;
    }

    real_t3  m1 = make_float3(mx[i_], my[i_], mz[i_]); // "neighbor" mag
    real_t  Ms1 = amul(Ms_, Ms_mul, i_);

    if (is0(m1) || (Ms1 == 0.0f)) {
            return;
    }

    if (Edens != NULL) {
        real_t dot1 = dot(m0, m1);
        Edens[I]  += (sig_eff + sig2_eff * (1.0f + dot1)) * (1.0f - dot1);
    }
}
__kernel void
setPhi(__global real_t* __restrict phi,
       __global real_t* __restrict  mx, __global real_t* __restrict my,
                               int  Nx,                         int Ny, int Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int  I = idx(ix, iy, iz);     // central cell index
    phi[I] = atan2(my[I], mx[I]);
}
__kernel void
setTheta(__global real_t* __restrict theta,
         __global real_t* __restrict    mz,
                                 int    Nx, int Ny, int Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int    I = idx(ix, iy, iz);    // central cell index
    theta[I] = acos(mz[I]);
}
// TODO: this could act on x,y,z, so that we need to call it only once.
__kernel void
settemperature2(__global real_t* __restrict      B, __global real_t* __restrict     noise, real_t kB2_VgammaDt,
                __global real_t* __restrict    Ms_,                      real_t    Ms_mul,
                __global real_t* __restrict  temp_,                      real_t  temp_mul,
                __global real_t* __restrict alpha_,                      real_t alpha_mul,
                                        int      N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t  temp = amul(temp_, temp_mul, i);
        real_t alpha = amul(alpha_, alpha_mul, i);

        B[i] = noise[i] * sqrt((kB2_VgammaDt * alpha * temp * invMs ));
    }
}
// Set s to the topological charge density.
// See topologicalcharge.go.
__kernel void
settopologicalcharge(__global real_t* __restrict     s,
                     __global real_t* __restrict    mx, __global real_t* __restrict my, __global real_t* __restrict mz,
                                          real_t icxcy,
                                             int    Nx,                         int Ny,                         int Nz,
                                         uint8_t   PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int I = idx(ix, iy, iz);                      // central cell index

    real_t3          m0 = make_float3(mx[I], my[I], mz[I]); // +0
    real_t3        dmdx = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // ∂m/∂x
    real_t3        dmdy = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);    // ∂m/∂y
    real_t3 dmdx_x_dmdy = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);       // ∂m/∂x ❌ ∂m/∂y
    int i_;                                                // neighbor index

    if (is0(m0)) {
        s[I] = 0.0f;
        return;
    }

    // x derivatives (along length)
    {
        real_t3 m_m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // -2
        i_ = idx(lclampx(ix-2), iy, iz);                 // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-2 >= 0) || PBCx)
        {
            m_m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // -1
        i_ = idx(lclampx(ix-1), iy, iz);                 // load neighbor m if inside grid, keep 0 otherwise
        if ((ix-1 >= 0) || PBCx)
        {
            m_m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // +1
        i_ = idx(hclampx(ix+1), iy, iz);
        if ((ix+1 < Nx) || PBCx)
        {
            m_p1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);     // +2
        i_ = idx(hclampx(ix+2), iy, iz);
        if ((ix+2 < Nx) || PBCx)
        {
            m_p2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        if (is0(m_p1) && is0(m_m1))                       //  +0
        {
            dmdx = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);         // --1-- zero
        }
        else if ((is0(m_m2) | is0(m_p2)) && !is0(m_p1) && !is0(m_m1))
        {
            dmdx = 0.5f * (m_p1 - m_m1);                  // -111-, 1111-, -1111 central difference,  ε ~ h^2
        }
        else if (is0(m_p1) && is0(m_m2))
        {
            dmdx =  m0 - m_m1;                            // -11-- backward difference, ε ~ h^1
        }
        else if (is0(m_m1) && is0(m_p2))
        {
            dmdx = -m0 + m_p1;                            // --11- forward difference,  ε ~ h^1
        }
        else if (!is0(m_m2) && is0(m_p1))
        {
            dmdx =  0.5f * m_m2 - 2.0f * m_m1 + 1.5f * m0; // 111-- backward difference, ε ~ h^2
        }
        else if (!is0(m_p2) && is0(m_m1))
        {
            dmdx = -0.5f * m_p2 + 2.0f * m_p1 - 1.5f * m0; // --111 forward difference,  ε ~ h^2
        }
        else
        {
            dmdx = (2.0f/3.0f)*(m_p1 - m_m1) + (1.0f/12.0f)*(m_m2 - m_p2); // 11111 central difference,  ε ~ h^4
        }
    }

    // y derivatives (along height)
    {
        real_t3 m_m2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-2), iz);
        if ((iy-2 >= 0) || PBCy)
        {
            m_m2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_m1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, lclampy(iy-1), iz);
        if ((iy-1 >= 0) || PBCy)
        {
            m_m1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p1 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+1), iz);
        if  ((iy+1 < Ny) || PBCy)
        {
            m_p1 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        real_t3 m_p2 = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);
        i_ = idx(ix, hclampy(iy+2), iz);
        if  ((iy+2 < Ny) || PBCy)
        {
            m_p2 = make_float3(mx[i_], my[i_], mz[i_]);
        }

        if (is0(m_p1) && is0(m_m1))                                        //  +0
        {
            dmdy = make_float3((real_t)0.0, (real_t)0.0, (real_t)0.0);                          // --1-- zero
        }
        else if ((is0(m_m2) | is0(m_p2)) && !is0(m_p1) && !is0(m_m1))
        {
            dmdy = 0.5f * (m_p1 - m_m1);                                   // -111-, 1111-, -1111 central difference,  ε ~ h^2
        }
        else if (is0(m_p1) && is0(m_m2))
        {
            dmdy =  m0 - m_m1;                                             // -11-- backward difference, ε ~ h^1
        }
        else if (is0(m_m1) && is0(m_p2))
        {
            dmdy = -m0 + m_p1;                                             // --11- forward difference,  ε ~ h^1
        }
        else if (!is0(m_m2) && is0(m_p1))
        {
            dmdy =  0.5f * m_m2 - 2.0f * m_m1 + 1.5f * m0;                 // 111-- backward difference, ε ~ h^2
        }
        else if (!is0(m_p2) && is0(m_m1))
        {
            dmdy = -0.5f * m_p2 + 2.0f * m_p1 - 1.5f * m0;                 // --111 forward difference,  ε ~ h^2
        }
        else
        {
            dmdy = (2.0f/3.0f)*(m_p1 - m_m1) + (1.0f/12.0f)*(m_m2 - m_p2); // 11111 central difference,  ε ~ h^4
        }
    }
    dmdx_x_dmdy = cross(dmdx, dmdy);

    s[I] = icxcy * dot(m0, dmdx_x_dmdy);
}
// Returns the topological charge contribution on an elementary triangle ijk
// Order of arguments is important here to preserve the same measure of chirality
// Note: the result is zero if an argument is zero, or when two arguments are the same
static inline real_t triangleCharge(real_t3 mi, real_t3 mj, real_t3 mk) {
    real_t numer   = dot(mi, cross(mj, mk));
    real_t denom   = 1.0f + dot(mi, mj) + dot(mi, mk) + dot(mj, mk);
    return 2.0f * atan2(numer, denom);
}

// Set s to the toplogogical charge density for lattices based on the solid angle 
// subtended by triangle associated with three spins: a,b,c
//
//       s = 2 atan[(a . b x c /(1 + a.b + a.c + b.c)] / (dx dy)
//
// After M Boettcher et al, New J Phys 20, 103014 (2018), adapted from
// B. Berg and M. Luescher, Nucl. Phys. B 190, 412 (1981), and implemented by
// Joo-Von Kim.
//
// A unit cell comprises two triangles, but s is a site-dependent quantity so we
// double-count and average over four triangles.
__kernel void
settopologicalchargelattice(__global real_t* __restrict     s,
                            __global real_t* __restrict    mx, __global real_t* __restrict my, __global real_t* __restrict mz,
                                                 real_t icxcy,
                                                    int    Nx,                         int Ny,                         int Nz,
                                                uint8_t   PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int     i0 = idx(ix, iy, iz);                     // central cell index
    real_t3 m0 = make_float3(mx[i0], my[i0], mz[i0]); // central cell magnetization

    if (is0(m0)) {
        s[i0] = 0.0f;
        return;
    }

    // indices of the 4 neighbors (counter clockwise)
    int i1 = idx(hclampx(ix+1), iy, iz); // (i+1,j)
    int i2 = idx(ix, hclampy(iy+1), iz); // (i,j+1)
    int i3 = idx(lclampx(ix-1), iy, iz); // (i-1,j)
    int i4 = idx(ix, lclampy(iy-1), iz); // (i,j-1)

    // magnetization of the 4 neighbors
    real_t3 m1 = make_float3(mx[i1], my[i1], mz[i1]);
    real_t3 m2 = make_float3(mx[i2], my[i2], mz[i2]);
    real_t3 m3 = make_float3(mx[i3], my[i3], mz[i3]);
    real_t3 m4 = make_float3(mx[i4], my[i4], mz[i4]);

    // local topological charge (accumulator)
    real_t topcharge = 0.0f;

    // charge contribution from the upper right triangle
    // if diagonally opposite neighbor is not zero, use a weight of 1/2 to avoid counting charges twice
    if (((ix+1<Nx) || PBCx) && ((iy+1<Ny) || PBCy)) { 
        int         i_ = idx(hclampx(ix+1), hclampy(iy+1), iz); // diagonal opposite neighbor in upper right quadrant
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m1, m2);
    }

    // upper left
    if (((ix-1>=0) || PBCx) && ((iy+1<Ny) || PBCy)) { 
        int         i_ = idx(lclampx(ix-1), hclampy(iy+1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m2, m3);
    }

    // bottom left
    if (((ix-1>=0) || PBCx) && ((iy-1>=0) || PBCy)) { 
        int         i_ = idx(lclampx(ix-1), lclampy(iy-1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m3, m4);
    }

    // bottom right
    if (((ix+1<Nx) || PBCx) && ((iy-1>=0) || PBCy)) { 
        int         i_ = idx(hclampx(ix+1), lclampy(iy-1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m4, m1);
    }

    s[i0] = icxcy * topcharge;
}
// Add uniaxial magnetocrystalline anisotropy field to B.
// http://www.southampton.ac.uk/~fangohr/software/oxs_uniaxial4.html
__kernel void
adduniaxialanisotropy(__global real_t* __restrict  Bx, __global real_t* __restrict     By, __global real_t* __restrict  Bz,
                      __global real_t* __restrict  mx, __global real_t* __restrict     my, __global real_t* __restrict  mz,
                      __global real_t* __restrict Ms_,                      real_t Ms_mul,
                      __global real_t* __restrict K1_,                      real_t K1_mul,
                      __global real_t* __restrict ux_,                      real_t ux_mul,
                      __global real_t* __restrict uy_,                      real_t uy_mul,
                      __global real_t* __restrict uz_,                      real_t uz_mul,
                                              int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3     u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        real_t  invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t     K1 = amul(K1_, K1_mul, i);

        K1  *= invMs;

        real_t3  m = {mx[i], my[i], mz[i]};
        real_t  mu = dot(m, u);
        real_t3 Ba = (real_t)2.0*K1*(mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
// Add uniaxial magnetocrystalline anisotropy field to B.
// http://www.southampton.ac.uk/~fangohr/software/oxs_uniaxial4.html
__kernel void
adduniaxialanisotropy2(__global real_t* __restrict  Bx, __global real_t* __restrict     By, __global real_t* __restrict  Bz,
                       __global real_t* __restrict  mx, __global real_t* __restrict     my, __global real_t* __restrict  mz,
                       __global real_t* __restrict Ms_,                      real_t Ms_mul,
                       __global real_t* __restrict K1_,                      real_t K1_mul,
                       __global real_t* __restrict K2_,                      real_t K2_mul,
                       __global real_t* __restrict ux_,                      real_t ux_mul,
                       __global real_t* __restrict uy_,                      real_t uy_mul,
                       __global real_t* __restrict uz_,                      real_t uz_mul,
                                               int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3     u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        real_t  invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t     K1 = amul(K1_, K1_mul, i);
        real_t     K2 = amul(K2_, K2_mul, i);

        K1  *= invMs;
        K2  *= invMs;

        real_t3 m  = {mx[i], my[i], mz[i]};
        real_t  mu = dot(m, u);
        real_t3 Ba = (real_t)2.0*K1*    (mu)*u+
                     (real_t)4.0*K2*pow3(mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
// Add voltage-controlled magnetic anisotropy field to B.
// https://www.nature.com/articles/s42005-019-0189-6.pdf
__kernel void
addvoltagecontrolledanisotropy2(__global real_t* __restrict         Bx, __global real_t* __restrict            By, __global real_t* __restrict Bz,
                                __global real_t* __restrict         mx, __global real_t* __restrict            my, __global real_t* __restrict mz,
                                __global real_t* __restrict        Ms_,                      real_t        Ms_mul,
                                __global real_t* __restrict vcmaCoeff_,                      real_t vcmaCoeff_mul,
                                __global real_t* __restrict   voltage_,                      real_t   voltage_mul,
                                __global real_t* __restrict        ux_,                      real_t        ux_mul,
                                __global real_t* __restrict        uy_,                      real_t        uy_mul,
                                __global real_t* __restrict        uz_,                      real_t        uz_mul,
                                                        int          N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3         u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        real_t      invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t  vcmaCoeff = amul(vcmaCoeff_, vcmaCoeff_mul, i) * invMs;
        real_t    voltage = amul(voltage_, voltage_mul, i) * invMs;
        real_t3         m = {mx[i], my[i], mz[i]};
        real_t         mu = dot(m, u);
        real_t3        Ba = (real_t)2.0*vcmaCoeff*voltage*    (mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
__kernel void
vecnorm(__global real_t* __restrict dst,
        __global real_t* __restrict  ax, __global real_t* __restrict ay, __global real_t* __restrict az,
                                int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t3 A = {ax[i], ay[i], az[i]};
        dst[i] = sqrt(dot(A, A));
    }
}
// set dst to zero in cells where mask != 0
__kernel void
zeromask(__global real_t* __restrict dst, __global real_t* maskLUT, __global uint8_t* regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        if (maskLUT[regions[i]] != 0){
            dst[i] = (real_t)0.0;
        }
    }
}
#define PREFACTOR ((MUB) / (2 * QE * GAMMA0))

// spatial derivatives without dividing by cell size
#define deltax(in) (in[idx(hclampx(ix+1), iy, iz)] - in[idx(lclampx(ix-1), iy, iz)])
#define deltay(in) (in[idx(ix, hclampy(iy+1), iz)] - in[idx(ix, lclampy(iy-1), iz)])
#define deltaz(in) (in[idx(ix, iy, hclampz(iz+1))] - in[idx(ix, iy, lclampz(iz-1))])

__kernel void
addzhanglitorque2(__global real_t* __restrict     tx, __global real_t* __restrict        ty, __global real_t* __restrict tz,
                  __global real_t* __restrict     mx, __global real_t* __restrict        my, __global real_t* __restrict mz,
                  __global real_t* __restrict    Ms_,                      real_t    Ms_mul,
                  __global real_t* __restrict    jx_,                      real_t    jx_mul,
                  __global real_t* __restrict    jy_,                      real_t    jy_mul,
                  __global real_t* __restrict    jz_,                      real_t    jz_mul,
                  __global real_t* __restrict alpha_,                      real_t alpha_mul,
                  __global real_t* __restrict    xi_,                      real_t    xi_mul,
                  __global real_t* __restrict   pol_,                      real_t   pol_mul,
                                       real_t     cx,                      real_t        cy,                      real_t cz,
                                          int     Nx,                        int        Ny,                        int Nz,
                                      uint8_t    PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int i = idx(ix, iy, iz);

    real_t  alpha = amul(alpha_, alpha_mul, i);
    real_t     xi = amul(xi_, xi_mul, i);
    real_t    pol = amul(pol_, pol_mul, i);
    real_t  invMs = inv_Msat(Ms_, Ms_mul, i);
    real_t      b = invMs * PREFACTOR / (1.0f + xi*xi);
    real_t3  Jvec = vmul(jx_, jy_, jz_, jx_mul, jy_mul, jz_mul, i);
    real_t3     J = pol*Jvec;
    real_t3 hspin = make_float3(0.0f, 0.0f, 0.0f); // (u·∇)m

    if (J.x != (real_t)0.0) {
        hspin += (b/cx)*J.x * make_float3(deltax(mx), deltax(my), deltax(mz));
    }
    if (J.y != (real_t)0.0) {
        hspin += (b/cy)*J.y * make_float3(deltay(mx), deltay(my), deltay(mz));
    }
    if (J.z != (real_t)0.0) {
        hspin += (b/cz)*J.z * make_float3(deltaz(mx), deltaz(my), deltaz(mz));
    }

    real_t3      m = make_float3(mx[i], my[i], mz[i]);
    real_t3 torque = ((real_t)-1.0/((real_t)1.0 + alpha*alpha)) * (
                         ((real_t)1.0+xi*alpha) * cross(m, cross(m, hspin))
                                  +(  xi-alpha) * cross(m, hspin)           );

    // write back, adding to torque
    tx[i] += torque.x;
    ty[i] += torque.y;
    tz[i] += torque.z;
}
/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
State buffer stores 6*N uint, where N is the total number of RNGs
The first contiguous N entries are the x[0] of the xorwow states
The second contiguous N entries are the x[1] of the xorwow states
The third contiguous N entries are the x[2] of the xorwow states
The fourth contiguous N entries are the x[3] of the xorwow states
The fifth contiguous N entries are the x[4] of the xorwow states
The sixth contiguous N entries are the d (Weyl sequence number) of the xorwow states
*/

/**
Seeds xorwow RNG.

@param state_buf Variable, that holds state of the generator to be seeded.
@param seed Value used for seeding. Should be randomly generated for each instance of generator (thread).
*/
__kernel void
xorwow_seed(
    __global uint* __restrict       state_buf,
    __global uint* __restrict g_jump_matrices,
    ulong seed) {
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Using local registers to compute state from seed	
    uint x[XORWOW_N];
    uint d;

    // Initialize state buffer and Weyl sequence number
    x[0] = 123456789U;
    x[1] = 362436069U;
    x[2] = 521288629U;
    x[3] = 88675123U;
    x[4] = 5783321U;
    d = 6615241U;

    // Update RNG state with seed value
    // Constants are arbitrary prime numbers
    const uint s0 = (uint)(seed) ^ 0x2c7f967fU;
    const uint s1 = (uint)(seed >> 32) ^ 0xa03697cbU;
    const uint t0 = 1228688033U * s0;
    const uint t1 = 2073658381U * s1;
    x[0] += t0;
    x[1] ^= t0;
    x[2] += t1;
    x[3] ^= t1;
    x[4] += t0;
    d += t1 + t0;

    // discarding subsequences to obtain non-overlapping random bit streams via parallelism...
    if (global_idx != 0) {
        xorwow_discard_subsequence(global_idx, x, g_jump_matrices);
    }

    // Write out state to global buffer
    int idx = global_idx;
    state_buf[idx] = x[0];
    idx += grp_offset;
    state_buf[idx] = x[1];
    idx += grp_offset;
    state_buf[idx] = x[2];
    idx += grp_offset;
    state_buf[idx] = x[3];
    idx += grp_offset;
    state_buf[idx] = x[4];
    idx += grp_offset;
    state_buf[idx] = d;
}
/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
Generates a random 32-bit unsigned integer using xorwow RNG.

@param state_buf State of the RNG to use.
@param d_data Output.
*/
__kernel void
xorwow_uint(
    __global uint* __restrict state_buf,
    __global uint* __restrict   d_data,
    int count){
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Only threads witin the count bounds will generate the random number
    if (global_idx < count) {
        // Using local registers to compute next state
        uint x[5];
        uint d;

        // Get state from global buffer
        int idx = global_idx;
        x[0] = state_buf[idx];
        idx += grp_offset;
        x[1] = state_buf[idx];
        idx += grp_offset;
        x[2] = state_buf[idx];
        idx += grp_offset;
        x[3] = state_buf[idx];
        idx += grp_offset;
        x[4] = state_buf[idx];
        idx += grp_offset;
        d = state_buf[idx];

        // For each thread that is launched, iterate until the index is out of bounds
        for (uint pos = global_idx; pos < count; pos += grp_offset) {
            const uint t = x[0] ^ (x[0] >> 2);
            x[0] = x[1];
            x[1] = x[2];
            x[2] = x[3];
            x[3] = x[4];
            x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

            d += 362437;

            d_data[pos] = (d + x[4]); // output value
        }

        // update the state buffer with the latest state
        idx = global_idx;
        state_buf[idx] = x[0];
        idx += grp_offset;
        state_buf[idx] = x[1];
        idx += grp_offset;
        state_buf[idx] = x[2];
        idx += grp_offset;
        state_buf[idx] = x[3];
        idx += grp_offset;
        state_buf[idx] = x[4];
        idx += grp_offset;
        state_buf[idx] = d;
    }
}
/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
Generates a random uniformly distributed float using xorwow RNG.

@param state State of the RNG to use.
@param d_data Output.
*/
__kernel void
xorwow_uniform(
    __global  uint* __restrict state_buf,
    __global float* __restrict    d_data,
    int count){
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Only threads witin the count bounds will generate the random number
    if (global_idx < count) {
        // Using local registers to compute next state
        uint x[5];
        uint d;

        // Get state from global buffer
        int idx = global_idx;
        x[0] = state_buf[idx];
        idx += grp_offset;
        x[1] = state_buf[idx];
        idx += grp_offset;
        x[2] = state_buf[idx];
        idx += grp_offset;
        x[3] = state_buf[idx];
        idx += grp_offset;
        x[4] = state_buf[idx];
        idx += grp_offset;
        d = state_buf[idx];

        // For each thread that is launched, iterate until the index is out of bounds
        for (uint pos = global_idx; pos < count; pos += grp_offset) {
            // generate a pair of uint32 (one uint64)
            // first number...
            uint t = x[0] ^ (x[0] >> 2);
            x[0] = x[1];
            x[1] = x[2];
            x[2] = x[3];
            x[3] = x[4];
            x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

            d += 362437;

            uint num1 = d+x[4];

            // second number...
            t = x[0] ^ (x[0] >> 2);
            x[0] = x[1];
            x[1] = x[2];
            x[2] = x[3];
            x[3] = x[4];
            x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

            d += 362437;

            uint num2 = d+x[4];

            d_data[pos] = uint2float(num1, num2); // output value
        }

        // update the state buffer with the latest state
        idx = global_idx;
        state_buf[idx] = x[0];
        idx += grp_offset;
        state_buf[idx] = x[1];
        idx += grp_offset;
        state_buf[idx] = x[2];
        idx += grp_offset;
        state_buf[idx] = x[3];
        idx += grp_offset;
        state_buf[idx] = x[4];
        idx += grp_offset;
        state_buf[idx] = d;
    }
}
/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
Generates a random normally distributed float using xorwow RNG.

@param state State of the RNG to use.
@param d_data Output.
*/
__kernel void
xorwow_normal(
    __global  uint* __restrict state_buf,
    __global float* __restrict    d_data,
    int count){
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Only threads witin the count bounds will generate the random number
    if (global_idx < count) {
        // Using local registers to compute next state
        uint x[5];
        uint d;

        // Get state from global buffer
        int idx = global_idx;
        x[0] = state_buf[idx];
        idx += grp_offset;
        x[1] = state_buf[idx];
        idx += grp_offset;
        x[2] = state_buf[idx];
        idx += grp_offset;
        x[3] = state_buf[idx];
        idx += grp_offset;
        x[4] = state_buf[idx];
        idx += grp_offset;
        d = state_buf[idx];
        bool generate = true;
        float z0 = 0.0f;
        float z1 = 0.0f;

        // For each thread that is launched, iterate until the index is out of bounds
        for (uint pos = global_idx; pos < count; pos += grp_offset) {
            if (generate) {
                // generate a pair of uint32 (one uint64)
                // first number...
                uint t = x[0] ^ (x[0] >> 2);
                x[0] = x[1];
                x[1] = x[2];
                x[2] = x[3];
                x[3] = x[4];
                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                d += 362437;

                uint num1 = d+x[4];

                // second number...
                t = x[0] ^ (x[0] >> 2);
                x[0] = x[1];
                x[1] = x[2];
                x[2] = x[3];
                x[3] = x[4];
                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                d += 362437;

                uint num2 = d+x[4];

                float tmpRes1 = uint2float(num1, num2); // output value

                // Repeat for second float...
                // generate a pair of uint32 (one uint64)
                // first number...
                t = x[0] ^ (x[0] >> 2);
                x[0] = x[1];
                x[1] = x[2];
                x[2] = x[3];
                x[3] = x[4];
                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                d += 362437;

                num1 = d+x[4];

                // second number...
                t = x[0] ^ (x[0] >> 2);
                x[0] = x[1];
                x[1] = x[2];
                x[2] = x[3];
                x[3] = x[4];
                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                d += 362437;

                num2 = d+x[4];

                float tmpRes2 = uint2float(num1, num2); // output value

                z0 = sqrt( -2.0f * log(tmpRes1)) * cospi(2.0f * tmpRes2);
                z1 = sqrt( -2.0f * log(tmpRes1)) * sinpi(2.0f * tmpRes2);
                d_data[pos] = z0; // output normal random value
                generate = !generate;
            } else {
                d_data[pos] = z1; // output normal random value
            }
        }

        // update the state buffer with the latest state
        idx = global_idx;
        state_buf[idx] = x[0];
        idx += grp_offset;
        state_buf[idx] = x[1];
        idx += grp_offset;
        state_buf[idx] = x[2];
        idx += grp_offset;
        state_buf[idx] = x[3];
        idx += grp_offset;
        state_buf[idx] = x[4];
        idx += grp_offset;
        state_buf[idx] = d;
    }
}
/**
@file
Implements threefry RNG.
*******************************************************
 * Modified version of Random123 library:
 * https://www.deshawresearch.com/downloads/download_random123.cgi/
 * The original copyright can be seen here:
 *
 * RANDOM123 LICENSE AGREEMENT
 *
 * Copyright 2010-2011, D. E. Shaw Research. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * * Redistributions of source code must retain the above copyright notice,
 *   this list of conditions, and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions, and the following disclaimer in the
 *   documentation and/or other materials provided with the distribution.
 *
 * Neither the name of D. E. Shaw Research nor the names of its contributors
 * may be used to endorse or promote products derived from this software
 * without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *********************************************************/

/**
State of threefry RNG. We will store in global buffer as a set of uint
**
counter: uint[4]
key:     uint[4]
state:   uint[4]
index:   uint
typedef struct{
        uint counter[4];
        uint result[4];
        uint key[4];
        uint tracker;
} threefry_state;
**/

/**
Seeds threefry RNG.
@param state Variable, that holds state of the generator to be seeded.
@param seed Value used for seeding. Should be randomly generated for each instance of generator (thread).
**/
__kernel void
threefry_seed(
    __global uint* __restrict     state_key,
    __global uint* __restrict state_counter,
    __global uint* __restrict  state_result,
    __global uint* __restrict state_tracker,
    __global uint* __restrict          seed) {
    uint gid = get_global_id(0);
    uint rng_count = get_global_size(0);
    uint idx = gid;
    uint localJ = seed[gid];
    state_key[idx] = localJ;
    state_counter[idx] = 0x00000000;
    state_result[idx] = 0x00000000;
    state_tracker[idx] = 4;
    idx += rng_count;
    state_key[idx] = 0x00000000;
    state_counter[idx] = 0x00000000;
    state_result[idx] = 0x00000000;
    idx += rng_count;
    state_key[idx] = gid;
    state_counter[idx] = 0x00000000;
    state_result[idx] = 0x00000000;
    idx += rng_count;
    state_key[idx] = 0x00000000;
    state_counter[idx] = 0x00000000;
    state_result[idx] = 0x00000000;
}
/**
@file
Implements threefry RNG.
*******************************************************
 * Modified version of Random123 library:
 * https://www.deshawresearch.com/downloads/download_random123.cgi/
 * The original copyright can be seen here:
 *
 * RANDOM123 LICENSE AGREEMENT
 *
 * Copyright 2010-2011, D. E. Shaw Research. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * * Redistributions of source code must retain the above copyright notice,
 *   this list of conditions, and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions, and the following disclaimer in the
 *   documentation and/or other materials provided with the distribution.
 *
 * Neither the name of D. E. Shaw Research nor the names of its contributors
 * may be used to endorse or promote products derived from this software
 * without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *********************************************************/

/**
State of threefry RNG. We will store in global buffer as a set of uint
**
typedef struct{
        uint counter[4];
        uint result[4];
        uint key[4];
        uint tracker;
} threefry_state;
**/

/**
Generates a random 32-bit unsigned integer using threefry RNG.
@param state State of the RNG to use.
**/
__kernel void
threefry_uint(
    __global uint* __restrict     state_key,
    __global uint* __restrict state_counter,
    __global uint* __restrict  state_result,
    __global uint* __restrict state_tracker,
    __global uint* __restrict        output,
    int data_size) {
    uint gid = get_global_id(0);
    uint rng_count = get_global_size(0);
    uint tmpIdx = gid;
    threefry_state state_;
    threefry_state *state = &state_;

    // For first out of four sets...
    // Read in counter
    state->counter[0] = state_counter[tmpIdx];
    // Read in result
    state->result[0] = state_result[tmpIdx];
    // Read in key
    state->key[0] = state_key[tmpIdx];
    // Read in tracker
    state->tracker = state_tracker[tmpIdx];

    // For second out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[1] = state_counter[tmpIdx];
    // Read in result
    state->result[1] = state_result[tmpIdx];
    // Read in key
    state->key[1] = state_key[tmpIdx];

    // For third out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[2] = state_counter[tmpIdx];
    // Read in result
    state->result[2] = state_result[tmpIdx];
    // Read in key
    state->key[2] = state_key[tmpIdx];

    // For last out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[3] = state_counter[tmpIdx];
    // Read in result
    state->result[3] = state_result[tmpIdx];
    // Read in key
    state->key[3] = state_key[tmpIdx];

    for (uint outIndex = gid; outIndex < data_size; outIndex += rng_count) {
        if (state->tracker > 3) {
            threefry_round(state);
            state->tracker = 1;
            output[outIndex] = state->result[0];
        } else if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            output[outIndex] = tmp;
        } else {
            output[outIndex] = state->result[state->tracker++];
        }
    }
    
    // For first out of four sets...
    // Write out counter
    tmpIdx = gid;
    state_counter[tmpIdx] = state->counter[0];
    // Write out result
    state_result[tmpIdx] = state->result[0];
    // Write out key
    state_key[tmpIdx] = state->key[0];
    // Write out tracker
    state_tracker[tmpIdx] = state->tracker;

    // For second out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[1];
    // Write out result
    state_result[tmpIdx] = state->result[1];
    // Write out key
    state_key[tmpIdx] = state->key[1];

    // For third out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[2];
    // Write out result
    state_result[tmpIdx] = state->result[2];
    // Write out key
    state_key[tmpIdx] = state->key[2];

    // For last out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[3];
    // Write out result
    state_result[tmpIdx] = state->result[3];
    // Write out key
    state_key[tmpIdx] = state->key[3];

}
/**
@file
Implements threefry RNG.
*******************************************************
 * Modified version of Random123 library:
 * https://www.deshawresearch.com/downloads/download_random123.cgi/
 * The original copyright can be seen here:
 *
 * RANDOM123 LICENSE AGREEMENT
 *
 * Copyright 2010-2011, D. E. Shaw Research. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * * Redistributions of source code must retain the above copyright notice,
 *   this list of conditions, and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions, and the following disclaimer in the
 *   documentation and/or other materials provided with the distribution.
 *
 * Neither the name of D. E. Shaw Research nor the names of its contributors
 * may be used to endorse or promote products derived from this software
 * without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *********************************************************/

/**
State of threefry RNG. We will store in global buffer as a set of uint
**
typedef struct{
        uint counter[4];
        uint result[4];
        uint key[4];
        uint tracker;
} threefry_state;
**/

/**
Generates a random uniformly distributed float using threefry RNG.

@param state State of the RNG to use.
**/
__kernel void
threefry_uniform(
    __global  uint* __restrict     state_key,
    __global  uint* __restrict state_counter,
    __global  uint* __restrict  state_result,
    __global  uint* __restrict state_tracker,
    __global float* __restrict        output,
    int data_size) {
    uint gid = get_global_id(0);
    uint rng_count = get_global_size(0);
    uint tmpIdx = gid;
    threefry_state state_;
    threefry_state *state = &state_;

    // For first out of four sets...
    // Read in counter
    state->counter[0] = state_counter[tmpIdx];
    // Read in result
    state->result[0] = state_result[tmpIdx];
    // Read in key
    state->key[0] = state_key[tmpIdx];
    // Read in tracker
    state->tracker = state_tracker[tmpIdx];

    // For second out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[1] = state_counter[tmpIdx];
    // Read in result
    state->result[1] = state_result[tmpIdx];
    // Read in key
    state->key[1] = state_key[tmpIdx];

    // For third out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[2] = state_counter[tmpIdx];
    // Read in result
    state->result[2] = state_result[tmpIdx];
    // Read in key
    state->key[2] = state_key[tmpIdx];

    // For last out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[3] = state_counter[tmpIdx];
    // Read in result
    state->result[3] = state_result[tmpIdx];
    // Read in key
    state->key[3] = state_key[tmpIdx];

    for (uint outIndex = gid; outIndex < data_size; outIndex += rng_count) {
        uint num1[2];
        uint lidx = 0;
        if (state->tracker > 3) {
            threefry_round(state);
            state->tracker = 1;
            num1[lidx++] = state->result[0];
        } else if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx++] = tmp;
        } else {
            num1[lidx++] = state->result[state->tracker++];
        }
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx] = tmp;
        } else {
            num1[lidx] = state->result[state->tracker++];
        }
        output[outIndex] = uint2float(num1[0], num1[1]);
    }
    
    // For first out of four sets...
    // Write out counter
    tmpIdx = gid;
    state_counter[tmpIdx] = state->counter[0];
    // Write out result
    state_result[tmpIdx] = state->result[0];
    // Write out key
    state_key[tmpIdx] = state->key[0];
    // Write out tracker
    state_tracker[tmpIdx] = state->tracker;

    // For second out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[1];
    // Write out result
    state_result[tmpIdx] = state->result[1];
    // Write out key
    state_key[tmpIdx] = state->key[1];

    // For third out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[2];
    // Write out result
    state_result[tmpIdx] = state->result[2];
    // Write out key
    state_key[tmpIdx] = state->key[2];

    // For last out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[3];
    // Write out result
    state_result[tmpIdx] = state->result[3];
    // Write out key
    state_key[tmpIdx] = state->key[3];

}
/**
@file
Implements threefry RNG.
*******************************************************
 * Modified version of Random123 library:
 * https://www.deshawresearch.com/downloads/download_random123.cgi/
 * The original copyright can be seen here:
 *
 * RANDOM123 LICENSE AGREEMENT
 *
 * Copyright 2010-2011, D. E. Shaw Research. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * * Redistributions of source code must retain the above copyright notice,
 *   this list of conditions, and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions, and the following disclaimer in the
 *   documentation and/or other materials provided with the distribution.
 *
 * Neither the name of D. E. Shaw Research nor the names of its contributors
 * may be used to endorse or promote products derived from this software
 * without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *********************************************************/

/**
State of threefry RNG. We will store in global buffer as a set of uint
**
typedef struct{
        uint counter[4];
        uint result[4];
        uint key[4];
        uint tracker;
} threefry_state;
**/

/**
Generates a random normally distributed float using threefry RNG.

@param state State of the RNG to use.
**/
__kernel void
threefry_normal(
    __global  uint* __restrict     state_key,
    __global  uint* __restrict state_counter,
    __global  uint* __restrict  state_result,
    __global  uint* __restrict state_tracker,
    __global float* __restrict        output,
    int data_size) {
    uint gid = get_global_id(0);
    uint rng_count = get_global_size(0);
    uint tmpIdx = gid;
    threefry_state state_;
    threefry_state *state = &state_;

    // For first out of four sets...
    // Read in counter
    state->counter[0] = state_counter[tmpIdx];
    // Read in result
    state->result[0] = state_result[tmpIdx];
    // Read in key
    state->key[0] = state_key[tmpIdx];
    // Read in tracker
    state->tracker = state_tracker[tmpIdx];

    // For second out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[1] = state_counter[tmpIdx];
    // Read in result
    state->result[1] = state_result[tmpIdx];
    // Read in key
    state->key[1] = state_key[tmpIdx];

    // For third out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[2] = state_counter[tmpIdx];
    // Read in result
    state->result[2] = state_result[tmpIdx];
    // Read in key
    state->key[2] = state_key[tmpIdx];

    // For last out of four sets...
    tmpIdx += rng_count;
    // Read in counter
    state->counter[3] = state_counter[tmpIdx];
    // Read in result
    state->result[3] = state_result[tmpIdx];
    // Read in key
    state->key[3] = state_key[tmpIdx];

    for (uint outIndex = gid; outIndex < data_size / 2; outIndex += rng_count) {
        uint num1[2];
        float res1[4];
        uint lidx = 0;
        if (state->tracker > 3) {
            threefry_round(state);
            state->tracker = 1;
            num1[lidx++] = state->result[0];
        } else if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx++] = tmp;
        } else {
            num1[lidx++] = state->result[state->tracker++];
        }
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx] = tmp;
        } else {
            num1[lidx] = state->result[state->tracker++];
        }
        res1[0] = uint2float(num1[0], num1[1]);
        lidx = 0;
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx++] = tmp;
        } else {
            num1[lidx++] = state->result[state->tracker++];
        }
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx] = tmp;
        } else {
            num1[lidx] = state->result[state->tracker++];
        }
        res1[1] = uint2float(num1[0], num1[1]);
        res1[2] = sqrt( -2.0f * log(res1[0])) * cospi(2.0f * res1[1]);
        res1[3] = sqrt( -2.0f * log(res1[0])) * sinpi(2.0f * res1[1]);
        output[outIndex] = res1[2];
        output[outIndex + (data_size/2)] = res1[3];
    }
    
    // For first out of four sets...
    // Write out counter
    tmpIdx = gid;
    state_counter[tmpIdx] = state->counter[0];
    // Write out result
    state_result[tmpIdx] = state->result[0];
    // Write out key
    state_key[tmpIdx] = state->key[0];
    // Write out tracker
    state_tracker[tmpIdx] = state->tracker;

    // For second out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[1];
    // Write out result
    state_result[tmpIdx] = state->result[1];
    // Write out key
    state_key[tmpIdx] = state->key[1];

    // For third out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[2];
    // Write out result
    state_result[tmpIdx] = state->result[2];
    // Write out key
    state_key[tmpIdx] = state->key[2];

    // For last out of four sets...
    // Write out counter
    tmpIdx += rng_count;
    state_counter[tmpIdx] = state->counter[3];
    // Write out result
    state_result[tmpIdx] = state->result[3];
    // Write out key
    state_key[tmpIdx] = state->key[3];

}


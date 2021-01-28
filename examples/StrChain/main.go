package main

import (
	"fmt"
	"unicode"

	"github.com/davecgh/go-spew/spew"

	"github.com/shyang107/paw"
)

func main() {

	str := paw.NewStrChain("abcde Abcde 12中文3 4.?!abcde")
	t := str
	fmt.Println("t=", t)
	fmt.Printf("[]byte = %v, len = %d\n", t.Bytes(), len(t.Bytes()))
	fmt.Printf("[]rune = %v, len = %d\n", t.Runes(), len(t.Runes()))

	s := str
	fmt.Printf("1. GetAbbrString(10,») = %s\n", s.GetAbbrString(10, "»"))

	s = str
	fmt.Println("2. StringWidth() = ", s.StringWidth())

	s = str
	h, a := s.CountPlaceHolder()
	fmt.Println("3. CountPlaceHolder() =", h, a)

	s = str
	fmt.Println("4. HasChineseChar() =", s.HasChineseChar())

	fmt.Println("                    t =", t)
	s = str
	fmt.Println("5. NumberBanner()     =", s.NumberBanner())

	// s = str
	// fmt.Println("4.1. NumberBannerRune() =", s.NumberBannerRune())

	s = str
	fmt.Println("6. Reverse()          =", s.Reverse())

	s = str
	fmt.Println("7. HasPrefix(abc)     =", s.HasPrefix("abc"))

	s = str
	fmt.Println("8. HasSuffix(abc)     =", s.HasSuffix("abc"))

	s = str
	fmt.Println("9. Contains(中)       =", s.Contains("中"))

	s = str
	fmt.Println("10.1. ContainsAny(中) =", s.ContainsAny("中"))
	fmt.Println("10.2. ContainsAny(90) =", s.ContainsAny("90"))

	s = str
	fmt.Println("11. Fields()          =", s.Fields())
	spew.Dump(s.Fields())

	s = str
	f := func(r rune) bool {
		if r == '中' {
			return true
		}
		return false
	}
	fmt.Println("12. FieldsFunc(f)          =", s.FieldsFunc(f))
	fmt.Println(`   f := func(r rune) bool {
        if r == '中' {
            return true
        }
        return false
        }`)
	spew.Dump(s.FieldsFunc(f))

	s = str
	fmt.Println("13. ContainsRune('文') =", s.ContainsRune('文'))

	t1 := "ABcde"
	t2 := "abCDE"
	fmt.Println("14. t1 =", t1, "\n    t2 =", t2)
	s1 := paw.NewStrChain(t1)
	fmt.Println("    t1.EqualFold(t2) =", s1.EqualFold(t2))

	fmt.Println("\nt=", t)
	s = str
	fmt.Println("15. Index(\"c\") =", s.Index("c"))
	s = str
	fmt.Println("16.1. IndexAny(\"cd\") =", s.IndexAny("cd"))
	fmt.Println("16.2. IndexAny(\"dc\") =", s.IndexAny("dc"))
	s = str
	fmt.Println("17. IndexByte(byte(\"c\") =", s.IndexByte(byte('c')))

	s = str
	f = func(c rune) bool {
		return unicode.Is(unicode.Han, c)
	}
	fmt.Println(`    f = func(c rune) bool {
        return unicode.Is(unicode.Han, c)
        }`)
	fmt.Println("18. IndexFunc(f) =", s.IndexFunc(f))

	s = str
	fmt.Println("19. IndexRune('中') =", s.IndexRune('中'))

	s = str
	fmt.Println("20. LastIndex(\"34.\") =", s.LastIndex("34."), " here use the numbers of placeholders")
	s = str
	fmt.Println("21. LastIndexAny(\"12?\") =", s.LastIndexAny("12?"), " here use len of bytes not runewidth")

	s = str
	fmt.Println("22. LastIndexByte('c')  =", s.LastIndexByte(byte('c')))
	fmt.Println("23. LastIndexFunc(f)    =", s.LastIndexFunc(f), " here use len of bytes not runewidth")
	fmt.Printf("24. Split(\" \") = %#v, sizes: %d\n", s.Split(" "), len(s.Split(" ")))
	fmt.Printf("25. SplitN(\" \",2) = %#v, sizes: %d\n", s.SplitN(" ", 2), len(s.SplitN(" ", 2)))
	fmt.Printf("26. SplitAfter(\" \") = %#v, sizes: %d\n", s.SplitAfter(" "), len(s.SplitAfter(" ")))

	s = str
	fmt.Println(s, "len =", s.Len(), "byteWidth =", len(s.Bytes()), "runeLen =", len(s.Runes()), "StringWidth =", s.StringWidth())
	fmt.Println("", s.NumberBanner())
	s = str
	fmt.Printf("27. Trim(\"abcde\")  = %q\n", s.Trim("abcde"))

	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Println("\n", s, "len =", s.Len(), "byteWidth =", len(s.Bytes()), "runeLen =", len(s.Runes()), "StringWidth =", s.StringWidth())
	fmt.Println("", s.NumberBanner())
	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	f = func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }
	fmt.Println(`f = func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r)}`)
	fmt.Printf("28. TrimFunc(f)  = %q\n", s.TrimFunc(f))

	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("29. TrimLeft(\"¡¡¡\")  = %q\n", s.TrimLeft("¡¡¡"))

	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("30. TrimLeftFunc(f)  = %q\n", s.TrimLeftFunc(f))

	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("31. TrimPrefix(\"¡¡\")  = %q\n", s.TrimPrefix("¡¡"))

	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("32. TrimRight(\"!!\")  = %q\n", s.TrimRight("!!"))
	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("33. TrimRightFunc(f)  = %q\n", s.TrimRightFunc(f))
	s = paw.NewStrChain("¡¡¡" + str.String() + "!!!")
	fmt.Printf("34. TrimSuffix(\"!!\")  = %q\n", s.TrimSuffix("!!"))

	s = paw.NewStrChain("\r ¡¡¡" + str.String() + "!!!\r ")
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("35. TrimSpace()  = %q\n", s.TrimSpace())

	s = str
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("36. ToUpper()  = %q\n", s.ToUpper())
	s = str
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("37. ToTitle()  = %q\n", s.ToTitle())
	s = str
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("38. ToLower()  = %q\n", s.ToLower())
	s = str
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("39. Title()  = %q\n", s.Title())

	s = str
	fmt.Printf("\n%#v\n", s)
	m := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Println(`map := func(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return 'A' + (r-'A'+13)%26
	case r >= 'a' && r <= 'z':
		return 'a' + (r-'a'+13)%26
	}
	return r
}`)
	fmt.Printf("40. Map(m)  = %q\n", s.Map(m))

	s = str
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("41. Repeat(2)  = %q\n", s.Repeat(2))

	s = str
	fmt.Printf("\n%#v\n", s)
	gs, _ := s.Utf8ToGbkString()
	fmt.Printf("41. Utf8ToGbkString()  = %q\n", gs)
	us, _ := gs.GbkToUtf8String()
	fmt.Printf("42. GbkToUtf8String()  = %q\n", us)

	s = str
	fmt.Printf("\n%#v\n", s)
	bs, _ := s.Utf8ToBig5String()
	fmt.Printf("41. Utf8ToBig5String()  = %q\n", bs)
	us, _ = bs.Big5ToUtf8String()
	fmt.Printf("42. Big5ToUtf8String()  = %q\n", us)

	s = paw.NewStrChain("abcDE")
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("43. IsEqualString(\"ABCDE\",true)  = %v\n", s.IsEqualString("ABCDE", true))
	fmt.Printf("44. IsEqualString(\"ABCDE\",false)  = %v\n", s.IsEqualString("ABCDE", false))

	s = paw.NewStrChain("abcDE")
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("43. FillLeft(10)  = %q\n", s.FillLeft(10))
	fmt.Printf("                     %s\n", s.NumberBanner())
	fmt.Printf("   StringWidth()  = %v\n", s.StringWidth())
	s = paw.NewStrChain("abcDE")
	fmt.Printf("44. FillRight(10)  = %q\n", s.FillRight(10))
	fmt.Printf("                      %s\n", s.NumberBanner())
	fmt.Printf("   StringWidth()  = %v\n", s.StringWidth())

	s = paw.NewStrChain("abc中defghijklmnopqrstuvwxyz")
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("45. Truncate(10,\"-»\")  = %q\n", s.Truncate(10, "-»"))
	fmt.Printf("                          %s\n", s.NumberBanner())

	s = paw.NewStrChain("abc中defghijklmnopqrstuvwxyz")
	fmt.Printf("\n%#v\n", s)
	fmt.Printf("46. Wrap(10)  =\n%v\n", s.Wrap(10))
	fmt.Printf("%s\n", "0123456789")
}

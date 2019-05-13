// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	i18n "github.com/qioalice/i18n"
)

//
type tCtxKeyboardGenerator struct {
	loc *i18n.Locale
}

//
func (kbgen *tCtxKeyboardGenerator) Delete() *tCtxKeyboardGenerator {

}

// 'Text' defines what text will be displayed on keyboard (inline or not)
// buttons and provides a lots of ways to declare desired behaviour.
//
// Generally, there are 3 groups of behaviour declaration:
// 1. Display arguments as is (with casting to text).
// 2. Using i18n and l10n, with possibility to deep generation and
// use a whole translate node as text source.
// 3. Use proxy generator as addition to any of the ways above.
// You can combine all these categories!
//
// 1. Display argument as is:
// Just pass arguments of any type except 'func' or any strings starts
// with '$' or '$$' symbol.
// Each non-string argument will be casted to the string using:
//  - String() method, if object implements 'Stringer' interface
//  - fmt.Sprintf("%+v",) call, in all other cases
//
// 2. Using i18n and l10n (translate way):
// As you know, i18n library of this project works through receiving
// "translation key" and "translation args".
// So, you can generate button's text in this way or use all phrases from
// a whole "translation node", or only specified or all except some.
// Very flexible, right?
//
// 2.0. Used locale object (i18n.Locale):
// By default the locale object of the 'user.User' will be used
// to the recognize the translation.
// But if you want to use some another locale, just pass it AS FIRST ARGUMENT!
// It's important! If you willn't pass it as first but as ?'th arg,
// it will be ignored!
// For example, if you want to use Default locale to generate button's text:
// Example:
//     g.Text(i18n.DefLC(), ...)
//
// 2.1. Use values of some phrases using i18n lib:
// Just pass "translation key" but with '$' symbol before.
// You can pass as many these arguments as you want. Each that argument
// will be treated as a separate button (text for a separate button).
// Example:
//     // Text for two buttons; using i18n phrases as source with
//     // keys "root.cat1.val1" and "root.cat1.val2".
//     g.Text("$root.cat1.val1", "$root.cat1.val2")
//
// 2.2. With "translation args":
// Just pass 'i18n.Args' object after "translation key" and that object
// will be used to generate translation.
// Example:
//     // '123' will be substituted instead '{{key}}' placeholder in
//     // the "root1.cat1.valargs" value phrase.
//     g.Text("$root.cat1.valargs", i18n.Args{ "key": 123 })
//
// 2.3. The short way to specify many phrases from one category
// Just pass '.' after "translation node" name (category), before
// "translation key" and then you can pass just a "translation key"s names
// (w/o "translation node" head)
// if you want to specify what phrases should be used.
// Example:
//     // Text for two buttons; using i18n phrases as source with
//     // keys "root.cat1.val1" and "root.cat1.val2".
//     g.Text("$root.cat1.", "$val1", "$val2")
//
// Of course, you can use more than one category:
//     // Text for two buttons; using i18n phrases as source with
//     // keys "root.cat1.val1", "root.cat1.val2" and "root.cat2.val1".
//     g.Text("$root.cat1.", "$val1", "$val2", "$root.cat2.", "$val1")
//
// 2.4. Want to use absolutely all phrases from "translation node"? Do it.
// Just pass a special value 'ALL' as "translation key" of desired
// "translation node".
// Example:
//     // All phrases from "root.cat1" will be used as buttons.
//     g.Text("$root.cat1.ALL")
//
// 2.4.1. Want arguments? No problem. Pass not more than one 'i18n.Args'
// object (another will be ignored) which will be used as arguments for
// each phrases.
// Example:
//     // '42' will be substituted instead '{{key}}' placeholder in
//     // each "root.cat1" phrases.
//     g.Text("$root.cat1.ALL", i18n.Args{ "key": 42 })
//
// 2.4.2. Want more flexible way for arguments?
// For example, arguments for each phrase? Pass a callback that takes
// a two args: int as number of phrase in category, and string as phrase,
// and return an 'i18n.Args' object.
// Example:
//     // If "root1.cat1" node has the following phrases:
//     // { "key1": "val with {{key}}", "key2": "val with {{key}}" },
//     // then the following button text will be generated:
//     // "val with 0", "val with 1".
//     // (as you see, the numbering of 'idx' starts from zero).
//     g.Text("$root.cat1.ALL", func(idx int, phrase string) i18n.Args {
//         return i18n.Args{ "key": idx }
//     })
//
// 2.4.3. Combine 2.4.1 and 2.4.2? Yes, 2.4.1 - default args, 2.4.2 - special.
// Example:
//     // If "root1.cat1" node has the following phrases:
//     // { "key1": "{{text}} aaa {{key}}", "key2": "{[text}} bbb {{key}}" },
//     // then the following button text will be generated:
//     // "val with aaa 0", "val with bbb 1".
//     g.Text("$root.cat1.ALL", i18n.Args{ "text": "val with" },
//            func(idx int, phrase string) i18n.Args {
//                return i18n.Args { "key": idx }
//            },
//     )
//
// 2.5. Want to use absolutely all phrases from "translation node"
// but except for a few?
// Just pass a special value 'ALL_EXCEPT' as "translation key" of desired
// "translation node".
// Example:
//     // All phrases from "root.cat1" will be used as buttons,
//     // because there is no exceptions
//     g.Text("$root.cat1.ALL_EXCEPT")
//     // But if "root.cat1" node has the phrases with following keys:
//     // ["key1", "key2", "key3", key4"] and you want all except the "key1":
//     g.Text("$root.cat1.ALL_EXCEPT", "$key1")
//
// 2.5.1, 2.5.2, 2.5.3: The same as 2.4.1, 2.4.2, 2.4.3.
//
// 3. Proxy generators.
// Button text generated using 1 or 2 category rule but you want to change it?
// You can pass a special callback: that takes and returns only one argument:
// the string. Takes the string which would be used as button text and
// you should return a string with really button text that you want to use.
//
// 3.1. Proxy generator for one button:
// Pass callback after
func (kbgen *tCtxKeyboardGenerator) Text(buttons ...interface{}) *tCtxKeyboardGenerator {
	// todo: check if object is valid
	if len(buttons) == 0 {
		return kbgen
	}
	// Check whether the first argument is desired locale
	if loc, ok := buttons[0].(*i18n.Locale); ok && loc != nil {
		kbgen.loc, buttons = loc, buttons[1:]
	}
	// Declare a variables for "translation node", "translation key"
	trnode, trkey := "", ""
	// Parse the remaining arguments until the argument list is empty
	for len(buttons) != 0 {
		// If it's just a text, add it
		if txt := kbgen.extractText(buttons[0]); txt != "" {
			kbgen.texts = append(kbgen.texts, txt)
			continue // todo: i++
		}
		// If it's a something of "translation key" - part or node name or
		// full translation key
		if key := kbgen.extractTrKey(buttons[0]); key != "" {

		}
	}
	oneBtnText, restRules := kbgen.parseText(buttons)
	for oneBtnText != "" && len(restRules) > 0 {
		kbgen.texts = append(kbgen.texts, oneBtnText)
		oneBtnText, restRules = kbgen.parseText(buttons)
	}
	return kbgen
}

//
func (kbgen *tCtxKeyboardGenerator) LocNode(key string, values ...interface{}) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) EachLocNode(key string, values ...interface{}) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) LocTr(key string) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) EachLocTr(key string) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) URL(links ...interface{}) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) Actions(args ...interface{}) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) AsTpl(key interface{}, args ...bool) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) FromTpl(key interface{}, args ...bool) *tCtxKeyboardGenerator {

}

//
func (kbgen *tCtxKeyboardGenerator) parseText(in []interface{}) (out string, rest []interface{}) {
	if len(in) == 0 {
		return "", nil
	}

}

//
func (*tCtxKeyboardGenerator) extractText(i interface{}) (text string) {

}

//
func (*tCtxKeyboardGenerator) extractTrKey(i interface{}) (key string) {

}
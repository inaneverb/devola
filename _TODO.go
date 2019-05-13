package tgbot

// // 'tRegistrableEvent' is the part of 'tReger'.
// // When you calling one of registration methods 'Text', 'Button', 'Command' or
// // 'InlineButton', these calls creates 'tRegistrableEvent' object and store
// // it in the 'tReger' object.
// // Later, when you will call 'Handler' or 'Middleware' methods,
// // these methods will apply callback from argument to each stored
// // 'tRegistrableEvent' event object.
// type tRegistrableEvent struct {
// 	typ  tEventType
// 	data string
// }

// refactor reflect.Call feature to unsafe.Pointer()() in tViewBaseFinisher

// add feature "Big brother"
// (see all conference with bot for chats)

// add feature: "Take control"
// (emulate user interaction with bot for some user)

// add feature: Make tRegistrator thread-safe

// tReceiver should decode ikb data, tEvent.Data should be tViewID

// add AutoRegister param to tViewIDConverter
// (auto register view id while trying encode if view id isn't registered)

// rename all "parent" to more readable names

// snippet of tRegistrator.match calling:
// if ctx != nil {
// 	typ = ctx.Event.Type
// 	data = ctx.Event.Data
// 	if ctx.sess != nil {
// 		viewID = ctx.sess.ViewIDEncoded
// 	}
// }

// all errors should be an object and moved to separate package

// rework error mechanism
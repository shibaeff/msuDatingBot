package bot

//import "echoBot/pkg/models"
//
//func (b *bot) next(user *models.User) (reply interface{}) {
//	// check if register is over
//	// TODO registration controller
//	if user.RegiStep < regOver {
//		raw := replyWithText(notRegistered)
//		raw.ChatID = user.Id
//		reply = raw
//		return
//	}
//	// get unseen users for current user
//	unseen, e := b.store.GetUnseen(user.Id)
//	if len(unseen) == 0 || e != nil {
//		reply = replyWithText("Вы просмотрели всех пользователей на данный момент")
//		return reply
//	}
//	unseen_user, _ := b.store.GetUser(unseen[0].Whome)
//	b.actionsLog.Printf("%d VIEWED %d\n", user.Id, unseen_user.Id)
//	card := replyWithCard(unseen_user, user.Id)
//	card.ParseMode = "html"
//	reply = card
//	b.actionsLog.Printf("%d VIEWED %d", user.Id, unseen_user.Id)
//	return
//}
//
//func (b *bot) dislike(user *models.User) (reply interface{}) {
//	// remove last user from the seen directory
//	unseen, _ := b.store.GetUnseen(user.Id)
//	b.store.GetUnseenRegistry().DeleteItem(user.Id, unseen[0].Whome)
//	reply = b.next(user)
//	return
//}

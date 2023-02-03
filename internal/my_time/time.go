//package my_Moment
//
//type Moment struct {
//	seconds, minutes, hours, days, years int
//}
//
//func (Moment *Moment) CheckMomentFormat() {
//	Moment.minutes += Moment.seconds / 60
//	Moment.seconds %= 60
//
//	Moment.hours += Moment.minutes / 60
//	Moment.minutes %= 60
//
//	Moment.days += Moment.hours / 24
//	Moment.hours %= 24
//
//	Moment.years += Moment.days % 365 // Можно добавить обработку високосных годов
//	Moment.days %= 365
//}
//
//func (Moment *Moment) To() int
//
//func (Moment *Moment) ToSeconds() int {
//	return Moment.seconds + Moment.minutes*60
//}
//
//func NewMoment(seconds, minutes, hours, days, years int) Moment {
//	var Moment = Moment{seconds, minutes, hours, days, years}
//	Moment.CheckMomentFormat()
//	return Moment
//}

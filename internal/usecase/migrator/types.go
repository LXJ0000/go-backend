package migrator

type Entity interface {
	ID() int64
	Compare(item Entity) bool
}

// func Compare() {
// 	if !reflect.DeepEqual(b, t) {
// 		to do something
// 	}
// }
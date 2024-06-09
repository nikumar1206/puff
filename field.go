package puff

type Param[T any] struct {
	Description string
	Kind        string //Header, Path, Query, Body
	Value       T
}

func HeaderParam[T any](description string) Param[T] {
	return Param[T]{
		Description: description,
		Kind:        "Header",
	}
}

func PathParam[T any](description string) Param[T] {
	return Param[T]{
		Description: description,
		Kind:        "Path",
	}
}

func QueryParam[T any](description string) Param[T] {
	return Param[T]{
		Description: description,
		Kind:        "Query",
	}
}

func BodyParam[T any](description string) Param[T] {
	return Param[T]{
		Description: description,
		Kind:        "Header",
	}
}

//DREAM CODE WITH THIS:

// func HandleRoute(params Context[HelloWorld]) {
// 	params.drinks.Value[0] //has correct type of string
// }
//
// type HelloWorld struct {
// 	food   Param[string]
// 	drinks Param[[]string]
// 	price  Param[int]
// }
//
// func main() {
// 	params_on_my_route := HelloWorld{
// 		food:   QueryParam[string]("The food you have ordered."),
// 		drinks: BodyParam[[]string]("The drinks you have ordered."),
// 		price:  QueryParam[int]("The price that you would like to charge out of your card."),
// 	}
// }

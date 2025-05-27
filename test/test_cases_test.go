package routes__test

import "net/http"

type endPointTest struct {
	name     string
	method   string
	path     string
	headers  kvp
	cookies  kvp
	body     string
	eStatus  int
	eHeaders kvp
	eBody    string
	dbres    []string
}

var tests = []endPointTest{
	// {"TEST NAME   ",
	// 	"METHOD", "/URL/PATH&querry=value", []kvp{{"Header Name": "Header Value"}}, kvp{"Cookie Name": "Cookie Value"},`Request Body`,
	// 	http.ExpectedStatus, []kvp{{"Expected Header Name": "Expected Header Value"}}, `Expected Body`} ,
	//  []string{"Database Query", "Database Result Row1", "Database Result Row2"}},
	// }
	{
		"GET root 1   ", "GET", "/",
		kvp{},
		kvp{},
		``,
		http.StatusOK,
		kvp{},
		`hello world`,
		nil,
	},
	// TODO: CHANGE AUTH COOKIE to JWT
	// {
	// 	"Isloged in false no uId  ",
	// 	"POST", "/api/v1/auth",
	// 	kvp{},
	// 	kvp{},
	// 	``,
	// 	http.StatusUnauthorized,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"isLoggedIn":false}` + "\n",
	// 	nil,
	// },

	// {
	// 	"Isloged in false bad cookie name ",
	// 	"POST", "/api/v1/auth",
	// 	kvp{},
	// 	kvp{"userid": "6352337196a2449cb772b524818bea36"},
	// 	``,
	// 	http.StatusUnauthorized,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"isLoggedIn":false}` + "\n",
	// 	nil,
	// },

	// {
	// 	"Isloged in false no uid value ",
	// 	"POST", "/api/v1/auth",
	// 	kvp{},
	// 	kvp{"userId": ""},
	// 	``,
	// 	http.StatusUnauthorized,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"isLoggedIn":false}` + "\n",
	// 	nil,
	// },

	// {
	// 	"Isloged in false bad uid value ",
	// 	"POST", "/api/v1/auth",
	// 	kvp{},
	// 	kvp{"userId": "6352337196a2449cb772b524818bea37"},
	// 	``,
	// 	http.StatusUnauthorized,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"isLoggedIn":false}` + "\n",
	// 	nil,
	// },

	// {
	// 	"Isloged in true   ",
	// 	"POST", "/api/v1/auth",
	// 	kvp{},
	// 	kvp{"userId": "6352337196a2449cb772b524818bea36"},
	// 	``,
	// 	http.StatusOK,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"isLoggedIn":true}` + "\n",
	// 	nil,
	// },

	// {
	// 	"RegisterHandler bad allready logged in  ",
	// 	"POST", "/api/v1/auth/register",
	// 	kvp{},
	// 	kvp{"userId": "6352337196a2449cb772b524818bea36"},
	// 	`{
	// 		"username" : "","email" : "","password" : "","age" : "","gender" : "",
	// 		"fname" : "","lname" : "","birthdate" : "","avatar" : "",
	// 		"aboutme" : "","status" : ""
	// 		}`,
	// 	http.StatusForbidden,
	// 	kvp{"Content-Type": "application/json"},
	// 	`{"error": "Already logged in"}`,
	// 	nil,
	// },
	// TODO: CHANGE AUTH COOKIE to JWT

	{
		"RegisterHandler bad empty body ",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		``,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "Request body cannot be empty"}`,
		nil,
	},
	{
		"RegisterHandler bad empty feilds ",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{ "username": "", "email": "", "birthdate": "", "fname": "", "lname": "", "password": "", "gender": "" }`,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "empty feild (birthdate) expected format (YYYY-MM-DD)"}`,
		nil,
	},
	{
		"RegisterHandler bad extra feild ",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{ "extra": "" }`,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "json: unknown field \"extra\""}`,
		nil,
	},
	{
		"RegisterHandler bad feild  value 1",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{ "username": "a",
		"email": "example@web.site","birthdate": "2001-11-09","fname": "Firstname","lname": "LastName","password": "password","gender": "DFK"}`,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "field 'username' is invalid: invalid charachters used"}`,
		nil,
	},
	{
		"RegisterHandler bad feild  value 2",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{"username": "User_Name",
		"email": "example@website",
		"birthdate": "2001-11-09","fname": "Firstname","lname": "LastName","password": "password","gender": "DFK" }`,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "field 'email' is invalid: malformed or missing data"}`,
		nil,
	},
	{
		"RegisterHandler bad feild  value 3",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{"username": "User_Name","email": "example@website",
"birthdate": "2001-13-09",
		"fname": "Firstname", "lname": "LastName", "password": "password", "gender": "DFK"}`,
		http.StatusBadRequest,
		kvp{"Content-Type": "application/json"},
		`{"error": "invalid birthdate format (expected YYYY-MM-DD): parsing time \"2001-13-09\": month out of range"}`,
		nil,
	},
	{
		"RegisterHandler good ",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{"username": "User_Name","email": "example@web.site","birthdate": "2001-11-09","fname": "Firstname", "lname": "LastName", "password": "password", "gender": "DFK"}`,
		http.StatusOK,
		kvp{"Content-Type": "application/json"},
		`{"message":"User registered successfully"}`,
		[]string{
			"select username, email, birthdate, password, gender, fname, lname from users where username = 'User_Name';",
			"username=User_Name, email=example@web.site, birthdate=2001-11-09 00:00:00 +0000 UTC, password=password, gender=2, fname=Firstname, lname=LastName",
		},
	},
	{
		"RegisterHandler bad all ready exist 1",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{"username": "User_Name","email": "example@web.site","birthdate": "2001-11-09","fname": "Firstname", "lname": "LastName", "password": "password", "gender": "DFK"}`,
		http.StatusConflict,
		kvp{"Content-Type": "application/json"},
		`{"error": "sorry email already exists"}`,
		nil,
	},
	{
		"RegisterHandler bad all ready exist 2",
		"POST", "/api/v1/auth/register",
		kvp{},
		kvp{},
		`{"username": "User_Name","email": "email@web.site","birthdate": "2001-11-09","fname": "Firstname", "lname": "LastName", "password": "password", "gender": "DFK"}`,
		http.StatusConflict,
		kvp{"Content-Type": "application/json"},
		`{"error": "sorry username already exists"}`,
		nil,
	},
}

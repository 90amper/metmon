package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// func (s *MemStorage) Stringer() string {}

var MemStor MemStorage = MemStorage{
	gauge:   make(map[string]float64),
	counter: make(map[string]int64),
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", collectorHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "POST" && r.Header.Get("Content-Type") == "text/plain" {
	if r.Method == "POST" {
		// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		path := strings.Split(r.URL.Path, "/")

		if len(path) < 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if path[2] == "counter" {
			vali, err := strconv.ParseInt(path[4], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			MemStor.counter[path[3]] += vali
		} else if path[2] == "gauge" {
			valf, err := strconv.ParseFloat(path[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			MemStor.gauge[path[3]] = valf
		} else {
			println(path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("%+v\r\n", MemStor)

		// test := strings.Join(path, " | ")

		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("test: " + test))
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// func JSONHandler(w http.ResponseWriter, req *http.Request) {
// 	// собираем данные
// 	subj := Subj{"Milk", 50}
// 	// кодируем в JSON
// 	resp, err := json.Marshal(subj)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	// устанавливаем заголовок Content-Type
// 	// для передачи клиенту информации, кодированной в JSON
// 	w.Header().Set("content-type", "application/json")
// 	// устанавливаем код 200
// 	w.WriteHeader(http.StatusOK)
// 	// пишем тело ответа
// 	w.Write(resp)
// }

// const form = `<html>
//     <head>
//     <title></title>
//     </head>
//     <body>
//         <form action="/" method="post">
//             <label>Логин</label><input type="text" name="login">
//             <label>Пароль<input type="password" name="password">
//             <input type="submit" value="Login">
//         </form>
//     </body>
// </html>`

// func Auth(login, password string) bool {
// 	return login == `guest` && password == `demo`
// }

// func mainPage(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		login := r.FormValue("login")
// 		password := r.FormValue("password")
// 		// http.Header
// 		if Auth(login, password) {
// 			io.WriteString(w, "Добро пожаловать!")
// 		} else {
// 			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
// 		}
// 		return
// 	} else {
// 		io.WriteString(w, form)
// 	}
// }

// func apiPage(res http.ResponseWriter, req *http.Request) {
// 	// res.Write([]byte("Это страница /api."))
// 	body := fmt.Sprintf("Method: %s\r\n", req.Method)
// 	body += "Header ===============\r\n"
// 	for k, v := range req.Header {
// 		body += fmt.Sprintf("%s: %v\r\n", k, v)
// 	}
// 	body += "Query parameters ===============\r\n"
// 	if err := req.ParseForm(); err != nil {
// 		res.Write([]byte(err.Error()))
// 		return
// 	}
// 	for k, v := range req.Form {
// 		body += fmt.Sprintf("%s: %v\r\n", k, v)
// 	}
// 	res.Write([]byte(body))

// }

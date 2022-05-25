package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/hello" {
		http.Error(w, "OOOO", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
	}

	fmt.Println("heelo handler")

}

func formHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Parse Form err: %v", err)
	}
	fmt.Fprintf(w, "Post Request successful\n")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprint(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(1 << 20)

	file, handler, err := r.FormFile("myFile")
	fmt.Println("gogogogoogogogogo")
	fmt.Println(reflect.TypeOf(file))
	fmt.Println(reflect.TypeOf(handler))
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("./uploads", "upload-*.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	var encoded []byte = HuffmanEncode(fileBytes)

	fmt.Printf(string(encoded))
	// return encoded
}

// func main() {
// 	fmt.Println("Hello, World!")
// 	fileServer := http.FileServer(http.Dir("./static"))
// 	http.Handle("/", fileServer)
// 	http.HandleFunc("/upload", uploadFile)
// 	http.HandleFunc("/form", formHandler)
// 	http.HandleFunc("/hello", helloHandler)

// 	fmt.Printf("Starting server at port 8080\n")

// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

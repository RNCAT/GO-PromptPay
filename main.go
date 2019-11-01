package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/carlescere/scheduler"
	"log"
	"net/http"
	"os"
	"strconv"

	pp "github.com/Frontware/promptpay"
	"github.com/skip2/go-qrcode"
	"html/template"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))
var fs = http.FileServer(http.Dir("static/"))

func main()  {
	_, _ = scheduler.Every(7).Day().Run(delFile)
	log.Println("Server started at : http://localhost")
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/",home)
	http.HandleFunc("/promptpay",promtpay)
	_ = http.ListenAndServe(":80", nil)
}

func home(w http.ResponseWriter, r *http.Request)  {
	_ = tmpl.ExecuteTemplate(w, "index", nil)
}

func promtpay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST"{
		http.Redirect(w,r,"/",301)
	}
	id := r.FormValue("id")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"),64)
	payment := pp.PromptPay{
		PromptPayID: id,
		Amount:      amount,
	}
	promtPayCode, err := payment.Gen()
	if err != nil {
		_ = tmpl.ExecuteTemplate(w, "error", nil)
	}else {
		hash := md5.New()
		hash.Write([]byte(id))
		filename := hex.EncodeToString(hash.Sum(nil))
		log.Print(promtPayCode)
		_ = qrcode.WriteFile(promtPayCode, qrcode.Medium, 256, "./static/results/"+filename+".png")

		_ = tmpl.ExecuteTemplate(w,"qrshow",filename)
	}
}

func delFile()  {
	_ = os.RemoveAll("./static/results")
	log.Println("Delete results dir complete")
	_ = os.MkdirAll("./static/results",0700)
	log.Println("Create results dir complete")
}

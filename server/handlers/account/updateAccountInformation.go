package account

import "net/http"

func UpdateAccountInformation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bka"))
}

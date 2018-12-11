package main

type (
	// login / auth
	ReqAuth struct {
		UserID      string `json:"user_id"`
		MacAddress1 string `json:"mac_address_1"`
		MacAddress2 string `json:"mac_address_2"`
		ReleaseSn   string `json:"release_sn"`
	}

	RespItem struct {
		Me    Terminal `json:"me"`
		Token string   `json:"token"`
	}

	Terminal struct {
		UserID       string `json:"user_id"`
		RegisterDate string `json:"register_date"`
		Expiration   string `json:"expiration"`
	}
)

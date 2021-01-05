package look_and_feel

type Model struct {
	LoginNameProject string `json:"Login_name_project"`
	LoginVersion     string `json:"Login_version"`
	LoginSlogan      string `json:"Login_slogan"`
	LoginLogo        string `json:"Login_logo"`
	MenuLogo         string `json:"menu_url_logo"`
	FooterLogo       string `json:"footer_url_logo"`
	Primary          string `json:"primary"`
	Secondary        string `json:"secondary"`
	Tertiary         string `json:"tertiary"`
	ID               string `json:"id"`
	Key              string `json:"key"`
}

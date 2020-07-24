package appConfig

type Mail struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	User string   `json:"user"`
	Pass string   `json:"pass"`
	To   []string `json:"to"` // 收件人
}

type SLSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Project         string `json:"project"`
	LogStore        string `json:"log_store"`
	Topic           string `json:"topic"`
}

type DataBase struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DataBase string `json:"database"`
}

type Redis struct {
	NetWork  string `json:"net_work"`
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DataBase int    `json:"data_base"`
	Prefix   string `json:"prefix"` // 前缀
}

type Oss struct {
	EndPoint        string `json:"end_point"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	BucketUrl       string `json:"bucket_url"`
	CallbackAddr    string `json:"callback_addr"`
	UploadDir       string `json:"upload_dir"`
}

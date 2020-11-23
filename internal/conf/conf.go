package conf

// This is supposed to validate, test and parse the configuration file

type Conf struct {
	user          string
	group         string
	defaultServer *ServerConf
	vhosts        []ServerConf
}

type ServerConf struct {
	name       string
	root       string
	ports      []int
	indexPages []string
	errorPages []ErrorPage
	errorLog   string
	accessLog  string
}

type ErrorPage struct {
	code int
	page string
}

func Load() error {
	// TODO
	return nil
}
